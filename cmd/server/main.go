package main

import (
	"fmt"
	"log"

	"github.com/Navl-bm/todo-final-project/config"
	"github.com/Navl-bm/todo-final-project/database"
	"github.com/Navl-bm/todo-final-project/handlers/api"
	"github.com/Navl-bm/todo-final-project/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("файл .env не найден, используются стандартные значения")
	}

	// Создание общего конфига для сервера используя переменные окружения
	// В случае отсутствия переменных используются стандартные значения
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка при инициализации конфига: %v", err)
	}

	// Инициализация базы данных
	db := database.NewDataBase(cfg)
	err = db.Init()
	defer db.DB.Close()

	if err != nil {
		log.Fatalf("ошибка подключения к базе данных: %v", err)
	}

	router := gin.Default() // Создание общего роутера

	router.StaticFS("/js/", gin.Dir(fmt.Sprintf("%s/js", cfg.WebDir), false))
	router.StaticFS("/css/", gin.Dir(fmt.Sprintf("%s/css", cfg.WebDir), false))
	router.StaticFile("/", fmt.Sprintf("%s/index.html", cfg.WebDir))
	router.StaticFile("/index.html", fmt.Sprintf("%s/index.html", cfg.WebDir))
	router.StaticFile("/login.html", fmt.Sprintf("%s/login.html", cfg.WebDir))
	router.StaticFile("/favicon.ico", fmt.Sprintf("%s/favicon.ico", cfg.WebDir))

	apiHandler := api.NewApiHandler(cfg, &db) // создание API интерфейса
	apiRouter := router.Group("/api")         // роутер для API приложения

	apiRouter.POST("/signin", apiHandler.SignIn)    // POST /api/signin запрос на авторизацию пользователя
	apiRouter.GET("/nextdate", apiHandler.NextDate) // GET /api/nextdate запрос на получение ближайшей даты выполнения задания

	middleware := middlewares.NewMiddleware(cfg) // миддлварь для проверки авторизации
	apiRouter.Use(middleware.CheckAuthToken)     // проверка авторизации для использования защищенного API
	{
		apiRouter.GET("/tasks", apiHandler.GetTasks) // GET /api/tasks запрос на общий список задач

		apiRouter.GET("/task", apiHandler.GetTask)       // GET /api/task запрос на задачу по id
		apiRouter.POST("/task", apiHandler.AddTask)      // POST /api/task запрос на добавление новой задачи
		apiRouter.PUT("/task", apiHandler.UpdateTask)    // PUT /api/task запрос на обновление задачи
		apiRouter.DELETE("/task", apiHandler.DeleteTask) // DELETE /api/task запрос на удаление задачи

		apiRouter.POST("/task/done", apiHandler.DoneTask) // POST запрос на выполнение задачи
	}

	// Запуск сервера
	if err := router.Run(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
		log.Fatalf("ошибка при запуске сервера: %v", err)
	}
}
