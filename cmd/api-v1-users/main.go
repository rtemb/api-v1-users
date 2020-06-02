package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gitlab.com/rtemb/api-v1-users/internal/api"
	"gitlab.com/rtemb/api-v1-users/internal/auth"
	"gitlab.com/rtemb/api-v1-users/internal/config"
	"gitlab.com/rtemb/api-v1-users/internal/middleware"
	service "gitlab.com/rtemb/api-v1-users/internal/service"
	"gitlab.com/rtemb/api-v1-users/pkg/version"
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
