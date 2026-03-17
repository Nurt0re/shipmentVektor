package pkg

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresDB(conn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	db.Ping()
	return db, nil
}
