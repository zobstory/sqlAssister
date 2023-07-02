/*
Package sqlAssister provides the QueryAssister interface which provides cleaner sql CRUD operations along with query & error logging

Example:

	package main

	import (
		"database/sql"
		_ "github.com/lib/pq"
		"github.com/zobstory/sqlAssister"
		"log"
	)

	var statementAssister *sqlAssister.Assister

	type Book struct {
		ID   string
		Name string
	}

	func init() {
		db, err := sql.Open("postgres", "DB info placeholder")
		if err != nil {
			log.Fatal(err)
		}

		statementAssister = sqlAssister.New(db)
	}

	func SelectBook(bookId string) (*Book, error) {
		book := &Book{}
		const statement = `
			SELECT
				"ID",
				"cpu_temp",
				"fan_speed",
				"hdd_space",
				"last_logged_in",
				"sys_time"
			FROM "Network"."vw_device"
			WHERE "ID" = $1;`

		row, err := statementAssister.SingleRowScanner(statement, bookId)
		if err != nil {
			return nil, err
		}

		err = row.Scan(book)
		if err != nil {
			return nil, err
		}

		return book, nil
	}

	func main() {
		book, err := SelectBook("1")
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(book)
	}

See https://pkg.go.dev/database/sql for documentation on the standard sql library
*/
package sqlAssister

import (
	"database/sql"
	"github.com/zobstory/sqlAssister/utils"
)

type Assister struct {
	DB *sql.DB
}

// New returns a new instance of Assister to access the QueryAssister interface
func New(db *sql.DB) *Assister {
	config := &Assister{
		DB: db,
	}
	return config
}

// UpdateSingleRow executes any CRUD operation EXCEPT Read for a single record
/*

Example:

	err := Assister.UpdateSingleRow(statement, args)
	if err != nil {
		return nil, err
	}
*/
func (ac Assister) UpdateSingleRow(query string, args ...any) error {
	stmt, err := ac.DB.Prepare(query)
	if err != nil {
		return err
	}
	results, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	err = utils.GetRowsAffected(results, 1)
	if err != nil {
		return err
	}

	return nil
}

// SingleRowScanner Executes Read operation on a single record & scans a single record into a struct.
// Expects ONLY a single record to be returned
/*

Example:

	yourStruct := &YourStruct{}
	row, err := Assister.SingleRowScanner(statement)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&yourStruct)
	if err != nil {
		return nil, err
	}
*/
func (ac Assister) SingleRowScanner(query string) (*sql.Row, error) {
	err := utils.QueryChecker(query)
	if err != nil {
		return nil, err
	}

	row := ac.DB.QueryRow(query)
	return row, nil
}

// SingleRowScannerWithArgs Executes Read operation on a single record & scans a single record into a struct.
// Expects ONLY a single record to be returned
/*

Example:

	yourStruct := &YourStruct{}
	row, err := Assister.SingleRowScanner(query, args)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&yourStruct)
	if err != nil {
		return nil, err
	}
*/
func (ac Assister) SingleRowScannerWithArgs(query string, args ...any) (*sql.Row, error) {
	err := utils.QueryCheckerWithArgs(query, args)
	if err != nil {
		return nil, err
	}

	row := ac.DB.QueryRow(query, args...)
	return row, nil
}

// MultipleRowScanner Executes Read operation on multiple records & scans them into a slice of a struct
// NOTE: MultipleRowScanner can work with a single record BUT please use SingleRowScanner if you are only expecting a single record to be found
/*

Example:

	var yourStructSlice []*YourStruct
	rows, err := Assister.MultipleRowScanner(query, args)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		yourStruct := &YourStruct{}
		err := rows.Scan(&yourStruct)
		if err != nil {
			return nil, err
		}
		yourStructSlice = append(yourStructSlice, yourStruct)
	}
*/
func (ac Assister) MultipleRowScanner(query string) (*sql.Rows, error) {
	err := utils.QueryCheckerWithArgs(query)
	if err != nil {
		return nil, err
	}

	rows, err := ac.DB.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// MultipleRowScannerWithArgs Executes Read operation on multiple records & scans them into a slice of a struct
// NOTE: MultipleRowScannerWithArgs can work with a single record BUT please use SingleRowScannerWithArgs if you are only expecting a single record to be found
/*

Example:

	var yourStructSlice []*YourStruct
	rows, err := Assister.MultipleRowScannerWithArgs(statement, args)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		yourStruct := &YourStruct{}
		err := rows.Scan(&yourStruct)
		if err != nil {
			return nil, err
		}
		yourStructSlice = append(yourStructSlice, yourStruct)
	}
*/
func (ac Assister) MultipleRowScannerWithArgs(query string, args ...any) (*sql.Rows, error) {
	err := utils.QueryCheckerWithArgs(query, args)
	if err != nil {
		return nil, err
	}

	rows, err := ac.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
