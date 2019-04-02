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
func UpdateBotUptime(uploadedBot *bot.Bot) error {
	collection := db.Collection
	filter := bson.D{{"botID", uploadedBot.BotID}}
	update := bson.D{{"$set",
		bson.D{{
			"currentSession.upTime", uploadedBot.CurrentSession.UpTime,
		}},
	}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return err
}

// RegisterOrUpdate updates a bot if already in database or registers a new bot called after few minutes a bot is running
func RegisterOrUpdate(uploadedBot *bot.Bot) error {
	collection := db.Collection
	// check if bot exists
	var foundBot bot.Bot
	filter := bson.D{{"botID", uploadedBot.BotID}}
	count, err := collection.Count(context.Background(), filter)
	if err != nil {
		return err
	}
	if count > 0 {
		// update existing bot
		err := collection.FindOne(context.Background(), filter).Decode(&foundBot)
		if err != nil {
			return err
		}
		return registerNewBotSession(uploadedBot, &foundBot)
	} else {
		// create new bot
		return createNewBot(uploadedBot)
	}
}

func createNewBot(uploadedBot *bot.Bot) error {
	_, err := db.Collection.InsertOne(context.Background(), uploadedBot)
	if err != nil {
		return err
	}
	log.Println("Added bot with id:", uploadedBot.BotID)
	return nil
}

func registerNewBotSession(uploadedBot *bot.Bot, foundBot *bot.Bot) error {
	// add current session to session history (start new session)
	foundBot.SessionHistory = append(foundBot.SessionHistory, foundBot.CurrentSession)
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
	filter := bson.D{{"botID", uploadedBot.BotID}}
	_, err := db.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	log.Println("Registed new session for bot with id:", uploadedBot.BotID)
	return nil
}
