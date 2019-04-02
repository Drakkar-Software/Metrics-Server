package model

import (
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

// Session stores data about a bot session
type Session struct {
	StartedAt  int      `json:"startedat,omitempty" bson:"startedAt"`
	UpTime     int      `json:"uptime,omitempty" bson:"upTime"`
	Simulator  bool     `json:"simulator,omitempty" bson:"simulator"`
	Trader     bool     `json:"trader,omitempty" bson:"trader"`
	EvalConfig []string `json:"evalconfig,omitempty" bson:"evalConfig"`
	Exchanges  []string `json:"exchanges,omitempty" bson:"exchanges"`
}

// Bot stores usage info about a specific bot identified by BotID
type Bot struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	BotID          string             `json:"botid" bson:"botID"`
	CurrentSession Session            `json:"currentsession,omitempty" bson:"currentSession"`
	SessionHistory []Session          `json:"sessionsistory,omitempty" bson:"sessionHistory"`
}

// GenerateBotID generates a new bot id
func GenerateBotID() string {
	return strconv.Itoa(int(time.Now().UnixNano() / 1000000))
}

// Bots is a slice of Bot
type Bots []Bot
