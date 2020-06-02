package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	err := os.Setenv("LOG_LEVEL", "error")
	require.NoError(t, err)

	err = os.Setenv("APP_PORT", "777")
	require.NoError(t, err)
	err = os.Setenv("GRACEFUL_SHUTDOWN_TIMEOUT", "30s")
	require.NoError(t, err)
	err = os.Setenv("WRITE_TIMEOUT", "30s")
	require.NoError(t, err)
	err = os.Setenv("READ_TIMEOUT", "30s")
	require.NoError(t, err)
	err = os.Setenv("IDLE_TIMEOUT", "30s")
	require.NoError(t, err)

	testDuration, err := time.ParseDuration("30s")
	require.NoError(t, err)

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "8081", cfg.Server.GatewayPort)
	assert.Equal(t, testDuration, cfg.Server.GracefulShutdownTimeout)
	assert.Equal(t, testDuration, cfg.Server.WriteTimeout)
	assert.Equal(t, testDuration, cfg.Server.ReadTimeout)
	assert.Equal(t, testDuration, cfg.Server.IdleTimeout)
}
