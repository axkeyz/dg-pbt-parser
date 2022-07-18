// datautils_test.go contains tests for datautils.go
package main

import "testing"

// TestGetCustomerRef calls GetCustomerRef and checks if the
// customer reference is returned in the correct format.
func TestGetCustomerRef(t *testing.T) {
	case1 := GetCustomerRef("BD987", "n987")
	if case1 != "BD987 (N987)" {
		// Test failed
		t.Fatalf(
			`TestGetCustomerRef did not return BD987 (N987), got %s`,
			case1,
		)
	}

	case2 := GetCustomerRef("ux12345", "")
	if case2 != "UX12345" {
		// Test failed
		t.Fatalf(
			`TestGetCustomerRef did not return UX12345, got %s`,
			case2,
		)
	}
}

// TestGetRegion calls GetRegion and checks if the
// region is extracted correctly from a string.
func TestGetRegion(t *testing.T) {
	// Test cases
	cases := map[string]string{
		"PBT Couriers Christchurch Depot":     "CHRISTCHURCH",
		"PBT Couriers Palmerston North Depot": "PALMERSTON NORTH",
		"PBT Couriers Dunedin Depot":          "DUNEDIN",
		"PBT Couriers Auckland Central Depot": "AUCKLAND",
		"PBT Couriers Auckland South Depot":   "AUCKLAND",
		"PBT Couriers Auckland Depot":         "AUCKLAND",
	}

	// Repeat test for each test case
	for test, expected := range cases {
		actual := GetRegion(test)

		if actual != expected {
			// Test failed
			t.Fatalf(
				`TestGetRegion did not return %s, got %s`,
				expected, actual,
			)
		}
	}
}

// TestGetSortbyCode calls GetSortbyCode and checks if the
// sortby code is returned in the correct format.
func TestGetSortbyCode(t *testing.T) {
	// Test cases
	cases := map[string]string{
		"S&T 1873":      "ST",
		"EC19602":       "EC",
		"NTT4501151003": "NTT",
	}

	// Repeat test for each test case
	for test, expected := range cases {
		actual := GetSortbyCode(test)

		if actual != expected {
			// Test failed
			t.Fatalf(
				`TestGetSortbyCode did not return %s, got %s`,
				expected, actual,
			)
		}
	}
}