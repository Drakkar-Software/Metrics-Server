package dao

import (
	"context"

	"github.com/Drakkar-Software/Metrics-Server/api/model"
	"github.com/Drakkar-Software/Metrics-Server/database"
	"github.com/mongodb/mongo-go-driver/bson"
)

// IsAuthorizedUser returns true if the given user (identified by API-Key) is authorized to get this level of info
func IsAuthorizedUser(apiKey string, accessLevel int8) bool {
	if len(apiKey) > 6 {
		collection := database.Database.UserCollection
		filter := bson.M{"apiKey": apiKey}
		var foundUser model.User
		err := collection.FindOne(context.Background(), filter).Decode(&foundUser)
		if err != nil {
			return false
		}
		return foundUser.AccessLevel >= accessLevel
	}
	return false
}
