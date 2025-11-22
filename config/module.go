package config

import (
	"go.uber.org/fx"
)

type configSupply struct {
	fx.Out

	Mode     ModeEnum
	HTTP     *HttpConfig
	Database *DbConfig
	Logger   *Logger
}

func provideNestedConfigs(cfg *Config) configSupply {
	return configSupply{
		Mode:     cfg.Mode,
		HTTP:     cfg.Http,
		Database: cfg.Database,
		Logger:   cfg.Logger,
	}
}

var Module = fx.Module(
	"config",
	fx.Provide(provideNestedConfigs),
)
