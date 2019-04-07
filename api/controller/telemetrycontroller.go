package controller

import (
	"log"

	"github.com/Drakkar-Software/Metrics-Server/api/dao"
	"github.com/Drakkar-Software/Metrics-Server/api/model"

	"net/http"

	"github.com/labstack/echo"
)

// GetBots returns a json representation of all the bots
func PublicGetBots(c echo.Context) error {
	bots, err := dao.PublicGetBots()
	if err != nil {
		log.Panic(err)
		return c.JSON(http.StatusBadRequest, bots)
	}

	return c.JSON(http.StatusOK, bots)
}

// GenerateBotID returns a new bot ID
func GenerateBotID(c echo.Context) error {
	id, err := dao.GenerateBotID()
	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, id)
	}
	return c.JSON(http.StatusOK, id)
}

// UpdateBotUptime updates a bot uptime
func UpdateBotUptime(c echo.Context) error {
	bot := new(model.Bot)
	c.Bind(bot)
	id, err := dao.UpdateBotUptime(bot)
	if err != nil {
		if err == dao.ErrBotNotFound {
			return c.JSON(http.StatusNotFound, id)
		} else {
			log.Println(err, bot.ID)
			return c.JSON(http.StatusBadRequest, id)
		}
	}

	return c.JSON(http.StatusOK, id)
}

// RegisterBot registers a bot as started (creates a new bot if necessary)
func RegisterBot(c echo.Context) error {
	bot := new(model.Bot)
	c.Bind(bot)
	id, err := dao.RegisterOrUpdate(bot)
	if err != nil {
		if err == dao.ErrBotNotFound {
			return c.JSON(http.StatusNotFound, id)
		} else {
			log.Println(err, bot.ID)
			return c.JSON(http.StatusBadRequest, id)
		}
	}
	return c.JSON(http.StatusOK, id)
}
