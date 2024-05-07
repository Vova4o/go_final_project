package main

import (
	"github.com/vova4o/go_final_project/internal/server"
)

func main() {
	// Initialize the server
	app := server.NewApp()

	// Start the server
	app.StartServer()

	// Wait for an interrupt signal to shutdown the server
	app.ShutdownServer(app.NewServer())
}