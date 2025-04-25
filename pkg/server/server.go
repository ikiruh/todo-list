package server

import (
	"log"
	"net/http"
	"os"
)

func StartServer() {
	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	port := os.Getenv("TODO_PORT")
	log.Println("Start server http://localhost:7540")
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Start server error: ", err)
	}
}
