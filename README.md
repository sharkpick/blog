# blog


# Requirements: 
authentication (https://github.com/sharkpick/authentication)
Go 1.17+ 
sqlite3

## About:
package blog offers simple functions to manage a blog. requires authentication package (https://github.com/sharkpick/authentication) in order to manage entry ownership.

entries may be classified as a Blog or News entry, which can make managing a news feed easy by allowing you to send more important blog entries to your news feed, or use both separately.

## Usage:
simply open your sqlite3 database and send it along to the functions to manage the blog. run the GenerateBlogTable() function at least once to build the table so your entries have somewhere to go.

use NewBlog() with a *sql.DB to start up the Blog struct. 

use GetEntries() to get a []Entry of all of your entries. This is useful for presentation, you can use html templates and that []Entry to present all of your posts.

DeleteEntry deletes an entry by ID.

InsertEntry takes a title, body and userID and generates its own timestamp. 
