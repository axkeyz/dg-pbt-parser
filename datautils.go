// datautils.go contains utility functions to extract and format
// data from strings for the final accounting spreadsheet
package main

import (
	"fmt"
	"strings"
	"time"
)

// GetCustomerRef combines the main reference (mainRef) with the
// secondary reference (subRef) if applicable.
func GetCustomerRef(mainRef string, subRef string) string {
	if subRef != "" {
		// Combine main reference with secondary reference
		mainRef = mainRef + " (" + subRef + ")"
	}
	// Return output
	return strings.ToUpper(mainRef)
}

// GetSortbyCode attempts to extract the sortby code from the
// customer reference (customerRef). If the corresponding value
// cannot be obtained, then it is returned as "UNKNOWN".
func GetSortbyCode(
	customerRef string, receiverName string,
	customers map[string][]string, salesCustomers map[string][]string,
) string {
	sortbyCode := strings.ToUpper(StripNonLetters(customerRef))

	if ok, sortbyCode := HasDGSortbyCode(sortbyCode, customers); ok {
		fmt.Println(customerRef, receiverName, sortbyCode)
		// Retrun standardised DG sortby code
		return sortbyCode
	}

	// Try getting ecom/sales code from the receiver name
	ecomSalesCode := TryDGSalesEComSortbyCode(receiverName, salesCustomers)
	if ecomSalesCode == "UNKNOWN" &&
		StripNonLetters(customerRef) == "EC" {
		// Check if ecom/sales code by customer reference
		return sortbyCode
	}

	// Return code
	return ecomSalesCode
}

// HasDGSortbyCode returns true if the given sortby code (sortbyCode) has a
// DG-specific version in the given customers map.
func HasDGSortbyCode(
	sortbyCode string, customers map[string][]string,
) (bool, string) {
	if sortbyCode == "" {
		return false, ""
	}
	for correctedCode, customer := range customers {
		for _, c := range customer {
			if strings.Contains(c, sortbyCode) ||
				strings.Contains(sortbyCode, c) {
				return true, correctedCode
			}
		}
	}
	// Does not have a DG sortby code
	return false, ""
}

// TryDGSalesEComSortbyCode tries to detect if the receiving customer's name
// (receiverName) is a sales customer outlined in the sales customers map.
// It generalises other customers with two words as an e commerce customer
// and everything else as unknown. TryDGSalesEComSortbyCode will always return
// a sortby code of SALES (or whatever codes given in the salesCustomers map),
// E COMMERCE or UNKNOWN depending on the situation.
func TryDGSalesEComSortbyCode(
	receiverName string, salesCustomers map[string][]string,
) string {
	receiverName = strings.ToUpper(receiverName)
	// Iterate through sales customers to check if the receiver name
	// matches any of the sales customer aliases.
	for sortbyCode, names := range salesCustomers {
		for _, name := range names {
			if strings.Contains(name, receiverName) ||
				strings.Contains(receiverName, name) {
				return sortbyCode
			}
		}
	}

	// Allow name look-alikes to pass as E Commerce
	if len(strings.Split(receiverName, " ")) == 2 {
		return "E COMMERCE"
	}

	// Failed to detect sortbyCode
	return "UNKNOWN"
}

// GetRegion attempts to extract the region when given the PBT
// depot name (depotName).
func GetRegion(depotName string) string {
	if strings.Contains(depotName, "Auckland") {
		// All Auckland depots are returned as Auckland
		return "AUCKLAND"
	}

	// Remove first two words (PBT Couriers) and last word (Depot)
	name := strings.Split(depotName, " ")[2:]
	return strings.ToUpper(strings.Join(name[:len(name)-1], " "))
}

// GetInvoiceCostTypeAndConsignment attempts to get the cost type
// and the consignment number.
// Cost type codes:
//	- RUR: Rural
//	- ADJ: Adjustment
//	- CL: CL-type items
//	- UT: Underticket
//	- NOR: Normal
func GetInvoiceCostTypeAndConsignment(reference string,
	description string) (costtype string, consignment string) {
	reference = strings.ToUpper(strings.Split(reference, " ")[0])
	description = strings.ToUpper(description)

	if strings.Contains(reference, "RU") {
		return "RUR", strings.Split(description, "-")[0]
	} else if strings.Contains(reference, "RD") {
		return "ADJ", strings.Split(description, " ")[0]
	} else if strings.Contains(reference, "CL") {
		return "CL", reference
	} else if strings.Contains(reference, "UT") ||
		reference == "" {
		return "UT", strings.Split(description, " ")[1]
	} else {
		return "NOR", reference
	}
}

// GetInvoiceDate gets the invoice date from the A1 cell.
func GetInvoiceDate(a1 string) string {
	cell := strings.Split(a1, " ")
	date := cell[len(cell)-3:]
	a1 = strings.Join(date, " ")
	return FormatDate(a1, "02 Jan 2006")
}

// GetItemDate returns an item date in an invoice spreadsheet
// in the format 2006-01-02.
func GetItemDate(date string, year string) string {
	return FormatDate(date+" "+year[:4], "02 Jan 2006")
}

// FormatDate standardises a date string of the given format
// to the format 2006-01-02
func FormatDate(date string, format string) string {
	t, _ := time.Parse(format, date)
	return t.Format("2006-01-02")
}
