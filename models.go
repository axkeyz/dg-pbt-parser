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
	FFItems         string
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
		"ff_items":         item.FFItems,
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
