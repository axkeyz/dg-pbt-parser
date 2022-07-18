// utils_test.go contains tests for utils.go
package main

import "testing"

// TestStripNonLetters calls StripNonLetters and checks if all
// non-letters are correctly removed.
func TestStripNonLetters(t *testing.T) {
	// Test cases
	cases := map[string]string{
		"Ux2348SHDFx":    "UxSHDFx",
		"MuiSGH":         "MuiSGH",
		"az893 U":        "azU",
		"sd5 ek@90-sdIO": "sdeksdIO",
	}

	// Repeat test for each test case
	for test, expected := range cases {
		actual := StripNonLetters(test)

		if actual != expected {
			// Test failed
			t.Fatalf(
				`TestStripNonLetters did not return %s, got %s`,
				expected, actual,
			)
		}
	}
}

// TestOpenConfigJSON calls OpenConfigJSON and checks if the config
// file is successfully mapped and data can be opened correctly.
func TestOpenConfigJSON(t *testing.T) {
	config := OpenConfigJSON("customers_test")

	if config["RTF"][2] != "IN" {
		// Test failed
		t.Fatalf(
			`TestOpenConfigJSON did not return %v, got %v`,
			[]string{"IO", "VSD", "IN", "IXU"}, config["RTF"],
		)
	}
}
