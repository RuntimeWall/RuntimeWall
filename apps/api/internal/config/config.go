package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds runtime configuration for the API server.
type Config struct {
	Addr                string
	ReadHeaderTimeout   time.Duration
	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	ShutdownTimeout     time.Duration
	DockerHost          string
	SandboxImage        string
	UbuntuSandboxImage  string
	SandboxNetwork      string
	SandboxLabel        string
	EnableDocker        bool
	SandboxTTL          time.Duration
	SandboxSweepEvery   time.Duration
}

// Load reads configuration from environment variables.
func Load() Config {
	return Config{
		Addr:               envOr("HTTP_ADDR", ":8080"),
		ReadHeaderTimeout:  envDuration("HTTP_READ_HEADER_TIMEOUT", 10*time.Second),
		ReadTimeout:        envDuration("HTTP_READ_TIMEOUT", 0),
		WriteTimeout:       envDuration("HTTP_WRITE_TIMEOUT", 0),
		ShutdownTimeout:    envDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
		DockerHost:         envOr("DOCKER_HOST", ""),
		SandboxImage:       envOr("SANDBOX_IMAGE", "runtimewall/sandbox:latest"),
		UbuntuSandboxImage: envOr("UBUNTU_SANDBOX_IMAGE", "ubuntu:22.04"),
		SandboxNetwork:     envOr("SANDBOX_NETWORK", "runtimewall"),
		SandboxLabel:       envOr("SANDBOX_LABEL", "runtimewall.managed"),
		EnableDocker:       envBool("ENABLE_DOCKER", true),
		SandboxTTL:         envDuration("SANDBOX_TTL", 30*time.Minute),
		SandboxSweepEvery:  envDuration("SANDBOX_SWEEP_INTERVAL", time.Minute),
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func envDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
