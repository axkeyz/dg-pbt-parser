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
	FirstInvoice    string
	LastInvoice     string
}
