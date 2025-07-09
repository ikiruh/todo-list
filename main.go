package main

import (
	"log"
	"os"

	"github.com/ikiruh/go_final_project/pkg/db"
	"github.com/ikiruh/go_final_project/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	err = db.Init(os.Getenv("TODO_DBFILE"))
	if err != nil {
		log.Fatal("Error during database initialization: ", err)
	}
	defer db.Close()

	server.StartServer()
}
