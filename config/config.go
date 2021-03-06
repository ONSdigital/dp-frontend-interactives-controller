package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-frontend-interactives-controller
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	ServeFromEmbeddedContent   bool          `envconfig:"SERVE_FROM_EMBEDDED_CONTENT"`
	APIRouterURL               string        `envconfig:"API_ROUTER_URL"`
	ServiceAuthToken           string        `envconfig:"SERVICE_AUTH_TOKEN" json:"-"`
	//todo other ctlrs all go through api-router - but this doesnt pass on auth headers so fails when passed thru to download-svc
	DownloadAPIURL string `envconfig:"DOWNLOAD_API_URL"`
}

// Get returns the default config with any modifications through environment vars
func Get() (*Config, error) {
	cfg := &Config{
		BindAddr:                   ":27300",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		ServeFromEmbeddedContent:   false,
		APIRouterURL:               "http://localhost:23200",
		DownloadAPIURL:             "http://localhost:23600",
	}

	return cfg, envconfig.Process("", cfg)
}
