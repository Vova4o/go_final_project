package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vova4o/go_final_project/internal/database"
)

// NextDate возвращает следующую дату, когда нужно выполнить задачу в соответствии с заданной датой и периодичностью
func NextDate(c *gin.Context) {
	nowStr := c.Query("now")
	date := c.Query("date")
	repeat := c.Query("repeat")

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	next, err := database.NextDate(now, date, repeat)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, next)
}
