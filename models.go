package main

import "strings"

// PBTItem represents a row of PBT data (one row per freight
// item).
type PBTItem struct {
	ID              string
	ConsignmentDate string
	ManifestNum     string
	Consignment     string
	CustomerRef     string
	ReceiverName    string
	AreaTo          string
	TrackingNumber  string
	Weight          string
	Cubic           string
	ItemCost        string
	SortbyCode      string
	RuralDelivery   string
	UnderTicket     string
	Adjustment      string
	Other           string
	FFItem          string
	FirstInvoice    string
	LastInvoice     string
}

// *PBTItem.SetCost sets PBT a cost depending on the
// given cost type.
func (item *PBTItem) SetCost(
	cost string, costtype string) {
	switch costtype {
	case "NOR":
		item.ItemCost = cost
	case "RUR":
		item.RuralDelivery = cost
	case "ADJ":
		item.Adjustment = cost
	case "UT":
		item.UnderTicket = cost
	case "CL":
		item.Other = cost
	case "CT":
		item.Other = cost
	default:
		item.ItemCost = cost
	}
}

// *PBTItem.GetNonEmptyCols returns a []string of all
// non-null fields in the given PBTItem.
func (item *PBTItem) GetNonEmptyCols() (
	columns []string, data []string,
) {
	fields := map[string]string{
		"consignment_date": item.ConsignmentDate,
		"manifest_num":     item.ManifestNum,
		"consignment":      item.Consignment,
		"customer_ref":     item.CustomerRef,
		"receiver_name":    item.ReceiverName,
		"area_to":          item.AreaTo,
		"tracking_number":  item.TrackingNumber,
		"weight":           item.Weight,
		"cubic":            item.Cubic,
		"item_cost":        item.ItemCost,
		"sortby_code":      item.SortbyCode,
		"rural_delivery":   item.RuralDelivery,
		"under_ticket":     item.UnderTicket,
		"adjustment":       item.Adjustment,
		"other":            item.Other,
		"ff_item":          item.FFItem,
		"first_invoice":    item.FirstInvoice,
		"last_invoice":     item.LastInvoice,
	}

	for col, val := range fields {
		if val != "" {
			columns = append(columns, col)
			data = append(data, "\""+val+"\"")
		}
	}

	return columns, data
}

// *PBTItem.GetNonEmptyCols returns a []string of all
// non-null fields in the given PBTItem.
func (item *PBTItem) GetNonEmptyColsEqualised() string {
	columns, data := item.GetNonEmptyCols()
	var equalised []string

	for index, col := range columns {
		equalised = append(equalised, col+" = "+data[index])
	}

	return strings.Join(equalised, ", ")
}

// *PBTItem.SwapInvoiceDates swaps the FirstInvoice date with
// the LastInvoice date.
func (item *PBTItem) SwapInvoiceDates() {
	temp := item.FirstInvoice
	item.FirstInvoice = item.LastInvoice
	item.LastInvoice = temp
}

// *PBTItem.GetCLDetails gets the details of a CL-type item
// when given two consecutive rows as a [][]string.
func (item *PBTItem) GetCLDetails(row [][]string) {
	customers := OpenConfigJSON("customers")
	sales := OpenConfigJSON("sales")

	item.ConsignmentDate = GetItemDate(row[0][0], item.FirstInvoice)
	item.Consignment = item.TrackingNumber
	item.ReceiverName = strings.ToUpper(row[0][3][3:])
	item.AreaTo = strings.ToUpper(row[1][3])
	item.CustomerRef = strings.ToUpper(row[0][2])
	item.SortbyCode = GetSortbyCode(
		item.CustomerRef, item.ReceiverName,
		customers, sales,
	)
	item.ManifestNum = strings.ToUpper(row[1][1])
	item.Weight = strings.ToUpper(row[0][7])
	item.Cubic = strings.ToUpper(row[0][8])
	item.FFItem = strings.ToUpper(row[0][6])
}

func (item *PBTItem) GetCTDetails(row []string) {
	// item.ConsignmentDate = GetItemDate(row[0], item.FirstInvoice)
	item.SortbyCode = "DIVERS"
	item.CustomerRef = "ADMIN CHARGE"
	item.ManifestNum = "ADMIN CHARGE"
	item.ReceiverName = "ADMIN CHARGE"
	item.SortbyCode = "DIVERS"
	item.ConsignmentDate = GetItemDate(row[0], item.FirstInvoice)
}
