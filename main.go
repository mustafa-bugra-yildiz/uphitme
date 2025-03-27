package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/env"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
	"github.com/mustafa-bugra-yildiz/uphitme/router"

	_ "github.com/lib/pq"
)

func main() {
	db := connectDB()
	defer db.Close()

	taskRepo := task.NewRepo(db)

	log.Printf("Starting server on port %s", env.Port)
	log.Fatal(http.ListenAndServe(
		":"+env.Port,
		router.New(taskRepo),
	))
}

func connectDB() *sql.DB {
	log.Println("Connecting the database")
	db, err := sql.Open("postgres", env.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected the database")

	var version string
	err = db.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database version is", version)

	return db
}
