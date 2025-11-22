package utils

import (
	"errors"
	"log"
	"net/http"

	"agentic/commerce/pkg/apperror"
	"agentic/commerce/pkg/specs/api"

	"github.com/labstack/echo/v4"
)

const ApiStatusSuccess = "success"
const ApiStatusError = "error"

var InternalServerErrorMessage = "something went wrong"

func SuccessResponse[TData any](ctx echo.Context, data TData) error {
	return SuccessResponseWithMessage(ctx, "", data)
}

func SuccessResponseWithMessage[TData any](ctx echo.Context, message string, data TData) error {
	return ctx.JSON(http.StatusOK, api.APIResponse[TData]{
		BaseResponse: api.BaseResponse{
			Status:  ApiStatusSuccess,
			Message: []string{message},
		},
		Data: data,
	})
}

func ErrorResponse(c echo.Context, err error, message ...string) error {
	var appErr *apperror.ErrorWithStatus

	if errors.As(err, &appErr) {
		var m []string
		if message != nil {
			m = message
		} else {
			m = []string{appErr.Err.Message}
		}

		return c.JSON(appErr.StatusCode, api.BaseResponse{
			Status:  "error",
			Message: m,
		})
	}

	log.Printf("Unhandled error: %v", err)
	return c.JSON(http.StatusInternalServerError, api.BaseResponse{
		Status:  ApiStatusError,
		Message: []string{InternalServerErrorMessage},
	})
}
