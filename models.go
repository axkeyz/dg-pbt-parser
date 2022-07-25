package main

import (
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// PBTItem represents a row of PBT data (one row per freight
// item).
type PBTItem struct {
	ID              string
	ConsignmentDate string
	ManifestNumber  string
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
	Account         string
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
		"manifest_number":  item.ManifestNumber,
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
		"account":          item.Account,
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
	item.ManifestNumber = strings.ToUpper(row[1][1])
	item.Weight = strings.ToUpper(row[0][7])
	item.Cubic = strings.ToUpper(row[0][8])
	item.FFItem = strings.ToUpper(row[0][6])
	item.LastInvoice = item.FirstInvoice
}

// *PBTItem.GetCTDetails creates the details of a CT-type
// item (admin charges)
func (item *PBTItem) GetCTDetails(row []string) {
	item.CustomerRef = "ADMIN CHARGE"
	item.ManifestNumber = "ADMIN CHARGE"
	item.ReceiverName = "ADMIN CHARGE"
	item.ConsignmentDate = GetItemDate(row[0], item.FirstInvoice)
	item.LastInvoice = item.FirstInvoice
}

// *PBTItem.ToExcel writes the details of the PBTItem in the excel
// file and sheet in the given row.
func (item *PBTItem) ToExcel(
	file *excelize.File, sheet string, row int,
) {
	r := strconv.Itoa(row)

	// Set default columns
	file.SetCellValue(sheet, "A"+r, item.ConsignmentDate)
	file.SetCellValue(sheet, "B"+r, item.ManifestNumber)
	file.SetCellValue(sheet, "C"+r, item.Consignment)
	file.SetCellValue(sheet, "D"+r, item.CustomerRef)
	file.SetCellValue(sheet, "E"+r, item.ReceiverName)
	file.SetCellValue(sheet, "F"+r, item.AreaTo)
	file.SetCellValue(sheet, "G"+r, item.TrackingNumber)
	SetCellFloat(file, sheet, "I"+r, item.Weight)
	SetCellFloat(file, sheet, "J"+r, item.Cubic)
	file.SetCellValue(sheet, "N"+r, item.SortbyCode)
	SetCellFloat(file, sheet, "M"+r, item.ItemCost)
	SetCellFloat(file, sheet, "P"+r, item.RuralDelivery)
	SetCellFloat(file, sheet, "R"+r, item.UnderTicket)
	SetCellFloat(file, sheet, "S"+r, item.Adjustment)
	file.SetCellFormula(sheet, "O"+r, "=IF(M"+r+">0, M"+r+", \"\")")
	SetCellFloat(file, sheet, "T"+r, item.Other)
	SetCellFloat(file, sheet, "X"+r, item.FFItem)
	file.SetCellFormula(sheet, "U"+r, "=SUM(O"+r+":T"+r+")")
	file.SetCellFormula(sheet, "V"+r, "=U"+r+"*$W$1")
	file.SetCellFormula(sheet, "W"+r, "=V"+r+"+U"+r)
	file.SetCellValue(sheet, "AE"+r,
		item.Account+" "+FormatDBDate(item.FirstInvoice),
	)

	// Set conditional columns based on invoice cost type
	invoicetype, _ := GetInvoiceCostTypeAndConsignment(
		item.TrackingNumber, item.CustomerRef)

	if invoicetype == "NOR" {
		file.SetCellFormula(sheet, "Z"+r,
			"=VLOOKUP(F"+r+",'new rates 29feb16'!K:L,2,FALSE)")
		file.SetCellFormula(sheet, "AA"+r, "=200*J"+r)
		file.SetCellFormula(sheet, "AB"+r, "=ROUND(MAX(I"+r+",AA"+r+"),0)")
		file.SetCellFormula(sheet, "AC"+r,
			"=INDEX(NewRates,MATCH(AB"+r+
				",Weight,-1),MATCH(Z"+r+",Service,0))")
		file.SetCellFormula(sheet, "AD"+r,
			"=O"+r+"-AC"+r)
	}
}
