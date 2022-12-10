package controller

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Drakkar-Software/Metrics-Server/api/caches"

	"github.com/Drakkar-Software/Metrics-Server/api/dao"
	"github.com/Drakkar-Software/Metrics-Server/api/model"

	"net/http"

	"github.com/labstack/echo/v4"
)

// PublicGetBots returns a json representation of all the bots
func PublicGetBots(c echo.Context) error {
	bots, err := dao.PublicGetBots(0, true)
	if err != nil {
		log.Panic(err)
		return c.JSON(http.StatusBadRequest, bots)
	}

	return c.JSON(http.StatusOK, bots)
}

func getTraderTypeParam(c echo.Context) string {
	traderTypeParam := c.QueryParam("traderType")
	if traderTypeParam == "simulated" {
		return model.SimulatedTraders
	}
	if traderTypeParam == "real" {
		return model.RealTraders
	}
	return model.AllTraders
}

// TopExchanges return the top of bot current sessions exchanges
func TopExchanges(c echo.Context) error {
	sinceParam, err := strconv.ParseInt(c.Param("since"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	traderType := getTraderTypeParam(c)

	value, exists := caches.GetExchangeTop(traderType, sinceParam)
	if exists {
		return c.JSON(http.StatusOK, value)
	}
	topExchanges := dao.FetchTop("exchanges", traderType, sinceParam, 1, 100, false)
	caches.SetExchangeTop(traderType, sinceParam, topExchanges)

	return c.JSON(http.StatusOK, topExchanges)
}

// TopPairs return the top of bot current sessions pairs
func TopPairs(c echo.Context) error {
	sinceParam, err := strconv.ParseInt(c.Param("since"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	traderType := getTraderTypeParam(c)

	value, exists := caches.GetPairTop(traderType, sinceParam)
	if exists {
		return c.JSON(http.StatusOK, value)
	}
	topPairs := dao.FetchTop("pairs", traderType, sinceParam, 1, 1000, false)
	caches.SetPairTop(traderType, sinceParam, topPairs)

	return c.JSON(http.StatusOK, topPairs)
}

// TopTradingModes return the top of bot current sessions trading modes
func TopTradingModes(c echo.Context) error {
	sinceParam, err := strconv.ParseInt(c.Param("since"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	traderType := getTraderTypeParam(c)

	value, exists := caches.GetTradingModeTop(traderType, sinceParam)
	if exists {
		return c.JSON(http.StatusOK, value)
	}
	topTradingModes := dao.FetchTop("evalConfig", traderType, sinceParam, 3, 100, true)
	caches.SetTradingModeTop(traderType, sinceParam, topTradingModes)

	return c.JSON(http.StatusOK, topTradingModes)
}

// TopEvaluationConfigs return the top of bot current sessions trading modes and evaluators
func TopEvaluationConfigs(c echo.Context) error {
	sinceParam, err := strconv.ParseInt(c.Param("since"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	traderType := getTraderTypeParam(c)

	value, exists := caches.GetEvaluationConfigTop(traderType, sinceParam)
	if exists {
		return c.JSON(http.StatusOK, value)
	}
	topEvalConfig := dao.FetchTop("evalConfig", traderType, sinceParam, 3, 100, false)

	caches.SetEvaluationConfigTop(traderType, sinceParam, topEvalConfig)
	return c.JSON(http.StatusOK, topEvalConfig)
}

// AuthenticatedGetBots returns a json representation of all the bots without filters
func AuthenticatedGetBots(c echo.Context) error {
	return getAuthBots(c, false, model.CurrentAccess)
}

// AuthenticatedGetBotsHistory returns a json representation of all the bots with historical data without filters
func AuthenticatedGetBotsHistory(c echo.Context) error {
	return getAuthBots(c, true, model.HistoricalDataAccess)
}

func getAuthBots(c echo.Context, history bool, minAccessRight int8) error {
	log.Println(c.Request().Header)
	if IsIPAllowed(c) {
		if dao.IsAuthorizedUser(c.Request().Header.Get("Api-Key"), minAccessRight) {
			bots, err := dao.CompleteGetBots(history)
			if err != nil {
				log.Panic(err)
				return c.JSON(http.StatusBadRequest, bots)
			}

			return c.JSON(http.StatusOK, bots)
		}
		return c.JSON(http.StatusUnauthorized, nil)
	}
	return c.JSON(http.StatusTooManyRequests, nil)
}

// PublicGetCount returns the number of total / yearly / monthly / daily active bots
func PublicGetCount(c echo.Context) error {
	jsonMap := make(map[string]interface{})
	yearsParam := c.Param("years")
	monthsParam := c.Param("months")
	daysParam := c.Param("days")

	years, err := strconv.Atoi(yearsParam)
	if err != nil {
		log.Panic(err)
	}

	months, err := strconv.Atoi(monthsParam)
	if err != nil {
		log.Panic(err)
	}

	days, err := strconv.Atoi(daysParam)
	if err != nil {
		// try removing request params
		splitDays := strings.Split(daysParam, "&")
		days, err = strconv.Atoi(splitDays[0])
		if err != nil {
			log.Panic(err)
		}
	}

	untilTime := int64(0)
	// returns full total if params are all zeros
	if !(years == 0 && months == 0 && days == 0) {
		lastMonthTimeStamp := time.Now().AddDate(years, months, days)
		untilTime = lastMonthTimeStamp.Unix()
	}
	totalBots, exists := caches.GetBotCount(untilTime)
	if !exists {
		totalBots, err = dao.PublicGetCountBots(untilTime)
		if err != nil {
			log.Panic(err)
		}
		caches.SetBotCount(untilTime, totalBots)
	}
	jsonMap["total"] = totalBots
	return c.JSON(http.StatusOK, jsonMap)
}

// GenerateBotID returns a new bot ID
func GenerateBotID(c echo.Context) error {
	if IsIPAllowed(c) {
		id, err := dao.GenerateBotID()
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusBadRequest, id)
		}
		return c.JSON(http.StatusOK, id)
	}
	return c.JSON(http.StatusTooManyRequests, nil)
}

// UpdateBotUptimeAndProfitability updates a bot uptime
func UpdateBotUptimeAndProfitability(c echo.Context) error {
	bot := new(model.Bot)
	_ = c.Bind(bot)
	id, err := dao.UpdateBotUptimeAndProfitability(bot)
	if err != nil {
		if err == dao.ErrBotNotFound {
			return c.JSON(http.StatusNotFound, id)
		}
		log.Println(err, bot.ID)
		return c.JSON(http.StatusBadRequest, id)
	}
	return c.JSON(http.StatusOK, id)
}

// RegisterBot registers a bot as started (creates a new bot if necessary)
func RegisterBot(c echo.Context) error {
	if IsIPAllowed(c) {
		bot := new(model.Bot)
		_ = c.Bind(bot)
		id, err := dao.RegisterOrUpdate(bot)
		if err != nil {
			if err == dao.ErrBotNotFound {
				return c.JSON(http.StatusNotFound, id)
			} else if err == dao.ErrInvalidData {
				return c.JSON(http.StatusBadRequest, id)
			}
			log.Println(err, bot.ID)
			return c.JSON(http.StatusBadRequest, id)
		}
		return c.JSON(http.StatusOK, id)
	}
	return c.JSON(http.StatusTooManyRequests, nil)
}
