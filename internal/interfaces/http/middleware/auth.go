package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func WithAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") == "" && !strings.Contains(c.Request().URL.String(), "swagger") {
			c.Error(echo.NewHTTPError(http.StatusUnauthorized, "please login first"))
			return nil
		}

		userID := int64(111)

		ctx := context.WithValue(c.Request().Context(), "userId", userID)

		req := c.Request().WithContext(ctx)
		c.SetRequest(req)
		return next(c)
	}
}

func GetUserID(ctx context.Context) int64 {
	if v := ctx.Value("userId"); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}
