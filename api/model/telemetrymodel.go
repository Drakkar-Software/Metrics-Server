package model

import (
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Bot struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BotID     string             `json:"botID,omitempty" bson:"botID"`
	CreatedAt int                `json:"createdAt,omitempty" bson:"createdAt"`
	UpTime    int                `json:"upTime,omitempty" bson:"upTime"`
}

func (t Bot) IsValid() bool {
	if l := len(t.BotID); l > 4 {
		return true
	}

	return false
}

func GenerateBotID() string {
	return strconv.Itoa(int(time.Now().UnixNano() / 1000000))
}

type Bots []Bot
