package utils

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type DB struct {
	DB *sqlx.DB

	// Configs
	Url             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
}

func (_db *DB) Connect() {
	db, err := sqlx.Connect("postgres", _db.Url)
	if err != nil {
		log.Fatal(err.Error())
	}

	db.SetMaxOpenConns(_db.MaxOpenConns)
	db.SetMaxIdleConns(_db.MaxIdleConns)
	db.SetConnMaxLifetime(time.Second * time.Duration(_db.ConnMaxLifetime))

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("DB connected")

	_db.DB = db
}
