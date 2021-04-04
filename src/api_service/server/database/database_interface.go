package database

type Database interface {
	InsertIntoDB(key string, value string, date int64) error
	FetchOriginalUrl(shortenURL string) error
}