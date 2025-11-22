package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/vladgrskkh/pr_service/internal/application"
	"github.com/vladgrskkh/pr_service/internal/server"
)

func main() {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "config/config.toml", "path to a config file")

	flag.Parse()

	app := application.NewAppllication(cfgFile)

	defer app.DB.Close()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Cfg.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	server := server.NewServer(srv, app.Logger)

	app.Logger.Info("Starting server")
	err := server.Serve()
	if err != nil {
		app.Logger.Log(context.Background(), application.LevelFatal, err.Error())
	}
}
