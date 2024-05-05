package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/vova4o/go_final_project/internal/database"
	"github.com/vova4o/go_final_project/internal/server"
)

// DeleteTask удаляет задачу по id
func DeleteTask(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан id"})
		return
	}

	err := database.DeleteTask(server.DB, id)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
