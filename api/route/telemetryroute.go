package route

import (
	"github.com/Drakkar-Software/Metrics-Server/api/controller"
	"github.com/labstack/echo/v4"
)

// Init registers the server routes
func Init(e *echo.Echo) {
	e.GET("/gen-bot-id", controller.GenerateBotID)
	e.GET("/metrics/community", controller.PublicGetBots)
	e.GET("/metrics/recent_bots/:since", controller.PublicGetRecentBotsWithProfitability)
	e.GET("/metrics/full_data", controller.AuthenticatedGetBots)
	e.GET("/metrics/full_data/history", controller.AuthenticatedGetBotsHistory)
	e.GET("/metrics/community/count/:years/:months/:days", controller.PublicGetCount)
	e.POST("/metrics/uptime", controller.UpdateBotUptimeAndProfitability)
	e.POST("/metrics/register", controller.RegisterBot)
}
