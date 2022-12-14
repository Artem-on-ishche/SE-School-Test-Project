package main

import (
	"gses2.app/api/pkg"
	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/infrastructure/logger"
	"gses2.app/api/pkg/presentation/http/routes"
)

func main() {
	config.LoadEnv()

	loggerService := logger.NewRabbitMQLogger()
	defer loggerService.Close()

	err := routes.SetupRouter(pkg.InitServices(loggerService)).Run(config.NetworkPort)
	if err != nil {
		loggerService.Error("Router Error")
	}
}
