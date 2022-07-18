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
// sortby code is returned in the correct format. Also se
func TestGetSortbyCode(t *testing.T) {
	// Test cases
	cases := map[string][]string{
		"MX": {
			"IO908896690", "sdfkjh IO",
		},
		"RTF": {
			"IXU 2348", "VSD INOB",
		},
		"MF": {
			"MF & 2390", "Edf EE",
		},
		"OP": {
			"OPURO@2388", "OPdjjjh",
		},
	}

	// Repeat test for each test case
	for expected, inputs := range cases {
		actual := GetSortbyCode(inputs[0], inputs[1])
		if actual != expected {
			// Test failed
			t.Fatalf(
				`TestGetSortbyCode did not return %s, got %s`,
				expected, actual,
			)
		}
	}
}

// TestHasDGSortbyCode tests HasDGSortbyCode and checks if the
// sortby code is in the customers outlined in config/customers_test.json.
func TestHasDGSortbyCode(t *testing.T) {
	customers := OpenConfigJSON("customers_test")

	cases := map[string]bool{
		"RTF":   false,
		"MX":    true,
		"MXUDF": true,
		"OPURO": true,
		"OP":    false,
		"NXU":   false,
	}

	// Repeat test for each test case
	for test, expected := range cases {
		actual := HasDGSortbyCode(test, customers)

		if actual != expected {
			// Test failed
			t.Fatalf(
				`TestStripNonLetters("%s") did not return %v, got %v`,
				test, expected, actual,
			)
		}
	}
}
