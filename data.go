// data.go contains functions to format items into a struct and
// general management of PBT data
package main

import "strings"

func Create200779Rows(worksheetRows [][]string) []PBTItem {
	var pbtRows []PBTItem

	for _, row := range worksheetRows[1:] {
		// Create new pbtItem from each row
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
