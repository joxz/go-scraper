// https://stackoverflow.com/questions/35090436/struct-time-property-doesnt-load-from-go-sqlx-library
// https://www.compose.com/articles/accessing-relational-databases-using-go/
// http://jmoiron.github.io/sqlx/
// https://www.thepolyglotdeveloper.com/2017/04/using-sqlite-database-golang-application/

package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Prefix struct {
	Id     int64
	Region string
	Prefix string
}

func conndb() *sqlx.DB {
	db := sqlx.MustConnect("sqlite3", ":memory:")
	return db
}

func createschema(db *sqlx.DB) {
	schema := `
		CREATE TABLE IF NOT EXISTS zscaler (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			region TEXT, 
			hostname TEXT, 
			location  TEXT, 
			prefix TEXT
			);

		CREATE TABLE IF NOT EXISTS sfdc (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			region TEXT, 
			prefix TEXT
			);
	`
	db.MustExec(schema)
}

func main() {

}
