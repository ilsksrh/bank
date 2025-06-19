package db

import (
	"jusan_demo/pkg/config"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func Init() {
	var err error
	DB, err = sqlx.Connect("postgres", config.AppConfig.GetDBConn())

	if err != nil {
		log.Fatalln("Не удалось подключиться к базе данных:", err)
	}
	log.Println("Успешное подключение к базе данных")
}

