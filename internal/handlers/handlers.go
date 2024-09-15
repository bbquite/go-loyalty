package handlers

import (
	"database/sql"
	"html/template"

	"github.com/bbquite/go-loyalty/internal/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	// services      *service.MetricService
	indexTemplate *template.Template
	logger        *zap.SugaredLogger
	db            *sql.DB
}

func NewHandler(logger *zap.SugaredLogger, db *sql.DB) (*Handler, error) {
	return &Handler{
		// services: services,
		logger: logger,
	}, nil
}

func (h *Handler) initRoutes() *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.RequestsLoggingMiddleware(h.logger))
	chiRouter.Use(middleware.GzipMiddleware)

	chiRouter.Route("/", func(r chi.Router) {
		r.Get("/", h.renderMetricsPage)
		r.Get("/ping/", h.databasePing)
		r.Route("/value/", func(r chi.Router) {
			r.Post("/", h.valueMetricJSON)
			r.Get("/{m_type}/{m_name}", h.valueMetricURI)
		})
		r.Route("/update/", func(r chi.Router) {
			r.Post("/", h.updateMetricJSON)
			r.Post("/{m_type}/{m_name}/{m_value}", h.updateMetricURI)
		})
	})

	return chiRouter
}
