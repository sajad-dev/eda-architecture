package checksystem

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/sajad-dev/eda-architecture/internal/database/migration"
	testutils "github.com/sajad-dev/eda-architecture/internal/test_utils"
)

// TestMigrationTables verifies that all expected tables exist in the database
func TestMigrationTables(t *testing.T) {
	var db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("USERNAME_DB"), os.Getenv("PASSWORD_DB"),
		os.Getenv("IP_DB"), os.Getenv("PORT_DB"), os.Getenv("DATABASE_DB")))

	if err == nil {
		// Query the list of tables from the database
		qu, err := db.Query("SHOW TABLES")
		if err != nil {
			return
		}
		x := 0

		// Retrieve expected migration table list
		table, _ := testutils.MiggarionListAppend()

		// Check if all tables exist
		for qu.Next() {
			x++
			var name = ""
			qu.Scan(&name)

			// If a table exists in the database but not in migration, report an error
			if !slices.Contains(table, name) {
				t.Fatalf("Database %s not deleted", name)
			}
		}

		// Ensure the number of tables in the database matches the expected number
		if len(table) != x {
			t.Fatal("You have problem in tables")
		}
	}
}

// TestMigrationTablesParams verifies that all table columns match expected definitions
func TestMigrationTablesParams(t *testing.T) {
	var db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("USERNAME_DB"), os.Getenv("PASSWORD_DB"),
		os.Getenv("IP_DB"), os.Getenv("PORT_DB"), os.Getenv("DATABASE_DB")))

	if err != nil {
		mig := migration.Migration{}
		mig.AppendToMigrationList()

		// Iterate over each migration table definition
		for _, v := range mig.MigrateList {
			rqfunc := v.Function()

			// Ensure that the migration function is not empty
			if len(rqfunc) == 0 {
				t.Fatal("Your migration not be empty")
			}

			// Query column details of the table
			rqdb, err := db.Query(fmt.Sprintf("SHOW FULL COLUMNS FROM %s", v.Table))
			if err != nil {
				t.Fatal(err.Error())
			}

			arr := []string{}

			// Scan the database schema for comparison
			for rqdb.Next() {
				var name, tp, null, extra, key, privileges, comment string
				var collation, df sql.NullString

				// Extract column attributes from the database
				if err := rqdb.Scan(&name, &tp, &collation, &null, &key, &df, &extra, &privileges, &comment); err != nil {
					t.Fatal("Error scanning row:", err)
				}

				// Normalize NULL constraints
				if null == "YES" {
					null = "NULL"
				} else {
					null = "NOT NULL"
				}

				// Normalize DEFAULT values
				dfStr := "DEFAULT ''"
				if df.Valid {
					dfStr = fmt.Sprintf("DEFAULT '%s'", df.String)
				}

				// Normalize PRIMARY and UNIQUE keys
				if key == "PRI" {
					key = "PRIMARY"
				}
				if key == "UNI" {
					key = "UNIQUE"
				}

				// Append extra constraints
				if extra != "" {
					extra = fmt.Sprintf(" %s", extra)
				}

				// Construct column definition string
				str := fmt.Sprintf("%s %s %s %s%s %s",
					name, strings.ToUpper(tp), null, dfStr, strings.ToUpper(extra), key)

				// Handle special case for "id" column
				if name == "id" {
					str = str + " KEY"
					str = strings.ReplaceAll(str, "INT(11)", "INT")
					str = strings.ReplaceAll(str, "DEFAULT '' ", "")
				}

				// Add the constructed column definition to the array
				arr = append(arr, str)
			}

			// Compare the extracted database schema with the expected migration schema
			for ind, val := range arr {
				if !strings.Contains(rqfunc[ind], val) {
					t.Fatal("Migration create not like migration exist")
				}
			}
		}
	}
}
