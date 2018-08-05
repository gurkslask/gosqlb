package main

import (
	//"database/sql"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	// "os"
	//"strconv"
)

var debug bool = true

func main() {
	if debug {
		fmt.Println("In main")
	}
	var path = flag.String("path", "/home/alex/data.db", "Path to .db file")
	var table = flag.String("table", "", "Name of the table in db")
	flag.Parse()
	db, err := connectDB(path)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
	}

	if *table == "" {
		tables, err := printTables(db)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("No table provided, listed tables: \n", tables)
	} else {
		result, err := printTable(db, table)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}

}

func connectDB(path *string) (*sqlx.DB, error) {
	if debug {
		fmt.Println("In connect")
	}
	db, err := sqlx.Open("sqlite3", *path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

type tablestruct struct {
	colname1 string
	colname2 string
	colname3 string
	colname4 string
	coldata1 string
	coldata2 string
	coldata3 string
	coldata4 string
}

func printTable(db *sqlx.DB, tablename *string) (string, error) {
	if debug {
		fmt.Println("*************In printTable")
	}

	rows, err := db.Query("PRAGMA table_info(TEST)")
	defer rows.Close()
	var i int = 0
	for rows.Next() {



	}

	fmt.Println("Table")

	rows, err = db.Query("select * from TEST")
	if err != nil {
		return "", err
	}
	dcols, err := rows.Columns()
	if err != nil {
		return "", err
	}

	rawResult := make([][]byte, len(dcols))
	result := make([]string, len(dcols))

	dest := make([]interface{}, len(dcols))
	for i, _ := range rawResult {
		dest[i] = &rawResult[i]
	}

	if err != nil {
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return "", err
		}
		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}
		fmt.Println(result)
	}
	for i, v range columns

	return "", nil
}

func printTables(db *sqlx.DB) (string, error) {
	if debug {
		fmt.Println("In printTables")
	}
	var result string
	rows, err := db.Query(`SELECT name FROM sqlite_master WHERE type = "table"`)
	if err != nil {
		return "", err
	}
	for rows.Next() {

		var name string
		err = rows.Scan(&name)
		if err != nil {
			return "", err
		}
		result += fmt.Sprint(name, "\n")
	}
	return result, nil
}
