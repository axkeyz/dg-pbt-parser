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
		manifest_number TEXT NULL,
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
		last_invoice DATE NULL,
		account TEXT NULL
	);`, table)

	_, err := database.Exec(query)
	FormatError(err)
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

// IsRowInDBCurrent returns a bool of whether there is a
// record in the DB with the same tracking number and in the current
// month or has no month. It also returns the latest rowID if there
// is a record (regardless of month).
func IsRowInDBCurrent(database *sql.DB, table string,
	item *PBTItem) (exists bool, id string) {
	var date sql.NullString
	month := strings.Split(item.FirstInvoice, "-")[1]
	query := fmt.Sprintf(
		`SELECT id, first_invoice from %s WHERE tracking_number = '%s'
		AND (first_invoice LIKE "%%-%s-%%" OR first_invoice IS NULL)`,
		table, item.TrackingNumber, month,
	)

	err := database.QueryRow(query).Scan(&id, &date)
	if err != nil && err != sql.ErrNoRows {
		FormatError(err)
	}

	if id != "" {
		d := date.String

		if d == item.FirstInvoice || d == "" {
			// In the same invoice
			item.LastInvoice = item.FirstInvoice
			exists = true
		} else if d != "" && strings.Split(d, "-")[1] ==
			month {
			// In the same month
			item.SwapInvoiceDates()
			exists = true
		}
	}

	return exists, id
}

// UpdateDBRow updates a row if the item.first_invoice is in the same
// month as the first_invoice in the database. Otherwise, it will create
// a new row.
func UpdateDBForInvoices(database *sql.DB, table string, item PBTItem) {
	isCurrentRow, rowID := IsRowInDBCurrent(database, table, &item)
	isRow := IsRowInDB(database, table, item.TrackingNumber)

	var query string

	if isCurrentRow {
		// Update the row with non-empty rows
		query = fmt.Sprintf(
			"UPDATE %s SET %s WHERE id = %s",
			table, item.GetNonEmptyColsEqualised(),
			rowID,
		)
		// Execute the query
		_, err := database.Exec(query)
		FormatError(err)
	} else if isRow {
		query = fmt.Sprintf(
			`INSERT INTO %s (consignment_date, manifest_number,
				consignment, customer_ref, receiver_name,
				area_to, tracking_number, weight, cubic,
				sortby_code) SELECT consignment_date, manifest_number,
				consignment, customer_ref, receiver_name,
				area_to, tracking_number, weight, cubic,
				sortby_code FROM %s WHERE tracking_number = "%s"`,
			table, table, item.TrackingNumber,
		)

		// Execute the query
		_, err := database.Exec(query)
		FormatError(err)

		// Refresh the data
		UpdateDBForInvoices(database, table, item)
	} else {
		// Try to get sortby code
		customers := OpenConfigJSON("customers")
		sales := OpenConfigJSON("sales")

		item.SortbyCode = GetSortbyCode(
			item.CustomerRef, item.ReceiverName,
			customers, sales,
		)

		// Create a new row with non-empty rows
		NewDBRow(database, table, item, false)
	}
}

// GetDBRowsByMonth returns all rows and columns in the database (db)
// of the given table, where the first_invoice date is in the given month.
// The data is returned as a []PBTItem.
func GetDBRowsByMonth(db *sql.DB, table string, month string) []PBTItem {
	var items []PBTItem
	var item PBTItem

	// Create query
	query := fmt.Sprintf(
		`SELECT id, ifnull(consignment_date, ''),
		ifnull(manifest_number, ''), ifnull(consignment, ''),
		ifnull(customer_ref, ''), ifnull(receiver_name, ''),
		ifnull(area_to, ''), ifnull(tracking_number, ''),
		ifnull(weight, ''), ifnull(cubic, ''),
		ifnull(item_cost, ''), ifnull(sortby_code, ''),
		ifnull(rural_delivery, ''), ifnull(under_ticket, ''),
		ifnull(adjustment, ''), ifnull(other, ''),
		ifnull(ff_item, ''), ifnull(first_invoice, ''),
		ifnull(last_invoice, ''), ifnull(account, '')
		from %s where first_invoice LIKE "%%-%s-%%"`,
		table, month,
	)

	rows, err := db.Query(query)
	FormatError(err)
	defer rows.Close()

	for rows.Next() {
		// Place data into PBTItem struct and display any errors
		err = rows.Scan(&item.ID, &item.ConsignmentDate,
			&item.ManifestNumber, &item.Consignment,
			&item.CustomerRef, &item.ReceiverName,
			&item.AreaTo, &item.TrackingNumber,
			&item.Weight, &item.Cubic,
			&item.ItemCost, &item.SortbyCode,
			&item.RuralDelivery, &item.UnderTicket,
			&item.Adjustment, &item.Other, &item.FFItem,
			&item.FirstInvoice, &item.LastInvoice,
			&item.Account,
		)
		FormatError(err)

		item.ConsignmentDate = FormatExcelDate(item.ConsignmentDate, "2006-01-02")

		// Place into []PBTItem
		items = append(items, item)
	}

	return items
}
