package main

import (
	"log"

	"github.com/ikiruh/go_final_project/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	server.StartServer()
}
