package models

import (
	"testing"
	"time"
	
	"github.com/shopspring/decimal"
)


func TestValidate(t *testing.T) {
	buyer := 1
	seller := 2
	units := decimal.NewFromInt(10)
	
	cap := Cap{
		ID: 1,
		Fund: nil,
		Buyer: &buyer,
		Seller: &seller,
		Units: &units,
		CreatedOn: time.Now(),
	}


	// FUND TESTS

	// nil fund
	err := cap.Validate()
	if err == nil {
		t.Fatalf("Cap validation failed to catch nil fund")
	}

	// fund <= 0
	fund := -1
	cap.Fund = &fund

	err = cap.Validate()
	if err == nil {
		t.Fatalf("Cap validation failed to catch fund <= 0")
	}

	// Set to valid value for next tests
	fund = 1
	cap.Fund = &fund


	// BUYER TESTS

	// nil buyer
	cap.Buyer = nil

	err = cap.Validate()
	if err == nil {
		t.Fatalf("Cap validation failed to catch nil buyer")
	}
	
	// buyer <= 0
	buyer = 0
	cap.Buyer = &buyer

	err = cap.Validate()
	if err == nil {
		t.Fatalf("Cap validation failed to catch buyer <= 0")
	}

	// Set to valid value for next tests
	buyer = 1
	cap.Buyer = &buyer


	// SELLER TESTS

	// nil seller
	cap.Seller = nil

	err = cap.Validate()
	if err == nil {
		t.Fatalf("Cap validation failed to catch nil seller")
	}

	// seller <= 0
	seller = 0
	cap.Seller = &seller

	err = cap.Validate()
	if err == nil {
		t.Fatalf("Cap validation failed to catch seller <= 0")
	}

	// Set to valid value for next tests
	seller = 2
	cap.Seller = &seller


	// UNITS TESTS

	// nil units
	cap.Units = nil

	err = cap.Validate()
	if err == nil {
		t.Fatalf("Cap validation failed to catch nil units")
	}

	// units <= 0
	zeroUnits := decimal.NewFromInt(0)
	cap.Units = &zeroUnits

	err = cap.Validate()
	if err == nil {
		t.Fatalf("Cap validation failed to catch units <= 0")
	}

	// Set to valid value for next tests
	cap.Units = &units


	// BUYER == SELLER TEST

	cap.Seller = &buyer

	err = cap.Validate()
	if err == nil {
		t.Fatalf("Cap validation failed to catch buyer == seller")
	}

	// Reset seller
	cap.Seller = &seller


	// VALID CAP TEST

	err = cap.Validate()
	if err != nil {
		t.Fatalf("Cap validation failed on valid cap: %v", err)
	}
}