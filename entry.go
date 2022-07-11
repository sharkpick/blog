package blog

type EntryType int64

const (
	Blog EntryType = iota
	News
)

type Entry struct {
	ID                int64
	Title             string
	Body              string
	Timestamp         string
	ModifiedTimestamp string
	UserID            int64
	Username          string
	Type              EntryType
	Comments          []Comment
}
