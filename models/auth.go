package models

import "github.com/golang-jwt/jwt/v5"

// LoginData модель передаваемых данных при авторизации
type LoginData struct {
	Password string `json:"password"`
}

// JWTToken модель токена пользователя
type JWTToken struct {
	PasswordHash string
	jwt.RegisteredClaims
}
