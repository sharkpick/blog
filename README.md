# blog is the basic stuff you need for a blog using sqlite3 and Go.

# Requirements: 
Go 1.17+ 
sqlite3

# Notes:
- you'll need to configure a table for users, then entries (and use userID as a foreign key in entries to link them together). The Entry struct is a good reference, but keep in mind that the username string is identified by the DB Query and doesn't need to a column in your table.

use NewBlog() with a *sql.DB to start up the Blog struct. 

use GetEntries() to get a []Entry of all of your entries. This is useful for presentation, you can use html templates and that []Entry to present all of your posts.

DeleteEntry deletes an entry by ID.

InsertEntry takes a title, body and userID and generates its own timestamp. 
