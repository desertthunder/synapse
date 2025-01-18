package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const DataSourceName string = "./synapse.db"

type Connection struct {
	Db *sql.DB
}

func CreateConnection() *Connection {
	var err error
	c := Connection{}

	if c.Db, err = sql.Open("sqlite3", DataSourceName); err != nil {
		DefaultLogger().Fatal(
			fmt.Sprintf("unable to connect to database %v %v", DataSourceName, err.Error()),
		)
	}
	return &c
}
