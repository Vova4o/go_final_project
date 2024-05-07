package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(mux *gin.Engine, h *Handler) {
	mux.StaticFS("/web", gin.Dir("./web", true))
	mux.StaticFile("/favicon.ico", "./web/favicon.ico")
	mux.StaticFile("/index.html", "./web/index.html")
	mux.Static("/css", "./web/css")
	mux.Static("/js", "./web/js")
	mux.StaticFile("/login.html", "./web/login.html")

	mux.GET("/", Index)
	mux.POST("/api/signin", SignIn)
	mux.GET("/api/nextdate", NextDate)

	//	hendlers will go here
	api := mux.Group("/api")
	api.Use(AuthMiddleware())
	api.POST("/task", h.AddTask)       // to midleware
	api.GET("/task", h.FindTask)       // to midleware
	api.PUT("/task", h.UpdateTask)     // to midleware
	api.DELETE("/task", h.DeleteTask)  // to midleware
	api.POST("/task/done", h.DoneTask) // to midleware
	api.GET("/tasks", h.Tasks)         // to midleware
}

func Index(c *gin.Context) {
	c.File("./web/index.html")
}
