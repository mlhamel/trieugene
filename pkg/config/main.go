package config

import (
	"os"
	"strconv"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/rs/zerolog"
)

// Config is the main configuration structure
type Config struct {
	httpPort int
	logger   *zerolog.Logger
	statsd   *statsd.Client
}

func NewConfig() *Config {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()

	httpPort, err := strconv.Atoi(GetEnv("PORT", "8080"))

	statsd, err := statsd.New("127.0.0.1:8125", statsd.WithNamespace("trieugene."))

	if err != nil {
		panic(err)
	}

	return &Config{
		httpPort: httpPort,
		logger:   &logger,
		statsd:   statsd,
	}
}

func (c *Config) ProjectID() string {
	return GetEnv("GOOGLE_CLOUD_PROJECT", "trieugene")
}

func (c *Config) BucketName() string {
	return "trieugene-storage"
}

func (c *Config) Logger() *zerolog.Logger {
	return c.logger
}

func (c *Config) PubSubURL() string {
	return "trieugene.myshopify.io:8432"
}

func (c *Config) GCSURL() string {
	return "trieugene.myshopify.io:4443"
}

func (c *Config) HTTPPort() int {
	return c.httpPort
}
