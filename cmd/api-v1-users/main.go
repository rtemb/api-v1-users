package main

import (
	"github.com/pkg/errors"
	"github.com/rtemb/api-v1-users/internal/api"
	"github.com/rtemb/api-v1-users/internal/auth"
	"github.com/rtemb/api-v1-users/internal/config"
	"github.com/rtemb/api-v1-users/internal/middleware"
	"github.com/rtemb/api-v1-users/internal/service"
	"github.com/rtemb/api-v1-users/pkg/version"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New().WithFields(logrus.Fields{
		"gitSha":  version.GitSha,
		"version": version.ServiceVersion,
		"logger":  "cmd/api-v1-users",
	})
	logger.Println("loading service configurations")
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal(errors.Wrap(err, "could not load service config"))
	}

	lvl, err := logrus.ParseLevel(cfg.AppConfig.LogLevel)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "could parse log level"))
	}
	logger.Logger.SetLevel(lvl)

	simpleAuth := auth.NewSimpleAuth("1234")
	mw := middleware.NewMiddleware(simpleAuth, logger)

	s := service.NewService(logger)
	apiHandler := api.NewHandler(s, logger)

	api.StartService(cfg.Server, mw, apiHandler, logger)
}
