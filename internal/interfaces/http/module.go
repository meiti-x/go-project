package http

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"http",
	fx.Provide(NewServer),
	fx.Invoke(CommonRoutes),
	fx.Invoke(RegisterHTTPServer),
)
