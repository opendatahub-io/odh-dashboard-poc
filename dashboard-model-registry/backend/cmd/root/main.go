package main

import (
	"flag"
	"fmt"
	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/cmd/config"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/opendatahub-io/odh-dashboard-poc/dashboard-model-registry/internal/application"
)

func main() {

	var cfg config.Config
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	//ederign this is not working
	flag.Var(&cfg.Env, "env", "Environment (development|staging|production)")
	cfg.Env = "development"
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := application.NewApp(cfg, logger)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      app.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)

}
