package config

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "monitor.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
		return nil, err
	}

	sqlFile, err := ioutil.ReadFile("db/init.sql")
	if err != nil {
		log.Fatal("Failed to read init.sql:", err)
		return nil, err
	}

	_, err = db.Exec(string(sqlFile))
	if err != nil {
		log.Fatal("Failed to execute init.sql:", err)
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}
