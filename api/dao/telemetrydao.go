package dao

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Drakkar-Software/Metrics-Server/api/model"
	bot "github.com/Drakkar-Software/Metrics-Server/api/model"
	"github.com/Drakkar-Software/Metrics-Server/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrBotNotFound is returned new a bot is not found in database
var ErrBotNotFound = errors.New("bot not found")

// ErrInvalidData is returned new a bot is invalid (ex: CurrentSession.StartedAt == 0)
var ErrInvalidData = errors.New("invalid data")

// CompleteGetBots returns all data about all bots
func CompleteGetBots(history bool) (bot.Bots, error) {
	return fetchBots(false, history, 0)
}

// PublicGetCountBots returns the number of active bot until time
func PublicGetCountBots(untilTime int64) (int64, error) {
	filter := bson.M{"$expr": bson.M{"$gt": bson.A{
		bson.M{"$add": bson.A{
			"$currentSession.startedAt",
			"$currentSession.upTime"}},
		untilTime}}}

	count, err := database.Database.Collection.CountDocuments(context.Background(), filter)

	if err != nil {
		return 0, err
	}

	return count, err
}

// UpdateBotUptimeAndProfitability updates bot in argument upTime and profitability (using BotID)
func UpdateBotUptimeAndProfitability(uploadedBot *bot.Bot) (interface{}, error) {
	collection := database.Database.Collection
	filter := bson.M{"_id": uploadedBot.ID}
	update := bson.M{"$set": bson.M{
		"currentSession.upTime":        uploadedBot.CurrentSession.UpTime,
		"currentSession.profitability": uploadedBot.CurrentSession.Profitability,
		"currentSession.tradedVolumes": uploadedBot.CurrentSession.TradedVolumes,
	},
	}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount != 1 {
		return nil, ErrBotNotFound
	}
	return uploadedBot.ID, err
}

// RegisterOrUpdate updates a bot if already in database or registers a new bot called after few minutes a bot is running
func RegisterOrUpdate(uploadedBot *bot.Bot) (interface{}, error) {
	// Allow a post for 24h in the future but not more
	maxTime := time.Now().AddDate(0, 0, 1).Unix()
	if uploadedBot.CurrentSession.StartedAt == 0 || uploadedBot.CurrentSession.UpTime == 0 ||
		int64(uploadedBot.CurrentSession.StartedAt+uploadedBot.CurrentSession.UpTime) > maxTime {
		return nil, ErrInvalidData
	}
	collection := database.Database.Collection
	// check if bot exists
	var foundBot bot.Bot
	filter := bson.D{{"_id", uploadedBot.ID}}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	if count < 1 {
		return nil, ErrBotNotFound
	}
	err = collection.FindOne(context.Background(), filter).Decode(&foundBot)
	if err != nil {
		return nil, err
	}
	return registerNewBotSession(uploadedBot, &foundBot)
}

// GenerateBotID generates a new bot id
func GenerateBotID() (interface{}, error) {
	newID, err := database.Database.Collection.InsertOne(context.Background(), bot.Bot{})
	if err != nil {
		return nil, err
	}
	log.Println("Registered bot with id:", newID.InsertedID)
	return newID.InsertedID, nil
}

func registerNewBotSession(uploadedBot *bot.Bot, foundBot *bot.Bot) (interface{}, error) {
	newHistoricalSession := foundBot.CurrentSession
	foundBot.CurrentSession = uploadedBot.CurrentSession
	update := bson.M{
		"$set": bson.D{{
			"currentSession", foundBot.CurrentSession,
		}},
		"$push": bson.D{{
			"sessionHistory", newHistoricalSession,
		}},
	}
	if len(foundBot.SessionHistory) == 0 {
		if newHistoricalSession.UpTime > 0 {
			foundBot.SessionHistory = append(foundBot.SessionHistory, newHistoricalSession)
		}
		update = bson.M{
			"$set": bson.D{
				{
					"currentSession", foundBot.CurrentSession,
				},
				{
					"sessionHistory", foundBot.SessionHistory,
				},
			},
		}
	}
	filter := bson.D{{"_id", uploadedBot.ID}}
	updateResult, err := database.Database.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount != 1 {
		return nil, ErrBotNotFound
	}
	log.Println("Registed new session for bot with id:", uploadedBot.ID)
	return uploadedBot.ID, err
}

