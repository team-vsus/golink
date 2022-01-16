package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
	"github.com/team-vsus/golink/handler"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	// Load .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load .env file!")
	}

	r := handler.InitHandler()
	r.Run()
}
