package main

import (
	"os"

	"github.com/Drakkar-Software/Metrics-Server/api/dao"
	"github.com/Drakkar-Software/Metrics-Server/api/route"
	"github.com/labstack/echo"
)

func main() {
	err := dao.Initialize()
	if err != nil {
		panic(err)
	}
	e := echo.New()
	route.Init(e)

	e.Logger.Fatal(e.Start(":" + getPort()))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
