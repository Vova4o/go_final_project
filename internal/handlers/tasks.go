package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log"
	"github.com/vova4o/go_final_project/internal/database"
	"github.com/vova4o/go_final_project/internal/models"
	"github.com/vova4o/go_final_project/internal/server"
)

// Tasks возвращает последниее 10 задач из базы данных. Оставил возможность указать смещение, но не использую его.
func Tasks(c *gin.Context) {
	search, searchExists := c.GetQuery("search")
	var tasks []models.DBTask
	var err error

	if !searchExists {
		offset := 0
		tasks, err = database.Tasks(server.DB, offset)
		if err != nil {
			log.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		parsedDate, err := time.Parse("02.01.2006", search)
		if err != nil {
			// The search query is not a date, so perform a string search.
			tasks, err = database.SearchTasks(server.DB, search)
			if err != nil {
				log.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			// The search query is a date.
			tasks, err = database.TasksByDate(server.DB, parsedDate.Format("20060102"))
			if err != nil {
				log.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}

	if tasks == nil {
		tasks = []models.DBTask{}
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// FindTask возвращает задачу по id
func FindTask(c *gin.Context) {
	search := c.Query("id")
	task, err := database.FindTask(server.DB, search)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if (models.DBTask{}) == task {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
		return
	}

	c.JSON(http.StatusOK, task)
}
