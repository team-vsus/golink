package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/team-vsus/golink/handler"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load .env file!")
	}

	r := handler.InitHandler()
	r.Run()
}
