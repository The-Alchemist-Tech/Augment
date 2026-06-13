package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	models "augment/models"

	"github.com/shopspring/decimal"
)

// THE DB MUST BE RUNNING FOR THESE INTEGRATION TESTS TO RUN
// USE `docker compose up --build` to start the app

// Units transferred in test setup — used in both insertTestTransfers and test assertions.
const (
	transfer1Units = 100 // testUser1 sells to testUser2
	transfer2Units = 25  // testUser2 sells to testUser3
)

// testCapIDs holds the IDs of cap rows inserted during test setup so we can
// delete them after all tests finish, leaving the DB in its original state.
var testCapIDs []int64

// Baselines capture the state of the DB before test transfers are inserted,
// so tests can assert deltas rather than absolute values.
// Without this we would have to have a clean DB to run tests.
// Now you can run tests after making calls to update the DB.
var (
	baselineUnitsInvestor2 decimal.Decimal
	baselineUnitsInvestor3 decimal.Decimal
	baselineUnitsInvestor4 decimal.Decimal
	baselineHistoryCount   int
)

// TestMain is the entry point for all tests in this package. It runs once.
func TestMain(m *testing.M) {
	// Tests run with cwd set to the package directory — move up to project root
	// so the migration path (database/migrations) resolves correctly.
	if err := os.Chdir(".."); err != nil {
		log.Fatalf("Failed to set working directory: %v", err)
	}

	// Connect to the DB and run migrations (which includes the seed data).
	// I found this approach after doing the baseline logic, so I don't need to worry about additional
	// transfers to the db - it's clean for testing.
	_, err := CreateDatabase("augment_test")
	if err != nil {
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	// Capture baseline values before inserting so tests can assert deltas
	// rather than absolute values — this way a non-empty DB won't break tests.
	if err := captureBaselines(); err != nil {
		log.Fatalf("Failed to capture baselines: %v", err)
	}

	// Insert two additional transfers so we have real adds and subtracts to test against.
	if err := insertTestTransfers(); err != nil {
		log.Fatalf("Failed to insert test transfers: %v — investors may have insufficient units. Reset the DB with: docker compose down -v && docker compose up --build", err)
	}

	// Run all tests via m.Run()
	code := m.Run()

	// Delete the inserted transfers to restore the DB to its original state.
	if err := cleanupTestTransfers(); err != nil {
		log.Printf("Warning: failed to clean up test transfers: %v", err)
	}

	os.Exit(code)
}

// captureBaselines records the current unit balances and history count before
// test transfers are inserted, so tests can assert on deltas rather than
// absolute values and work regardless of existing data in the DB.
func captureBaselines() error {
	var err error

	baselineUnitsInvestor2, err = GetFundUnitsForInvestor(1, 2)
	if err != nil {
		return fmt.Errorf("failed to get baseline for investor 2: %v", err)
	}

	baselineUnitsInvestor3, err = GetFundUnitsForInvestor(1, 3)
	if err != nil {
		return fmt.Errorf("failed to get baseline for investor 3: %v", err)
	}

	baselineUnitsInvestor4, err = GetFundUnitsForInvestor(1, 4)
	if err != nil {
		return fmt.Errorf("failed to get baseline for investor 4: %v", err)
	}

	history, err := GetCapHistoryForFund(1)
	if err != nil {
		return fmt.Errorf("failed to get baseline history count: %v", err)
	}
	baselineHistoryCount = len(history)

	return nil
}

// insertTestTransfers adds two transfers on top of the seed data so tests
// can verify that the aggregation correctly handles multiple transfers.
// NOTE: IF TRANSFERS FAIL DUE TO INSUFFICIENT UNITS, RESET THE DB WITH:
// `docker compose down -v && docker compose up --build`
func insertTestTransfers() error {
	fund := 1

	// Transfer 1: testUser1 (ID 2) sells units to testUser2 (ID 3)
	buyer1 := 3
	seller1 := 2
	units1 := decimal.NewFromInt(transfer1Units)

	id1, err := InsertResource(&models.Cap{
		Fund:   &fund,
		Buyer:  &buyer1,
		Seller: &seller1,
		Units:  &units1,
	})
	if err != nil {
		return err
	}
	testCapIDs = append(testCapIDs, id1)

	// Transfer 2: testUser2 (ID 3) sells units to testUser3 (ID 4)
	buyer2 := 4
	seller2 := 3
	units2 := decimal.NewFromInt(transfer2Units)

	id2, err := InsertResource(&models.Cap{
		Fund:   &fund,
		Buyer:  &buyer2,
		Seller: &seller2,
		Units:  &units2,
	})
	if err != nil {
		return err
	}
	testCapIDs = append(testCapIDs, id2)

	return nil
}

// cleanupTestTransfers deletes the rows inserted by insertTestTransfers,
// restoring the DB to its original state after tests finish.
func cleanupTestTransfers() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for _, id := range testCapIDs {
		_, err := database.ExecContext(ctx, "DELETE FROM cap WHERE id = ?", id)
		if err != nil {
			return err
		}
	}
	return nil
}

