package datastore

import (
	"fmt"
	//	"database/sql"
	"github.com/jmoiron/sqlx"
	// Blank import comment
	_ "github.com/lib/pq"
)


// Postgre access code
var Postgre *sqlx.DB

// ConnectPostgre interface function, exportable
func ConnectPostgre() {
	var err error

	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DbUser, DbPassword, DbName)
	//"user=postgres dbname=udemy_fileserver sslmode=disable"
	Postgre, err = sqlx.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	Postgre.SetMaxIdleConns(1)
	Postgre.SetMaxOpenConns(8)
}
