package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ := sql.Open("sqlite3", "./dgpbt.db")

	table := "pbt_main"

	CreateDB(database, table)

	CreateAll200779DBRows(database, table)
	CreateAllInvoiceCosts(database, table)
}
