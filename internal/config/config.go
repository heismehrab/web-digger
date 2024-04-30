package config

import (
	"log/slog"
	"web-digger/pkg/logger"
)

type Config struct {
	App    AppSettings   `mapstructure:",squash"`
	Logger logger.Config `mapstructure:",squash"`
}

type AppSettings struct {
	Port int `mapstructure:"APP_PORT"`
}

type Logger struct {
	LogLevel      slog.Level `mapstructure:"LOGGER_LOG_LEVEL"`
	GrayLogActive bool       `mapstructure:"LOGGER_GRAYLOG_ACTIVE"`
	GrayLogServer string     `mapstructure:"LOGGER_GRAYLOG_SERVER"`
	GrayLogStream string     `mapstructure:"LOGGER_GRAYLOG_STREAM"`
}
