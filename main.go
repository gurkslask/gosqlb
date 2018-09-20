package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var debug bool = false
var debugp *bool = &debug

func main() {

	if len(os.Args) < 2 {
		printHelp()
		os.Exit(1)
	}
	if string(os.Args[1][0]) == "-" {
		printHelp()
		os.Exit(1)
	}

	var table = flag.String("table", "", "Name of the table in db")
	var limit = flag.Int("limit", 10, "Number of lines to print")
	var offset = flag.Int("offset", 1, "Offset from where to start printing")
	var debugf = flag.Bool("debug", false, "Prints extra debug info")
	var help = flag.Bool("help", false, "Prints help dialog")
	var dbtype = flag.String("dbtype", "sqlite", `What kind of database?\nSupported databases: sqlite and postgres`)
	flag.CommandLine.Parse(os.Args[2:])

	if *help {
		printHelp()
		os.Exit(0)
	}

	path := &os.Args[1]
	var err error
	if *dbtype == "sqlite" {
		_, err = os.Stat(*path)
		if os.IsNotExist(err) {
			fmt.Println("Supplied path does not exist")
			os.Exit(1)
		}
	}

	debugp = debugf
	if *debugp {
		fmt.Println("In main")
		fmt.Println(table)
	}

	if *path == "" {
		err = errors.New("No path specified")
		if err != nil {
			log.Fatal(err)
		}
	}
	dbSpecifics, err := getDbSpecifics(*dbtype)
	if err != nil {
		log.Fatal(err)
	}
	err = connectDB(path, dbSpecifics)
	defer dbSpecifics.db.Close()
	if err != nil {
		fmt.Println(err)
	}
	if *table == "" {
		tables, err := dbSpecifics.lister.availableTables(dbSpecifics.db)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("No table provided, listed tables: \n", tables)
	} else {
		result, err := printTable(*dbSpecifics, table, limit, offset)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
	}

}

type dblister interface {
	columnInfo(*string, *sql.DB) ([][]string, error)
	availableTables(*sql.DB) (string, error)
	statement() string
}
type dsa struct {
	dbtype string
	db     *sql.DB
	lister dblister
}

func printTable(d dsa, tablename *string, limit, offset *int) (string, error) {
	// Get data with queries and print it nicely with padding
	if *debugp {
		fmt.Println("*************In printTable")
		fmt.Println("In Printtable with tablename: %s", *tablename)

	}
	fmt.Printf("Table %s\n\n", *tablename)
	//stmt, err := d.db.Prepare(fmt.Sprintf("SELECT * FROM %s LIMIT $1 OFFSET $2 ", *tablename))
	//stmt, err := d.db.Prepare(fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ", *tablename))
	stmt, err := d.db.Prepare(fmt.Sprintf(d.lister.statement(), *tablename))
	if err != nil {
		fmt.Println("stmt err")
		return "", err
	}
	rows, err := stmt.Query(*limit, *offset)
	if err != nil {
		return "", err
	}
	data, err := getData(rows)
	if err != nil {
		return "", err
	}
	if data == nil {
		err := errors.New("No data from the query")
		return "", err
	}
	heads, err := d.lister.columnInfo(tablename, d.db)
	if err != nil {
		log.Fatal("heads ", err)
	}
	if *debugp {
		fmt.Println("This is the data:")
		fmt.Println(data)
		fmt.Println("This is the heads:")
		fmt.Println(heads)
	}
	columnlengths := maxColumnLength(data, heads)

	resultstring := ""
	for i := 0; i < len(data[0]); i++ {
		padString(heads[i][1], columnlengths[i], &resultstring)
		//resultstring += "\t"
	}
	resultstring += "\n"
	for i := 0; i < len(data[0]); i++ {
		////resultstring += heads[i][2]
		padString(heads[i][2], columnlengths[i], &resultstring)
		//resultstring += "\t"
	}

	resultstring += "\n"
	resultstring += "\n"
	for _, row := range data {
		for i, col := range row {
			padString(col, columnlengths[i], &resultstring)
		}
		resultstring += "\n"
	}

	if *debugp {
		fmt.Println("Out Printtable with tablename: %s", *tablename)
	}
	return resultstring, nil
}
