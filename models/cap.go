package models

import (
	"fmt"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

type Cap struct {
	ID           int              `json:"id"`
	Fund         *int             `json:"fund"`
	Buyer        *int             `json:"buyer"`
	Seller       *int             `json:"seller"`
	Units        *decimal.Decimal `json:"units"`
	CreatedOn    time.Time        `json:"created_on"`
	LastModified time.Time        `json:"last_modified"`
}

func (c *Cap) String() string {
	// nil checks to avoid nil pointer dereference
	fund := "<nil>"
	if c.Fund != nil {
		fund = fmt.Sprintf("%d", *c.Fund)
	}
	buyer := "<nil>"
	if c.Buyer != nil {
		buyer = fmt.Sprintf("%d", *c.Buyer)
	}
	seller := "<nil>"
	if c.Seller != nil {
		seller = fmt.Sprintf("%d", *c.Seller)
	}
	units := "<nil>"
	if c.Units != nil {
		units = fmt.Sprintf("%f", *c.Units)
	}

	return fmt.Sprintf("Cap[ID=%d, FundID=%s, BuyerID=%s, SellerID=%s, Units=%s, CreatedOn=%s, LastModified=%s]",
		c.ID, fund, buyer, seller, units, c.CreatedOn.Format(time.RFC3339), c.LastModified.Format(time.RFC3339))
}

func (c *Cap) Validate() error {
	log.Println("Validating request inputs")
	// Values within acceptable ranges and types
	if c.Fund == nil || *c.Fund <= 0 {
		return fmt.Errorf("Fund must be present and greater than zero.")
	}

	if c.Buyer == nil || *c.Buyer <= 0 {
		return fmt.Errorf("Buyer must be present and greater than zero.")
	}

	if c.Seller == nil || *c.Seller <= 0 {
		return fmt.Errorf("Seller must be present and greater than zero.")
	}

	if c.Units == nil || c.Units.Cmp(decimal.NewFromInt(0)) <= 0 {
		return fmt.Errorf("Units must be present and greater than zero.")
	}

	//Cannot buy from yourself
	log.Printf("In validate - Buyer ID: %d; Seller ID: %d", *c.Buyer, *c.Seller)
	if *c.Buyer == *c.Seller {
		return fmt.Errorf("Buyer and Seller cannot be the same person.")
	}

	return nil
}
