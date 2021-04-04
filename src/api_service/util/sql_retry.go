package util

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/prometheus/common/log"
	"strings"
	"time"
)

const (
	retrySqlMaxCount    = 5
	retrySqlDuration    = 200
	wsrepStatusNotReady = 1047
	sqlTimeoutExceeded  = 1205
	sqlDeadlockRetry    = 1213
)

type RetrySql struct {
	*sql.DB
}

type Row struct {
	rs *RetrySql
	*sql.Row
	query string
	args  []interface{}
}

func NewRetrySql(sql *sql.DB) *RetrySql {
	return &RetrySql{sql}
}

func (rs *RetrySql) Exec(query string, args ...interface{}) (sql.Result, error) {
	tokens := strings.Fields(query)
	var res sql.Result
	var err error
	if len(tokens) > 0 && len(tokens[0]) > 0 {
		if strings.ToUpper(tokens[0]) == "UPDATE" {
			log.Debugf("In RetrySql.Exec(%s)", tokens[0])

			retryCounter := retrySqlMaxCount
			for retryCounter > 0 {
				res, err = rs.DB.Exec(query, args...)
				if err == nil {
					return res, err
				}

				retryCounter--
				log.Errorf("Retrying sql(%d): %v", retryCounter, err)
				time.Sleep(retrySqlDuration * time.Millisecond)
			}

			return res, err
		}
	}

	return rs.DB.Exec(query, args...)
}

func (rs *RetrySql) QueryRow(query string, args ...interface{}) UrlRow {
	return &Row{rs, rs.DB.QueryRow(query, args...),
		query,
		args}
}

func (rsRow *Row) Scan(dest ...interface{}) error {
	err := rsRow.Row.Scan(dest...)
	if err == nil {
		return nil
	}

	sqlErr, ok := err.(*mysql.MySQLError)

	// Check for the node  wsrep not ready status
	// in the case of network glitch we are mitigating the issue.
	if ok && (sqlErr.Number == wsrepStatusNotReady || sqlErr.Number == sqlTimeoutExceeded || sqlErr.Number == sqlDeadlockRetry) {
		log.Errorf("Encounter error for retry: %v", err)
		for i := 0; i < retrySqlMaxCount; i++ {
			err = rsRow.rs.DB.QueryRow(rsRow.query, rsRow.args...).
				Scan(dest...)
			if err == nil {
				break
			}

			log.Errorf("Wsrep status not ready retrying:(%v)", err)
			time.Sleep(retrySqlDuration * time.Millisecond)
		}
	}

	return err
}

