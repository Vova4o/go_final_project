package main

import (
	"github.com/gin-gonic/gin"

	"github.com/vova4o/go_final_project/internal/database"
	"github.com/vova4o/go_final_project/internal/handlers"
	"github.com/vova4o/go_final_project/internal/server"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	mux := gin.Default()

	config := server.NewApp(mux)

	handlers.SetupRoutes(mux)

	mux.Use(config.Log.GinLogger(), gin.Recovery())

	srv := config.NewServer()

	config.StartServer()

	defer database.CloseDB(server.DB)

	config.ShutdownServer(srv)
}
