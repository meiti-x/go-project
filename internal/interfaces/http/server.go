package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"agentic/commerce/config"
	"agentic/commerce/internal/interfaces/http/middleware"

	m "github.com/labstack/echo/v4/middleware"

	"github.com/TickLabVN/tonic/core/docs"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

type Server struct {
	Router *echo.Echo
	Spec   *docs.OpenApi
}

func NewServer() *Server {
	engine := echo.New()
	engine.JSONSerializer = &middleware.JsonV2{}
	engine.Use(m.RemoveTrailingSlash())
	engine.Use(middleware.WithRecoverMiddleware)
	engine.Use(middleware.WithAuthMiddleware)

	apiDoc := &docs.OpenApi{
		OpenAPI: "3.0.1",
		Info: docs.InfoObject{
			Version: "1.0.0",
			Title:   "Swagger API documentation",
			Contact: &docs.ContactObject{
				Name:  "Author",
				URL:   "https://github.com/meiti-x",
				Email: "mahdimomeni6@gmail.com",
			},
		},
		Components: docs.ComponentsObject{
			SecuritySchemes: map[string]docs.SecuritySchemeOrReference{
				"BearerAuth": {
					SecuritySchemeObject: &docs.SecuritySchemeObject{
						Type: "apiKey",
						In:   "header",
						Name: "Authorization",
					},
				},
			},
		},
		Security: []docs.SecurityRequirement{
			{"BearerAuth": {}},
		},
	}

	return &Server{
		Router: engine,
		Spec:   apiDoc,
	}
}

func (s *Server) Start(cfg *config.HttpConfig) error {
	address := cfg.UrlString()
	if err := s.Router.Start(address); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	fmt.Println("🛑 Gracefully shutting down HTTP server...")
	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	return s.Router.Shutdown(ctxTimeout)
}

func RegisterHTTPServer(lc fx.Lifecycle, cfg *config.HttpConfig, srv *Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func(srv *Server, cfg *config.HttpConfig) {
				err := srv.Start(cfg)
				if err != nil {
					_ = fmt.Errorf("HTTP server error: %w", err)

				}
			}(srv, cfg)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			err := srv.Shutdown(ctx)
			if err != nil {
				return err
			}
			return nil
		},
	})
}
