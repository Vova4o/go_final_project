package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/vova4o/go_final_project/internal/database"
	"github.com/vova4o/go_final_project/internal/server"
)

// DoneTask помечает задачу как выполненную по id, если задача не повторяющаяся, то удаляет ее из базы данных,
// в противном случае устанавливает дату следующего выполнения и записывает в базу данных.
func DoneTask(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан id"})
		return
	}

	err := database.DoneTask(server.DB, id)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
