package main

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "log"
)

func InitDB(dataSourceName string) (*sql.DB, error) {
    fmt.Printf("entra %s", dataSourceName)
    db, err := sql.Open("mysql", dataSourceName)

    if err != nil {
        fmt.Printf("error db %s", err)
        log.Panic(err)
    }

    if err = db.Ping(); err != nil {
        fmt.Printf("error db 2 %s", err)
        log.Panic(err)
    }

    return db, nil
}
