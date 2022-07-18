package main

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

// SetPBTItemCost sets PBT a cost depending on the
// given cost type.
func SetPBTItemCost(cost string, costtype string,
	item *PBTItem) {
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
	}
}
