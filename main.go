package main

import (
	"log"
	"net/http"

	app "augment/augment"
	db "augment/database"

	"github.com/gorilla/mux"
)

func main() {
	// DB setup
	database, err := db.CreateDatabase()
	if err != nil {
		log.Fatalf("Database connection failed: %s", err.Error())
	}
	
	defer database.Close()

	// Init app and router
	app := &app.App{
		Router:   mux.NewRouter().StrictSlash(true),
		Database: database,
	}

	app.SetupRouter()

	log.Fatal(http.ListenAndServe(":8080", app.Router))
}