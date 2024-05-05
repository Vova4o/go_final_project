package server

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vova4o/go_final_project/internal/config"
	"github.com/vova4o/go_final_project/internal/database"
	"github.com/vova4o/go_final_project/internal/logger"
)

type ServerConfig struct {
	Addr    string
	Handler *gin.Engine
	DB      *sql.DB
	Log     *logger.Logger
}

var DB *sql.DB

func NewApp(handler *gin.Engine) *ServerConfig {
	addr := config.Address()

	var err error
	DB, err = database.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	log := logger.New()

	return &ServerConfig{
		Addr:    addr,
		Handler: handler,
		DB:      DB,
		Log:     log,
	}
}

func (c *ServerConfig) NewServer() *http.Server {
	return &http.Server{
		Addr:    c.Addr,
		Handler: c.Handler,
	}
}

func (c *ServerConfig) StartServer() {
	go func() {
		log.Printf("Starting server on %s\n", c.Addr)
		if err := c.NewServer().ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func (c *ServerConfig) ShutdownServer(srv *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
