package main

import (
	"fmt"
	"log"
	"net/http"

	dockercontainermanagement "example.com/dockercontainermanagement"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	store, err := dockercontainermanagement.InitDB("internal.db")
	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	//http.HandleFunc("/", handler)
	//log.Fatal(http.ListenAndServe(":8080", nil))

	//fmt.Println("Database oppened successfully")

	// Test for dockerDeployment
	var test_container dockercontainermanagement.DockerContainer
	test_container.Container = "AlpineTest"
	test_container.Name = "alpine:latest"
	test_container.ID = "alpineTest"
	dockercontainermanagement.DeployContainerBackground(test_container)
	fmt.Println("Docker Container lauched ! 🚀")
}
