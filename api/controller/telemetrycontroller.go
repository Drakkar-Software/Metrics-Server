package controller

import (
	"log"

	"github.com/Drakkar-Software/Metrics-Server/api/dao"
	"github.com/Drakkar-Software/Metrics-Server/api/model"

	"net/http"

	"github.com/labstack/echo"
)

func GetAll(c echo.Context) error {
	ts, err := dao.All()
	if err != nil {
		log.Panic(err)
	}

	return c.JSON(http.StatusOK, ts)
}

func GenerateBotID(c echo.Context) error {
	newID := model.GenerateBotID()

	return c.JSON(http.StatusOK, newID)
}

func UpdateBotUptime(c echo.Context) error {
	var bot model.Bot

	c.Bind(bot)
	err := dao.UpdateBotUptime(bot)
	if err != nil {
		log.Panic(err)
	}

	return c.JSON(http.StatusOK, bot.BotID)
}
