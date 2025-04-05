# TODO-APP
Приложение предназначено для формирования и отслеживания списка задач.
При выполнении повторяющейся задачи она автоматически переносится на следующую ближайшую дату в соответствии с настойками повтора выполнения

## Выполненные дополнительные задачи
- [x]  Настройка порта сервера с помощью переменных окружения (`PORT`)
- [x]  Настройка пути к БД с помощью переменных окружения (`DBFILE`)
- [x]  Поиск ближайшей даты в случае повтора по неделе и месяцу
- [x]  Поиск задач по заданному фильтру
- [x]  Аутентификация пользователя (`PASSWORD` и `JWT_SECRET`)
- [x]  Управление каталога фронтенда через переменные окружения (`WEB_DIR`)
- [x]  Фильтрация запросов для файлов и страниц с помощью переменных окружения (`ALLOWED_EXT_FILES`)
- [x]  Возможность изменения имени стартовой страницы приложения с помощью пременных окружения (`INDEX_FILE_NAME`)
- [x]  Сборка docker контейнера

## Локальный запуск
### Пример файла .env  

```env
PORT=7540
WEB_DIR=./web
INDEX_FILE_NAME=index.html
ALLOWED_EXT_FILES=.html,.js,.css,.ico
DBFILE=scheduler.db
PASSWORD=12345
JWT_SECRET=jwt_secret_string
```

### Запуск приложения
```sh
go run ./cmd/server/main.go
```
Адрес приложения ```http://localhost:7540/```

## Тестирование приложения

### Файл `./tests/settings.go`
```GO
package tests

var Port = 7540
var DBFile = "../scheduler.db"
var FullNextDate = true
var Search = true
var Token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJQYXNzd29yZEhhc2giOiJXWlJIR3JzQkVTcjh3WUZaOXN4MHRQVVJ1WmdHMmxtenl2V3B3WFBLejhVPSIsImlzcyI6InRvZG8tZmluYWwifQ.J6AybyDrzNyOmnvtPTZ9fMObv-KXJmmg6SLmbam59kc`
```

### Запуск тестов
```sh
go test ./tests
```

## Запуск docker конейтнера
```docker
docker build --tag todo-final-app:v1 .
docker run -d -p 7540:7540 todo-final-app:v1
```
Адрес приложения ```http://localhost:7540/```