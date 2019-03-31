package route

import (
	"github.com/Drakkar-Software/Metrics-Server/api/controller"
	"github.com/labstack/echo"
)

func Init(e *echo.Echo) {
	e.GET("/", controller.GetAll)
	e.GET("/gen-bot-id", controller.GenerateBotID)
	e.POST("/metrics/uptime", controller.UpdateBotUptime)
}
