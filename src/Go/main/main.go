package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	dockercontainermanagement "example.com/dockercontainermanagement"
)

func main() {
	store, err := dockercontainermanagement.InitDB("internal.db")
	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	// Test for dockerDeployment
	var test_container dockercontainermanagement.DockerContainer
	test_container.Container = "AlpineTest"
	test_container.Name = "alpine:latest"
	test_container.ID = "alpineTest"
	dockercontainermanagement.DeployContainerBackground(test_container, store)
	fmt.Println("Docker Container lauched ! 🚀")

	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Routes for GIN :
	// TODO : Change this to a better path
	r.LoadHTMLGlob("../../../web-interface/HTML/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
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

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
