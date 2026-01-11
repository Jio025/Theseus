package dockercontainermanagement

import (
	"bytes"
	"log"
	"os"
	"text/template"
)

// This file contains the methods for launching a webtop docker container

// Making the custom compose.yaml file for docker compose
func YamlWriter(container DockerContainer) {
	tmpl, _ := template.ParseFiles("webtop.tmpl")
	var tpl bytes.Buffer
	tmpl.Execute(&tpl, container)
	data := tpl.Bytes()

	err := os.WriteFile("docker-compose.yml", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("docker-compose.yml created successfully")
}
