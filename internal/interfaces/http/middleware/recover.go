package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
)

const TemporaryDebugMode = true //TODO

// WithRecoverMiddleware is a custom panic recovery middleware for Echo
func WithRecoverMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC RECOVERED] %v\nStack Trace:\n%s", r, string(debug.Stack()))

				if TemporaryDebugMode {
					c.Error(echo.NewHTTPError(http.StatusInternalServerError, r))
				}
				c.Error(echo.NewHTTPError(http.StatusInternalServerError, "internal server error"))
			}
		}()

		return next(c)
	}
}
