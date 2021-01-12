package controller

import (
	"log"
	"strconv"
	"time"

	"github.com/Drakkar-Software/Metrics-Server/api/dao"
	"github.com/Drakkar-Software/Metrics-Server/api/model"

	"net/http"

	"github.com/labstack/echo/v4"
)

// PublicGetBots returns a json representation of all the bots
func PublicGetBots(c echo.Context) error {
	bots, err := dao.PublicGetBots(true)
	if err != nil {
		log.Panic(err)
		return c.JSON(http.StatusBadRequest, bots)
	}

	return c.JSON(http.StatusOK, bots)
}

// PublicGetRecentBotsWithProfitability returns a json representation of all the recent bots
// with their profitability
func PublicGetRecentBotsWithProfitability(c echo.Context) error {
	sinceParam, err := strconv.ParseInt(c.Param("since"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	bots, err := dao.PublicGetBots(false)
	if err != nil {
		log.Panic(err)
		return c.JSON(http.StatusBadRequest, getOnlyRecentBots(bots, int(sinceParam)))
	}

	return c.JSON(http.StatusOK, getOnlyRecentBots(bots, int(sinceParam)))
}

// AuthenticatedGetBots returns a json representation of all the bots without filters
func AuthenticatedGetBots(c echo.Context) error {
	return getAuthBots(c, false, model.CurrentAccess)
}

// AuthenticatedGetBotsHistory returns a json representation of all the bots with historical data without filters
func AuthenticatedGetBotsHistory(c echo.Context) error {
	return getAuthBots(c, true, model.HistoricalDataAccess)
}

func getOnlyRecentBots(bots model.Bots, since int) model.Bots {
	// Create a zero-length slice with the same underlying array
	filteredBots := bots[:0]
	for _, bot := range bots {
		if bot.CurrentSession.UpTime >= since {
			// Append desired values to slice
			filteredBots = append(filteredBots, bot)
		}
	}
	return filteredBots
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
		log.Panic(err)
	}

	// returns full total if params are all zeros
	if years == 0 && months == 0 && days == 0 {
		totalBots, err := dao.PublicGetCountBots(0)
		if err != nil {
			log.Panic(err)
		}
		jsonMap["total"] = totalBots
	} else {
		lastMonthTimeStamp := time.Now().AddDate(years, months, days)
		totalBots, err := dao.PublicGetCountBots(lastMonthTimeStamp.Unix())
		if err != nil {
			log.Panic(err)
		}
		jsonMap["total"] = totalBots
	}
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
