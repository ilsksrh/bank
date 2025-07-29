package main

import (
	"log"
	"net/http"

	"jusan_demo/pkg/app"
	"jusan_demo/pkg/db"
	"jusan_demo/pkg/config"
)

// @title Jusan Demo API
// @version 1.0
// @description API для банковской системы
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	config.LoadConfig()
	db.Init()
	

	router := app.SetupRoutes()

	log.Println("Сервер запущен на :8080")
	if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}




//20.04
//server go
//config show
//auth struct
//refresh token(30 min) access token 5 min

//23.04
//role 
//after regist create card auto
//api swagger
//swag init -g cmd/main.go чтобы обновилось

//25.04
//sql remove
//transactions fix
//separate api loans
//roles superadmin
