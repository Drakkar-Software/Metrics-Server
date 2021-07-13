package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Session stores data about a bot session
type Session struct {
	StartedAt      int      `json:"startedat,omitempty" bson:"startedAt"`
	UpTime         int      `json:"uptime,omitempty" bson:"upTime"`
	Version        string   `json:"version,omitempty" bson:"version"`
	Simulator      bool     `json:"simulator,omitempty" bson:"simulator"`
	Trader         bool     `json:"trader,omitempty" bson:"trader"`
	EvalConfig     []string `json:"evalconfig,omitempty" bson:"evalConfig"`
	Pairs          []string `json:"pairs,omitempty" bson:"pairs"`
	Exchanges      []string `json:"exchanges,omitempty" bson:"exchanges"`
	Notifications  []string `json:"notifications,omitempty" bson:"notifications"`
	Type           string   `json:"type,omitempty" bson:"type"`
	Platform       string   `json:"platform,omitempty" bson:"platform"`
	ReferenceMaket string   `json:"referencemarket,omitempty" bson:"referenceMarket"`
	PortfolioValue float32  `json:"portfoliovalue,omitempty" bson:"portfolioValue"`
	Profitability  float32  `json:"profitability,omitempty" bson:"profitability"`
	TradedVolumes  map[string]float32  `json:"tradedvolumes,omitempty" bson:"tradedVolumes"`
	Supports  map[string][]string  `json:"supports,omitempty" bson:"supports"`
}

// Bot stores usage info about a specific bot identified by BotID
type Bot struct {
	ID             primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CurrentSession Session            `json:"currentsession,omitempty" bson:"currentSession"`
	SessionHistory []Session          `json:"sessionhistory,omitempty" bson:"sessionHistory"`
}

// Bots is a slice of Bot
type Bots []Bot

// FilterPublicInfo Resets non public info
func (bot *Bot) FilterPublicInfo() {
	bot.ID = primitive.NilObjectID
	bot.CurrentSession.FilterPublicInfo()
	bot.SessionHistory = nil
}

// FilterPublicInfo Resets non public info
func (session *Session) FilterPublicInfo() {
	session.UpTime = session.StartedAt + session.UpTime
	session.StartedAt = 0
	session.Notifications = nil
	session.Type = ""
	session.Platform = ""
	session.PortfolioValue = 0
	session.Profitability = 0
	session.TradedVolumes = nil
	session.ReferenceMaket = ""
	session.Supports = nil
}
