package migration

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/sajad-dev/eda-architecture/internal/app/exception"
	connectiondb "github.com/sajad-dev/eda-architecture/internal/database/connection_db"
)

// MigrateFunc defines a function type that returns a slice of strings for migration operations
type MigrateFunc func() []string

// Migrate represents a single migration entry
type Migrate struct {
	Table     string      // Name of the database table
	Function  MigrateFunc // Migration function returning column definitions
}

// Migration holds a list of migration entries
type Migration struct {
	MigrateList []*Migrate // List of migrations to be processed
}

// MigrationInterface defines the interface for migration operations
type MigrationInterface interface {
	AppendToMigrationList()   // Append migrations to the migration list
	CheckTable()              // Check and process table changes
	AddTable(strSlice string, table string) // Add a new table or column
	checkDeletedTable()       // Remove deleted tables from the database
	UpdateTable(table string, column_name_old string, parametr string) // Update table column properties
	HandelUpdate(field string, fieldtype string, collection sql.NullString,
		null string, defult sql.NullString, key string, extera string, privileges string, comment string, function MigrateFunc, table string, x int) // Handle updates to table columns
}

// AppendToMigrationList adds migration entries to the migration list
func (migrate *Migration) AppendToMigrationList() {
	migrate.MigrateList = ArrMigrations
}

// CheckTable validates the database schema against the migration list
func (migrate *Migration) CheckTable() {
	database := os.Getenv("DATABASE_DB") // Get database name from environment
	migrate.checkDeletedTable() // Remove obsolete tables

	sqlqu := fmt.Sprintf(` 
SELECT table_name
FROM information_schema.tables 
WHERE table_schema = '%s';
 `, database)

	row, err := connectiondb.Database.Query(sqlqu)
	exception.Log(err)

	for row.Next() {
		var tb string
		err := row.Scan(&tb)
		exception.Log(err)
	}

	// Process each migration entry
	for _, m := range migrate.MigrateList {
		sql_qu := fmt.Sprintf("SHOW FULL COLUMNS FROM %s", m.Table)
		qu, err := connectiondb.Database.Query(sql_qu)
		exception.Log(err)

		x := 0
		for qu.Next() {
			var (
				field, fieldtype, null, key, extra, privileges, comment string
				collection, defult                                      sql.NullString
			)
			err = qu.Scan(&field, &fieldtype, &collection, &null, &key, &defult, &extra, &privileges, &comment)
			exception.Log(err)

			migrate.HandelUpdate(field, fieldtype, collection, null, defult, key, extra, privileges, comment, m.Function, m.Table, x)
			x++
		}

		// Add missing columns
		strSlice := m.Function()
		if len(strSlice) > x {
			for i := x; i < len(strSlice); i++ {
				migrate.AddTable(strSlice[i], m.Table)
			}
		}
	}
}

// HandelUpdate processes updates to table columns
func (migrate *Migration) HandelUpdate(field string, fieldtype string, collection sql.NullString,
	null string, defult sql.NullString, key string, extera string, privileges string, comment string, function MigrateFunc, table string, x int) {
	
	strSlice := function()

	// Drop column if not in migration list
	if x >= len(strSlice) {
		sql_qu := fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table, field)
		_, err := connectiondb.Database.Query(sql_qu)
		exception.Log(err)
		return
	}

	// Handle unique key constraints
	if key == "UNI" && strings.Contains(strSlice[x], "UNIQUE") {
		sql_qu := fmt.Sprintf("SELECT DISTINCT COLUMN_NAME, INDEX_NAME FROM information_schema.statistics WHERE table_name = '%s'", table)
		qu, err := connectiondb.Database.Query(sql_qu)
		exception.Log(err)

		unilist := []string{}
		for qu.Next() {
			var column, key string
			qu.Scan(&column, &key)
			if column == field {
				unilist = append(unilist, key)
			}
		}
		for _, v := range unilist {
			sql_qu := fmt.Sprintf("ALTER TABLE %s DROP INDEX %s", table, v)
			connectiondb.Database.Query(sql_qu)
			exception.Log(err)
		}
	}

	// Remove primary key constraints
	if key == "PRI" {
		strSlice[x] = strings.ReplaceAll(strSlice[x], "PRIMARY KEY", "")
	}
	migrate.UpdateTable(table, field, strSlice[x])
}

// checkDeletedTable removes tables that are no longer in the migration list
func (migrate *Migration) checkDeletedTable() {
	row, err := connectiondb.Database.Query("SHOW TABLES")
	exception.Log(err)

	tableArr := []string{}
	for _, v := range migrate.MigrateList {
		tableArr = append(tableArr, v.Table)
	}

	for row.Next() {
		var table string
		row.Scan(&table)

		if !slices.Contains(tableArr, table) {
			connectiondb.Database.Query(fmt.Sprintf("DROP TABLE %s", table))
		}
	}
}

// Handel initializes the migration process
func Handel() {
	var migra MigrationInterface
	migra = &Migration{}
	migra.AppendToMigrationList()
	migra.CheckTable()
}

// AddTable adds a new column or foreign key constraint
func (migrate *Migration) AddTable(strSlice string, table string) {
	if strings.Contains(strSlice, "FOREIGN") {
		re := regexp.MustCompile(`\((.*?)\)`)
		matches := re.FindStringSubmatch(strSlice)

		// Add foreign key column
		sql_qu := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s INT", table, matches[1])
		_, err := connectiondb.Database.Query(sql_qu)
		exception.Log(err)

		// Add foreign key constraint
		sql_qu = fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s %s", table, matches[1], strSlice)
		_, err = connectiondb.Database.Query(sql_qu)
		exception.Log(err)
		return
	}

	// Add new column
	sql_qu := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", table, strSlice)
	_, err := connectiondb.Database.Query(sql_qu)
	exception.Log(err)
}

// UpdateTable modifies an existing column definition
func (migrate *Migration) UpdateTable(table string, column_name_old string, parametr string) {
	sql_qu := fmt.Sprintf("ALTER TABLE %s CHANGE %s %s ", table, column_name_old, parametr)

	// Handle foreign key constraints
	if strings.Contains(sql_qu, "FOREIGN") {
		sql_du2 := fmt.Sprintf("ALTER TABLE %s DROP FOREIGN KEY %s", table, table+"_"+column_name_old+"_fg")
		connectiondb.Database.Query(sql_du2)

		sql_qu = fmt.Sprintf("ALTER TABLE %s CHANGE COLUMN %s %s ", table, column_name_old, parametr)
		sql_qu = strings.Replace(sql_qu, "CONSTRAINT", "ADD CONSTRAINT ", 1)
	}

	_, err := connectiondb.Database.Query(sql_qu)
	exception.Log(err)
}
