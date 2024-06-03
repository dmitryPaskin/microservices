package main

import (
	"microservices/notify/config"
	"microservices/notify/internal/infrastructure/logs"
	"microservices/notify/run"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	conf := config.NewAppConf()

	logger := logs.NewLogger(conf, os.Stdout)

	conf.Init(logger)

	app := run.NewApp(conf, logger)

	if code := app.Run(); code != 0 {
		logger.Info("app run error")
		os.Exit(code)
	}
}
