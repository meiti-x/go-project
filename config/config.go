package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type ModeEnum string

const (
	ModeDev  ModeEnum = "dev"
	ModeProd ModeEnum = "prod"
)

type Config struct {
	Mode     ModeEnum `yaml:"mode"`
	Http     *HttpConfig
	Database *DbConfig
	Logger   *Logger
}

type LogLevel string

const (
	LogDisabled = "disabled"
	LevelNone   = "none"
	LevelTrace  = "trace"
	LevelDebug  = "debug"
	LevelInfo   = "info"
	LevelWarn   = "warn"
	LevelError  = "error"
	LevelFatal  = "fatal"
	LevelPanic  = "panic"
)

type Logger struct {
	Level LogLevel `yaml:"level"`
}

type DbConfig struct {
	Host                 string        `yaml:"host"`
	Database             string        `yaml:"database"`
	User                 string        `yaml:"user"`
	Password             string        `yaml:"password"`
	Port                 int           `yaml:"port"`
	DialTimeout          time.Duration `yaml:"dialTimeout"`
	Timeout              time.Duration `yaml:"timeout"`
	DialRetry            int           `yaml:"dialRetry"`
	MaxConn              int           `yaml:"maxConn"`
	IdleConn             int           `yaml:"idleConn"`
	DisableWithReturning bool          `yaml:"disableWithReturning"`
	Location             string        `yaml:"location"`
}

func (c *DbConfig) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.Location,
	)
}

type HttpClientConfig struct {
	DialTimeout *time.Duration `yaml:"dialTimeout"`
	TlsTimeout  *time.Duration `yaml:"tlsTimeout"`
	Timeout     *time.Duration `yaml:"timeout"`
}

type HttpConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (c *Config) IsDev() bool {
	return c.Mode == ModeDev
}

func (c *Config) IsProd() bool {
	return c.Mode == ModeProd
}

func (c *HttpConfig) UrlString() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func ReadConfig(configPath string) (*Config, error) {
	var cfg Config

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	mode := v.GetString("mode")
	overrideFile := "config." + mode + ".yml"
	override := filepath.Join(filepath.Dir(configPath), overrideFile)

	if _, err := os.Stat(override); err == nil {
		v.SetConfigFile(override)
		if err := v.MergeInConfig(); err != nil {
			return nil, err
		}
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
