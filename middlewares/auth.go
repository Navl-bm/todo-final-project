package middlewares

import (
	"net/http"

	"github.com/Navl-bm/todo-final-project/errors"
	"github.com/Navl-bm/todo-final-project/models"
	"github.com/Navl-bm/todo-final-project/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// CheckAuthToken проверяет токен пользователя
//
// Берет токен из cookie и проверяет его на валидность
func (m *Middleware) CheckAuthToken(ctx *gin.Context) {
	if m.Config.UserPassword == "" {
		ctx.Next()
		return
	}

	tokenString, err := ctx.Cookie("token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized.Error()})
		return
	}

	claims := models.JWTToken{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(m.Config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized.Error()})
		return
	}

	if claims.PasswordHash != utils.GetSha256Password(m.Config.UserPassword) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.ErrUnauthorized.Error()})
		return
	}

	ctx.Next()
}
