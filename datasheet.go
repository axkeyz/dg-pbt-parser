// datautils.go contains functions to extract and import data into
// the given spreadsheets
package main

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

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
func CreateAll200779DBRows(database *sql.DB, table string) {
	// Get the PBT runsheets
	runsheets := GetMatchingRunsheets("uploads/*runsheet_exporting*.xls*")

	for _, runsheet := range runsheets {
		// Create a PBT item row for each line in spreadsheet
		pbtRows := Create200779Rows(ExtractSheet(runsheet, "xlsx"))

		// Create a new row in the database for each PBT item row
		for _, row := range pbtRows {
			NewDBRow(database, table, row, true)
		}
	}
}

// CreateAllInvoiceCosts creates all invoice costs by
// matching corresponding lines in the spreadsheet to
// the database.
func CreateAllInvoiceCosts(db *sql.DB, table string) {
	// Get the runsheets
	invoices := GetMatchingRunsheets("uploads/WebCSVStmt_*")

	for _, invoice := range invoices {
		// Convert rows to PBTItem structs
		items := CreateInvoiceRows(ExtractSheet(invoice, "csv"))

		for _, item := range items {
			// Update database with items
			UpdateDBForInvoices(db, table, item)
		}
	}
}

// Create23635Rows creates rows according to the matching
// spreadsheet under the PBT 23635 account.
func CreateAll23635Rows(db *sql.DB, table string) {
	// Get the invoices
	invoices := GetMatchingRunsheets("uploads/PBT_Invoice_*")

	for _, invoice := range invoices {
		// Convert rows to PBTItem structs
		items := Create23635Items(ExtractSheet(invoice, "csv"))
		// fmt.Println(ExtractSheet(invoice, "csv"))

		for _, item := range items {
			// Update database with items
			UpdateDBForInvoices(db, table, item)
		}
	}
}

// ExportDB exports the database for the given month
func ExportDB(db *sql.DB, table string, month string) {
	// Create new workbook and sheet
	f := excelize.NewFile()
	f.SetSheetName("Sheet1", "PBT")

	// Get data
	items := GetDBRowsByMonth(db, table, month)
	currentTime := time.Now()

	SetHeader(f, "PBT")

	for row, item := range items {
		item.ToExcel(f, "PBT", row+4)
	}

	if err := f.SaveAs("exports/PBT" +
		currentTime.Format("2006-02-01") +
		".xlsx",
	); err != nil {
		FormatError(err)
	}
}

// SetHeader sets a header of an export spreadsheet.
func SetHeader(f *excelize.File, sheet string) {
	f.SetCellValue(sheet, "V1", "FAF")
	SetCellFloat(f, sheet, "W1", "0.1189")

	columns := []string{
		"consignment_date", "manifest_number", "consignment",
		"customer_consignment", "receivers_name", "area_to",
		"tracking_number", "product_type", "weight", "cubic",
		"dg_charge", "sat_del_charge", "item_cost", "sortby_code",
		"Total Item Cost", "Rural Delivery",
		"Accrue Rural Delivery", "U/Ticket", "Adjs/sheet",
		"Other", "Sub-total", "FAF", "Cost", "# of FF items",
		"# of tickets per parcel", "RATING",
		"Cubic converted to weight", "Highest weight",
		"Expected Charge", "Variance", "Invoice #", "Note",
	}

	for i, col := range columns {
		cell, err := excelize.CoordinatesToCellName(i+1, 3)
		FormatError(err)
		f.SetCellValue(sheet, cell, col)
	}
}

// SetCellFloat converts a string to a float and then adds the value
// into the given cell, of the sheet in the file.
func SetCellFloat(
	file *excelize.File, sheet string, cell string, value string,
) {
	if value == "" {
		// Do not add anything to cell
		return
	}

	// Add value to cell
	file.SetCellFloat(sheet, cell, StringToFloat(value), 2, 64)
}
