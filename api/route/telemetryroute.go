package route

import (
	"github.com/Drakkar-Software/Metrics-Server/api/controller"
	"github.com/labstack/echo"
)

// Init registers the server routes
func Init(e *echo.Echo) {
	e.GET("/gen-bot-id", controller.GenerateBotID)
	e.GET("/metrics/community", controller.PublicGetBots)
	e.GET("/metrics/community/count/:years/:months/:days", controller.PublicGetCount)
	e.POST("/metrics/uptime", controller.UpdateBotUptime)
	e.POST("/metrics/register", controller.RegisterBot)
}
