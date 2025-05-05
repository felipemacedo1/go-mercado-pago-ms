package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"

	"go-mercado-pago-ms/config"
	"go-mercado-pago-ms/routes"
	"go-mercado-pago-ms/utils"
)

func main() {
	utils.Logger.Info("Starting application...")

	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	cfg := config.Config{}
	err = cfg.LoadEnv()
	if err != nil {
		utils.Logger.Fatal("failed to load env variables")
	}

	e := echo.New()

	routes.SetupRoutes(e)

	utils.Logger.Info("Server started", "port", cfg.ServerPort)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.ServerPort)))
}