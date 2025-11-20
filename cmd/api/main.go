package main

import (
	"context"
	"net/http"
	"time"

	"github.com/vladgrskkh/pr_service/internal/application"
	"github.com/vladgrskkh/pr_service/internal/server"
)

func main() {
	app := application.NewAppllication()

	srv := &http.Server{
		Addr:         ":8080",
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
