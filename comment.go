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

func InsertComment(db *sql.DB, body, name string, postid int64) (Comment, error) {
	timestamp := time.Now().Format(TimestampFormat)
	sql_string := "INSERT INTO " + CommentsTable + "(body, name, timestamp, postid) VALUES (?, ?, ?, ?)"
	prepared, err := db.Prepare(sql_string)
	if err != nil {
		return Comment{}, fmt.Errorf("InsertComment %w: %s", ErrUnableToPrepare, err)
	}
	n, err := prepared.Exec(body, name, timestamp, postid)
	if err != nil {
		return Comment{}, fmt.Errorf("InsertComment %w: %s", ErrUnableToExecute, err)
	}
	id, err := n.LastInsertId()
	if err != nil {
		return Comment{}, fmt.Errorf("InsertComment %w: %s", ErrCantGetLastInsertID, err)
	}
	return GetCommentByID(db, id)
}

func GetCommentByID(db *sql.DB, id int64) (Comment, error) {
	sql_string := "SELECT id, body, name, timestamp, postid FROM " + CommentsTable + " WHERE id=?"
	prepared, err := db.Prepare(sql_string)
	if err != nil {
		return Comment{}, fmt.Errorf("GetCommentByID %w: %s", ErrUnableToPrepare, err)
	}
	row := prepared.QueryRow(id)
	var comment Comment
	err = row.Scan(&comment.ID, &comment.Body, &comment.Name, &comment.Timestamp, &comment.PostID)
	if err != nil {
		return Comment{}, fmt.Errorf("GetCommentByID %w: %s", ErrUnableToScan, err)
	}
	return comment, nil
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
