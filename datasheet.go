// datautils.go contains functions to extract and import data into
// the given spreadsheets
package main

import (
	"fmt"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

// ExtractPBTONERunsheet extracts the data from the PBTONE
// runsheet. This is usually a file that looks like
// runsheet_exporting.xlsx
func ExtractPBTONERunsheet() [][]string {
	// Search for a file that matches the runsheet file name
	pattern := "*runsheet_exporting*.xls*"
	var rows [][]string

	matches, err := filepath.Glob(pattern)

	if err != nil {
		// Display the error
		fmt.Println(err)
		return rows
	}

	runsheet := matches[0]

	// Open the runsheet
	f, err := excelize.OpenFile(runsheet)
	if err != nil {
		fmt.Println(err)
		return rows
	}
	// Get all the rows in the Sheet1.
	rows, err = f.GetRows("RunSheet_exporting")
	if err != nil {
		fmt.Println(err)
		return rows
	}

	return rows
}
