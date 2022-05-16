package models

import "github.com/golang-jwt/jwt"

type AppClaims struct {
	jwt.StandardClaims
	UserId string `json:"userId"`
}
