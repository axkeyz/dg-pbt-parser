// datautils.go contains utility functions to extract and format
// data from strings for the final accounting spreadsheet
package main

import "strings"

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
		return sortbyCode
	}
	return TryDGSalesEComSortbyCode(receiverName, salesCustomers)
}

// HasDGSortbyCode returns true if the given sortby code (sortbyCode) has a
// DG-specific version in the given customers map.
func HasDGSortbyCode(
	sortbyCode string, customers map[string][]string,
) (bool, string) {
	for correctedCode, customer := range customers {
		for _, c := range customer {
			if c == sortbyCode {
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
