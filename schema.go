package blog

import (
	"database/sql"
	"fmt"

	"github.com/sharkpick/authentication"
)

func GenerateBlogTable(db *sql.DB) error {
	sql_string := `CREATE TABLE IF NOT EXISTS "` + BlogEntriesTable + `"(
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT,
		"body" TEXT,
		"timestamp" TEXT NOT NULL,
		"updated" TEXT,
		"userid" integer NOT NULL, type int,
		FOREIGN KEY(userid) REFERENCES "` + authentication.UsersTable + `"(id));`
	prepared, err := db.Prepare(sql_string)
	if err != nil {
		return fmt.Errorf("GenerateBlogTable %w: %s", ErrUnableToPrepare, err)
	}
	_, err = prepared.Exec()
	if err != nil {
		return fmt.Errorf("GenerateBlogTable: %w: %s", ErrUnableToExecute, err)
	}
	return nil
}
