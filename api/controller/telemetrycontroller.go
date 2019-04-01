package controller

import (
	"log"

	"github.com/Drakkar-Software/Metrics-Server/api/dao"
	"github.com/Drakkar-Software/Metrics-Server/api/model"

	"net/http"

	"github.com/labstack/echo"
)

// GetBots returns a json representation of all the bots
func GetBots(c echo.Context) error {
	bots, err := dao.GetBots()
	if err != nil {
		log.Panic(err)
	}

	return c.JSON(http.StatusOK, bots)
}

// GenerateBotID returns a new bot ID
func GenerateBotID(c echo.Context) error {
	return c.JSON(http.StatusOK, model.GenerateBotID())
}

// UpdateBotUptime updates a bot uptime
func UpdateBotUptime(c echo.Context) error {
	var bot model.Bot
	c.Bind(bot)
	err := dao.UpdateBotUptime(bot)
	if err != nil {
		log.Panic(err)
	}

	return c.JSON(http.StatusOK, bot.BotID)
}

// RegisterBot registers a bot as started (creates a new bot if necessary)
func RegisterBot(c echo.Context) error {
	var bot model.Bot
	c.Bind(bot)
	err := dao.RegisterOrUpdate(bot)
	if err != nil {
		log.Panic(err)
	}
	return c.JSON(http.StatusOK, bot.BotID)
}
