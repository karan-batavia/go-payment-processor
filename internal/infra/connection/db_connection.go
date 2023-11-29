package connection

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func DBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s", dsn))
	if err != nil {
		return nil, err
	}

	return db, nil
}
