package blog

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const TimestampFormat = time.RFC1123

type Blog struct {
	database *sql.DB
}

func (b *Blog) GetDatabase() *sql.DB {
	return b.database
}

func NewBlog(file ...string) *Blog {
	dbFile := func() string {
		if len(file) == 0 {
			return "./blog-database.db"
		} else {
			return file[0]
		}
	}()
	b := Blog{}
	if database, err := sql.Open("sqlite3", dbFile); err != nil {
		log.Fatalln("Fatal Error: Unable to read", dbFile)
	} else {
		b.database = database
	}
	return &b

}

func (b *Blog) CleanupBlog() {
	b.database.Close()
}

func (b *Blog) InsertEntry(title, body string, userID int) {
	timestamp := fmt.Sprintf("%v", time.Now().Format(TimestampFormat))
	sqlStatement := `INSERT INTO entries(timestamp, title, body, userID) VALUES (?, ?, ?, ?)`
	if statement, err := b.database.Prepare(sqlStatement); err != nil {
		log.Fatalln("Fatal Error in b.InsertEntry():", err)
	} else {
		if _, err := statement.Exec(timestamp, title, body, userID); err != nil {
			log.Fatalln("Fatal Error in b.InsertEntry():", err)
		}
	}
	log.Println("finished inserting blog entry", title)
}

func (b *Blog) DeleteEntry(id int) {
	sqlStatement := `DELETE FROM entries WHERE id=?`
	if statement, err := b.database.Prepare(sqlStatement); err != nil {
		log.Fatalln("Fatal Error in b.DeleteEntry()", err)
	} else {
		if _, err := statement.Exec(id); err != nil {
			log.Fatalln("Fatal Error executing b.DeleteEntry()")
		}
	}
	log.Println("Finished deleting blog entry", id)
}

func (b *Blog) GetEntries() []Entry {
	entries := make([]Entry, 0)
	if row, err := b.database.Query(`SELECT entries.id, entries.timestamp, 
	entries.title, entries.body, entries.userid, 
	users.username FROM entries INNER JOIN users 
	ON entries.userid=users.id ORDER BY entries.id DESC;`); err != nil {
		log.Fatalln("Fatal Error in b.GetEntries():", err)
	} else {
		defer row.Close()
		for row.Next() {
			var entry Entry
			row.Scan(&entry.ID, &entry.Timestamp, &entry.Title, &entry.Body, &entry.UserID, &entry.Username)
			entry.Username = strings.Split(entry.Username, "@")[0]
			entries = append(entries, entry)
		}
	}
	return entries
}
