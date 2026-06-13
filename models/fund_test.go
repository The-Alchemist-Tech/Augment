package models

import (
	"testing"
)

func TestFundValidate(t *testing.T) {
	units := 100

	fund := Fund{
		ID:    1,
		Name:  "Test Fund",
		Units: &units,
	}


	// NAME TESTS

	// empty name
	fund.Name = ""

	err := fund.Validate()
	if err == nil {
		t.Fatalf("Fund validation failed to catch empty name")
	}

	// whitespace-only name
	fund.Name = "   "

	err = fund.Validate()
	if err == nil {
		t.Fatalf("Fund validation failed to catch whitespace-only name")
	}

	// Set to valid value for next tests
	fund.Name = "Test Fund"


	// UNITS TESTS

	// nil units
	fund.Units = nil

	err = fund.Validate()
	if err == nil {
		t.Fatalf("Fund validation failed to catch nil units")
	}

	// units <= 0
	zeroUnits := 0
	fund.Units = &zeroUnits

	err = fund.Validate()
	if err == nil {
		t.Fatalf("Fund validation failed to catch units <= 0")
	}

	// Set to valid value for next tests
	fund.Units = &units


	// VALID FUND TEST

	err = fund.Validate()
	if err != nil {
		t.Fatalf("Fund validation failed on valid fund: %v", err)
	}
}
