package api

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// ProcessTransactions processes all transactions from CSV files in the specified directory
func (c *Client) ProcessTransactions(dataDir string, processType string) error {
	// Initialize entity counters based on process type
	var entityCounters map[string]int
	if processType == "organisation" {
		entityCounters = map[string]int{
			"minister":   0,
			"department": 0,
		}
	} else if processType == "person" {
		entityCounters = map[string]int{
			"citizen": 0,
		}
	} else {
		return fmt.Errorf("invalid process type: %s", processType)
	}

	// Get all CSV files in the directory
	files, err := os.ReadDir(dataDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dataDir, err)
	}

	// Collect all transactions from all files
	var allTransactions []map[string]interface{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			// Extract file type from filename (e.g., "ADD" from "2403-38_ADD.csv" or "ADD.csv")
			fileName := strings.TrimSuffix(file.Name(), ".csv")
			fileType := "ADD" // Default to ADD
			if strings.Contains(fileName, "TERMINATE") {
				fileType = "TERMINATE"
			} else if strings.Contains(fileName, "MOVE") {
				fileType = "MOVE"
			} else if strings.Contains(fileName, "MERGE") {
				fileType = "MERGE"
			} else if strings.Contains(fileName, "RENAME") {
				fileType = "RENAME"
			}

			// Load transactions from the CSV file
			transactions, err := loadTransactions(filepath.Join(dataDir, file.Name()), fileType)
			if err != nil {
				return fmt.Errorf("failed to load transactions from %s: %w", file.Name(), err)
			}
			allTransactions = append(allTransactions, transactions...)
		}
	}

	// Sort transactions by transaction_id, handling numeric parts correctly
	sort.Slice(allTransactions, func(i, j int) bool {
		idI := allTransactions[i]["transaction_id"].(string)
		idJ := allTransactions[j]["transaction_id"].(string)

		// Split the IDs into parts
		partsI := strings.Split(idI, "_")
		partsJ := strings.Split(idJ, "_")

		// Compare the first part (e.g., "2153/12")
		if partsI[0] != partsJ[0] {
			return partsI[0] < partsJ[0]
		}

		// Compare the second part (e.g., "tr")
		if partsI[1] != partsJ[1] {
			return partsI[1] < partsJ[1]
		}

		// Compare the numeric part by converting to integers
		numI := strings.TrimPrefix(partsI[2], "tr_")
		numJ := strings.TrimPrefix(partsJ[2], "tr_")

		// Convert to integers for numeric comparison
		valI, _ := strconv.Atoi(numI)
		valJ, _ := strconv.Atoi(numJ)
		return valI < valJ
	})

	// Process transactions in order
	for _, transaction := range allTransactions {
		fmt.Printf("Processing transaction: %s (Type: %s)\n", transaction["transaction_id"], transaction["file_type"])

		switch transaction["file_type"] {
		case "ADD":
			// Check if the transaction type matches the process type
			childType := transaction["child_type"].(string)
			if (processType == "organisation" && (childType == "minister" || childType == "department")) ||
				(processType == "person" && childType == "citizen") {
				var err error

				if processType == "person" && childType == "citizen" {
					entityCounters[childType], err = c.AddPersonEntity(transaction, entityCounters)
				} else {
					entityCounters[childType], err = c.AddOrgEntity(transaction, entityCounters)
				}

				if err != nil {
					return fmt.Errorf("failed to process add transaction %s: %w", transaction["transaction_id"], err)
				}
				//entityCounters[childType] = newCounter
				fmt.Printf("Processed Add transaction: %s\n", transaction["transaction_id"])
			} else {
				fmt.Printf("Skipping transaction %s: type %s does not match process type %s\n",
					transaction["transaction_id"], childType, processType)
			}

		case "TERMINATE":
			if processType == "organisation" {
				err := c.TerminateOrgEntity(transaction)
				if err != nil {
					return fmt.Errorf("failed to process terminate transaction %s: %w", transaction["transaction_id"], err)
				}
				fmt.Printf("Processed Terminate transaction: %s\n", transaction["transaction_id"])
			} else if processType == "person" {
				err := c.TerminatePersonEntity(transaction)
				if err != nil {
					return fmt.Errorf("failed to process terminate transaction %s: %w", transaction["transaction_id"], err)
				}
				fmt.Printf("Processed Terminate transaction: %s\n", transaction["transaction_id"])
			}

		case "MOVE":
			if processType == "organisation" {
				err := c.MoveDepartment(transaction)
				if err != nil {
					return fmt.Errorf("failed to process move transaction %s: %w", transaction["transaction_id"], err)
				}
				fmt.Printf("Processed Move transaction: %s\n", transaction["transaction_id"])
			} else if processType == "person" {
				err := c.MovePerson(transaction)
				if err != nil {
					return fmt.Errorf("failed to process move transaction %s: %w", transaction["transaction_id"], err)
				}
				fmt.Printf("Processed Move transaction: %s\n", transaction["transaction_id"])
			}

		case "MERGE":
			if processType == "organisation" {
				newCounter, err := c.MergeMinisters(transaction, entityCounters)
				if err != nil {
					return fmt.Errorf("failed to process merge transaction %s: %w", transaction["transaction_id"], err)
				}
				entityCounters["minister"] = newCounter
				fmt.Printf("Processed Merge transaction: %s\n", transaction["transaction_id"])
			}

		case "RENAME":
			if processType == "organisation" {
				newCounter, err := c.RenameMinister(transaction, entityCounters)
				if err != nil {
					return fmt.Errorf("failed to process rename transaction %s: %w", transaction["transaction_id"], err)
				}
				entityCounters["minister"] = newCounter
				fmt.Printf("Processed Rename transaction: %s\n", transaction["transaction_id"])
			}

		default:
			fmt.Printf("Skipping unknown transaction type: %s\n", transaction["file_type"])
		}
	}

	return nil
}

