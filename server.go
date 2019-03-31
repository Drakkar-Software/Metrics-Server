package main

import (
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

	e.Logger.Fatal(e.Start(":8080"))
}
