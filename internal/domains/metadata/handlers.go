package metadata

import (
	"agentic/commerce/internal/utils"
	"agentic/commerce/pkg/apperror"
	"agentic/commerce/pkg/logger"
	"agentic/commerce/pkg/specs/api"

	"github.com/labstack/echo/v4"
)

type IContentResource interface {
	CreateMetadata() echo.HandlerFunc
	GetMetadata() echo.HandlerFunc
	ListMetadata() echo.HandlerFunc
}

type contentResource struct {
	ContentService IContentService
	Logger         *logger.AppLogger
}

func NewContentResource(service IContentService, logger *logger.AppLogger) IContentResource {
	return &contentResource{
		ContentService: service,
		Logger:         logger.WithScope(contentResource{}),
	}
}

func (v *contentResource) CreateMetadata() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var req api.MetadataRequest

		reqCtx := ctx.Request().Context()

		v.Logger.Info("contentService.CreateMetadata called")

		err := ctx.Bind(&req)
		if err != nil {
			return utils.ErrorResponse(ctx, apperror.ErrBadRequest, "Check your input: "+err.Error())
		}

		resp, err := v.ContentService.CreateMetaData(reqCtx, &req)
		if err != nil {
			return utils.ErrorResponse(ctx, err, "Cant create the metadata")
		}

		return utils.SuccessResponse(ctx, resp)
	}
}

func (v *contentResource) GetMetadata() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var req api.MetadataIDAwareRequest
		reqCtx := ctx.Request().Context()

		err := ctx.Bind(&req)
		if err != nil {
			return utils.ErrorResponse(ctx, apperror.ErrBadRequest, "Check your input: "+err.Error())
		}
		v.Logger.Info("contentService.ListMetadata called")

		resp, err := v.ContentService.GetMetaData(reqCtx, &req)
		if err != nil {
			return utils.ErrorResponse(ctx, err, "Cant create the metadata")
		}

		return utils.SuccessResponse(ctx, resp)
	}
}

func (v *contentResource) ListMetadata() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		reqCtx := ctx.Request().Context()

		v.Logger.Info("contentService.ListMetadata called")

		resp, err := v.ContentService.ListMetaData(reqCtx)
		if err != nil {
			return utils.ErrorResponse(ctx, err, "Cant list the metadata")
		}

		return utils.SuccessResponse(ctx, resp)
	}
}
