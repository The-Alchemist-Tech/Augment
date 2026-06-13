package models

import (
	"testing"
)

func TestInvestorValidate(t *testing.T) {
	investor := Investor{
		ID:        1,
		Username:  "testName",
		Email:     "test@testemail.com",
		FirstName: "Tester",
		LastName:  "One",
	}


	// USERNAME TESTS

	// empty username
	investor.Username = ""

	err := investor.Validate()
	if err == nil {
		t.Fatalf("Investor validation failed to catch empty username")
	}

	// whitespace-only username
	investor.Username = "   "

	err = investor.Validate()
	if err == nil {
		t.Fatalf("Investor validation failed to catch whitespace-only username")
	}

	// Set to valid value for next tests
	investor.Username = "testName"


	// EMAIL TESTS

	// empty email
	investor.Email = ""

	err = investor.Validate()
	if err == nil {
		t.Fatalf("Investor validation failed to catch empty email")
	}

	// whitespace-only email
	investor.Email = "   "

	err = investor.Validate()
	if err == nil {
		t.Fatalf("Investor validation failed to catch whitespace-only email")
	}

	// invalid email format
	investor.Email = "thisWillNotWork"

	err = investor.Validate()
	if err == nil {
		t.Fatalf("Investor validation failed to catch invalid email format")
	}

	// Set to valid value for next tests
	investor.Email = "test@testemail.com"


	// FIRST NAME TESTS

	// empty first name
	investor.FirstName = ""

	err = investor.Validate()
	if err == nil {
		t.Fatalf("Investor validation failed to catch empty first name")
	}

	// whitespace-only first name
	investor.FirstName = "   "

	err = investor.Validate()
	if err == nil {
		t.Fatalf("Investor validation failed to catch whitespace-only first name")
	}

	// Set to valid value for next tests
	investor.FirstName = "Tester"


	// LAST NAME TESTS

	// empty last name
	investor.LastName = ""

	err = investor.Validate()
	if err == nil {
		t.Fatalf("Investor validation failed to catch empty last name")
	}

	// whitespace-only last name
	investor.LastName = "   "

	err = investor.Validate()
	if err == nil {
		t.Fatalf("Investor validation failed to catch whitespace-only last name")
	}

	// Set to valid value for next tests
	investor.LastName = "One"


	// VALID INVESTOR TEST

	err = investor.Validate()
	if err != nil {
		t.Fatalf("Investor validation failed on valid investor: %v", err)
	}
}
