package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EClaesson/go-luhn"
	"github.com/bbquite/go-loyalty/internal/middleware"
	"github.com/bbquite/go-loyalty/internal/models"
	"github.com/bbquite/go-loyalty/internal/services"
	"github.com/bbquite/go-loyalty/internal/utils"
	"github.com/go-chi/chi/v5"
	ChiMiddleware "github.com/go-chi/chi/v5/middleware"
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

	chiRouter.Use(ChiMiddleware.Logger)
	// chiRouter.Use(middleware.GzipMiddleware)

	chiRouter.Route("/api/user/", func(r chi.Router) {
		r.Post("/register/", h.registerAccount)
		r.Post("/login/", h.loginAccount)
		r.Route("/orders/", func(r chi.Router) {
			r.Post("/", middleware.TokenAuthMiddleware(h.sendPurchase))
			r.Get("/", middleware.TokenAuthMiddleware(h.purchasesList))
		})
		r.Route("/balance/", func(r chi.Router) {
			r.Get("/", middleware.TokenAuthMiddleware(h.accountBalance))
			r.Post("/withdraw/", middleware.TokenAuthMiddleware(h.accountBalanceWithdraw))
		})
		r.Get("/withdrawals/", middleware.TokenAuthMiddleware(h.accountWithdrawHistory))
	})

	return chiRouter
}

func (h *Handler) registerAccount(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer
	var reqData models.UserLoginData

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Debug(err)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &reqData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.logger.Debug(err)
		return
	}

	if reqData.Username == "" || reqData.Password == "" {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		h.logger.Debug(err)
		return
	}

	token, err := h.services.RegisterUser(&reqData)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg, _ := json.Marshal(token)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token.Token)

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (h *Handler) loginAccount(w http.ResponseWriter, r *http.Request) {
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

	if reqData.Username == "" || reqData.Password == "" {
		http.Error(w, "invalid request format", http.StatusBadRequest)
		h.logger.Debug(err)
		return
	}

	token, err := h.services.LoginUser(&reqData)
	if err != nil {
		if errors.Is(err, services.ErrIncorrectLoginData) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	msg, _ := json.Marshal(token)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", "Bearer "+token.Token)

	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (h *Handler) accountBalance(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(utils.AccountIDContextKey).(uint32)
	balance, err := h.services.GetAccountBalance(accountID)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	msg, _ := json.Marshal(balance)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (h *Handler) sendPurchase(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accountID := r.Context().Value(utils.AccountIDContextKey).(uint32)
	orderID := buf.String()

	if orderID == "" {
		http.Error(w, "invalid request format", http.StatusUnprocessableEntity)
		return
	}

	orderValid, err := luhn.IsValid(orderID)
	if err != nil {
		http.Error(w, "invalid request format", http.StatusUnprocessableEntity)
		h.logger.Debug(err)
		return
	}

	if !orderValid {
		http.Error(w, "invalid request format", http.StatusUnprocessableEntity)
		return
	}

	err = h.services.SendPurchase(accountID, orderID)
	if err != nil {
		h.logger.Debug(err)
		if errors.Is(err, services.ErrPurchaseAlreadySend) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, services.ErrPurchaseConflict) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	return
}

func (h *Handler) purchasesList(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(utils.AccountIDContextKey).(uint32)
	purchaseList, err := h.services.PurchasesList(accountID)
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(purchaseList) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	msg, _ := json.Marshal(purchaseList)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}

func (h *Handler) accountBalanceWithdraw(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) accountWithdrawHistory(w http.ResponseWriter, r *http.Request) {
	accountID := r.Context().Value(utils.AccountIDContextKey).(uint32)
	balanceHistory, err := h.services.GetAccountBalanceHistory(accountID, "OUT")
	if err != nil {
		h.logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(balanceHistory) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	msg, _ := json.Marshal(balanceHistory)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(msg)
}
