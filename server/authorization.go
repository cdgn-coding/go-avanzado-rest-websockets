package server

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"go-rest-websockets/models"
	"time"
)

type Authorization struct{}

func NewAuthorization() Authorization {
	return Authorization{}
}

func (a Authorization) SignToken(secretKey string, userId string) (string, error) {
	claims := models.AppClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
		UserId: userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

func (a Authorization) ParseAndVerifyToken(secretKey string, tokenString string) (*models.AppClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*models.AppClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	return claims, nil
}
