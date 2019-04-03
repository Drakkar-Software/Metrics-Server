package dao

import (
	"context"
	"errors"
	"log"

	bot "github.com/Drakkar-Software/Metrics-Server/api/model"
	database "github.com/Drakkar-Software/Metrics-Server/database"

	"github.com/mongodb/mongo-go-driver/bson"
)

var db = database.DB{}

// Init initializes the database connection
func Init() error {
	return db.Initialize()
}

// GetBots returns all data about all bots
func GetBots() (bot.Bots, error) {
	bots := bot.Bots{}
	cur, err := db.Collection.Find(context.Background(), bson.D{{}})
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
		bots = append(bots, decodedBot)
	}
	return bots, err
}

// UpdateBotUptime updates bot in argument upTime (using BotID)
func UpdateBotUptime(uploadedBot *bot.Bot) (interface{}, error) {
	collection := db.Collection
	filter := bson.D{{"_id", uploadedBot.ID}}
	update := bson.D{{"$set",
		bson.D{{
			"currentSession.upTime", uploadedBot.CurrentSession.UpTime,
		}},
	}}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount != 1 {
		return nil, errors.New("No matching document with given id")
	}
	return updateResult.UpsertedID, err
}

// RegisterOrUpdate updates a bot if already in database or registers a new bot called after few minutes a bot is running
func RegisterOrUpdate(uploadedBot *bot.Bot) (interface{}, error) {
	collection := db.Collection
	// check if bot exists
	var foundBot bot.Bot
	filter := bson.D{{"_id", uploadedBot.ID}}
	err := collection.FindOne(context.Background(), filter).Decode(&foundBot)
	if err != nil {
		return nil, err
	}
	return registerNewBotSession(uploadedBot, &foundBot)
}

// GenerateBotID generates a new bot id
func GenerateBotID() (interface{}, error) {
	newID, err := db.Collection.InsertOne(context.Background(), bot.Bot{})
	if err != nil {
		return nil, err
	}
	log.Println("Registered bot with id:", newID.InsertedID)
	return newID.InsertedID, nil
}

func registerNewBotSession(uploadedBot *bot.Bot, foundBot *bot.Bot) (interface{}, error) {
	// add current session to session history (start new session)
	if foundBot.CurrentSession.UpTime > 0 {
		foundBot.SessionHistory = append(foundBot.SessionHistory, foundBot.CurrentSession)
	}
	foundBot.CurrentSession = uploadedBot.CurrentSession
	update := bson.D{{"$set",
		bson.D{
			{
				"currentSession", foundBot.CurrentSession,
			},
			{
				"sessionHistory", foundBot.SessionHistory,
			},
		},
	}}
	filter := bson.D{{"_id", uploadedBot.ID}}
	updateResult, err := db.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}
	if updateResult.MatchedCount != 1 {
		return nil, errors.New("No matching document with given id")
	}
	log.Println("Registed new session for bot with id:", uploadedBot.ID)
	return updateResult.UpsertedID, err
}
