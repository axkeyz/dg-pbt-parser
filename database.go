// database.go contains functions related to SQL CRUD operations.
package main

import (
	"database/sql"
	"fmt"
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
		weight FLOAT NULL,
		cubic FLOAT NULL,
		item_cost FLOAT NULL,
		sortby_code TEXT NULL,
		rural_delivery FLOAT NULL,
		under_ticket FLOAT NULL,
		adjustment FLOAT NULL,
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

	// Display the error
	if err != nil {
		fmt.Println(err)
	}
}

// NewDBRow creates a new PBT row in the table, with the details
// in a PBTItem struct.
func NewDBRow(database *sql.DB, table string, item PBTItem) {
	// Check if row exists
	if IsRowInDB(database, table, item.TrackingNumber) {
		// Do not insert data
		return
	}

	// Add PBT item to table in database
	query := fmt.Sprintf(`INSERT INTO %s(
			consignment_date, manifest_num, consignment,
			customer_ref, receiver_name, area_to,
			tracking_number, weight, cubic, item_cost,
			sortby_code, rural_delivery, under_ticket,
			adjustment, first_invoice, last_invoice
		) VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s',
			'%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s'
		)`, table, item.ConsignmentDate, item.ManifestNum,
		item.Consignment, item.CustomerRef, item.ReceiverName,
		item.AreaTo, item.TrackingNumber, item.Weight, item.Cubic,
		item.ItemCost, item.SortbyCode, item.RuralDelivery,
		item.UnderTicket, item.Adjustment, item.FirstInvoice,
		item.LastInvoice,
	)

	_, err := database.Exec(query)

	// Display the error
	if err != nil {
		fmt.Println(err)
	}
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
		fmt.Println(err)
	}

	return exists
}
