// data.go contains functions to format items into a struct and
// general management of PBT data
package main

import (
	"strings"
)

func Create200779Rows(worksheetRows [][]string) []PBTItem {
	var pbtRows []PBTItem
	customers := OpenConfigJSON("customers")
	sales := OpenConfigJSON("sales")

	for _, row := range worksheetRows[1:] {
		customerRef := GetCustomerRef(row[2], row[23])
		// Create new pbtItem from each row
		pbtItem := PBTItem{
			ConsignmentDate: row[4][0:10],
			ManifestNum:     row[3],
			Consignment:     row[0],
			CustomerRef:     customerRef,
			ReceiverName:    strings.ToUpper(row[7]),
			AreaTo:          GetRegion(row[18]),
			TrackingNumber:  row[1],
			Weight:          row[11],
			Cubic:           row[12],
			SortbyCode:      GetSortbyCode(customerRef, row[7], customers, sales),
		}

		// add pbtItem to array of pbtDBRows
		pbtRows = append(pbtRows, pbtItem)
	}

	return pbtRows
}
