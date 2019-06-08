package main

import (
	"os"

	"github.com/Drakkar-Software/Metrics-Server/database"
	"github.com/Drakkar-Software/Metrics-Server/routes"
	"github.com/labstack/echo"
)

func main() {
	err := database.Init()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	routes.Init(e)

	e.Logger.Fatal(e.Start(":" + getPort()))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
