package controller

import (
	"github.com/Drakkar-Software/Metrics-Server/api/caches"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

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

// TopTradingModes return the top of bot current sessions exchanges
func TopExchanges(c echo.Context) error {
	sinceParam, err := strconv.ParseInt(c.Param("since"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	value, exists := caches.GetExchangeTop(sinceParam)
	if exists{
		return c.JSON(http.StatusOK, value)
	}
	bots, err := dao.PublicGetBots(sinceParam, false)
	if err != nil {
		log.Panic(err)
		return c.JSON(http.StatusBadRequest, bots)
	}
	top := make(map[string]int)
	for _, bot := range bots {
		for _, exchange := range bot.CurrentSession.Exchanges {
			val, present := top[exchange]
			if present {
				top[exchange] = val + 1
			}else{
				top[exchange] = 1
			}
		}
	}
	sortedTop := getSortedTop(top)
	caches.SetExchangeTop(sinceParam, sortedTop)
	return c.JSON(http.StatusOK, sortedTop)
}

// TopTradingModes return the top of bot current sessions trading modes
func TopTradingModes(c echo.Context) error {
	sinceParam, err := strconv.ParseInt(c.Param("since"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	value, exists := caches.GetTradingModeTop(sinceParam)
	if exists{
		return c.JSON(http.StatusOK, value)
	}
	bots, err := dao.PublicGetBots(sinceParam, false)
	if err != nil {
		log.Panic(err)
		return c.JSON(http.StatusBadRequest, bots)
	}
	top := make(map[string]int)
	for _, bot := range bots {
		for _, evalOrTradingMode := range bot.CurrentSession.EvalConfig {
			if strings.HasSuffix(evalOrTradingMode, "TradingMode") {
				val, present := top[evalOrTradingMode]
				if present {
					top[evalOrTradingMode] = val + 1
				}else{
					top[evalOrTradingMode] = 1
				}
			}
		}
	}
	sortedTop := getSortedTop(top)
	caches.SetTradingModeTop(sinceParam, sortedTop)
	return c.JSON(http.StatusOK, sortedTop)
}

// TopProfitabilities return the top of bot current sessions profitability with a limit of maxValues
func TopProfitabilities(c echo.Context) error {
	maxValues := 3
	sinceParam, err := strconv.ParseInt(c.Param("since"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	countParam, err := strconv.ParseInt(c.Param("count"), 10, 0)
	if err != nil || int(countParam) > maxValues {
		return c.JSON(http.StatusBadRequest, nil)
	}
	bots, err := dao.PublicGetBots(sinceParam, false)
	profitabilities := make([]float32, len(bots))
	i := 0
	for _, bot := range bots {
		profitabilities[i] = bot.CurrentSession.Profitability
		i++
	}
	sort.SliceStable(profitabilities, func(i, j int) bool {
		return profitabilities[i] > profitabilities[j]
	})

	returnedProfitabilities := profitabilities
	if len(profitabilities) >= int(countParam){
		returnedProfitabilities = profitabilities[:countParam]
	}
	return c.JSON(http.StatusOK, returnedProfitabilities)
}

// AuthenticatedGetBots returns a json representation of all the bots without filters
func AuthenticatedGetBots(c echo.Context) error {
	return getAuthBots(c, false, model.CurrentAccess)
}

// AuthenticatedGetBotsHistory returns a json representation of all the bots with historical data without filters
func AuthenticatedGetBotsHistory(c echo.Context) error {
	return getAuthBots(c, true, model.HistoricalDataAccess)
}

func getSortedTop(top map[string]int) []model.Top {
	topList := make([]model.Top, len(top))
	i := 0
	for key, value := range top {
		topList[i] = model.Top{Name: key, Count: value}
		i++
	}
	sort.SliceStable(topList, func(i, j int) bool {
		return topList[i].Count > topList[j].Count
	})
	return topList
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
		untilTime  = lastMonthTimeStamp.Unix()
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
