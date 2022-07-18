// database.go contains functions related to SQL CRUD operations.
package main

import (
	"database/sql"
	"fmt"
)

// CreateDB creates a new "pbt" table in the database if it
// does not exist already.
func CreateDB(database *sql.DB, tableName string) {
	// Create the database table and columns
	createTableSQL := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		consignment_date DATE,
		manifest_num TEXT NULL,
		consignment TEXT NULL,
		customer_ref TEXT NULL,
		receiver_name TEXT NULL,
		area_to TEXT NULL,
		tracking_number TEXT UNIQUE,
		weight FLOAT NULL,
		cubic FLOAT NULL,
		item_cost FLOAT NULL,
		sortby_code TEXT NULL,
		rural_delivery FLOAT NULL,
		under_ticket FLOAT NULL,
		adjustment FLOAT NULL,
		first_invoice DATE,
		last_invoice DATE
	);`, tableName)

	_, err := database.Exec(createTableSQL)

	// Display the error
	if err != nil {
		fmt.Println(err)
	}
}
