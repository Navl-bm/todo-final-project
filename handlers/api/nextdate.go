package api

import (
	"net/http"
	"time"

	"github.com/Navl-bm/todo-final-project/pkg/dates"
	"github.com/gin-gonic/gin"
)

// NextDate формирует ближайшую дату для выполнения задачи
func (h *ApiHandler) NextDate(ctx *gin.Context) {
	now, exists := ctx.GetQuery("now")
	if !exists {
		now = time.Now().Format("20060102")
	}

	dstart, exists := ctx.GetQuery("date")
	if !exists {
		ctx.Data(http.StatusBadRequest, "text/html", []byte("date required"))
		return
	}

	repeat, exists := ctx.GetQuery("repeat")
	if !exists {
		ctx.Data(http.StatusBadRequest, "text/html", []byte("repeat required"))
		return
	}
	nextdate, err := dates.NextDate(now, dstart, repeat)
	if err != nil {
		ctx.Data(http.StatusBadRequest, "text/html", []byte(err.Error()))
		return
	}
	ctx.Data(http.StatusOK, "text/html", []byte(nextdate))
}
