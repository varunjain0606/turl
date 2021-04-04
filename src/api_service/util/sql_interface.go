package util

import (
	"database/sql"
)

type UrlRow interface {
	Scan(...interface{}) error
}

type Sql interface {
	Begin() (*sql.Tx, error)
	Close() error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) UrlRow
}
