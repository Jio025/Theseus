package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	dockercontainermanagement "example.com/dockercontainermanagement"
)

func main() {
	store, err := dockercontainermanagement.InitDB("internal.db")
	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	// Test for DB content

	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()
	wd, _ := os.Getwd()

	// Clean static file path
	staticPath := filepath.Join(
		filepath.Dir(filepath.Dir(filepath.Dir(wd))), // Theseus/
		"web-interface",
		"static",
	)
	r.Static("/static", staticPath)

	// Routes for GIN :
	// TODO : Change this to a better path
	r.LoadHTMLGlob("../../../web-interface/HTML/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{})
	})
	r.GET("/settings", func(c *gin.Context) {
		c.HTML(http.StatusOK, "settings.html", gin.H{})
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	})

	// API Endpoints :
	// Define a simple GET endpoint
	r.GET("/ping", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/api/containers/running", func(c *gin.Context) {
		containers, err := store.GetAllActiveDockerContainers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Return JSON response
		c.JSON(http.StatusOK, containers)
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
