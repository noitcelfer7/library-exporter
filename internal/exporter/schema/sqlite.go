package schema

import "database/sql"

type Record struct {
	AuthorFirstName string
	AuthorLastName  string

	BookIsbn  string
	BookTitle string

	GenreTitle string

	IssueDate       string
	IssuePeriod     string
	IssueReturnDate sql.NullString

	ReaderFirstName   string
	ReaderLastName    string
	ReaderPhoneNumber string
}
