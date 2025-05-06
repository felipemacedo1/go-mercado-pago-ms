package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"

	"go-mercado-pago-ms/config"
	"go-mercado-pago-ms/routes"
	"go-mercado-pago-ms/utils"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production"
	}

	logger := utils.GetLogger(env)

	logger.Info("Starting application...")

	cfg, err := config.Load()
	if err != nil {
		utils.Logger.Fatal("Failed to load config", zap.Error(err))
	}

	e := echo.New()

	routes.SetupRoutes(e, cfg)

	logger.Info("Server started", zap.String("port", cfg.ServerPort))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.ServerPort)))
	logger.Info("Server stopped")
}
