package http

import (
	"go/types"

	"agentic/commerce/config"
	"agentic/commerce/internal/interfaces/http/handlers"

	echoAdapter "github.com/TickLabVN/tonic/adapters/echo"
	"gorm.io/gorm"
)

func CommonRoutes(
	s *Server,
	db *gorm.DB,
	mode config.ModeEnum,
) *Server {

	echoAdapter.UIHandle(s.Router, s.Spec, "/swagger-ui")

	healthResource := handlers.NewHealthResource(db, mode)
	healthGroup := s.Router.Group("/health")

	echoAdapter.AddRoute[types.Nil, types.Nil](s.Spec,
		healthGroup.GET("/ping", healthResource.Ping()),
	)
	echoAdapter.AddRoute[types.Nil, types.Nil](s.Spec,
		healthGroup.GET("/liveness", healthResource.Liveness()),
	)
	echoAdapter.AddRoute[types.Nil, types.Nil](s.Spec,
		healthGroup.GET("/readiness", healthResource.Readiness()),
	)

	return s
}
