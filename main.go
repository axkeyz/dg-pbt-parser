package main

import (
	"time"
	"fmt"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ := sql.Open("sqlite3", "./dgpbt.db")

	table := "pbt_main"

	CreateDB(database, table)

	CreateAll200779DBRows(database, table)
	CreateAllInvoiceCosts(database, table)
	CreateAll23635Rows(database, table)

	// fmt.Println(GetDBRowsByMonth(database, table, "07"))

	month := fmt.Sprintf("%02d", int(time.Now().Month()) - 1)

	ExportDB(database, table, month)
}
