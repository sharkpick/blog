package blog

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

const TimestampFormat = time.RFC1123

type Blog struct {
	database *sql.DB
	mutex    sync.Mutex
}

func (b *Blog) GetDatabase() *sql.DB {
	return b.database
}

func NewBlog() *Blog {
	b := Blog{}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if file, err := os.OpenFile("./blog-database.db", os.O_CREATE|os.O_RDWR, 0777); err != nil {
		log.Println("Error in NewBlog:", err)
	} else {
		file.Close()
		log.Println("./blog-database exists")
	}
	if database, err := sql.Open("sqlite3", "./blog-database.db"); err != nil {
		log.Fatalln("Fatal Error: Unable to read ./blog-database.db")
	} else {
		b.database = database
	}
	b.createTable()
	return &b

}

func (b *Blog) CleanupBlog() {
	b.database.Close()
}

func (b *Blog) createTable() {
	sqlStatement := `CREATE TABLE entries (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"timestamp" TEXT,
		"title" TEXT,
		"body" TEXT,
		"userid" integer,
		FOREIGN KEY(userid) REFERENCES users(id)
	);`
	if statement, err := b.database.Prepare(sqlStatement); err != nil {
		log.Println("Error in b.createTable():", err) // likely means we're re-launching but log it anyway
	} else {
		if _, err = statement.Exec(); err != nil {
			log.Fatalln("Failed to make table users: ", err)
		}
	}
}

func (b *Blog) InsertEntry(title, body string) {
	log.Println("inserting new blog entry")
	timestamp := fmt.Sprintf("%v", time.Now().Format(TimestampFormat))
	sqlStatement := `INSERT INTO entries(timestamp, title, body) VALUES (?, ?, ?)`
	if statement, err := b.database.Prepare(sqlStatement); err != nil {
		log.Fatalln("Fatal Error in b.InsertEntry():", err)
	} else {
		if _, err := statement.Exec(timestamp, title, body); err != nil {
			log.Fatalln("Fatal Error in b.InsertEntry():", err)
		}
	}
	log.Println("finished inserting blog entry", title)
}

func (b *Blog) GetEntries() []Entry {
	entries := make([]Entry, 0)
	if row, err := b.database.Query("SELECT * FROM entries ORDER BY id DESC"); err != nil {
		log.Fatalln("Fatal Error in b.GetEntries():", err)
	} else {
		defer row.Close()
		for row.Next() {
			var entry Entry
			var timestamp string
			row.Scan(&entry.ID, &timestamp, &entry.Title, &entry.Body)
			if t, err := time.Parse(TimestampFormat, timestamp); err != nil {
				log.Fatalln("Fatal Error - can't parse timestamp", err)
			} else {
				entry.Timestamp = t.Format(TimestampFormat)
			}
			entries = append(entries, entry)
		}
	}
	return entries
}
