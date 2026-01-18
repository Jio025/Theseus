package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

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
	var test1 dockercontainermanagement.DockerContainer
	test1.ID = "Id test1"
	test1.Name = "Name test1"
	test1.Container = "Container test1"
	test1.Status = "running"
	var test2 dockercontainermanagement.DockerContainer
	test2.ID = "Id test2"
	test2.Name = "Name test2"
	test2.Container = "Container test2"
	test2.Status = "restarting"
	var test3 dockercontainermanagement.DockerContainer
	test3.ID = "Id test3"
	test3.Name = "Name test3"
	test3.Container = "Container test3"
	test3.Status = "stopped"

	store.SaveDockerContainer(test1)
	store.SaveDockerContainer(test2)
	store.SaveDockerContainer(test3)
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
	r.GET("/container/launchpage", func(c *gin.Context) {
		c.HTML(http.StatusOK, "containerLaunchpage.html", gin.H{})
	})

	// API Endpoints :
	// API Enpoint for testing if the service is up
	r.GET("/status", func(c *gin.Context) {
		// Return JSON response
		c.JSON(http.StatusOK, gin.H{
			"message": "Theseus service is up!",
		})
	})
	// API to get all runing containers in BBolt
	r.GET("/api/containers/running", func(c *gin.Context) {
		containers, err := store.GetAllDockerContainers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		// Return JSON response
		c.JSON(http.StatusOK, containers)
	})

	// API to create a webtop container
	r.POST("/api/webtop/create", func(c *gin.Context) {
		var container dockercontainermanagement.DockerContainer
		if err := c.BindJSON(&container); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if container.ID == "" {
			bytes := make([]byte, 4) // 8 characters hex
			if _, err := rand.Read(bytes); err == nil {
				container.ID = hex.EncodeToString(bytes)
			} else {
				// Fallback if random fails
				container.ID = fmt.Sprintf("container-%d", time.Now().Unix())
			}
		}

		dockercontainermanagement.YamlWriter(container)
		c.JSON(http.StatusCreated, gin.H{"message": "webtop created successfully"})
	})

	mockContainer := dockercontainermanagement.DockerContainer{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Name:      "lscr.io/linuxserver/webtop:latest",
		Container: "webtop_guillaume_dev",
		HostMachine: dockercontainermanagement.HostMachines{
			ID:     "node-01",
			IP:     "192.168.1.50",
			Status: "online",
		},
		RestartPolicy: "unless-stopped",
		Port: []dockercontainermanagement.PortBinding{
			{Internal: 3000, External: 3000},
			{Internal: 3001, External: 3001},
		},
		EnvironmentVariable: map[string]string{
			"PUID":  "1000",
			"PGID":  "1000",
			"TZ":    "America/Toronto",
			"TITLE": "Guillaume-Webtop",
		},
		VolumeMounts: map[string]string{
			"/home/guillaume/webtop/config": "/config",
			"/var/run/docker.sock":          "/var/run/docker.sock",
		},
		ShmSize: "1gb",
		Status:  "Active",
	}
	dockercontainermanagement.YamlWriter(mockContainer)

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
