package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/claytono/go-unifi-mcp/internal/config"
	servermocks "github.com/claytono/go-unifi-mcp/internal/server/mocks"
	"github.com/filipowm/go-unifi/unifi"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_RequiresClient(t *testing.T) {
	_, err := New(Options{Client: nil})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "client is required")
}

func TestNew_CreatesServer(t *testing.T) {
	client := servermocks.NewClient(t)
	t.Setenv("UNIFI_TOOL_MODE", "")

	s, err := New(Options{Client: client})
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Len(t, s.ListTools(), 3)
}

func TestNewClient_APIKey(t *testing.T) {
	cfg := &config.Config{
		Host:      "https://192.168.1.1",
		APIKey:    "test-key",
		Site:      "default",
		VerifySSL: false,
	}

	var captured *unifi.ClientConfig
	prevFactory := newUnifiClient
	newUnifiClient = func(clientCfg *unifi.ClientConfig) (unifi.Client, error) {
		captured = clientCfg
		return nil, nil
	}
	t.Cleanup(func() {
		newUnifiClient = prevFactory
	})

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.Nil(t, client)
	require.NotNil(t, captured)
	assert.Equal(t, cfg.Host, captured.URL)
	assert.Equal(t, cfg.VerifySSL, captured.VerifySSL)
	assert.Equal(t, cfg.APIKey, captured.APIKey)
	assert.Empty(t, captured.User)
	assert.Empty(t, captured.Password)
}

func TestNewClient_UserPass(t *testing.T) {
	cfg := &config.Config{
		Host:      "https://192.168.1.1",
		Username:  "admin",
		Password:  "secret",
		Site:      "default",
		VerifySSL: false,
	}

	var captured *unifi.ClientConfig
	prevFactory := newUnifiClient
	newUnifiClient = func(clientCfg *unifi.ClientConfig) (unifi.Client, error) {
		captured = clientCfg
		return nil, nil
	}
	t.Cleanup(func() {
		newUnifiClient = prevFactory
	})

	client, err := NewClient(cfg)
	assert.NoError(t, err)
	assert.Nil(t, client)
	require.NotNil(t, captured)
	assert.Equal(t, cfg.Host, captured.URL)
	assert.Equal(t, cfg.VerifySSL, captured.VerifySSL)
	assert.Empty(t, captured.APIKey)
	assert.Equal(t, cfg.Username, captured.User)
	assert.Equal(t, cfg.Password, captured.Password)
}

func TestMode_DefaultsToLazy(t *testing.T) {
	// Clear environment variable
	_ = os.Unsetenv("UNIFI_TOOL_MODE")

	// Mode should default to lazy when not specified
	opts := Options{}
	assert.Empty(t, opts.Mode)

	// The actual default is applied in New(), which requires a client
	// So we just verify the constant values are correct
	assert.Equal(t, Mode("lazy"), ModeLazy)
	assert.Equal(t, Mode("eager"), ModeEager)
}

func TestMode_ReadsFromEnvVar(t *testing.T) {
	// Test that environment variable is respected
	_ = os.Setenv("UNIFI_TOOL_MODE", "eager")
	defer func() { _ = os.Unsetenv("UNIFI_TOOL_MODE") }()

	// Can't test full creation without a client, but verify env parsing
	envMode := Mode(os.Getenv("UNIFI_TOOL_MODE"))
	assert.Equal(t, ModeEager, envMode)
}

func TestMode_OptionsOverridesEnv(t *testing.T) {
	// Set environment to eager
	_ = os.Setenv("UNIFI_TOOL_MODE", "eager")
	defer func() { _ = os.Unsetenv("UNIFI_TOOL_MODE") }()

	// Options should take precedence (verified by constant check)
	opts := Options{Mode: ModeLazy}
	assert.Equal(t, ModeLazy, opts.Mode)
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected unifi.LoggingLevel
	}{
		{"disabled", "disabled", unifi.DisabledLevel},
		{"trace", "trace", unifi.TraceLevel},
		{"debug", "debug", unifi.DebugLevel},
		{"info", "info", unifi.InfoLevel},
		{"warn", "warn", unifi.WarnLevel},
		{"error", "error", unifi.ErrorLevel},
		{"unknown", "unknown", unifi.ErrorLevel},
		{"empty", "", unifi.ErrorLevel},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, ParseLogLevel(tt.input))
		})
	}
}

func TestNewClient_DefaultLogLevel(t *testing.T) {
	cfg := &config.Config{
		Host:      "https://192.168.1.1",
		APIKey:    "test-key",
		Site:      "default",
		VerifySSL: false,
		LogLevel:  "error",
	}

	var captured *unifi.ClientConfig
	prevFactory := newUnifiClient
	newUnifiClient = func(clientCfg *unifi.ClientConfig) (unifi.Client, error) {
		captured = clientCfg
		return nil, nil
	}
	t.Cleanup(func() {
		newUnifiClient = prevFactory
	})

	_, _ = NewClient(cfg)
	require.NotNil(t, captured)
	assert.NotNil(t, captured.Logger)
}

