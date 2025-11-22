package metadata

import (
	"go/types"

	"agentic/commerce/internal/interfaces/http"
	"agentic/commerce/pkg/logger"
	"agentic/commerce/pkg/specs/api"

	echoAdapter "github.com/TickLabVN/tonic/adapters/echo"
)

func RegisterRoutes(s *http.Server, shipmentService IContentService, logger *logger.AppLogger) *http.Server {
	contentResourceObj := NewContentResource(shipmentService, logger)

	apis := s.Router.Group("/metadata")

	echoAdapter.AddRoute[api.MetadataRequest, api.APIResponse[types.Nil]](s.Spec,
		apis.POST("", contentResourceObj.CreateMetadata()),
	)

	echoAdapter.AddRoute[api.MetadataIDAwareRequest, api.APIResponse[api.MetadataItemResponse]](s.Spec,
		apis.GET("/:id", contentResourceObj.GetMetadata()),
	)

	echoAdapter.AddRoute[types.Nil, api.ApiPaginateResponse[api.MetadataItemResponse]](s.Spec,
		apis.GET("/list", contentResourceObj.ListMetadata()),
	)

	return s
}
