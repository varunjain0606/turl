package database

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/prometheus/common/log"
	"os"
	"time"
	"api_service/util"
)

type SqlDatabase struct {
	sql util.Sql
}

func (sqlDb *SqlDatabase) GetSql() util.Sql {
	return sqlDb.sql
}

func New() (*SqlDatabase, error) {
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbusername := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sqlDb, err := util.InitializeMySQLDB(dbhost, dbport, dbusername, dbpassword, dbname)
	if err != nil {
		log.Errorf("Error returned is : %v", err)
		return nil, err
	}
	return &SqlDatabase{sql: util.NewRetrySql(sqlDb)}, nil
}

func (sqlDb *SqlDatabase) InsertIntoDB(key string, value string, date int64) error{
	t := time.Unix(date, 0)
	err := sqlDb.sql.QueryRow("insert into urlMap (original_url, shorten_url, expiry_date, created_at)values (?,?,?, now())", key, value, t)
	if err != nil {
		return errors.New("Error inserting into database")
	}
	return nil
}

func (sqlDb *SqlDatabase) FetchOriginalUrl(key string) error {
	var originalURL string

	query := "SELECT original_url from urlMap where shorten_url=?"

	row := sqlDb.sql.QueryRow(query, key)
	switch err := row.Scan(&originalURL); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println("Database returned" + originalURL)
	default:
		return err
	}
	return nil
}
