package util

import (
	"database/sql"
	"fmt"
	"github.com/prometheus/common/log"
	"time"
)

const defaultMySQLDBPort = "3306"

// InitializeMySQLDB initialize mySQL DB connection with hostname, username, password and defaullt dbname parameters.
func InitializeMySQLDB(host, port, username, password, dbname string) (*sql.DB, error) {
	if port == "" {
		port = defaultMySQLDBPort
	}

	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", username, password, host, port, dbname)

	// Retry DB connection
	var db *sql.DB
	var err error
	// Try to connect forever.
	for {
		log.Infof("Opening DB: %s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", username, host, port, dbname)

		db, err = sql.Open("mysql", dbURL)
		if err != nil {
			log.Errorf("Opening DB failed: %v", err)
		} else if err = db.Ping(); err != nil {
			// Ping to db to make sure the connection is valid.
			log.Errorf("Connecting to DB failed: %v", err)
			db.Close()
		}

		if err == nil {
			break
		}

		log.Info("Wait and will try to connect to DB again.")
		time.Sleep(5 * time.Second)
	}

	log.Debugf("Opened DB successfully.")
	return db, nil
}
