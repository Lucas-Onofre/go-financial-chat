package model

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}
