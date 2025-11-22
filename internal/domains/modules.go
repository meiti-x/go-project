package domains

import (
	"agentic/commerce/internal/domains/metadata"
	"agentic/commerce/internal/infrastructure/database"

	"go.uber.org/fx"
)

var Modules = fx.Module(
	"domains",

	fx.Provide(database.CreateGormDB),
	metadata.Module,
)
