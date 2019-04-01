package dao

import (
	"context"
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
func UpdateBotUptime(bot bot.Bot) error {
	collection := db.Collection
	filter := bson.D{{"botID", bot.BotID}}
	update := bson.D{{"$set",
		bson.D{{
			"currentSession.upTime", bot.CurrentSession.UpTime,
		}},
	}}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	if updateResult.MatchedCount != 1 {
		log.Println("None or more than one bot updated with BotID: ", bot.BotID, ": ", updateResult.MatchedCount)
	}
	return err
}

// RegisterOrUpdate updates a bot if already in database or registers a new bot
func RegisterOrUpdate(bot bot.Bot) error {
	// TODO
	return nil
}
