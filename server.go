package main

import (
	"net/http"
	"os"

	"github.com/Drakkar-Software/Metrics-Server/database"
	"github.com/Drakkar-Software/Metrics-Server/routes"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := database.Init()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	routes.Init(e)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	e.Logger.Fatal(e.Start(":" + getPort()))
}

func getPort() string {
	port, exists := os.LookupEnv("PORT")
	if !exists || port == "" {
		port = "8080"
	}
	return port
}
