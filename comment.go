package blog

import "time"

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
