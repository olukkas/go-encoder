package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/lib/pq"
	"github.com/olukkas/go-encoder/domain"
	"log"
)

type DataBase struct {
	Db          *gorm.DB
	Env         string
	Dsn         string
	Debug       bool
	DbType      string
	AutoMigrate bool
}

func NewDataBase() *DataBase {
	return &DataBase{}
}

func NewDataBaseTest() *gorm.DB {
	db := DataBase{
		Env:         "test",
		DbType:      "sqlite3",
		Dsn:         ":memory:",
		AutoMigrate: true,
		Debug:       true,
	}

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Test db error: %s\n", err)
	}

	return conn
}

func (d *DataBase) Connect() (*gorm.DB, error) {
	var err error

	d.Db, err = gorm.Open(d.DbType, d.Dsn)
	if err != nil {
		return nil, err
	}

	if d.Debug {
		d.Db.LogMode(true)
	}

	if d.AutoMigrate {
		d.Db.AutoMigrate(&domain.Video{}, &domain.Job{})
		d.Db.Model(domain.Job{}).
			AddForeignKey("video_id", "videos (id)", "CASCADE", "CASCADE")
	}

	return d.Db, nil
}
