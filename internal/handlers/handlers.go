package handlers

import (
	"net/http"

	"github.com/bbquite/go-loyalty/internal/middleware"
	"github.com/bbquite/go-loyalty/internal/services"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	services *services.AppService
	logger   *zap.SugaredLogger
}

func NewHandler(services *services.AppService, logger *zap.SugaredLogger) (*Handler, error) {
	return &Handler{
		services: services,
		logger:   logger,
	}, nil
}

func (h *Handler) InitRoutes() *chi.Mux {
	chiRouter := chi.NewRouter()

	chiRouter.Use(middleware.RequestsLoggingMiddleware(h.logger))
	chiRouter.Use(middleware.GzipMiddleware)

	chiRouter.Route("/api/user/", func(r chi.Router) {
		r.Post("/register/", h.registerUser)
		r.Post("/login/", h.loginUser)
		r.Route("/orders/", func(r chi.Router) {
			r.Post("/", h.orderSend)
			r.Get("/", h.ordersList)
		})
		r.Route("/balance/", func(r chi.Router) {
			r.Get("/", h.userBalance)
			r.Post("/withdraw/", h.userBalanceWithdraw)
		})
		r.Get("/withdrawals/", h.withdrawHistory)
	})

	return chiRouter
}

func (h *Handler) registerUser(w http.ResponseWriter, r *http.Request) {
	h.services.RegisterUser()
}

func (h *Handler) loginUser(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) orderSend(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) ordersList(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) userBalance(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) userBalanceWithdraw(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) withdrawHistory(w http.ResponseWriter, r *http.Request) {

}
