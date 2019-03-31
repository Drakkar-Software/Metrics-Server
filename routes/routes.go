package routes

import (
	"github.com/Drakkar-Software/Metrics-Server/api/route"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Init(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())

	route.Init(e)
}
