package app

import (
	"context"
	"fmt"

	"agentic/commerce/pkg/logger"
	"go.uber.org/fx"
)

func registerLoggerShutdown(lc fx.Lifecycle, r *logger.AppLogger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			fmt.Println("🛑 Gracefully stopping logger...")
			return r.Close()
		},
	})
}

var LoggerModule = fx.Module(
	"logger",
	fx.Provide(logger.NewAppLogger),
	fx.Invoke(func(l *logger.AppLogger) {
		l.InitLogger()
	}),
	fx.Invoke(registerLoggerShutdown),
)
