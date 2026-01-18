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
	tmpl, err := template.ParseFiles("webtop.tmpl")
	if err != nil {
		log.Printf("❗ Error parsing template: %v", err)
		return
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, container); err != nil {
		log.Printf("❗ Error executing template: %v", err)
		return
	}
	data := tpl.Bytes()

	err = os.WriteFile("docker-compose.yml", data, 0644)
	if err != nil {
		log.Printf("❗ Error writing file: %v", err)
		return
	}
	log.Println("docker-compose.yml created successfully")
}
