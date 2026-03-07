package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	ErrMissingHost        = errors.New("UNIFI_HOST environment variable is required")
	ErrMissingCredentials = errors.New("either UNIFI_API_KEY or both UNIFI_USERNAME and UNIFI_PASSWORD must be set")
	ErrInvalidLogLevel    = errors.New("UNIFI_LOG_LEVEL must be one of: disabled, trace, debug, info, warn, error")
	ErrInvalidTransport   = errors.New("UNIFI_TRANSPORT must be one of: stdio, http")
	ErrInvalidHTTPPort    = errors.New("UNIFI_HTTP_PORT must be a valid port number (1-65535)")
)

var validLogLevels = map[string]bool{
	"disabled": true,
	"trace":    true,
	"debug":    true,
	"info":     true,
	"warn":     true,
	"error":    true,
}

var validTransports = map[string]bool{
	"stdio": true,
	"http":  true,
}

// Config holds the MCP server configuration.
type Config struct {
	Host      string // UNIFI_HOST - UniFi controller URL
	APIKey    string // UNIFI_API_KEY - API key auth (preferred)
	Username  string // UNIFI_USERNAME - username/password auth
	Password  string // UNIFI_PASSWORD - username/password auth
	Site      string // UNIFI_SITE - site name (default: "default")
	VerifySSL bool   // UNIFI_VERIFY_SSL - verify SSL certs (default: true)
	LogLevel  string // UNIFI_LOG_LEVEL - go-unifi log level (default: "error")

	Transport string // UNIFI_TRANSPORT - transport mode: stdio or http (default: "stdio")
	HTTPHost  string // UNIFI_HTTP_HOST - HTTP listen address (default: "0.0.0.0")
	HTTPPort  int    // UNIFI_HTTP_PORT - HTTP listen port (default: 8080)
	HTTPPath  string // UNIFI_HTTP_PATH - MCP endpoint path (default: "/mcp")
}

// Load loads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Host:      os.Getenv("UNIFI_HOST"),
		APIKey:    os.Getenv("UNIFI_API_KEY"),
		Username:  os.Getenv("UNIFI_USERNAME"),
		Password:  os.Getenv("UNIFI_PASSWORD"),
		Site:      os.Getenv("UNIFI_SITE"),
		VerifySSL: true,
	}

	// Parse UNIFI_VERIFY_SSL
	if v := os.Getenv("UNIFI_VERIFY_SSL"); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return nil, errors.New("UNIFI_VERIFY_SSL must be a boolean (true/false)")
		}
		cfg.VerifySSL = parsed
	}

	// Parse UNIFI_LOG_LEVEL
	if v := os.Getenv("UNIFI_LOG_LEVEL"); v != "" {
		v = strings.ToLower(v)
		if !validLogLevels[v] {
			return nil, fmt.Errorf("%w: got %q", ErrInvalidLogLevel, v)
		}
		cfg.LogLevel = v
	} else {
		cfg.LogLevel = "error"
	}

	// Parse UNIFI_TRANSPORT
	if v := os.Getenv("UNIFI_TRANSPORT"); v != "" {
		v = strings.ToLower(v)
		if !validTransports[v] {
			return nil, fmt.Errorf("%w: got %q", ErrInvalidTransport, v)
		}
		cfg.Transport = v
	} else {
		cfg.Transport = "stdio"
	}

	// Parse HTTP transport settings
	cfg.HTTPHost = os.Getenv("UNIFI_HTTP_HOST")
	if cfg.HTTPHost == "" {
		cfg.HTTPHost = "0.0.0.0"
	}

	if v := os.Getenv("UNIFI_HTTP_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil || port < 1 || port > 65535 {
			return nil, fmt.Errorf("%w: got %q", ErrInvalidHTTPPort, v)
		}
		cfg.HTTPPort = port
	} else {
		cfg.HTTPPort = 8080
	}

	cfg.HTTPPath = os.Getenv("UNIFI_HTTP_PATH")
	if cfg.HTTPPath == "" {
		cfg.HTTPPath = "/mcp"
	}

	// Set default site
	if cfg.Site == "" {
		cfg.Site = "default"
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks required configuration.
func (c *Config) Validate() error {
	if c.Host == "" {
		return ErrMissingHost
	}

	if !c.UseAPIKey() && !c.UseUserPass() {
		return ErrMissingCredentials
	}

	return nil
}

// UseAPIKey returns true if API key auth should be used.
func (c *Config) UseAPIKey() bool {
	return c.APIKey != ""
}

// UseUserPass returns true if username/password auth should be used.
func (c *Config) UseUserPass() bool {
	return c.Username != "" && c.Password != ""
}
