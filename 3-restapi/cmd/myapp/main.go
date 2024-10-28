package main

import (
	"database/sql"
	"fmt"
	"log"
	api "restapi/cmd"
	"restapi/config"
	"restapi/internal/database"
	"time"
)

func main() {
	config := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Envs.USER, config.Envs.PASSWORD, config.Envs.DBNAME, config.Envs.HOST, config.Envs.PORT,
	)

	db, err := database.NewSQLStorage(config)
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)

	server := api.NewServerAPI(":3000", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	log.Println("DB: Successfully Connected")
}
