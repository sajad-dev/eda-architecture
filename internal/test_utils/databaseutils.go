package testutils

import (
	"bytes"
	"os/exec"
	"strings"
)

// Mig is an empty struct used for defining migration-related utilities
type Mig struct{}

// MiggarionListAppend retrieves a list of migration tables and their corresponding functions
func MiggarionListAppend() ([]string, []string) {
	// Execute the Go documentation command to retrieve migration interface details
	ou := exec.Command("go", "doc", "migration", "MigrationInterface")

	// Capture the output of the command
	output, _ := ou.Output()
	
	// Define slices to store table names and function names
	res_table := []string{}
	res_func := []string{}

	// Split output by new lines
	slice := bytes.Split(output, []byte("\n"))

	// Iterate through the extracted lines
	for _, s := range slice {
		// Convert byte slice to string and remove "func" keyword
		text := string(strings.ReplaceAll(string(s), "func", ""))
		index := strings.Index(text, "(") // Find the index of "("
		migration_index := strings.Index(text, "Migration") // Find the index of "Migration"

		// If both indices are valid, extract table and function names
		if index > 0 && migration_index > 0 {
			table_name := strings.ToLower(text[:index][:migration_index]) // Extract table name
			res_table = append(res_table, strings.TrimSpace(table_name)) // Store table name
			res_func = append(res_func, strings.TrimSpace(string(text[:index]))) // Store function name
		}
	}

	// Return slices containing table names and function names
	return res_table, res_func
}
