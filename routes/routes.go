package routes

import (
	"github.com/Drakkar-Software/Metrics-Server/api/route"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Init initializes the echo routes and registers the middleware
func Init(e *echo.Echo) {
	e.Pre(middleware.RemoveTrailingSlash())
	// allow a max request size of 4000 characters
	e.Use(middleware.BodyLimitWithConfig(middleware.BodyLimitConfig{Limit: "4K"}))
	route.Init(e)
}
