package server

import (
	"final-project/pkg/api"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func RunServer() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("TODO_PORT")

	if port == "" {
		port = "7540"
	}

	webDir, err := filepath.Abs("./web")
	if err != nil {
		log.Fatalf("Ошибка получения пути до ./web: %v", err)
	}

	api.Init()
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	log.Printf("Сервер запущен на порту %s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
