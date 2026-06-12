package investor

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	errors "augment/errors"
	db "augment/database"
	models "augment/models"
)

func CreateInvestor(w http.ResponseWriter, r *http.Request) {
	// Decode the request body directly into the struct
	investor := &models.Investor{}

	err := json.NewDecoder(r.Body).Decode(investor)
	if err != nil {
		errors.WriteJSONError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}
	defer r.Body.Close()

	err = investor.Validate()
	if err != nil {
		errors.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", err.Error()))
		return
	}

	// Insert the new investor into the database
	id, err := db.InsertResource(investor)
	// Check for duplicate email or username errors and return appropriate status and message for each, otherwise return generic internal server error for any other DB error
	// Note that we check for both email and username duplicates in the DB function, so we can return specific errors for each case here
	if err == db.ErrDuplicateEmail {
		errors.WriteJSONError(w, http.StatusConflict,  fmt.Sprintf("Investor with email %s already exists", investor.Email))
		return
	} else if err == db.ErrDuplicateUsername {
		errors.WriteJSONError(w, http.StatusConflict, fmt.Sprintf("Investor with username %s already exists", investor.Username))
		return
	} else if err == db.ErrDuplicate {
		// Fallthrough case that should not be hit, but exists just in case to catch any duplicate case that is not specifically an email or username duplicate - return generic message since we don't know which field is duplicated
		errors.WriteJSONError(w, http.StatusConflict, fmt.Sprintf("Investor with email %s or username %s already exists", investor.Email, investor.Username))
		return
	} else if err != nil && err != db.ErrDuplicate {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to create investor")
		return
	}

	log.Printf("New investor created with ID: %d", id)

	// Get new investor object to return in response - returning just ID will not show accurate timestamps
	// for the created investor
	newInvestor, err := GetInvestorByID(id)
	if err != nil {
		errors.WriteJSONError(w, http.StatusInternalServerError, "Failed to retrieve created investor")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newInvestor)
}

// TODO in real-world app: handler for API call to get investor by ID and validate
// Then call the below function to retrieve the investor and return in response

func GetInvestorByID(id int64) (investor *models.Investor, err error) {
	log.Printf("Retrieving investor with ID: %d", id)

	// Get the investor from the database using the generic GetResourceByID function
	// Returns any type for use with multiple tables, so we need to cast it to an Investor below
	res, err := db.GetResourceByID("investors", id)
	if err != nil {
		return nil, err
	}

	investor, ok := res.(*models.Investor)
	if !ok {
		return nil, fmt.Errorf("Failed to cast resource to Investor model for ID %d", id)
	}

	return investor, nil
}