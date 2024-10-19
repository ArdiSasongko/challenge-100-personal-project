package main

import (
	"basic-rest-api/cmd/api"
	"basic-rest-api/config"
	"basic-rest-api/db"
	"database/sql"
	"fmt"
	"log"
)

func main() {
	config := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Envs.USER, config.Envs.PASSWORD, config.Envs.DBNAME, config.Envs.HOST, config.Envs.PORT)

	db, err := db.NewSQLStorage(config)
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)

	server := api.NewAPIServer(":3000", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: Successfully Connected")
}
