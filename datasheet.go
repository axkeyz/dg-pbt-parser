// datautils.go contains functions to extract and import data into
// the given spreadsheets
package main

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"strings"

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
func ExtractSheet(runsheet string, format string) [][]string {
	var rows [][]string

	if format == "xlsx" {
		// Open the runsheet
		f, err := excelize.OpenFile(runsheet)
		FormatError(err)
		// Get all the rows in the Sheet1.
		rows, err = f.GetRows(f.GetSheetList()[0])
		FormatError(err)
	} else {
		f, err := ioutil.ReadFile(runsheet)
		FormatError(err)

		items := strings.Split(string(f), "\n")

		for _, item := range items {
			row := make([]string, 11)
			i := strings.Split(item, ",")
			for j, k := range i {
				i[j] = strings.TrimSpace(k)
			}
			copy(row, i)
			rows = append(rows, row)
		}
	}

	return rows
}

// CreateAll200779Rows creates new rows in the 200779 database
// where the data is extracted from the available PBTOne spreadsheets.
func CreateAll200779DBRows(database *sql.DB) {
	// Get the PBT runsheets
	runsheets := GetMatchingRunsheets("uploads/*runsheet_exporting*.xls*")

	for _, runsheet := range runsheets {
		// Create a PBT item row for each line in spreadsheet
		pbtRows := Create200779Rows(ExtractSheet(runsheet, "xlsx"))

		// Create a new row in the database for each PBT item row
		for _, row := range pbtRows {
			NewDBRow(database, "pbt_200779", row, true)
		}
	}
}

// Get all invoices and extract needed data
func CreateAllInvoiceCosts(database *sql.DB, table string) {
	// func CreateAllInvoiceCosts() {
	invoices := GetMatchingRunsheets("uploads/WebCSVStmt_*")
	// var costtype, consignment string

	for _, invoice := range invoices {
		// fmt.Println(invoice)
		// fmt.Println(ExtractSheet(invoice, "csv")[10][10])
		items := CreateInvoiceRows(ExtractSheet(invoice, "csv"))
		// fmt.Println(items)
		for _, item := range items {
			UpdateDBForInvoices(database, table, item)
		}
	}
}
