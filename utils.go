// utils.go contains miscellaneous utility functions.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
)

// StripNonLetters is a utility function that strips all non-letters
// from the given string (str).
func StripNonLetters(str string) string {
	s := []byte(str)
	n := 0
	for _, b := range s {
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') {
			s[n] = b
			n++
		}
	}
	return string(s[:n])
}

// OpenConfigJSON returns a map[string][]string when given a config
// file name of a file in config/.
func OpenConfigJSON(fileName string) map[string][]string {
	var configJSON map[string][]string

	// Get entire contents of file
	content, err := ioutil.ReadFile("config/" + fileName + ".json")
	FormatError(err)

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(content, &configJSON)

	// Return content
	return configJSON
}

// IsInArray returns true if a string is in the []string.
func IsInArray(item string, array []string) bool {
	for _, i := range array {
		if item == i {
			return true
		}
	}
	return false
}

// FormatError is a generic function to format an error if
// applicable.
func FormatError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// StringToFloat converts a string to a float64.
func StringToFloat(intstring string) float64 {
	if intstring == "" {
		return 0
	}
	intfloat, err := strconv.ParseFloat(intstring, 64)
	FormatError(err)
	return intfloat
}

// StringToFloat converts a string to an int.
func StringToInt(intstring string) int {
	if intstring == "" {
		return 0
	}

	intfloat, err := strconv.Atoi(intstring)
	FormatError(err)

	return intfloat
}