func fetchBots(filterBots bool, includeHistory bool, since int64) (bot.Bots, error) {
	bots := bot.Bots{}
	filter := bson.M{}
	if since > 0 {
		filter = bson.M{"$expr": bson.M{"$gt": bson.A{
			bson.M{"$add": bson.A{
				"$currentSession.startedAt",
				"$currentSession.upTime"}},
			since}}}
	}
	cur, err := database.Database.Collection.Find(context.Background(), filter)
	if err != nil {
		return bots, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var decodedBot bot.Bot
		err := cur.Decode(&decodedBot)
		if err != nil {
			panic(err)
		}
		if filterBots {
			decodedBot.FilterPublicInfo()
		}
		if !includeHistory {
			decodedBot.SessionHistory = nil
		}
		bots = append(bots, decodedBot)
	}
	return bots, err
}

func FetchTop(field string, traderType string, since int64, minimal int64, limit int64, onlyFirstOfArray bool) []bot.Top {
	// count and order elements from arrays in bots documents, example with field=exchanges, traderType=Simulator since=1669853830, minimal=3, limit=100
	// db.metrics.aggregate([
	// 	{ $match: { $expr: {$and: [{$eq: ["$currentSession.trader", false]}, {$eq: ["$currentSession.simulator", true]}, {$gt: [{ $add: ["$currentSession.startedAt", "$currentSession.upTime"] }, 1669853830] }]} } },
	//  If full array:
	// 		{ $unwind: "$currentSession.exchanges" },
	// 		{ $group: { _id: "$currentSession.exchanges", count: { $sum: 1 } } },
	//  If onlyFirstOfArray:
	//  	{ $group: { _id: { $first: "$currentSession.evalConfig" }, count: { $sum: 1 } } }
	// 	{ $match: {$and: [{count: {"$gte": 3}}, {_id: {"$ne": null} }]}},
	// 	{ $sort: { count: -1, _id: 1 } },
	// 	{ $limit: 100 }
	// ])

	fieldKey := "$currentSession." + field
	timeSelector := bson.D{{"$gt", bson.A{bson.D{{"$add", bson.A{"$currentSession.startedAt", "$currentSession.upTime"}}}, since}}}
	timeAndTraderSelector := timeSelector
	if traderType != model.AllTraders {
		realTraderSelector := bson.D{{"$eq", bson.A{"$currentSession.trader", traderType == model.RealTraders}}}
		simulatedTraderSelector := bson.D{{"$eq", bson.A{"$currentSession.simulator", traderType == model.SimulatedTraders}}}
		timeAndTraderSelector = bson.D{{"$and", bson.A{realTraderSelector, simulatedTraderSelector, timeSelector}}}
	}
	matchStage := bson.D{{"$match", bson.D{{"$expr", timeAndTraderSelector}}}}
	unwindStage := bson.D{{"$unwind", fieldKey}}
	groupStage := bson.D{{"$group", bson.D{{"_id", fieldKey}, {"count", bson.D{{"$sum", 1}}}}}}
	if onlyFirstOfArray {
		groupStage = bson.D{{"$group", bson.D{{"_id", bson.D{{"$first", fieldKey}}}, {"count", bson.D{{"$sum", 1}}}}}}
	}
	matchOnCountStage := bson.D{{"$match", bson.D{{"$and", bson.A{bson.D{{"count", bson.D{{"$gte", minimal}}}}, bson.D{{"_id", bson.D{{"$ne", nil}}}}}}}}}
	sortStage := bson.D{{"$sort", bson.D{{"count", -1}, {"_id", 1}}}}
	limitStage := bson.D{{"$limit", limit}}

	var pipeline mongo.Pipeline
	if onlyFirstOfArray {
		pipeline = mongo.Pipeline{matchStage, groupStage, matchOnCountStage, sortStage, limitStage}
	} else {
		pipeline = mongo.Pipeline{matchStage, unwindStage, groupStage, matchOnCountStage, sortStage, limitStage}
	}

	cursor, err := database.Database.Collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		panic(err)
	}
	var fetchedTop []bot.Top
	if err = cursor.All(context.TODO(), &fetchedTop); err != nil {
		panic(err)
	}
	return fetchedTop
}
