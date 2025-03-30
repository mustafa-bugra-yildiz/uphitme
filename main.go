package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"

	"github.com/mustafa-bugra-yildiz/uphitme/auth"
	"github.com/mustafa-bugra-yildiz/uphitme/env"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/task"
	"github.com/mustafa-bugra-yildiz/uphitme/repos/user"
	"github.com/mustafa-bugra-yildiz/uphitme/router"
	"github.com/mustafa-bugra-yildiz/uphitme/scheduler"
	"github.com/mustafa-bugra-yildiz/uphitme/view"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed bin/tailwind.css
var tailwind string

func main() {
	// setup views
	view.Setup(tailwind)

	// setup db
	db := connectDB()
	defer db.Close()

	// setup repos
	taskRepo := task.NewRepo(db)
	userRepo := user.NewRepo(db)

	// setup modules
	auth := auth.New(userRepo)

	scheduler, err := scheduler.New(context.Background(), taskRepo)
	if err != nil {
		log.Fatal(err)
	}

	// run server
	log.Printf("Starting server on port %s", env.Port)
	log.Fatal(http.ListenAndServe(
		":"+env.Port,
		router.New(auth, taskRepo, userRepo, scheduler),
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
