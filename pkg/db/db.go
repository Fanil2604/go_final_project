package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func Init(dbFile string) error {
	var db *sql.DB
	var err error
	const schema string = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
	);
	CREATE INDEX scheduler_date ON scheduler (date);`

	_, err = os.Stat(dbFile)
	if os.IsNotExist(err) {
		f, err := os.Create(dbFile)
		if err != nil {
			panic(err)
		}
		f.Close()

	}
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = db.Exec(schema)
	if err != nil {
		return err
	}
	return err
}
