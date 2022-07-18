// utils.go contains miscellaneous utility functions.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(content, &configJSON)

	// Return content
	return configJSON
}
