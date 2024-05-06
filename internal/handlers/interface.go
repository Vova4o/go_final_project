package handlers

import (
	"github.com/vova4o/go_final_project/internal/database"
	"github.com/vova4o/go_final_project/internal/models"
)

var _ Storager = &database.Storage{}

type Storager interface {
	InitDB() error
	CloseDB()
	AddTaskDB(date string, title string, comment string, repeat string) (int64, error)
	FindTask(id string) (models.DBTask, error)
	UpdateTask(task models.DBTask) error
	Tasks(offset int) ([]models.DBTask, error)
	SearchTasks(search string) ([]models.DBTask, error)
	TasksByDate(date string) ([]models.DBTask, error)
	DoneTask(id string) error
	DeleteTask(id string) error
}
