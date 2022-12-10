package route

import (
	"github.com/Drakkar-Software/Metrics-Server/api/controller"
	"github.com/labstack/echo/v4"
)

// Init registers the server routes
func Init(e *echo.Echo) {
	e.GET("/metrics/community", controller.PublicGetBots)
	e.GET("/metrics/community/top/exchanges/:since", controller.TopExchanges)
	e.GET("/metrics/community/top/pairs/:since", controller.TopPairs)
	e.GET("/metrics/community/top/trading_modes/:since", controller.TopTradingModes)
	e.GET("/metrics/community/top/evaluation_configs/:since", controller.TopEvaluationConfigs)
	e.GET("/metrics/community/count/:years/:months/:days", controller.PublicGetCount)

	e.GET("/metrics/full_data", controller.AuthenticatedGetBots)
	e.GET("/metrics/full_data/history", controller.AuthenticatedGetBotsHistory)

	e.GET("/gen-bot-id", controller.GenerateBotID)
	e.POST("/metrics/uptime", controller.UpdateBotUptimeAndProfitability)
	e.POST("/metrics/register", controller.RegisterBot)
}
