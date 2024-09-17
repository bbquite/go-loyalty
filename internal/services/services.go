package services

import (
	"database/sql"
	"errors"

	"github.com/bbquite/go-loyalty/internal/models"
	"github.com/bbquite/go-loyalty/internal/utils"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists = errors.New("the user already exists")
)

type StorageRepo interface {
	CreateAccount(username string, password string) (int64, error)
	GetAccountByUsername(username string) (models.Account, error)
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

func (service *AppService) RegisterUser(registerData *models.UserLoginData) (utils.JWT, error) {
	var token utils.JWT

	_, err := service.store.GetAccountByUsername(registerData.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			userID, err := service.store.CreateAccount(
				registerData.Username, utils.GenerateSHAString(registerData.Password))
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
