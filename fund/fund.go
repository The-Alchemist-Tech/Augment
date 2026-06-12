package fund

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	errors "augment/errors"
	db "augment/database"
	models "augment/models"
)

func CreateFund(w http.ResponseWriter, r *http.Request) {
	// Decode the request body directly into the struct
	fund := &models.Fund{}

	err := json.NewDecoder(r.Body).Decode(fund)
	if err != nil {
		errors.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	err = fund.Validate()
	if err != nil {
		errors.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", err.Error()))
		return
	}

	defer r.Body.Close()

	// Insert the new fund into the database
	id, err := db.InsertResource(fund)
	if err != nil && err != db.ErrDuplicate {
		errors.WriteJSONError(w, http.StatusInternalServerError,  fmt.Sprintf("Failed to create fund: %v", err))
		return
	} else if err == db.ErrDuplicate {
		errors.WriteJSONError(w, http.StatusConflict, fmt.Sprintf("Fund with name %s already exists", fund.Name))
		return
	}

	log.Printf("New fund created with ID: %d", id)

	// Get new fund object to return in response - returning just ID will not show accurate timestamps
	// for the created fund
	newFund, err := GetFundByID(id)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to retrieve created fund")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newFund)
}

// In real-world app: handler for API call to get fund by ID and validate, then call the below function
func GetFundByID(id int64) (fund *models.Fund, err error) {
	log.Printf("Retrieving fund with ID: %d", id)

	// Get the fund from the database using the generic GetResourceByID function
	// Returns any type for use with multiple tables, so we need to cast it to a Fund below
	res, err := db.GetResourceByID("funds", id)
	if err != nil {
		return nil, err
	}

	fund, ok := res.(*models.Fund)
	if !ok {
		return nil, fmt.Errorf("Failed to cast resource to Fund model for ID %d", id)
	}

	return fund, nil
}