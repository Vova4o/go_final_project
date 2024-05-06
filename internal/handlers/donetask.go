package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
)

// DoneTask помечает задачу как выполненную по id, если задача не повторяющаяся, то удаляет ее из базы данных,
// в противном случае устанавливает дату следующего выполнения и записывает в базу данных.
func (h *Handler) DoneTask(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не указан id"})
		return
	}

	err := h.Storage.DoneTask(id)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
