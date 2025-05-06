package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func GetProductionConfig() zap.Config {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config
}

func GetDevelopmentConfig() zap.Config {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return config
}

func GetLogger(environment string) (*zap.Logger, error) {
	var config zap.Config
	if environment == "development" {
		config = GetDevelopmentConfig()
	} else {
		config = GetProductionConfig()
	}

	return config.Build()
}

func init() {
	var err error
	Logger, err = GetLogger("production") // Default to production in init
	if err != nil {
		panic(err) // panicking here is acceptable for logger initialization
	}
}