package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// GetRowsAffected helper function that takes the actual number rows affected & compares it to expected number rows affected.
// Returns an error if the expected rows affected don't match the actual rows affected
func GetRowsAffected(results sql.Result, targetNumRowsAffected int64) error {
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
