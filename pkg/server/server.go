package server

import (
	"final-project/pkg/api"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func RunServer() {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	webDir, err := filepath.Abs("./web")
	if err != nil {
		panic(err)
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на порту %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
