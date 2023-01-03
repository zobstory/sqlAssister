// Ephmrl (Ephemeral) sqlAssister functions allow for the DB connection to be opened, the function to be used for an operation, & the connection to be closed.
// Use these when the DB connection is expected to be ephemeral.
// The methods exposed by the interface are expecting a persistent connection to the DB

package sqlAssister

import (
	"database/sql"
	"errors"
	"log"
)

// EphmrlUpdateSingleRow executes any CRUD operation EXCEPT Read for a single record
/*

Example:
	db, err := sql.Open("postgres", "DB info placeholder")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	err := AssisterConfig.EphmrlUpdateSingleRow(db, statement, args)
	if err != nil {
		return nil, err
	}
*/
func EphmrlUpdateSingleRow(db *sql.DB, statement string, params ...interface{}) error {
	log.Printf("Query: %s", statement)
	stmt, err := db.Prepare(statement)
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

// EphmrlSingleRowScanner Executes Read operation on a single record & scans a single record into a struct.
// Expects ONLY a single record to be returned
/*

Example:

	db, err := sql.Open("postgres", "DB info placeholder")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	yourStruct := &YourStruct{}
	row, err := AssisterConfig.EphmrlSingleRowScanner(db, statement, args)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&yourStruct)
	if err != nil {
		return nil, err
	}
*/
func EphmrlSingleRowScanner(db *sql.DB, statement string, params ...interface{}) (*sql.Row, error) {
	if len(params) < 1 {
		noParamsErr := errors.New("no params were passed")
		return nil, noParamsErr
	}
	log.Printf("Query: %s", statement)
	stmt, err := db.Prepare(statement)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return nil, err
	}
	row := stmt.QueryRow(params...)
	return row, nil
}

// EphmrlMultipleRowScanner Executes Read operation on multiple records & scans them into a slice of a struct
// NOTE: EphmrlMultipleRowScanner can work with a single record BUT please use EphmrlSingleRowScanner if you are only expecting a single record to be found
/*

Example:

	db, err := sql.Open("postgres", "DB info placeholder")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	var yourStructSlice []*YourStruct
	rows, err := AssisterConfig.EphmrlMultipleRowScanner(db, statement, args)
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
func EphmrlMultipleRowScanner(db *sql.DB, statement string, params ...interface{}) (*sql.Rows, error) {
	if len(params) < 1 {
		noParamsErr := errors.New("no params were passed")
		return nil, noParamsErr
	}
	log.Printf("Query: %s", statement)
	stmt, err := db.Prepare(statement)
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