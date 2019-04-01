package model

import (
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Session stores data about a bot session
type Session struct {
	StartedAt int `json:"startedAt,omitempty" bson:"startedAt"`
	UpTime    int `json:"upTime,omitempty" bson:"upTime"`
}

// Bot stores usage info about a specific bot identified by BotID
type Bot struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BotID          string             `json:"botID,omitempty" bson:"botID"`
	CurrentSession Session            `json:"currentSession,omitempty" bson:"currentSession"`
	SessionHistory []Session          `json:"sessionHistory,omitempty" bson:"sessionHistory"`
}

// GenerateBotID generates a new bot id
func GenerateBotID() string {
	return strconv.Itoa(int(time.Now().UnixNano() / 1000000))
}

// Bots is a slice of Bot
type Bots []Bot
