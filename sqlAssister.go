/*
Package sqlAssister provides the StatementAssister interface which provides cleaner sql CRUD operations along with query & error logging

Example:

	package main

	import (
		"database/sql"
		_ "github.com/lib/pq"
		"github.com/zobstory/sqlAssister"
		"log"
	)

	var statementAssister *sqlAssister.AssisterConfig

	type Book struct {
		ID   string
		Name string
	}

	func init() {
		db, err := sql.Open("postgres", "DB info placeholder")
		if err != nil {
			log.Fatalln(err)
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
			log.Fatalln(err)
		}
		log.Fatalln(book)
	}

See https://pkg.go.dev/database/sql for documentation on the standard sql library
*/
package sqlAssister

import (
	"database/sql"
	"errors"
	"log"
)

type AssisterConfig struct {
	DB *sql.DB
}

type StatementAssister interface {
	UpdateSingleRow(query string, params ...interface{}) error
	SingleRowScanner(db *sql.DB, query string, params ...interface{}) (*sql.Row, error)
	MultipleRowScanner(query string, params ...interface{}) (*sql.Rows, error)
	New() (config *AssisterConfig)
}

// New returns a new instance of AssisterConfig to access the StatementAssister interface
func New(db *sql.DB) *AssisterConfig {
	config := &AssisterConfig{
		DB: db,
	}
	return config
}

// UpdateSingleRow executes any CRUD operation EXCEPT Read for a single record
/*

Example:

	err := AssisterConfig.UpdateSingleRow(statement, args)
	if err != nil {
		return nil, err
	}
*/
func (ac AssisterConfig) UpdateSingleRow(statement string, params ...interface{}) error {
	log.Printf("Query: %s", statement)
	stmt, err := ac.DB.Prepare(statement)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}
	results, err := stmt.Exec(params...)
	err = GetRowsAffected(results, 1)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}
	return nil
}

// SingleRowScanner Executes Read operation on a single record & scans a single record into a struct.
// Expects ONLY a single record to be returned
/*

Example:

	yourStruct := &YourStruct{}
	row, err := AssisterConfig.SingleRowScanner(statement, args)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&yourStruct)
	if err != nil {
		return nil, err
	}
*/
func (ac AssisterConfig) SingleRowScanner(statement string, params ...interface{}) (*sql.Row, error) {
	if len(params) < 1 {
		noParamsErr := errors.New("no params were passed")
		return nil, noParamsErr
	}
	log.Printf("Query: %s", statement)
	stmt, err := ac.DB.Prepare(statement)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return nil, err
	}
	row := stmt.QueryRow(params...)
	return row, nil
}

// MultipleRowScanner Executes Read operation on multiple records & scans them into a slice of a struct
// NOTE: MultipleRowScanner can work with a single record BUT please use SingleRowScanner if you are only expecting a single record to be found
/*

Example:

	var yourStructSlice []*YourStruct
	rows, err := AssisterConfig.MultipleRowScanner(statement, args)
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
func (ac AssisterConfig) MultipleRowScanner(statement string, params ...interface{}) (*sql.Rows, error) {
	if len(params) < 1 {
		noParamsErr := errors.New("no params were passed")
		return nil, noParamsErr
	}
	log.Printf("Query: %s", statement)
	stmt, err := ac.DB.Prepare(statement)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query(params...)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return nil, err
	}
	return rows, nil
}

