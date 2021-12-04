package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"path/filepath"
	"github.com/xuri/excelize/v2"
	"strings"
)

// CreateDB creates a new "pbt" table in the database if it
// does not exist already.
func CreateDB(database *sql.DB) {
	// Execute the query
	_, err := database.Exec(`CREATE TABLE IF NOT EXISTS pbt_200779 (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		consignment_date DATE,
		manifest_num TEXT NULL,
		consignment TEXT NULL,
		customer_ref TEXT NULL,
		receivers_name TEXT NULL,
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
	);`)

	if err != nil {
		// Display the error
		fmt.Println(err)
	}

	return
}

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

func Create200779Rows() {
	var consignment_date, manifest_num, consignment string
	var customer_ref, receivers_name, area_to string
	var weight, cubic string
	// var under_ticket, adjustment float64
	var tracking_number, sortby_code string

	rows := ExtractPBTONERunsheet()

	for _, row := range rows[1:] {
		consignment = row[0]
		tracking_number = row[1]
		customer_ref = row[2]
		if row[23] != "" {
			customer_ref = customer_ref + " (" + row[23] + ")"
		}
		manifest_num = row[3]
		consignment_date = row[4][0:10]
		receivers_name = strings.ToUpper(row[7])
		weight = row[11]
		cubic = row[12]
		area_to = strings.ToUpper(row[26])
		sortby_code = ""

		fmt.Println(consignment_date, manifest_num, consignment,
		customer_ref, receivers_name, area_to, weight, cubic,
		tracking_number, sortby_code)
	}
}

func main() {
	database, _ := sql.Open("sqlite3", "./dgpbt.db")

	CreateDB(database)

	Create200779Rows()

	fmt.Println("Done")

	fmt.Println("This won't run")
}