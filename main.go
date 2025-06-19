package main

import (
	"log"
	"net/http"

	"jusan_demo/pkg/app"
	"jusan_demo/pkg/db"
	"jusan_demo/pkg/config"
)

func main() {
	config.LoadConfig()
	db.Init()
	

	router := app.SetupRoutes()

	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

//server go
//config show
//auth struct
//refresh token(30 min) access token 5 min
//role 
//after regist create card auto