package services

import (
	"database/sql"
	"errors"
	"github.com/bbquite/go-loyalty/internal/models"
	"github.com/bbquite/go-loyalty/internal/utils"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrIncorrectLoginData = errors.New("incorrect login or password")

	ErrPurchaseAlreadySend = errors.New("purchase already exists")
	ErrPurchaseConflict    = errors.New("purchase conflict with another user")
)

type StorageRepo interface {
	CreateAccount(username string, password string) (uint32, error)
	GetAccountByUsername(username string) (models.Account, error)
	GetAccountByLoginData(username string, password string) (models.Account, error)

	GetAccountBalance(accountID uint32) (models.Balance, error)
	GetAccountBalanceHistory(accountID uint32, trType string) ([]models.BalanceHistory, error)

	GetPurchase(purchaseID string) (models.Purchase, error)
	CreatePurchase(accountID uint32, purchaseID string) error
	GetPurchasesList(accountID uint32) ([]models.Purchase, error)
}

type AppService struct {
	store  StorageRepo
	logger *zap.SugaredLogger
}

func NewAppService(store StorageRepo, logger *zap.SugaredLogger) *AppService {
	return &AppService{
		store:  store,
		logger: logger,
	}
}

func (service *AppService) RegisterUser(userData *models.UserLoginData) (utils.JWT, error) {
	var token utils.JWT

	_, err := service.store.GetAccountByUsername(userData.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			userID, err := service.store.CreateAccount(
				userData.Username, utils.GenerateSHAString(userData.Password))
			if err != nil {
				return token, err
			}

			tokenString, err := utils.CreateAccessToken(userID)
			if err != nil {
				return token, err
			}

			token.Token = tokenString
			return token, nil
		}
		return token, err
	}
	return token, ErrUserAlreadyExists
}

func (service *AppService) LoginUser(userData *models.UserLoginData) (utils.JWT, error) {
	var token utils.JWT

	shaInputPassword := utils.GenerateSHAString(userData.Password)

	account, err := service.store.GetAccountByLoginData(userData.Username, shaInputPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return token, ErrIncorrectLoginData
		}
		return token, err
	}

	tokenString, err := utils.CreateAccessToken(account.Id)
	if err != nil {
		return token, err
	}

	token.Token = tokenString
	return token, nil
}

func (service *AppService) GetAccountBalance(accountID uint32) (models.Balance, error) {
	response, err := service.store.GetAccountBalance(accountID)
	if err != nil {
		return models.Balance{}, err
	}
	return response, nil
}

func (service *AppService) GetAccountBalanceHistory(accountID uint32, trType string) ([]models.BalanceHistory, error) {
	response, err := service.store.GetAccountBalanceHistory(accountID, trType)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (service *AppService) SendPurchase(accountID uint32, purchaseID string) error {

	purchase, err := service.store.GetPurchase(purchaseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err := service.store.CreatePurchase(accountID, purchaseID)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	if purchase.AccountID == accountID {
		return ErrPurchaseAlreadySend
	}
	return ErrPurchaseConflict
}

func (service *AppService) PurchasesList(accountID uint32) ([]models.Purchase, error) {
	response, err := service.store.GetPurchasesList(accountID)
	if err != nil {
		return nil, err
	}
	return response, nil
}
