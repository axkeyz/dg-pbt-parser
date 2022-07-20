// database.go contains functions related to SQL CRUD operations.
package main

import (
	"database/sql"
	"fmt"
	"strings"
)

// CreateDB creates a new "pbt" table in the database if it
// does not exist already.
func CreateDB(database *sql.DB, table string) {
	// Create the database table and columns
	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		consignment_date DATE,
		manifest_num TEXT NULL,
		consignment TEXT NULL,
		customer_ref TEXT NULL,
		receiver_name TEXT NULL,
		area_to TEXT NULL,
		tracking_number TEXT,
		weight TEXT NULL,
		cubic TEXT NULL,
		item_cost TEXT NULL,
		sortby_code TEXT NULL,
		rural_delivery TEXT NULL,
		under_ticket TEXT NULL,
		adjustment TEXT NULL,
		other TEXT NULL,
		ff_item TEXT NULL,
		first_invoice DATE NULL,
		last_invoice DATE NULL
	);`, table)

	_, err := database.Exec(query)

	// Display the error
	if err != nil {
		fmt.Println(err)
	}
}

// RefreshDBRows clears all rows in the given table name.
func RefreshDBRows(database *sql.DB, table string) {
	// Delete the rows in the table
	query := fmt.Sprintf(
		"DELETE FROM %s", table,
	)

	_, err := database.Exec(query)
	FormatError(err)
}

// NewDBRow creates a new PBT row in the table, with the details
// in a PBTItem struct.
func NewDBRow(
	database *sql.DB, table string, item PBTItem, unique bool,
) {
	// Do not insert data if row already exists
	if IsRowInDB(database, table, item.TrackingNumber) &&
		unique {
		return
	}

	columns, data := item.GetNonEmptyCols()

	// Add PBT item to table in database
	query := fmt.Sprintf(
		`INSERT INTO %s (%s) VALUES(%s)`,
		table, strings.Join(columns, ", "),
		strings.Join(data, ", "),
	)

	// Display the error if any
	_, err := database.Exec(query)
	FormatError(err)
}

// IsRowInDB returns true if the given tracking number exists in the
// named table of the database.
func IsRowInDB(database *sql.DB, table string, tracking string) bool {
	var exists bool
	query := fmt.Sprintf(
		`SELECT exists (SELECT id from %s
		WHERE tracking_number = '%s')`, table,
		tracking,
	)

	err := database.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		FormatError(err)
	}

	return exists
}

// IsRowInDBCurrentMonth returns a bool of whether there is a
// record in the DB with the same tracking number and in the current
// month or has no month. If this is true, it returns the row ID
// as a string.
func IsRowInDBCurrentMonth(database *sql.DB, table string,
	item PBTItem) (exists bool, first bool, id string) {
	var first_invoice sql.NullString
	query := fmt.Sprintf(
		`SELECT id, first_invoice from %s WHERE tracking_number = '%s' 
		AND (first_invoice IS NULL OR
			strftime('%%m', first_invoice) = 
			'%s')`, table,
		item.TrackingNumber, item.FirstInvoice,
	)

	err := database.QueryRow(query).Scan(&id, &first_invoice)
	if err != nil && err != sql.ErrNoRows {
		FormatError(err)
	}

	if id != "" {
		exists = true
	}

	if first_invoice.String != "" {
		fmt.Println(first_invoice.String)
		first = true
	}

	return exists, first, id
}

// UpdateDBRow updates a row if the item.first_invoice is in the same
// month as the first_invoice in the database. Otherwise, it will create
// a new row.
func UpdateDBForInvoices(database *sql.DB, table string, item PBTItem) {
	isRow := IsRowInDB(database, table, item.TrackingNumber)
	isCurrentRow, isFirst, rowID := IsRowInDBCurrentMonth(database, table, item)

	var query string

	if isCurrentRow {
		// Change invoice dates
		if isFirst {
			item.SwapInvoiceDates()
			// Do something
		}
		// Update the row with non-empty rows
		query = fmt.Sprintf(
			"UPDATE %s SET %s WHERE id = %s",
			table, item.GetNonEmptyColsEqualised(),
			rowID,
		)
		// Execute the query
		_, err := database.Exec(query)
		FormatError(err)
		// fmt.Println(query)
	} else if isRow {
		// Copy the data
		query = fmt.Sprintf(
			`INSERT INTO %s (consignment_date, manifest_num,
				consignment, customer_ref, receiver_name,
				area_to, tracking_number, weight, cubic,
				sortby_code) SELECT consignment_date, manifest_num,
				consignment, customer_ref, receiver_name,
				area_to, tracking_number, weight, cubic,
				sortby_code FROM %s WHERE tracking_number = '%s'`,
			table, table, item.TrackingNumber,
		)
		// Execute the query
		_, err := database.Exec(query)
		FormatError(err)
		// fmt.Println(query)

		// Refresh the data
		UpdateDBForInvoices(database, table, item)
	} else {
		// Create a new row with non-empty rows
		item.SortbyCode = "UNKNOWN"
		NewDBRow(database, table, item, false)
	}
}
