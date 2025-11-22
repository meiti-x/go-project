package metadata

import (
	"agentic/commerce/internal/interfaces/http/middleware"
	"agentic/commerce/pkg/apperror"
	"agentic/commerce/pkg/logger"
	"agentic/commerce/pkg/specs/api"

	"github.com/google/uuid"

	"context"
)

type IContentService interface {
	CreateMetaData(ctx context.Context, req *api.MetadataRequest) (*api.MetadataResponse, error)
	GetMetaData(ctx context.Context, req *api.MetadataIDAwareRequest) (*api.MetadataItemResponse, error)
	ListMetaData(ctx context.Context) ([]api.MetadataItemResponse, error)
}

type contentService struct {
	repository IContentRepository
	logger     *logger.AppLogger
	mappers    IContentMapper
}

func NewContentService(
	logger *logger.AppLogger,
	repository IContentRepository,
	mappers IContentMapper,

) IContentService {
	return &contentService{
		repository: repository,
		logger:     logger.WithScope(&contentService{}),
		mappers:    mappers,
	}
}

func (s *contentService) CreateMetaData(ctx context.Context, req *api.MetadataRequest) (*api.MetadataResponse, error) {
	model, err := s.mappers.mapContentRequestToModel(req, uuid.New().String())
	if err != nil {
		return nil, err
	}
	err = s.repository.Create(ctx, model)
	if err != nil {
		return nil, apperror.ErrServer
	}

	return &api.MetadataResponse{
		UUID: *model.UUid,
	}, err
}
func (s *contentService) GetMetaData(ctx context.Context, req *api.MetadataIDAwareRequest) (*api.MetadataItemResponse, error) {
	res, err := s.repository.GetByUserID(ctx, req.ID, middleware.GetUserID(ctx))
	if err != nil {
		return nil, apperror.ErrServer
	}

	return s.mappers.mapToMetadataItem(res), err
}

func (s *contentService) ListMetaData(ctx context.Context) ([]api.MetadataItemResponse, error) {
	res, err := s.repository.ListByUserID(ctx, middleware.GetUserID(ctx))
	if err != nil {
		return nil, apperror.ErrServer
	}

	return s.mappers.mapToMetadataList(res), err
}
