package services

import (
	"database/sql"
	"errors"

	"github.com/bbquite/go-loyalty/internal/models"
	"github.com/bbquite/go-loyalty/internal/utils"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists  = errors.New("the user already exists")
	ErrIncorrectLoginData = errors.New("incorrect login or password")
)

type StorageRepo interface {
	CreateAccount(username string, password string) (uint32, error)
	GetAccountByUsername(username string) (models.Account, error)
	GetAccountByLoginData(username string, password string) (models.Account, error)

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

func (service *AppService) RequestPurchase(accountID uint32, purchaseID string) error {
	err := service.store.CreatePurchase(accountID, purchaseID)
	if err != nil {
		return err
	}
	return nil
}

func (service *AppService) PurchasesList(accountID uint32) ([]models.Purchase, error) {
	response, err := service.store.GetPurchasesList(accountID)
	if err != nil {
		return nil, err
	}
	return response, nil
}
