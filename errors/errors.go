// Пакет errors содержит в себе константы для хранения текстов встречающихся ошибок

package errors

import (
	"errors"
)

var (
	ErrFileNotFound        = errors.New("файл не найден")
	ErrAccessDenied        = errors.New("доступ запрещен")
	ErrBadRepeatFormat     = errors.New("некорректный формат повторения")
	ErrBadTime             = errors.New("некорректный формат времени")
	ErrBadJSON             = errors.New("некорректный формат JSON")
	ErrBadTitle            = errors.New("не указан заголовок задачи")
	ErrInternalServerError = errors.New("внутренняя ошибка сервера")
	ErrTaskNotFound        = errors.New("задача не найдена")
	ErrIncorrectPassword   = errors.New("неправильный пароль")
	ErrUnauthorized        = errors.New("требуется авторизация")
)
