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
func GetSortbyCode(customerRef string, receiverName string) string {
	sortbyCode := strings.ToUpper(StripNonLetters(customerRef))
	return sortbyCode
}

// HasDGSortbyCode returns true if the given sortby code (sortbyCode) has a
// DG-specific version in the given customers map.
func HasDGSortbyCode(sortbyCode string, customers map[string][]string) bool {
	for _, customer := range customers {
		for _, c := range customer {
			if c == sortbyCode {
				return true
			}
		}
	}
	// Does not have a DG sortby code
	return false
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