// loadTransactions reads and processes transactions from a CSV file
func loadTransactions(filePath string, fileType string) ([]map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header from %s: %w", filePath, err)
	}

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read records from %s: %w", filePath, err)
	}

	var transactions []map[string]interface{}
	// Process each record
	for _, record := range records {
		transaction := make(map[string]interface{})
		for i, value := range record {
			transaction[header[i]] = value
		}

		// Add file type to transaction
		transaction["file_type"] = fileType

		// For ADD transactions, handle president column
		if fileType == "ADD" {
			// Check if parent_type is minister
			if parentType, ok := transaction["parent_type"].(string); ok && parentType == "minister" {
				president, ok := transaction["president"].(string)
				if !ok || president == "" {
					return nil, fmt.Errorf("president field is required for minister in ADD transaction")
				}
				// Append president to parent name
				if parent, ok := transaction["parent"].(string); ok {
					transaction["parent"] = parent + president
				}
			}

			// Check if child_type is minister
			if childType, ok := transaction["child_type"].(string); ok && childType == "minister" {
				president, ok := transaction["president"].(string)
				if !ok || president == "" {
					return nil, fmt.Errorf("president field is required for minister in ADD transaction")
				}
				// Append president to child name
				if child, ok := transaction["child"].(string); ok {
					transaction["child"] = child + president
				}
			}
		}

		// For MERGE transactions, handle president column
		if fileType == "MERGE" {
			if moveType, ok := transaction["type"].(string); ok && moveType == "minister" {
				president, ok := transaction["president"].(string)
				if !ok || president == "" {
					return nil, fmt.Errorf("president field is required for minister in MERGE transaction")
				}
				// Process old ministers array
				if old, ok := transaction["old"].(string); ok {
					// Remove the square brackets and split by comma
					old = strings.Trim(old, "[]")
					ministers := strings.Split(old, ",")

					// Process each minister name
					for i, minister := range ministers {
						// Trim quotes and spaces
						minister = strings.Trim(minister, "\" ")
						ministers[i] = minister + president
					}

					// Reconstruct the array string
					transaction["old"] = "[" + strings.Join(ministers, ",") + "]"
				}

				// Process new minister name
				if new, ok := transaction["new"].(string); ok {
					transaction["new"] = new + president
				}
			}
		}

		// For MOVE transactions, handle president column
		if fileType == "MOVE" {
			if moveType, ok := transaction["type"].(string); ok && moveType == "department" {
				// Get old parent president
				oldParentPres, ok := transaction["old_parent_pres"].(string)
				if !ok || oldParentPres == "" {
					return nil, fmt.Errorf("old_parent_pres field is required for department in MOVE transaction")
				}
				// Get new parent president
				newParentPres, ok := transaction["new_parent_pres"].(string)
				if !ok || newParentPres == "" {
					return nil, fmt.Errorf("new_parent_pres field is required for department in MOVE transaction")
				}
				// Append old_parent_pres to old_parent
				if oldParent, ok := transaction["old_parent"].(string); ok {
					transaction["old_parent"] = oldParent + oldParentPres
				}
				// Append new_parent_pres to new_parent
				if newParent, ok := transaction["new_parent"].(string); ok {
					transaction["new_parent"] = newParent + newParentPres
				}
			}
		}

		// For RENAME transactions, handle president column
		if fileType == "RENAME" {
			if renameType, ok := transaction["type"].(string); ok && renameType == "minister" {
				president, ok := transaction["president"].(string)
				if !ok || president == "" {
					return nil, fmt.Errorf("president field is required for minister in RENAME transaction")
				}
				// Append president to old name
				if old, ok := transaction["old"].(string); ok {
					transaction["old"] = old + president
				}
				// Append president to new name
				if new, ok := transaction["new"].(string); ok {
					transaction["new"] = new + president
				}
			}
		}

		// For TERMINATE transactions, handle president column
		if fileType == "TERMINATE" {
			needsPresident := false
			if parentType, ok := transaction["parent_type"].(string); ok && parentType == "minister" {
				needsPresident = true
			}
			if childType, ok := transaction["child_type"].(string); ok && childType == "minister" {
				needsPresident = true
			}

			if needsPresident {
				president, ok := transaction["president"].(string)
				if !ok || president == "" {
					return nil, fmt.Errorf("president field is required for minister in TERMINATE transaction")
				}
				// Check if parent_type is minister
				if parentType, ok := transaction["parent_type"].(string); ok && parentType == "minister" {
					if parent, ok := transaction["parent"].(string); ok {
						transaction["parent"] = parent + president
					}
				}

				// Check if child_type is minister
				if childType, ok := transaction["child_type"].(string); ok && childType == "minister" {
					if child, ok := transaction["child"].(string); ok {
						transaction["child"] = child + president
					}
				}
			}
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
