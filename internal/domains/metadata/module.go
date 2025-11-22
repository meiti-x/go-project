package metadata

import (
	"agentic/commerce/internal/infrastructure/database"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"content-request",
	fx.Provide(NewContentRepository),
	fx.Provide(NewContentMapper),
	fx.Provide(NewContentService),
	fx.Invoke(RegisterRoutes),
	database.AsModel(&MetaDataModel{}),
)
