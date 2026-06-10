package fund

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	db "augment/database"
)

type Fund struct {
	ID   			int       	`json:"id"`
	Name 			string    	`json:"name"`
	Units			int			`json:"units"`
	CreatedOn 		time.Time 	`json:"created_on"`
	LastModified	time.Time 	`json:"last_modified"`
}

func (fund *Fund) CreateFund(w http.ResponseWriter, r *http.Request) {
	// Decode the request body directly into the struct
	err := json.NewDecoder(r.Body).Decode(&fund)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Generic db.Insert
	rows, err := db.Insert("funds", "INSERT INTO `test` (name) VALUES ('myname')")
	if err != nil {
		log.Fatal("Database INSERT failed")
	} else if rows == 0 {
		
	}

	log.Printf("DB updated: rows affected: %d, ")
	w.WriteHeader(http.StatusOK)
}