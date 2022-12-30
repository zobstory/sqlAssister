/*

Package sqlAssister provides the StatementAssister interface which provides cleaner sql CRUD operations along with query & error logging

Example:

package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/zobstory/sqlAssist"
	"log"
)

var statementAssist *sqlAssist.AssisterConfig

type Book struct {
	ID   string
	Name string
}

func init() {
	db, err := sql.Open("postgres", "DB info placeholder")
	if err != nil {
		log.Fatalln(err)
	}

	statementAssist = sqlAssist.New(db)
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

	row, err := statementAssist.ScanSingleRow(statement, bookId)
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
	"fmt"
	"log"
)

type AssisterConfig struct {
	DB *sql.DB
}

type StatementAssister interface {
	UpdateSingleRow(query string, params ...interface{}) error
	ScanSingleRow(db *sql.DB, query string, params ...interface{}) (*sql.Row, error)
	ScanMultipleRows(query string, params ...interface{}) (*sql.Rows, error)
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

err := statementAssist.UpdateSingleRow(statement, args)
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
	err = getRowsAffected(results, 1)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return err
	}
	return nil
}

// ScanSingleRow Executes Read operation on a single record & scans a single record into a struct.
// Expects ONLY a single record to be returned
/*

Example:

yourStruct := &YourStruct{}
row, err := statementAssist.UpdateSingleRow(statement, args)
if err != nil {
	return nil, err
}

err = row.Scan(&yourStruct)
if err != nil {
	return nil, err
}

*/
func (ac AssisterConfig) ScanSingleRow(statement string, params ...interface{}) (*sql.Row, error) {
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

// ScanMultipleRows Executes Read operation on multiple records & scans them into a slice of a struct
// NOTE: ScanMultipleRows can work with a single record BUT please use ScanSingleRow if you are only expecting a single record to be found
/*

Example:

var yourStructSlice []*YourStruct
rows, err := statementAssist.UpdateSingleRow(statement, args)
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
func (ac AssisterConfig) ScanMultipleRows(statement string, params ...interface{}) (*sql.Rows, error) {
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

// getRowsAffected Non-exported helper function that takes the actual number rows affected & compares it to expected number rows affected.
// Returns an error if the expected rows affected don't match the actual rows affected
func getRowsAffected(results sql.Result, targetNumRowsAffected int64) error {
	rowsAffected, err := results.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != targetNumRowsAffected {
		sqlErr := errors.New(fmt.Sprintf("number of rows affected does not match the expected number of rows affected: %v / %v", rowsAffected, targetNumRowsAffected))
		log.Printf("ERROR: %s", sqlErr)
		return sqlErr
	}
	log.Printf("Rows affected: %v / %v", rowsAffected, targetNumRowsAffected)
	return nil
}

