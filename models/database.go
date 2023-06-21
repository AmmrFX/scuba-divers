package database

import (
        "database/sql"
        "log"

        _ "github.com/go-sql-driver/mysql"
)

var (
        db *sql.DB
)

func ConnectToDB(dataSourceName string) {
        var err error
        db, err = sql.Open("mysql", dataSourceName)
        if err != nil {
                log.Fatalf("Failed to connect to database: %v", err)
        }

        err = db.Ping()
        if err != nil {
                log.Fatalf("Failed to ping database: %v", err)
        }
}

func CloseDB() {
        err := db.Close()
        if err != nil {
                log.Fatalf("Failed to close database connection: %v", err)
        }
}
