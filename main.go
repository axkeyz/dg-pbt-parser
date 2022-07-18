package main

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
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

func Create200779Rows(worksheetRows [][]string) []PBTItem {
	var pbtRows []PBTItem

	for _, row := range worksheetRows[1:] {
		// Create new pbtDBItem
		pbtItem := PBTItem{
			ConsignmentDate: row[4][0:10],
			ManifestNum:     row[3],
			Consignment:     row[0],
			CustomerRef:     GetCustomerRef(row[2], row[23]),
			ReceiverName:    strings.ToUpper(row[7]),
			AreaTo:          strings.ToUpper(row[26]),
			TrackingNumber:  row[1],
			Weight:          row[11],
			Cubic:           row[12],
			SortbyCode:      "",
		}

		// add pbtItem to array of pbtDBRows
		pbtRows = append(pbtRows, pbtItem)
	}

	return pbtRows
}

// GetCustomerRef combines the main reference (mainRef) with the
// secondary reference (subRef) if applicable.
func GetCustomerRef(mainRef string, subRef string) string {
	if subRef != "" {
		// Combine main reference with secondary reference
		mainRef = mainRef + " (" + subRef + ")"
	}
	// Return output
	return mainRef
}

func main() {
	database, _ := sql.Open("sqlite3", "./dgpbt.db")

	CreateDB(database, "pbt_200779")

	Create200779Rows(ExtractPBTONERunsheet())

	fmt.Println("Done")
}
