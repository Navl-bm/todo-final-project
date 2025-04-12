package api

import (
	"log"
	"net/http"

	"github.com/Navl-bm/todo-final-project/errors"
	"github.com/Navl-bm/todo-final-project/models"
	"github.com/Navl-bm/todo-final-project/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SignIN формирует JWT токен для пользователя при авторизации
func (h *ApiHandler) SignIn(ctx *gin.Context) {
	loginData := models.LoginData{}
	err := ctx.ShouldBindBodyWithJSON(&loginData)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrIncorrectPassword.Error()})
		return
	}

	if loginData.Password != h.Config.UserPassword {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrIncorrectPassword.Error()})
		return
	}

	tokenData := jwt.NewWithClaims(jwt.SigningMethodHS256, models.JWTToken{
		PasswordHash: auth.GetSha256Password(h.Config.UserPassword),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "todo-final",
		},
	})

	token, err := tokenData.SignedString([]byte(h.Config.JWTSecret))
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServerError.Error()})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"token": token})
}
