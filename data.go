// data.go contains functions to format items into a struct and
// general management of PBT data
package main

import (
	"strconv"
	"strings"
)

// Create200779Rows formats a [][]string of rows from a PBTOne
// worksheet into a []PBTItem.
func Create200779Rows(worksheetRows [][]string) []PBTItem {
	var pbtRows []PBTItem
	customers := OpenConfigJSON("customers")
	sales := OpenConfigJSON("sales")

	for _, row := range worksheetRows[1:] {
		// Create new pbtItem from each row
		pbtItem := PBTItem{
			ConsignmentDate: row[4][0:10],
			ManifestNum:     row[3],
			Consignment:     row[0],
			CustomerRef:     GetCustomerRef(row[2], row[23]),
			ReceiverName:    strings.ToUpper(row[7]),
			AreaTo:          GetRegion(row[18]),
			TrackingNumber:  row[1],
			Weight:          row[11],
			Cubic:           row[12],
			SortbyCode:      GetSortbyCode(row[2], row[7], customers, sales),
		}

		// add pbtItem to array of pbtDBRows
		pbtRows = append(pbtRows, pbtItem)
	}

	return pbtRows
}

// CreateInvoiceRows formats a [][]string of rows from an invoice
// worksheet into a []PBTItem.
func CreateInvoiceRows(worksheetRows [][]string) []PBTItem {
	var pbtRows []PBTItem

	var item_cost int

	for _, row := range worksheetRows[1:] {
		item_cost, _ = strconv.Atoi(row[10])

		if item_cost > 0 {
			// Create new pbtItem from each row
			pbtItem := PBTItem{
				Consignment: row[0],
			}

			// add pbtItem to array of pbtDBRows
			pbtRows = append(pbtRows, pbtItem)
		}
	}

	return pbtRows
}
