package main

import (
	"fmt"

	"example.com/dockercontainermanagement"
)

func main() {
	fmt.Printf("Hello World")
	dockercontainermanagement.OpenBoltDB()
}
