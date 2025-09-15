package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Read all_depts.csv
	allDeptsFile, err := os.Open("all_depts.csv")
	if err != nil {
		fmt.Printf("Error opening all_depts.csv: %v\n", err)
		return
	}
	defer allDeptsFile.Close()

	allDeptsReader := csv.NewReader(allDeptsFile)
	allDepts, err := allDeptsReader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading all_depts.csv: %v\n", err)
		return
	}

	// Create a map of all departments for quick lookup
	allDeptsMap := make(map[string]bool)
	for i, record := range allDepts {
		if i == 0 {
			continue // Skip header
		}
		if len(record) > 0 {
			// Trim whitespace for exact case-sensitive comparison
			dept := strings.TrimSpace(record[0])
			allDeptsMap[dept] = true
			// Case-insensitive version (commented out):
			// allDeptsMap[strings.ToLower(dept)] = true
		}
	}

	// Read ranils_new_depts.csv
	ranilsFile, err := os.Open("ranils_new_depts.csv")
	if err != nil {
		fmt.Printf("Error opening ranils_new_depts.csv: %v\n", err)
		return
	}
	defer ranilsFile.Close()

	ranilsReader := csv.NewReader(ranilsFile)
	ranilsRecords, err := ranilsReader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading ranils_new_depts.csv: %v\n", err)
		return
	}

	// Find exact matches and non-matches
	var exactMatches [][]string
	var nonMatches [][]string
	exactMatches = append(exactMatches, []string{"minister", "department"}) // Add header
	nonMatches = append(nonMatches, []string{"minister", "department"})     // Add header

	for i, record := range ranilsRecords {
		if i == 0 {
			continue // Skip header
		}
		if len(record) >= 2 {
			minister := strings.TrimSpace(record[0])
			dept := strings.TrimSpace(record[1])

			// Check for exact match (case-sensitive)
			if allDeptsMap[dept] {
				exactMatches = append(exactMatches, []string{minister, dept})
			} else {
				// Department not found in all_depts.csv
				nonMatches = append(nonMatches, []string{minister, dept})
			}
			// Case-insensitive version (commented out):
			// if allDeptsMap[strings.ToLower(dept)] {
			//	exactMatches = append(exactMatches, []string{minister, dept})
			// } else {
			//	nonMatches = append(nonMatches, []string{minister, dept})
			// }
		}
	}

	// Write exact matches to new file
	exactMatchesFile, err := os.Create("exact_matches_departments.csv")
	if err != nil {
		fmt.Printf("Error creating exact matches file: %v\n", err)
		return
	}
	defer exactMatchesFile.Close()

	exactWriter := csv.NewWriter(exactMatchesFile)
	defer exactWriter.Flush()

	err = exactWriter.WriteAll(exactMatches)
	if err != nil {
		fmt.Printf("Error writing exact matches: %v\n", err)
		return
	}

	// Write non-matches to separate file
	nonMatchesFile, err := os.Create("non_matches_departments.csv")
	if err != nil {
		fmt.Printf("Error creating non-matches file: %v\n", err)
		return
	}
	defer nonMatchesFile.Close()

	nonMatchesWriter := csv.NewWriter(nonMatchesFile)
	defer nonMatchesWriter.Flush()

	err = nonMatchesWriter.WriteAll(nonMatches)
	if err != nil {
		fmt.Printf("Error writing non-matches: %v\n", err)
		return
	}

	fmt.Printf("Found %d exact matches. Results written to exact_matches_departments.csv\n", len(exactMatches)-1)
	fmt.Printf("Found %d non-matches. Results written to non_matches_departments.csv\n", len(nonMatches)-1)
}
