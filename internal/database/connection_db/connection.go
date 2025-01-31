package connectiondb

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sajad-dev/eda-architecture/internal/app/exception"
)

// Database is a global variable holding the database connection
var Database *sql.DB

// Connection initializes a MySQL database connection
func Connection() {
	// Open a connection to MySQL without specifying a database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		os.Getenv("USERNAME_DB"), os.Getenv("PASSWORD_DB"),
		os.Getenv("IP_DB"), os.Getenv("PORT_DB")))
	exception.Log(err)

	// Retrieve database name from environment variables
	databasename := os.Getenv("DATABASE_DB")

	// Create the database if it does not exist
	_, _ = db.Exec("CREATE DATABASE " + databasename)

	// Reopen the connection, now specifying the database
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("USERNAME_DB"), os.Getenv("PASSWORD_DB"),
		os.Getenv("IP_DB"), os.Getenv("PORT_DB"), databasename))
	exception.Log(err)

	// Verify the connection to the database
	if err := db.Ping(); err != nil {
		exception.Log(err)
	}

	// Assign the initialized database connection to the global variable
	Database = db
}
