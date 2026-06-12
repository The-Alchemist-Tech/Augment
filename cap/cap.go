package cap

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	db "augment/database"
	errors "augment/errors"
	fund "augment/fund"
	investor "augment/investor"
	models "augment/models"

	"github.com/shopspring/decimal"
)

// For this app we will only accept IDs for the fund, buyer, and seller in the request body
// This follows general UI and API practices of using IDs to reference objects.
// It may be advantageous to allow emails of the buyer and seller in a real app.
// Auth would also be necessary to determine that the user can make the transfer.
func CreateTransfer(w http.ResponseWriter, r *http.Request) {
	log.Println("Starting create transfer request")
	// Decode the request body directly into the struct
	cap := &models.Cap{}

	err := json.NewDecoder(r.Body).Decode(cap)
	if err != nil {
		log.Printf("Invalid JSON payload: %v", err)
		errors.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON payload: %v", err))
		return
	}

	log.Printf("Request decoded to cap: %v", cap)

	defer r.Body.Close()

	err = cap.Validate()
	if err != nil {
		log.Printf("Invalid request inputs: %v", err)
		errors.WriteJSONError(w, http.StatusBadRequest, fmt.Sprintf("Validation error: %s", err.Error()))
		return
	}

	// Deeper validations on the buyer, seller, and fund need to hit the database
	// These do not belong in the model valdidations and will cause circular imports with the db if done there
	err = validateInput(cap)
	if err != nil {
		log.Printf("Invalid inputs from deep validation: %v", err)
		errors.WriteJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Do transfer
	id, err := db.InsertResource(cap)
	if err != nil {
		log.Printf("Failed to insert transfer: %v", err)
		errors.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create transfer: %v", err))
		return
	}

	log.Printf("New cap transaction created with ID: %d", id)

	// Get new cap transaction object to return in response to show the full object created and show accurate timestamps
	newTransfer, err := getCapByID(id)
	if err != nil {
		log.Printf("Failed to record transfer: %v", err)
		errors.WriteJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve created transfer: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTransfer)
}

func GetFundCap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("")
}

func GetFundCapHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("")
}

func validateInput(c *models.Cap) error {
	// Due to time I was not able to optimize this into a batched query, but that would be important in a live app for DB performance

	// Check fund exists
	_, err := fund.GetFundByID(int64(*c.Fund))
	if err != nil {
		return fmt.Errorf("A Fund with ID '%d' does not exist.", *c.Fund)
	}

	// Check buyer exists
	_, err = investor.GetInvestorByID(int64(*c.Buyer))
	if err != nil {
		// I did not want to create a buyer if they do not exist as the seller would transfer to
		// an empty account with no real owner.
		// My ideal solution that does not make sense to implement here is to email the
		// buyer to tell them someone wants to transfer units and to create an account.
		// Maybe a table that holds the units until they create an account to claim them
		// and return to the seller after a period of time if unclaimed.
		// We would need to identify the buyer by email or more to confirm identity.
		return fmt.Errorf("A buyer with ID '%d' does not exist.", *c.Buyer)
	}

	// Check seller exists
	_, err = investor.GetInvestorByID(int64(*c.Seller))
	if err != nil {
		return fmt.Errorf("A seller with ID '%d' does not exist.", *c.Seller)
	}

	// Get number of units the seller has
	var sellerUnitsAvailable decimal.Decimal

	sellerUnitsAvailable, err = db.GetFundUnitsForInvestor(int64(*c.Fund), int64(*c.Seller))
	if err != nil {
		return fmt.Errorf("Failed to get units or no units are available for seller %d for fund %d: %v", *c.Fund, *c.Seller, err)
	}

	// Check that the seller has enough units in the fund to transfer
	if sellerUnitsAvailable.LessThan(*c.Units) {
		return fmt.Errorf(
			"Seller %d does not have %s units available in fund %d to transfer. %s available.",
			*c.Seller,
			(*c.Units).String(),
			*c.Fund,
			sellerUnitsAvailable.String(),
		)
	}

	return nil
}

func getCapByID(id int64) (cap *models.Cap, err error) {
	log.Printf("Retrieving fund with ID: %d", id)

	// Get the cap transaction from the database using the generic GetResourceByID function
	// Returns any type for use with multiple tables, so we need to cast it to a cap below
	res, err := db.GetResourceByID("cap", id)
	if err != nil {
		return nil, err
	}

	cap, ok := res.(*models.Cap)
	if !ok {
		return nil, fmt.Errorf("Failed to cast resource to Cap model for ID %d", id)
	}

	return cap, nil
}
