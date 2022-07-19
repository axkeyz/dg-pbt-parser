// datautils_test.go contains tests for datautils.go
package main

import "testing"

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
			"MXUDF908896690", "sdfkjh IO",
		},
		"RTF": {
			"IXU 2348", "VSD INOB",
		},
		"MF": {
			"MF & 2390", "Edf EE",
		},
		"OP": {
			"skjdhf2198", "MUXIO IO",
		},
		"ZZ": {
			"SNOOZE", "SNOOZE",
		},
		"UNKNOWN": {
			"MD", "FUIOX",
		},
		"SALES": {
			"SN", "SOMETHING ltd",
		},
	}
	customers := OpenConfigJSON("customers_test")
	sales := OpenConfigJSON("sales_test")

	// Repeat test for each test case
	for expected, inputs := range cases {
		actual := GetSortbyCode(inputs[0], inputs[1], customers, sales)
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
		actual, _ := HasDGSortbyCode(test, customers)

		if actual != expected {
			// Test failed
			t.Fatalf(
				`TestStripNonLetters("%s") did not return %v, got %v`,
				test, expected, actual,
			)
		}
	}
}

// TestTryDGSalesEComSortbyCode tests TryDGSalesEComSortbyCode
// and checks if the correct sortby code is returned as according to
// the customers outlined in config/sales_test.json.
func TestTryDGSalesEComSortbyCode(t *testing.T) {
	// Open config file for sales customers
	customers := OpenConfigJSON("sales_test")

	// Test cases
	cases := map[string]string{
		"RIOLEMY LTD":   "SALES",
		"SNOOZE":        "ZZ",
		"MF INDYOS LTD": "HH",
		"JOHN SMITH":    "E COMMERCE",
		"ASDFGHJLOOO":   "UNKNOWN",
	}

	// Repeat test for each test case
	for test, expected := range cases {
		actual := TryDGSalesEComSortbyCode(test, customers)

		if actual != expected {
			// Test failed
			t.Fatalf(
				`TestTryDGSalesEComSortbyCode did not return %v, got %v`,
				expected, actual,
			)
		}
	}
}

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

// TestGetInvoiceCostTypeAndConsignment tests
// GetInvoiceCostTypeAndConsignment to see if the correct cost type
// and consignment number is returned.
func TestGetInvoiceCostTypeAndConsignment(t *testing.T) {
	// Test cases
	cases := map[string][]string{
		"1": {
			"UT310666",
			"AKL OKA4502988001 0x0x0cm",
			"UT", "OKA4502988001",
		},
		"2": {
			"RD2441595",
			"OKA6502998011",
			"ADJ", "OKA6502998011",
		},
		"3": {
			"RU2390034",
			"OKA2692393001-2692393001",
			"RUR", "OKA2692393001",
		},
		"4": {
			"UT332616",
			"AKL OKA1502319001 0x0x0cm",
			"UT", "OKA1502319001",
		},
		"5": {
			"OKA3240408901",
			"Delivered-Wellington                MX5555",
			"NOR", "OKA3240408901",
		},
	}

	// Repeat test for each test case
	for test, expected := range cases {
		actualtype, actualconsignment := GetInvoiceCostTypeAndConsignment(
			expected[0], expected[1],
		)

		if actualtype != expected[2] || actualconsignment != expected[3] {
			// Test failed
			t.Fatalf(
				`TestGetInvoiceCostTypeAndConsignment2 #%s 
				did not return %s, %s got %s, %s`, test, expected[2],
				expected[3], actualtype, actualconsignment,
			)
		}
	}
}

func TestGetInvoiceDate(t *testing.T) {
	a1 := "Statement Invoice 3452349 for Account: 200779   -   Divers Group Trust - Warehouse                    FOR PERIOD ENDING: 10 Jul 2022"
	test := GetInvoiceDate(a1)

	if test != "11-07-2022" {
		// Test failed
		t.Fatalf(
			`TestGetInvoiceDate did not return 11-07-2022, got %s`,
			test,
		)
	}
}
