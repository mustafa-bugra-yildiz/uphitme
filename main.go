package main

import (
	"log"
	"net/http"
	"context"

	"github.com/mustafa-bugra-yildiz/uphitme/env"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
	"github.com/mustafa-bugra-yildiz/uphitme/router"
	"github.com/mustafa-bugra-yildiz/uphitme/scheduler"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	db := connectDB()
	defer db.Close()

	taskRepo := task.NewRepo(db)

	scheduler, err := scheduler.New(context.Background(), taskRepo)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on port %s", env.Port)
	log.Fatal(http.ListenAndServe(
		":"+env.Port,
		router.New(taskRepo, scheduler),
	))
}

func connectDB() *pgxpool.Pool {
	log.Println("Connecting the database")
	db, err := pgxpool.New(context.Background(), env.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected the database")

	var version string
	err = db.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database version is", version)

	return db
}
