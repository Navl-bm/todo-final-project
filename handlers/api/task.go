package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Navl-bm/todo-final-project/database"
	"github.com/Navl-bm/todo-final-project/errors"
	"github.com/Navl-bm/todo-final-project/models"
	"github.com/Navl-bm/todo-final-project/utils"
	"github.com/gin-gonic/gin"
)

// checkTaskForErrors проверяет текущую задачу на выполенение требований для создания и редактирования
func checkTaskForErrors(task *models.Task) (models.Task, error) {
	now := time.Now()

	if task.Title == "" {
		return models.Task{}, errors.ErrBadJSON
	}

	if task.Date == "" {
		task.Date = now.Format(utils.DateFormat)
	}

	date, err := time.Parse(utils.DateFormat, task.Date)
	if err != nil {
		return models.Task{}, errors.ErrBadTime

	}

	if task.Repeat != "" {
		_, err = utils.NextDate(now.Format(utils.DateFormat), task.Date, task.Repeat)
		if err != nil {
			return models.Task{}, err
		}
	}

	if !utils.AfterNow(date, now) {
		task.Date = now.Format(utils.DateFormat)
	}

	return *task, nil
}

// AddTask обработчик для добавления задачи в БД
//
// Проверяет корректность введенных данных и создает задачу
func (h *ApiHandler) AddTask(ctx *gin.Context) {
	var task models.Task

	err := ctx.ShouldBindBodyWithJSON(&task)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadJSON.Error()})
		return
	}

	task, err = checkTaskForErrors(&task)
	if err != nil {
		log.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := h.DB.AddTask(&task)
	if err != nil {
		log.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServerError.Error()})
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, gin.H{"id": id})
}

// GetTasks обработчик для передачи списка задач пользователю
//
// Получает отсортированный список задач пользователя
func (h *ApiHandler) GetTasks(ctx *gin.Context) {
	var tasks *models.TasksList
	var err error
	search, exists := ctx.GetQuery("search")
	if exists {
		date, errParse := time.Parse("02.01.2006", search)
		if errParse == nil {
			tasks, err = h.DB.GetTasks(50, date.Format("20060102"), database.SearchTypeDate)
		} else {
			tasks, err = h.DB.GetTasks(50, search, database.SearchTypeTitleOrComment)
		}
	} else {
		tasks, err = h.DB.GetTasks(50, "", database.SearchTypeNil)
	}

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errors.ErrInternalServerError.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// GetTask обработчик для передачи пользователю запрашиваемой задачи по ее идентификатору
//
// Получает задачу по идентификатору
func (h *ApiHandler) GetTask(ctx *gin.Context) {
	idString, exists := ctx.GetQuery("id")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	task, err := h.DB.GetTask(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// UpdateTask обработчик для сохранения изменений в задаче
//
// Проверяет корректность данных задачи и сохраняет ее
func (h *ApiHandler) UpdateTask(ctx *gin.Context) {
	var task models.Task
	err := ctx.ShouldBindBodyWithJSON(&task)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	task, err = checkTaskForErrors(&task)
	if err != nil {
		log.Println(err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.DB.UpdateTask(&task)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

// DoneTask обработчик выполнения задачи
//
// Проверяет поле repeat, в случае его наличия ищет новую ближайшую дату для задачи и обновляет ее
func (h *ApiHandler) DoneTask(ctx *gin.Context) {
	idString, exists := ctx.GetQuery("id")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	task, err := h.DB.GetTask(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	if task.Repeat != "" {
		task.Date, err = utils.NextDate(time.Now().Format(utils.DateFormat), task.Date, task.Repeat)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRepeatFormat.Error()})
			return
		}
		err := h.DB.UpdateTask(task)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrInternalServerError.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	} else {
		err = h.DB.DeleteTask(id)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{})
	}
}

// DeletTask обработчик для удаления задачи
//
// Удаляет задачу по id
func (h *ApiHandler) DeleteTask(ctx *gin.Context) {
	idString, exists := ctx.GetQuery("id")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	err = h.DB.DeleteTask(id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errors.ErrTaskNotFound.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
