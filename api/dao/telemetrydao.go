package dao

import (
	"context"
	"log"

	bot "github.com/Drakkar-Software/Metrics-Server/api/model"
	database "github.com/Drakkar-Software/Metrics-Server/database"

	"github.com/mongodb/mongo-go-driver/bson"
)

func All() (bot.Bots, error) {
	bots := bot.Bots{}

	db := database.DB{}
	err := db.Initialize()
	if err != nil {
		return bots, err
	}
	defer db.Close()

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

func UpdateBotUptime(bot bot.Bot) error {
	db := database.DB{}
	err := db.Initialize()
	if err != nil {
		return err
	}
	defer db.Close()
	collection := db.Collection
	filter := bson.D{{"botID", bot.BotID}}
	update := bson.D{{"$set",
		bson.D{{
			"upTime", bot.UpTime,
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
