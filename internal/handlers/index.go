package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(mux *gin.Engine) {
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
	api.POST("/task", AddTask)       // to midleware
	api.GET("/task", FindTask)       // to midleware
	api.PUT("/task", UpdateTask)     // to midleware
	api.DELETE("/task", DeleteTask)  // to midleware
	api.POST("/task/done", DoneTask) // to midleware
	api.GET("/tasks", Tasks)         // to midleware
}

func Index(c *gin.Context) {
	c.File("./web/index.html")
}
