package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type ServiceConfiguration struct {
	AppConfig *AppConfig
	Server    *Server
}

type AppConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"trace"`
}

type Server struct {
	GatewayPort             string        `envconfig:"GATEWAY_PORT" default:"8081"`
	GRPCHost                string        `envconfig:"GRPC_HOST" default:"0.0.0.0"`
	GRPCPort                string        `envconfig:"GRPC_PORT" default:"10000"`
	GracefulShutdownTimeout time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT" default:"15s"`
	WriteTimeout            time.Duration `envconfig:"WRITE_TIMEOUT" default:"5s"`
	ReadTimeout             time.Duration `envconfig:"READ_TIMEOUT" default:"5s"`
	IdleTimeout             time.Duration `envconfig:"IDLE_TIMEOUT" default:"300s"`
}

func Load() (*ServiceConfiguration, error) {
	cfg := &ServiceConfiguration{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse configuration")
	}

	return cfg, nil
}
