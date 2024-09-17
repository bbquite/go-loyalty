package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const ACCESS_TOKEN_EXPIRY_HOUR = time.Hour * 3
const ACCESS_TOKEN_SECRET = "goloyaltydiplomsecter"

type JWT struct {
	Token string `json:"token"`
}

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func CreateAccessToken(userid int64) (accessToken string, err error) {
	claims := &Claims{
		UserID: string(userid),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ACCESS_TOKEN_EXPIRY_HOUR)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(ACCESS_TOKEN_SECRET))
	if err != nil {
		return "", err
	}
	return t, err
}

func IsAuthorized(requestToken string) (bool, error) {
	_, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ACCESS_TOKEN_SECRET), nil
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func ExtractIDFromToken(requestToken string) (string, error) {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ACCESS_TOKEN_SECRET), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok && !token.Valid {
		return "", fmt.Errorf("invalid Token")
	}

	log.Print(claims)

	return claims["UserID"].(string), nil
}