// TestGetFundUnitsForInvestor verifies that net units are calculated correctly
// for each single investor across multiple transactions (buys and sells).
func TestGetFundUnitsForInvestor(t *testing.T) {
	tests := []struct {
		name       string
		investorID int64
		baseline   decimal.Decimal
		delta      int64
	}{
		{"testUser1 sold 100", 2, baselineUnitsInvestor2, -transfer1Units},
		{"testUser2 bought 100 sold 25", 3, baselineUnitsInvestor3, transfer1Units - transfer2Units},
		{"testUser3 bought 25", 4, baselineUnitsInvestor4, transfer2Units},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			units, err := GetFundUnitsForInvestor(1, tc.investorID)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			expected := tc.baseline.Add(decimal.NewFromInt(tc.delta))
			if !units.Equal(expected) {
				t.Fatalf("Expected %s units, got %s", expected, units)
			}
		})
	}
}

// TestGetCapTableForFund verifies the cap table aggregation for the whole fund.
// It checks that each investor's balance shifted by the expected delta from
// baseline, and that the seed 'fund fund' investor is excluded from results.
func TestGetCapTableForFund(t *testing.T) {
	entries, err := GetCapTableForFund(1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Build a name -> units map for easy lookup
	unitsByName := make(map[string]decimal.Decimal)
	for _, e := range entries {
		unitsByName[e.InvestorName] = e.Units
	}

	// Seed investor must never appear in the cap table
	if _, found := unitsByName["fund fund"]; found {
		t.Errorf("Seed investor 'fund fund' should be excluded from cap table")
	}

	// Assert each investor's balance is baseline + expected delta
	deltas := map[string]struct {
		baseline decimal.Decimal
		delta    int64
	}{
		"testFirst1 testLast1": {baselineUnitsInvestor2, -transfer1Units},
		"testFirst2 testLast2": {baselineUnitsInvestor3, transfer1Units - transfer2Units},
		"testFirst3 testLast3": {baselineUnitsInvestor4, transfer2Units},
	}

	for name, tc := range deltas {
		got, ok := unitsByName[name]
		if !ok {
			t.Errorf("Investor '%s' not found in cap table", name)
			continue
		}
		expected := tc.baseline.Add(decimal.NewFromInt(tc.delta))
		if !got.Equal(expected) {
			t.Errorf("Investor '%s': expected %s units, got %s", name, expected, got)
		}
	}
}

// TestGetCapHistoryForFund verifies the full transfer history for a fund.
// Records must be in chronological order and all belong to fund 1.
func TestGetCapHistoryForFund(t *testing.T) {
	history, err := GetCapHistoryForFund(1)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Expect exactly 2 more records than the baseline
	if len(history) != baselineHistoryCount+2 {
		t.Fatalf("Expected %d history entries, got %d", baselineHistoryCount+2, len(history))
	}

	// All records must belong to fund 1
	for i, entry := range history {
		if entry.Fund != 1 {
			t.Errorf("Entry %d: expected fund 1, got %d", i, entry.Fund)
		}
	}

	// Records must be in chronological order
	for i := 1; i < len(history); i++ {
		if history[i].CreatedOn.Before(history[i-1].CreatedOn) {
			t.Errorf("History out of order at index %d", i)
		}
	}

	// The first 3 records are the seed transfers from 'fund fund' to the real investors
	for i := 0; i < 3; i++ {
		if history[i].Seller != "fund fund" {
			t.Errorf("Seed record %d: expected seller 'fund fund', got '%s'", i, history[i].Seller)
		}
	}

	// The last 2 records are the transfers we inserted — they appear at the end
	// since they were inserted after any pre-existing records.
	last := len(history) - 1
	if history[last-1].Seller != "testFirst1 testLast1" || history[last-1].Buyer != "testFirst2 testLast2" {
		t.Errorf("Transfer 1: expected testFirst1 -> testFirst2, got %s -> %s", history[last-1].Seller, history[last-1].Buyer)
	}
	if history[last].Seller != "testFirst2 testLast2" || history[last].Buyer != "testFirst3 testLast3" {
		t.Errorf("Transfer 2: expected testFirst2 -> testFirst3, got %s -> %s", history[last].Seller, history[last].Buyer)
	}
}
