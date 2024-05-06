// Package postgresql - included functions for init SQL connections
package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewConn - function for opening new SQL connection
func NewConn(dsn string) (*sql.DB, error) {
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("DB open: %w", err)
	}
	return conn, nil
}
