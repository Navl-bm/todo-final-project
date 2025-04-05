// Пакет api используется для взаимодействия с API приложения

package api

import (
	"github.com/Navl-bm/todo-final-project/config"
	"github.com/Navl-bm/todo-final-project/database"
)

// Метод ApiHandler для хранения конфигурации приложения
type ApiHandler struct {
	DB     *database.DB
	Config *config.Config
}

// NewApiHandler создает новый экземпляр для взаимодействия с API
func NewApiHandler(cfg *config.Config, db *database.DB) *ApiHandler {
	return &ApiHandler{Config: cfg, DB: db}
}
