package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"

	"library_exporter/internal/exporter/schema"
)

type Database struct {
	*sql.DB
}

func NewDB(connStr string) (*Database, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, fmt.Errorf("sql.Open Error: %w", err)
	}

	err = db.Ping()

	if err != nil {
		return nil, fmt.Errorf("db.Ping Error: %w", err)
	}

	return &Database{db}, nil
}

func (db *Database) InsertRecord(record schema.Record) error {
	tx, err := db.Begin()

	if err != nil {
		return fmt.Errorf("db.Begin Error: %w", err)
	}

	defer tx.Rollback()

	var genreID string

	err = tx.QueryRow(`INSERT INTO genres (title) VALUES ($1) ON CONFLICT DO NOTHING RETURNING uuid`, record.GenreTitle).Scan(&genreID)

	if err == sql.ErrNoRows {
		err = tx.QueryRow(`SELECT uuid FROM genres WHERE title = $1`, record.GenreTitle).Scan(&genreID)
	}

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to insert genres: %w", err)
	}

	var readerID string

	err = tx.QueryRow(`INSERT INTO readers (first_name, last_name, phone_number) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING uuid`,
		record.ReaderFirstName, record.ReaderLastName, record.ReaderPhoneNumber).Scan(&readerID)

	if err == sql.ErrNoRows {
		err = tx.QueryRow(`SELECT uuid FROM readers WHERE first_name = $1 AND last_name = $2`, record.ReaderFirstName, record.ReaderLastName).Scan(&readerID)
	}

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to insert readers: %w", err)
	}

	var authorID string

	err = tx.QueryRow(`INSERT INTO authors (first_name, last_name) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING uuid`,
		record.AuthorFirstName, record.AuthorLastName).Scan(&authorID)

	if err == sql.ErrNoRows {
		err = tx.QueryRow(`SELECT uuid FROM authors WHERE first_name = $1 AND last_name = $2`, record.AuthorFirstName, record.AuthorLastName).Scan(&authorID)
	}

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to insert authors: %w", err)
	}

	var bookID string

	err = tx.QueryRow(`INSERT INTO books (isbn, title) VALUES ($1, $2) ON CONFLICT DO NOTHING RETURNING uuid`,
		record.BookIsbn, record.BookTitle).Scan(&bookID)

	if err == sql.ErrNoRows {
		err = tx.QueryRow(`SELECT uuid FROM books WHERE title = $1`, record.BookTitle).Scan(&bookID)
	}

	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to insert books: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO authors_to_books (author_uuid, book_uuid) VALUES ($1, $2) ON CONFLICT DO NOTHING`, authorID, bookID)

	if err != nil {
		return fmt.Errorf("failed to insert into authors_to_books: %w", err)
	}

	_, err = tx.Exec(`INSERT INTO books_to_genres (book_uuid, genre_uuid) VALUES ($1, $2) ON CONFLICT DO NOTHING`, bookID, genreID)

	if err != nil {
		return fmt.Errorf("failed to insert into books_to_genres: %w", err)
	}

	issueReturnDate := pq.NullTime{}

	if record.IssueReturnDate.Valid {
		date, err := time.Parse("0001-01-01 00:00:00 +0000 UTC", issueReturnDate.Time.String())

		if err != nil {
			return fmt.Errorf("failed to insert issues: %w", err)
		}

		issueReturnDate = pq.NullTime{Time: date, Valid: true}
	}

	_, err = tx.Exec(`INSERT INTO issues ("book_uuid", "reader_uuid", "issue_date", "period", "return_date") VALUES ($1, $2, $3, $4, $5) ON CONFLICT DO NOTHING`,
		bookID, readerID, record.IssueDate, record.IssuePeriod, issueReturnDate)

	if err != nil {
		return fmt.Errorf("failed to insert issues: %w", err)
	}

	err = tx.Commit()

	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
