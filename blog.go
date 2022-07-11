package blog

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

const TimestampFormat = time.RFC1123

const (
	UsersTable       = "tUsers"
	BlogEntriesTable = "tEntries"
	CommentsTable    = "tComments"
)

var (
	ErrEntryNotFound       = fmt.Errorf("blog entry not found")
	ErrUnableToScan        = fmt.Errorf("could not scan entry to struct")
	ErrUnableToPrepare     = fmt.Errorf("could not prepare sql statement")
	ErrUnableToExecute     = fmt.Errorf("could not execute sql statement")
	ErrUnableToQuery       = fmt.Errorf("could not query with sql statement")
	ErrCantFindComments    = fmt.Errorf("could not find comments for post")
	ErrCantGetLastInsertID = fmt.Errorf("could not get LastInsertID")
)

func InsertEntry(db *sql.DB, title, body string, userID int64, entryType EntryType) (Entry, error) {
	timestamp := time.Now().Format(TimestampFormat)
	sql_string := "INSERT INTO " + BlogEntriesTable + "(timestamp, title, body, userid, type) VALUES (?, ?, ?, ?, ?)"
	prepared, err := db.Prepare(sql_string)
	if err != nil {
		return Entry{}, fmt.Errorf("CreateEntryAndInsert %w: %s", ErrUnableToPrepare, err)
	}
	n, err := prepared.Exec(timestamp, title, body, userID, entryType)
	if err != nil {
		return Entry{}, fmt.Errorf("CreateEntryAndInsert %w: %s", ErrUnableToExecute, err)
	}
	id, err := n.LastInsertId()
	if err != nil {
		return Entry{}, fmt.Errorf("CreateEntryAndInsert %w: %s", ErrCantGetLastInsertID, err)
	}
	return GetEntry(db, id)
}

func entrySELECTStatement() string {
	return fmt.Sprintf("SELECT %s.id, %s.timestamp, %s.title, %s.body, %s.userid, %s.type, %s.username FROM %s INNER JOIN %s ON %s.userid=%s.id WHERE %s.id=?;",
		BlogEntriesTable, BlogEntriesTable, BlogEntriesTable, BlogEntriesTable, BlogEntriesTable, BlogEntriesTable, UsersTable, BlogEntriesTable, UsersTable, BlogEntriesTable, UsersTable, BlogEntriesTable)
}

func entriesSELECTStatement() string {
	return fmt.Sprintf("SELECT %s.id, %s.timestamp, %s.title, %s.body, %s.userid, %s.type, %s.username FROM %s INNER JOIN %s ON %s.userid=%s.id ORDER BY %s.id DESC;",
		BlogEntriesTable, BlogEntriesTable, BlogEntriesTable, BlogEntriesTable, BlogEntriesTable, BlogEntriesTable, UsersTable, BlogEntriesTable, UsersTable, BlogEntriesTable, UsersTable, BlogEntriesTable)
}

func GetEntry(db *sql.DB, id int64) (Entry, error) {
	sql_string := entrySELECTStatement()
	prepared, err := db.Prepare(sql_string)
	if err != nil {
		return Entry{}, fmt.Errorf("GetEntry %w: %s", ErrUnableToPrepare, err)
	}
	row := prepared.QueryRow(id)
	var entry Entry
	err = row.Scan(&entry.ID, &entry.Timestamp, &entry.Title, &entry.Body, &entry.UserID, &entry.Type, &entry.Username)
	if err != nil {
		return Entry{}, fmt.Errorf("GetEntry %w: %s", ErrUnableToQuery, err)
	}
	comments, err := GetComments(db, entry.ID)
	if err != nil {
		log.Println("GetEntries %w", err)
	}
	entry.Comments = comments
	return entry, nil
}

func GetEntries(db *sql.DB) ([]Entry, error) {
	sql_string := entriesSELECTStatement()
	prepared, err := db.Prepare(sql_string)
	if err != nil {
		return nil, fmt.Errorf("GetEntries %w: %s", ErrUnableToPrepare, err)
	}
	rows, err := prepared.Query()
	if err != nil {
		return nil, fmt.Errorf("GetEntries %w: %s", ErrUnableToQuery, err)
	}
	defer rows.Close()
	entries := make([]Entry, 0)
	for rows.Next() {
		var entry Entry
		err = rows.Scan(&entry.ID, &entry.Timestamp, &entry.Title, &entry.Body, &entry.UserID, &entry.Type, &entry.Username)
		if err != nil {
			log.Println("GetEntries %w: %s", ErrUnableToScan, err)
			continue
		}
		comments, err := GetComments(db, entry.ID)
		if err != nil {
			log.Println("GetEntries %w", err)
		}
		entry.Comments = comments
		entries = append(entries, entry)
	}
	return entries, nil
}

func UpdateEntry(db *sql.DB, title, body string, postID int64) (Entry, error) {
	sql_string := "UPDATE " + BlogEntriesTable + " SET title=?, body=? WHERE id=?"
	prepared, err := db.Prepare(sql_string)
	if err != nil {
		return Entry{}, fmt.Errorf("UpdateEntry %w: %s", ErrUnableToPrepare, err)
	}
	_, err = prepared.Exec(title, body, postID)
	if err != nil {
		return Entry{}, fmt.Errorf("UpdateEntry %w: %s", ErrUnableToExecute, err)
	}
	return GetEntry(db, postID)
}
