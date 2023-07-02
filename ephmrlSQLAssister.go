// Ephmrl (Ephemeral) sqlAssister functions allow for the DB connection to be opened, the function to be used for an operation, & the connection to be closed.
// Use these when the DB connection is expected to be ephemeral.
// The methods exposed by the interface are expecting a persistent connection to the DB

package sqlAssister

import (
	"database/sql"
	"github.com/zobstory/sqlAssister/utils"
)

// EphmrlUpdateSingleRow executes any CRUD operation EXCEPT Read for a single record
/*

Example:
	db, err := sql.Open("postgres", "DB info placeholder")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err := Assister.EphmrlUpdateSingleRow(db, query, args)
	if err != nil {
		return nil, err
	}
*/
func EphmrlUpdateSingleRow(db *sql.DB, statement string, args ...any) (*sql.Result, error) {
	err := utils.QueryCheckerWithArgs(statement, args)
	if err != nil {
		return nil, err
	}

	results, err := db.Exec(statement, args...)
	if err != nil {
		return nil, err
	}

	err = utils.GetRowsAffected(results, 1)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

// EphmrlSingleRowScannerWithArgs Executes Read operation on a single record & scans a single record into a struct.
// Expects ONLY a single record to be returned
/*

Example:

	db, err := sql.Open("postgres", "DB info placeholder")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	yourStruct := &YourStruct{}
	row, err := Assister.EphmrlSingleRowScannerWithArgs(db, query, args)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&yourStruct)
	if err != nil {
		return nil, err
	}
*/
func EphmrlSingleRowScannerWithArgs(db *sql.DB, query string, args ...any) (*sql.Row, error) {
	err := utils.QueryCheckerWithArgs(query, args)
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(query, args...)
	return row, nil
}

// EphmrlSingleRowScanner Executes Read operation on a single record & scans a single record into a struct.
// Expects ONLY a single record to be returned
/*

Example:

	db, err := sql.Open("postgres", "DB info placeholder")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	yourStruct := &YourStruct{}
	row, err := Assister.EphmrlSingleRowScanner(db, query)
	if err != nil {
		return nil, err
	}

	err = row.Scan(&yourStruct)
	if err != nil {
		return nil, err
	}
*/
func EphmrlSingleRowScanner(db *sql.DB, query string) (*sql.Row, error) {
	err := utils.QueryChecker(query)
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(query)
	return row, nil
}

// EphmrlMultipleRowScanner Executes Read operation on multiple records & scans them into a slice of a struct
// NOTE: EphmrlMultipleRowScanner can work with a single record BUT please use EphmrlSingleRowScanner if you are only expecting a single record to be found
/*

Example:

	db, err := sql.Open("postgres", "DB info placeholder")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var yourStructSlice []*YourStruct
	rows, err := Assister.EphmrlMultipleRowScanner(db, statement)
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
func EphmrlMultipleRowScanner(db *sql.DB, query string) (*sql.Rows, error) {
	err := utils.QueryChecker(query)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// EphmrlMultipleRowScannerWithArgs Executes Read operation on multiple records & scans them into a slice of a struct
// NOTE: EphmrlMultipleRowScanner can work with a single record BUT please use EphmrlSingleRowScannerWithArgs if you are only expecting a single record to be found
/*

Example:

	db, err := sql.Open("postgres", "DB info placeholder")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var yourStructSlice []*YourStruct
	rows, err := Assister.EphmrlMultipleRowScanner(db, statement, args)
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
func EphmrlMultipleRowScannerWithArgs(db *sql.DB, query string, args ...any) (*sql.Rows, error) {
	err := utils.QueryCheckerWithArgs(query, args)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(query, args)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
