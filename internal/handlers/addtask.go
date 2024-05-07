package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/labstack/gommon/log" // need to check it
	"github.com/vova4o/go_final_project/internal/models"
	"github.com/vova4o/go_final_project/internal/nextdate"
)

type task struct {
	Date    string `json:"date,omitempty"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type Handler struct {
	Storage Storager
}

// NewHandler создаёт новый объект Handler.
func NewHandler(storage Storager) *Handler {
	return &Handler{Storage: storage}
}

// AddTask добавляет задачу в базу данных.
func (h *Handler) AddTask(c *gin.Context) {
	var err error
	var t task
	if err = c.BindJSON(&t); err != nil {
		err := errors.New("ошибка десериализации JSON")
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = t.checkTask()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.Storage.AddTaskDB(t.Date, t.Title, t.Comment, t.Repeat)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// UpdateTask обновляет задачу по id в базе данных.
func (h *Handler) UpdateTask(c *gin.Context) {
	var t models.DBTask
	if err := c.ShouldBindJSON(&t); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	checkT := task{
		Date:    t.Date,
		Title:   t.Title,
		Comment: t.Comment,
		Repeat:  t.Repeat,
	}

	err := checkT.checkTask()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if t.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "не указан id задачи"})
		return
	}

	_, err = strconv.Atoi(t.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id задачи должен быть числом"})
		return
	}

	_, err = h.Storage.FindTask(t.ID)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	err = h.Storage.UpdateTask(t)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// checkTask проверяет корректность данных задачи и возвращает исправленную задачу и ошибку
func (t *task) checkTask() error {
	if t.Title == "" {
		return fmt.Errorf("не указан заголовок задачи")
	}
	if t.Date == "" {
		t.Date = time.Now().Format("20060102")
	}
	date, err := time.Parse("20060102", t.Date)
	if err != nil {
		return fmt.Errorf("дата представлена в формате, отличном от 20060102")
	}

	if t.Repeat != "" && t.Repeat[0] != 'd' && t.Repeat[0] != 'w' && t.Repeat[0] != 'm' && t.Repeat[0] != 'y' {
		return errors.New("неверное правило повторения")
	}

	if len(t.Repeat) > 0 {
		if t.Repeat[0] != 'd' && t.Repeat[0] != 'w' && t.Repeat[0] != 'm' && t.Repeat[0] != 'y' {
			return errors.New("неверное правило повторения")
		}
		if t.Repeat[0] == 'd' || t.Repeat[0] == 'w' || t.Repeat[0] == 'm' {
			if len(t.Repeat) < 3 {
				return errors.New("неверное правило повторения")
			}
		}
	}

	if date.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		if t.Repeat == "" {
			t.Date = time.Now().Format("20060102")
		}
	}

	if date.Truncate(24 * time.Hour).Before(time.Now().Truncate(24 * time.Hour)) {
		if t.Repeat != "" {
			t.Date, err = nextdate.NextDate(time.Now(), t.Date, t.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка при вычислении следующей даты: %v", err)
			}
		}
	}

	return nil
}
