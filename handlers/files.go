// Пакет handlers используется для описания обработчиков http запросов

package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Navl-bm/todo-final-project/config"
	"github.com/Navl-bm/todo-final-project/errors"
	"github.com/gin-gonic/gin"
)

// FilesHandler метод для взаимодействия с файлами
type FilesHandler struct {
	Config *config.Config
}

// Создание нового экземпляра FilesHandler
func NewFilesHandler(cfg *config.Config) *FilesHandler {
	return &FilesHandler{Config: cfg}
}

// checkFileExtAllowed используется дла проверки расширения запрашиваемого файла
// с разрешенным списком расширений из файла конфигурации
func (h *FilesHandler) checkFileExtAllowed(filePath string) bool {
	fileExt := filepath.Ext(filePath)

	lookup := make(map[string]bool)
	for _, s := range h.Config.AllowedExtFiles {
		lookup[s] = true
	}

	return lookup[fileExt]
}

// Get используется для отправки файлов клиенту
func (h *FilesHandler) Get(ctx *gin.Context) {
	var FileQuery string
	if ctx.Request.URL.Path == "/" {
		FileQuery = fmt.Sprintf("/%s", h.Config.IndexFileName)
	} else {
		FileQuery = ctx.Request.URL.Path
	}
	FullFilePath := fmt.Sprintf("%s%s", h.Config.WebDir, FileQuery)
	log.Println(FullFilePath)

	if !h.checkFileExtAllowed(FullFilePath) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": errors.ErrAccessDenied.Error()})
		return
	}

	_, err := os.Stat(FullFilePath)
	if os.IsNotExist(err) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": errors.ErrFileNotFound.Error()})
		return
	}

	ctx.File(FullFilePath)
}
