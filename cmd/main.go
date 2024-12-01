package main

import (
	"mymod/internal/app"
	"mymod/internal/config"
	"mymod/internal/database"
	"mymod/internal/services"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("log is loaded")

	var cfg config.MainConfig
	cfg.ConfigMustLoad("docker")
	services.GlobalEmail = cfg.Email
	log.Debug("config is loaded")

	var pgdb database.PostgresDatabase
	pgdb.Run(cfg)
	log.Debug("databases is loaded")

	var application app.App
	go application.NewServer(cfg.Server.Port)
	log.Debug("server is loaded")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	application.Stop()

}
