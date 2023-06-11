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
			ManifestNumber:  row[3],
			Consignment:     row[0],
			CustomerRef:     GetCustomerRef(row[2], row[23]),
			ReceiverName:    strings.ToUpper(row[7]),
			AreaTo:          GetRegion(row[18]),
			TrackingNumber:  row[1],
			Weight:          row[11],
			Cubic:           row[12],
			SortbyCode:      GetSortbyCode(row[2], row[7], customers, sales),
			Account:         "200779",
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
	var account = GetAccount(worksheetRows[0][0])

	for key, row := range worksheetRows {
		if IsDataStart(row) {
			continue
		} else if IsDataEnd(row) {
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
				Account:        account,
			}

			// Add details for CL-type items (which only contains
			// details on the invoice)
			if costtype == "CL" {
				item.GetCLDetails(worksheetRows[key : key+2])
			} else if costtype == "CT" {
				item.GetCTDetails(row)
			}

			// Set the cost
			item.SetCost(row[9], costtype)

			// add pbtItem to array of pbtDBRows
			pbtRows = append(pbtRows, item)
		}
	}

	return pbtRows
}

// Create23635Items returns a []PBTItem from a [][]string
// of PBT 23635 worksheet array.
func Create23635Items(worksheet [][]string) []PBTItem {
	var items []PBTItem
	a1 := strings.Split(worksheet[0][0], " ")
	var invdate = FormatDate(a1[len(a1)-1], "02/01/2006")

	for _, row := range worksheet {
		if Is23635DataDiscard(row) {
			continue
		}

		items = append(items, PBTItem{
			ConsignmentDate: FormatDate(row[0], "02/01/2006"),
			Consignment:     strings.ToUpper(row[1]),
			TrackingNumber:  strings.ToUpper(row[1]),
			CustomerRef:     strings.ToUpper(row[3]),
			ReceiverName:    strings.ToUpper(strings.Split(row[2], "  ")[0]),
			FirstInvoice:    invdate,
			LastInvoice:     invdate,
			Account:         "23635",
			FFItem:          row[4],
			Weight:          row[5],
			Cubic:           row[6],
			SortbyCode:      GetDefaultSortby(row[3], row[2]),
			Other:           row[7],
			AreaTo:          strings.ToUpper(row[10]),
		})
	}

	return items
}

// Checks if the data is a discarded comment before
// the start of the needed data.
func IsDataStart(row []string) bool {
	if strings.Contains(row[0], "Statement") ||
		strings.Contains(row[0], "GST") {
		return true
	}
	return false
}

// Checks if the data is a discarded comment after
// the end of the needed data.
func IsDataEnd(row []string) bool {
	if strings.Contains(row[3], "VFC") ||
		strings.Contains(row[3], "NETT AMT") ||
		strings.Contains(row[9], "----") {
		return true
	}

	return false
}

// Checks if the data is a discarded comment of 23635 invoices.
func Is23635DataDiscard(row []string) bool {
	return strings.Contains(row[0], "Invoice") ||
		strings.Contains(row[0], "GST") ||
		strings.Contains(row[0], "DATE") ||
		row[0] == ""
}
