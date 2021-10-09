package blog

import (
	"database/sql"
	"log"
	"testing"
)

func loadDB() *sql.DB {
	dbFile := "test.db"
	database, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalln("Fatal Error: Unable to read", dbFile)
	}
	return database

}
func TestNewBlog(t *testing.T) {
	myBlog := NewBlog(loadDB())
	if nil == myBlog.database {
		t.Fatalf("Error: database should not be nil")
	}
}

func TestInsertEntry(t *testing.T) {
	myBlog := NewBlog(loadDB())
	defer myBlog.CleanupBlog()
	myBlog.InsertEntry("Test", "TestBody", 1)
	myBlog.InsertEntry("Test2", "Test2Body", 1)
	entries := myBlog.GetEntries()
	if len(entries) != 2 {
		t.Fatalf("Error: wanted 2, got %d", len(entries))
	}
}

func TestDeleteEntry(t *testing.T) {
	myBlog := NewBlog(loadDB())
	defer myBlog.CleanupBlog()
	currentEntries := myBlog.GetEntries()
	for i := 0; i < len(currentEntries); i++ {
		myBlog.DeleteEntry(currentEntries[i].ID)
	}
	if len(myBlog.GetEntries()) != 0 {
		t.Fatalf("Error: myBlog shoud be empty but has %d entries", len(myBlog.GetEntries()))
	}
}
