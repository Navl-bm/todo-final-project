// Пакет middlewares предназначен для проверки авторизации
// пользователя перед выполнением защищенных запросов

package middlewares

import "github.com/Navl-bm/todo-final-project/config"

// Метод взаимодействия
type Middleware struct {
	Config *config.Config
}

// Создание нового экземпляра
func NewMiddleware(cfg *config.Config) *Middleware {
	return &Middleware{Config: cfg}
}
