package model

import (
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Bot stores usage info about a specific bot identified by BotID
type Bot struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BotID     string             `json:"botID,omitempty" bson:"botID"`
	CreatedAt int                `json:"createdAt,omitempty" bson:"createdAt"`
	UpTime    int                `json:"upTime,omitempty" bson:"upTime"`
}

// GenerateBotID generates a new bot id
func GenerateBotID() string {
	return strconv.Itoa(int(time.Now().UnixNano() / 1000000))
}

// Bots is a slice of Bot
type Bots []Bot
