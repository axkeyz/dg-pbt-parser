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
			ConsignmentDate: FormatDate(row[4][0:10], "02-01-2006"),
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

	var cost float64
	var invoicedate = GetInvoiceDate(worksheetRows[0][0])

	for key, row := range worksheetRows {
		if strings.Contains(row[0], "Statement") ||
			strings.Contains(row[0], "GST") {
			continue
		} else if strings.Contains(row[3], "VFC") ||
			strings.Contains(row[3], "NETT AMT") ||
			strings.Contains(row[9], "----") {
			break
		}

		cost, _ = strconv.ParseFloat(row[9], 32)

		if int(cost*100) > 0 {
			// Create new pbtItem from each row
			costtype, consignment := GetInvoiceCostTypeAndConsignment(row[1], row[3])

			// Set the consignment value
			item := PBTItem{
				TrackingNumber: consignment,
				FirstInvoice:   invoicedate,
			}

			// Add details for CL-type items (which only contains
			// details on the invoice)
			if costtype == "CL" {
				item.GetCLDetails(worksheetRows[key : key+2])
			}

			// Set the cost
			item.SetCost(row[9], costtype)

			// add pbtItem to array of pbtDBRows
			pbtRows = append(pbtRows, item)
		}
	}

	return pbtRows
}
