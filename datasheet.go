// datautils.go contains functions to extract and import data into
// the given spreadsheets
package main

import (
	"database/sql"
	"path/filepath"

	"github.com/xuri/excelize/v2"
)

// GetMatchingRunsheets returns a []string of all files that
// have a file name matching the givern pattern.
func GetMatchingRunsheets(pattern string) []string {
	// Search for a file that matches the runsheet file name
	matches, err := filepath.Glob(pattern)
	FormatError(err)

	return matches
}

// ExtractPBTONERunsheet extracts the data from the named
// PBTONE runsheet.
func ExtractPBTONERunsheet(runsheet string) [][]string {
	var rows [][]string

	// Open the runsheet
	f, err := excelize.OpenFile(runsheet)
	FormatError(err)
	// Get all the rows in the Sheet1.
	rows, err = f.GetRows("RunSheet_exporting")
	FormatError(err)

	return rows
}

// CreateAll200779Rows creates new rows in the 200779 database
// where the data is extracted from the available PBTOne spreadsheets.
func CreateAll200779DBRows(database *sql.DB) {
	// Get the PBT runsheets
	runsheets := GetMatchingRunsheets("uploads/*runsheet_exporting*.xls*")

	for _, runsheet := range runsheets {
		// Create a PBT item row for each line in spreadsheet
		pbtRows := Create200779Rows(ExtractPBTONERunsheet(runsheet))

		// Create a new row in the database for each PBT item row
		for _, row := range pbtRows {
			NewDBRow(database, "pbt_200779", row, true)
		}
	}
}

// TODO: Get all invoices and extract needed data
// func GetAllInvoices(database *sql.DB) {
// 	invoices := GetMatchingRunsheets("uploads/WebCSVStmt_*")
// }
