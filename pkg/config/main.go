package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

// Config is the main configuration structure
type Config struct {
	httpPort int
	logger   *zerolog.Logger
}

func NewConfig() *Config {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger()

	httpPort, err := strconv.Atoi(GetEnv("PORT", "8080"))

	if err != nil {
		panic(err)
	}

	return &Config{
		httpPort: httpPort,
		logger:   &logger,
	}
}

type environmentType string

const (
	Production  environmentType = "Production"
	Development environmentType = "Development"
)

func (c *Config) BucketName() string {
	return GetEnv("TRIEUGENE_BUCKET_NAME", "trieugene-storage")
}

func (c *Config) Environment() environmentType {
	return environmentType(GetEnv("TRIEUGENE_ENVIRONMENT", string(Development)))
}

func (c *Config) FaktoryURL() string {
	return GetEnv("TRIEUGENE_FAKTORY_URL", "tcp://faktory:7419")
}

func (c *Config) GCSURL() string {
	return GetEnv("TRIEUGENE_GCS_URL", "trieugene.myshopify.io:4443")
}

func (c *Config) GCSAccessKey() string {
	return GetEnv("TRIEUGENE_GCS_ACCESS_KEY", "")
}

func (c *Config) GCSAccessSecret() string {
	return GetEnv("TRIEUGENE_GCS_ACCESS_SECRET", "")
}

func (c *Config) HTTPPort() int {
	return c.httpPort
}

func (c *Config) LocalPrefix() string {
	return fmt.Sprintf("%s/.cfg/trieugene", GetEnv("HOME", "/tmp"))
}

func (c *Config) Logger() *zerolog.Logger {
	return c.logger
}

func (c *Config) ProjectID() string {
	return GetEnv("GOOGLE_CLOUD_PROJECT", "trieugene")
}

func (c *Config) PubSubURL() string {
	return GetEnv("TRIEUGENE_PUBSUB_URL", "trieugene.myshopify.io:8085")
}

func (c *Config) S3Bucket() string {
	return GetEnv("TRIEUGENE_S3_BUCKET", "trieugene")
}

func (c *Config) S3AccessKey() string {
	return GetEnv("TRIEUGENE_S3_ACCESS_KEY", "trieugene_key")
}

func (c *Config) S3SecretKey() string {
	return GetEnv("TRIEUGENE_S3_SECRET_KEY", "trieugene_secret_key")
}

func (c *Config) S3URL() string {
	return GetEnv("TRIEUGENE_S3_URL", "trieugene.myshopify.io:8000")
}

func (c *Config) S3Region() string {
	return GetEnv("TRIEUGENE_S3_REGION", "auto")
}

func (c *Config) StatsdURL() string {
	return GetEnv("TRIEUGENE_STATSD_URL", "127.0.0.1:8125")
}
