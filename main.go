package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, _ := sql.Open("sqlite3", "./dgpbt.db")

	CreateDB(database, "pbt_200779")

	Create200779Rows(ExtractPBTONERunsheet())

	fmt.Println("Done")
}
