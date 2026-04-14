package server

import (
	"log/slog"
	"net/http"

	"github.com/Strangebrewer/go-tracer/health"
	"github.com/Strangebrewer/go-tracer/middleware"
	"github.com/Strangebrewer/go-tracer/span"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Server struct {
	HTTPServer *http.Server
}

func New(addr string, allowedOrigins []string, store *span.Store, authMiddleware, serviceKeyMiddleware func(http.Handler) http.Handler) *Server {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Service-Key"},
		MaxAge:         300,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(slog.Default()))
	r.Use(chimiddleware.Recoverer)

	r.Get("/health", health.Handler)
	r.Mount("/", span.Routes(store, serviceKeyMiddleware, authMiddleware))

	return &Server{
		HTTPServer: &http.Server{
			Addr:    addr,
			Handler: r,
		},
	}
}
