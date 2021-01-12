package dao

import (
	"context"
	"errors"
	"log"
	"time"

	bot "github.com/Drakkar-Software/Metrics-Server/api/model"
	"github.com/Drakkar-Software/Metrics-Server/database"

	"go.mongodb.org/mongo-driver/bson"
)

// ErrBotNotFound is returned new a bot is not found in database
var ErrBotNotFound = errors.New("bot not found")

// ErrInvalidData is returned new a bot is invalid (ex: CurrentSession.StartedAt == 0)
var ErrInvalidData = errors.New("invalid data")

// PublicGetBots returns filtered data about all bots
func PublicGetBots(since int64) (bot.Bots, error) {
	return fetchBots(false, false, since)
}

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
	if uploadedBot.CurrentSession.StartedAt == 0 || uploadedBot.CurrentSession.UpTime == 0 {
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
