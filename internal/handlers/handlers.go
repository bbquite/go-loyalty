package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bbquite/go-loyalty/internal/middleware"
	"github.com/bbquite/go-loyalty/internal/models"
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
			r.Post("/", middleware.TestMW(h.orderSend))
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

	var buf bytes.Buffer
	var reqData models.UserLoginData

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Info(err)
		return
	}

	token, err := h.services.RegisterUser(&reqData)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		h.logger.Error(err)
	}

	msg, _ := json.Marshal(token)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
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
