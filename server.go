package main

import (
	"os"

	"github.com/Drakkar-Software/Metrics-Server/api/route"
	"github.com/Drakkar-Software/Metrics-Server/database"
	"github.com/labstack/echo"
)

func main() {
	database := new(database.DB)
	err := database.Initialize()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	route.Init(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
