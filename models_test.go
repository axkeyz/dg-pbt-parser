package main

import (
	"strings"
	"testing"
)

// TestSetCost calls *PBTItem.SetCost and checks if the
// given PBT item cost is correctly set in the PBTItem struct.
func TestSetCost(t *testing.T) {
	// Test cases
	var item PBTItem
	cases := map[string][]string{
		"1": {
			"4.5", "RUR",
			"4.5 |  |  | ",
		},
		"2": {
			"6.79", "UT",
			"4.5 |  | 6.79 | ",
		},
		"3": {
			"9.34", "ADJ",
			"4.5 | 9.34 | 6.79 | ",
		},
		"4": {
			"7.88", "NOR",
			"4.5 | 9.34 | 6.79 | 7.88",
		},
	}

	// Repeat test for each test case
	for key, test := range cases {
		// Update the PBT item with costs
		item.SetCost(test[0], test[1])

		// Get string of key cost values
		actual := strings.Join(
			[]string{
				item.RuralDelivery,
				item.Adjustment,
				item.UnderTicket,
				item.ItemCost,
			}, " | ")

		if actual != test[2] {
			// Test failed
			t.Fatalf(
				`TestSetPBTItemCost #%s 
				did not return %s, got %s`,
				key, test[2], actual,
			)
		}
	}
}

// TestGetNonEmptyCols calls *PBTItem.GetNonEmptyCols
// and tests if non-empty cols in the given PBT item
// are return correctly.
func TestGetNonEmptyCols(t *testing.T) {
	case1 := PBTItem{
		CustomerRef: "899OPUXS",
		Adjustment:  "9.23",
		AreaTo:      "Uranus",
	}
	actual := case1.GetNonEmptyCols()
	expected := []string{
		"customer_ref", "area_to", "adjustment",
	}

	for _, field := range expected {
		if !IsInArray(field, actual) {
			// Test failed
			t.Fatalf(
				`TestGetNonEmptyCols
				did not return %v, got %v`,
				expected, actual,
			)
		}
	}
	for _, field := range actual {
		if !IsInArray(field, expected) {
			// Test failed
			t.Fatalf(
				`TestGetPBTItemNonEmptyCols
				did not return %v, got %v`,
				expected, actual,
			)
		}
	}
}
