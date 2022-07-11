package blog

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Comment struct {
	ID        int64
	Body      string
	Name      string
	Timestamp string
	PostID    int64
}

func NewComment(postID int64, body, name string) Comment {
	return Comment{
		PostID:    postID,
		Body:      body,
		Name:      name,
		Timestamp: time.Now().Format(TimestampFormat),
	}
}

func GetComments(db *sql.DB, id int64) ([]Comment, error) {
	results := make([]Comment, 0)
	sql_string := "SELECT id, body, name, timestamp, postid FROM " + CommentsTable + " WHERE postid=?"
	prepared, err := db.Prepare(sql_string)
	if err != nil {
		return results, fmt.Errorf("GetComments %w: %s", ErrCantFindComments, err)
	}
	rows, err := prepared.Query(id)
	if err != nil {
		return results, fmt.Errorf("GetComments %w: %s", ErrUnableToQuery, err)
	}
	defer rows.Close()
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.ID, &comment.Body, &comment.Name, &comment.Timestamp, &comment.PostID)
		if err != nil {
			log.Printf("GetComments %w: %s\n", ErrUnableToScan, err)
			continue
		}
		results = append(results, comment)
	}
	return results, nil
}
