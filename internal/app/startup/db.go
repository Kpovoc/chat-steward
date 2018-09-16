package startup

import (
	"gitlab.com/Kpovoc/chat-steward/internal/app/sqlite"
	"log"
)

func InitializeDatabase(userDataDir string) {
	sqlite.InitDB(userDataDir)

	// Init Plugins
	initTitlesPlugin()
}

// Plugin Table Inits
func initTitlesPlugin() {
	// See if table exists, if not, create it
	tableName := "titles"
	columns :=
		"id INTEGER NOT NULL PRIMARY KEY, " +
			"title TEXT, " +
			"author TEXT NOT NULL, " +
			"created_on DATETIME, " +
			"total_votes INTEGER"
	err := sqlite.InitTable(tableName, columns)
	if err != nil {
		log.Fatal(err)
	}

	// Will only hold 1 name, the current show in case of crash
	tableName = "show"
	columns = "name TEXT"
	err = sqlite.InitTable(tableName, columns)
	if err != nil {
		log.Fatal(err)
	}
}
