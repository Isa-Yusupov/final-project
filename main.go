package main

import (
	"final-project/pkg/db"
	"final-project/pkg/server"
	"log"
)

func main() {
	err := db.Init("scheduler.db")
	if err != nil {
		log.Fatal("Ошибка БД", err)
	}
	defer db.Close()

	server.RunServer()
}
