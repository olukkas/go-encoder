package database

import (
	"database/sql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
	"log"
)

type DataBase struct {
	Db     *sql.DB
	Env    string
	Dsn    string
	DbType string
}

func NewDataBaseTest() *sql.DB {
	db := DataBase{
		Env:    "test",
		DbType: "sqlite3",
		Dsn:    ":memory:",
	}

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Test db error: %s\n", err)
	}

	return conn
}

func (d *DataBase) Connect() (*sql.DB, error) {
	var err error

	d.Db, err = sql.Open(d.DbType, d.Dsn)
	if err != nil {
		return nil, err
	}

	return d.Db, nil
}
