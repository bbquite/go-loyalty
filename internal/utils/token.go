package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const ACCESS_TOKEN_EXPIRY_HOUR = time.Hour * 3
const ACCESS_TOKEN_SECRET = "goloyaltydiplomsecter"

var ()

type JWT struct {
	Token string `json:"token"`
}

type Claims struct {
	jwt.RegisteredClaims
	Userid uint32
}

func CreateAccessToken(userid uint32) (accessToken string, err error) {
	log.Printf("create token func user_id: %d", userid)
	claims := &Claims{
		Userid: userid,
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

func ExtractIDFromToken(requestToken string) (uint32, error) {
	token, err := jwt.ParseWithClaims(requestToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ACCESS_TOKEN_SECRET), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		userID := claims.Userid
		return userID, nil
	}

	return 0, fmt.Errorf("invalid Token")
}
