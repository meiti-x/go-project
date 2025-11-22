package handlers

import (
	"net/http"

	"agentic/commerce/config"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type IHealthResource interface {
	Ping() echo.HandlerFunc
	Liveness() echo.HandlerFunc
	Readiness() echo.HandlerFunc
}

type healthResource struct {
	db   *gorm.DB
	mode config.ModeEnum
}

func NewHealthResource(
	db *gorm.DB,
	mode config.ModeEnum,
) IHealthResource {
	return &healthResource{
		db:   db,
		mode: mode,
	}
}

func (v *healthResource) Ping() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return ctx.HTML(http.StatusOK, "Pong")
	}
}

func (v *healthResource) Liveness() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tx := v.db.Exec("SELECT 1")
		if tx.Error != nil {
			return ctx.HTML(http.StatusInternalServerError, "Database failed to pong")
		}

		return ctx.HTML(http.StatusOK, "Success")
	}
}

func (v *healthResource) Readiness() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// Database Ping
		tx := v.db.Exec("SELECT 1")
		if tx.Error != nil {
			return ctx.HTML(http.StatusInternalServerError, "Database failed to pong")
		}

		return ctx.HTML(http.StatusOK, "Success")
	}
}