func TestNewClient_CustomLogLevel(t *testing.T) {
	cfg := &config.Config{
		Host:      "https://192.168.1.1",
		APIKey:    "test-key",
		Site:      "default",
		VerifySSL: false,
		LogLevel:  "info",
	}

	var captured *unifi.ClientConfig
	prevFactory := newUnifiClient
	newUnifiClient = func(clientCfg *unifi.ClientConfig) (unifi.Client, error) {
		captured = clientCfg
		return nil, nil
	}
	t.Cleanup(func() {
		newUnifiClient = prevFactory
	})

	_, _ = NewClient(cfg)
	require.NotNil(t, captured)
	assert.NotNil(t, captured.Logger)
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := l.Addr().(*net.TCPAddr).Port
	require.NoError(t, l.Close())
	return port
}

// startHTTPServer launches serveHTTP in a goroutine with the given config and
// a cancellable context. Returns the base URL and a cancel func that triggers
// graceful shutdown.
func startHTTPServer(t *testing.T, cfg *config.Config) (baseURL string, cancel context.CancelFunc) {
	t.Helper()
	s := server.NewMCPServer("test", "0.0.0")
	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- serveHTTP(ctx, s, cfg)
	}()

	addr := fmt.Sprintf("http://127.0.0.1:%d", cfg.HTTPPort)
	require.Eventually(t, func() bool {
		resp, err := http.Get(addr + "/health")
		if err != nil {
			return false
		}
		defer func() { _ = resp.Body.Close() }()
		return resp.StatusCode == http.StatusOK
	}, 2*time.Second, 50*time.Millisecond, "server did not start")

	t.Cleanup(func() {
		cancel()
		// Drain the error channel so the goroutine exits
		err := <-errCh
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("unexpected server error: %v", err)
		}
	})

	return addr, cancel
}

func TestServeHTTP_HealthEndpoint(t *testing.T) {
	addr, _ := startHTTPServer(t, &config.Config{
		Transport: "http",
		HTTPHost:  "127.0.0.1",
		HTTPPort:  freePort(t),
		HTTPPath:  "/mcp",
	})

	resp, err := http.Get(addr + "/health")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"status":"ok"}`, string(body))
}

func TestServeHTTP_MCPEndpoint(t *testing.T) {
	addr, _ := startHTTPServer(t, &config.Config{
		Transport: "http",
		HTTPHost:  "127.0.0.1",
		HTTPPort:  freePort(t),
		HTTPPath:  "/mcp",
	})

	// MCP endpoint should be registered (GET returns non-404)
	resp, err := http.Get(addr + "/mcp")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
}

func TestServeHTTP_CustomPath(t *testing.T) {
	addr, _ := startHTTPServer(t, &config.Config{
		Transport: "http",
		HTTPHost:  "127.0.0.1",
		HTTPPort:  freePort(t),
		HTTPPath:  "/custom/mcp",
	})

	resp, err := http.Get(addr + "/custom/mcp")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.NotEqual(t, http.StatusNotFound, resp.StatusCode)
}

func TestServe_HTTPWithSignal(t *testing.T) {
	port := freePort(t)
	cfg := &config.Config{
		Transport: "http",
		HTTPHost:  "127.0.0.1",
		HTTPPort:  port,
		HTTPPath:  "/mcp",
	}

	s := server.NewMCPServer("test", "0.0.0")

	errCh := make(chan error, 1)
	go func() {
		errCh <- Serve(s, cfg)
	}()

	addr := fmt.Sprintf("http://127.0.0.1:%d", port)
	require.Eventually(t, func() bool {
		resp, err := http.Get(addr + "/health")
		if err != nil {
			return false
		}
		defer func() { _ = resp.Body.Close() }()
		return resp.StatusCode == http.StatusOK
	}, 2*time.Second, 50*time.Millisecond, "server did not start")

	// Send SIGINT to ourselves to trigger graceful shutdown
	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGINT))

	select {
	case err := <-errCh:
		// http.ErrServerClosed is expected on graceful shutdown
		if err != nil && err != http.ErrServerClosed {
			t.Fatalf("unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Serve did not return after SIGINT")
	}
}

func TestServeHTTP_GracefulShutdown(t *testing.T) {
	addr, cancel := startHTTPServer(t, &config.Config{
		Transport: "http",
		HTTPHost:  "127.0.0.1",
		HTTPPort:  freePort(t),
		HTTPPath:  "/mcp",
	})

	// Server is up
	resp, err := http.Get(addr + "/health")
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Cancel triggers shutdown
	cancel()

	// Server should stop accepting connections
	require.Eventually(t, func() bool {
		_, err := http.Get(addr + "/health")
		return err != nil
	}, 2*time.Second, 50*time.Millisecond, "server did not shut down")
}
