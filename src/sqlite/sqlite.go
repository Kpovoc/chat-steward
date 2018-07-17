package sqlite

import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"log"
	"fmt"
)

var db *sql.DB

func InitDB(userDataDir string) {
	tempDB, err := sql.Open("sqlite3", userDataDir + "/main.db")
	if err != nil {
		log.Fatal(err)
	}
	db = tempDB
}

func GetInstance() *sql.DB {
	return db
}

/**
  Checks if table exists. If not, then creates a new table with provided tableName
 */
func InitTable(tableName string, columnNames string) error {
	queryForTable := fmt.Sprintf(
		"SELECT name " +
			"FROM sqlite_master " +
			"WHERE type='table' " +
			"AND name='%s';", tableName)
	rows, err := db.Query(queryForTable)
	if err != nil {
		return err
	}
	defer rows.Close()

	// should only be one result. If rows.Next, then table exists
	for rows.Next() {
		return nil
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	// Was no result, so need to create table
	createTable := fmt.Sprintf(
		"CREATE TABLE %s (%s);", tableName, columnNames)
	_, err = db.Exec(createTable)
	if err != nil {
		return err
	}

	return nil
}