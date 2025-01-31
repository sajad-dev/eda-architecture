package model

import (
	"fmt"
	"strings"

	"github.com/sajad-dev/eda-architecture/internal/app/exception"
	connectiondb "github.com/sajad-dev/eda-architecture/internal/database/connection_db"
)

// Sqltype represents a custom SQL type
type Sqltype string

// Where_st represents a condition in a SQL WHERE clause
type Where_st struct {
	Key      string // Column name
	Value    string // Value to match
	After    string // Additional condition
	Operator string // SQL operator (e.g., '=', '<', '>')
}

// GetOutput defines a type for query result sets
type GetOutput []map[string]string

// getStar returns a selection clause for SQL queries
func getStar(selection []string) string {
	if len(selection) == 0 {
		return fmt.Sprintf("*") // Select all columns if none are specified
	}
	return fmt.Sprintf("%s", strings.Join(selection, ",")) // Join selected columns
}

// getWhere constructs the WHERE clause of a SQL query
func getWhere(where []Where_st) string {
	str := ""

	for _, val := range where {
		str += fmt.Sprintf(" %s %s ? %s ", val.Key, val.Operator, val.After)
	}
	return str
}

// getOrder constructs an ORDER BY clause
func getOrder(by string, asc bool) string {
	str := fmt.Sprintf("ORDER BY %s ", by)
	if asc {
		str += "ASC"
	} else {
		str += "DESC"
	}
	return str
}

// Where generates a WHERE clause from a slice of conditions
func Where(where []Where_st) string {
	return getWhere(where)
}

// Selection constructs a SELECT SQL query
func Selection(selection []string, table string) string {
	return fmt.Sprintf("SELECT %s FROM %s ", getStar(selection), table)
}

// Get executes a SELECT query with optional filtering and ordering
func Get(selection []string, table string, where []Where_st, order string, asc bool) GetOutput {
	sql_qu := Selection(selection, table)

	if len(where) != 0 {
		sql_qu = sql_qu + " WHERE " + Where(where)
	}

	if order != "" {
		sql_qu += getOrder(order, asc)
	}
	arr := []any{}

	for _, val := range where {
		arr = append(arr, val.Value)
	}
	
	values := make([]interface{}, len(selection))
	query, err := connectiondb.Database.Query(sql_qu, arr...)
	
	exception.Log(err)
	
	var outputList GetOutput
	valuePtrs := make([]interface{}, len(selection))

	for i := range values {
		valuePtrs[i] = &values[i]
	}
	for query.Next() {
		query.Scan(valuePtrs...)
		var arr = map[string]string{}
		for i, _ := range selection {
			switch v := values[i].(type) {
			case []byte:
				arr[selection[i]] = string(v)
			case int64:
				arr[selection[i]] = fmt.Sprintf("%d", v)
			case int:
				arr[selection[i]] = fmt.Sprintf("%d", v)
			case string:
				arr[selection[i]] = v
			default:
				arr[selection[i]] = fmt.Sprintf("%v", v)
			}
		}
		outputList = append(outputList, arr)
	}
	return outputList
}

// Insert adds a new record to the specified table
func Insert(data map[string]string, table string) {
	value := []any{}
	key := []string{}
	for ke, val := range data {
		value = append(value, val)
		key = append(key, ke)
	}
	value_str := ""
	for i := 0; i < len(value); i++ {
		value_str += "?"
		if i != len(value)-1 {
			value_str += ","
		}
	}
	key_str := strings.Join(key, ",")
	sql_qu := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, key_str, value_str)
	_, err := connectiondb.Database.Query(sql_qu, value...)
	exception.Log(err)
}

// Update modifies an existing record in the specified table
func Update(data map[string]string, table string, where [2]string) {
	value := []string{}
	arr := []any{}
	for ke, val := range data {
		value = append(value, fmt.Sprintf("%s= ?", ke))
		arr = append(arr, val)
	}
	value_str := strings.Join(value, " , ")
	sql_qu := fmt.Sprintf("UPDATE %s SET %s WHERE %s = %s", table, value_str, where[0], where[1])

	_, err := connectiondb.Database.Query(sql_qu, arr...)
	exception.Log(err)
}

// Delete removes a record from the specified table
func Delete(table string, where [2]string) {
	sql_qu := fmt.Sprintf("DELETE FROM %s WHERE %s = ?", table, where[0])

	_, err := connectiondb.Database.Query(sql_qu, where[1])
	exception.Log(err)
}
