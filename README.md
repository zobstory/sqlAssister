# sqlAssister
A simple package to make working with the standard go sql library easier. 

`Assister` provides the methods:
- `UpdateSingleRow()`
  - Updates a single record or returns `error`
- `SingleRowScanner()`
  - Reruns `*sql.Row` or `error`
- `SingleRowScannerWithArgs()`
  - Requires at least a single argument to be passed with the query & returns `*sql.Row` or `error`
- `MultipleRowScanner()`
  - Returns `*sql.Rows` or `error`
- `MultipleRowScannerWithArgs()`
  - Requires at least a single argument to be passed with the query & returns `*sql.Rows` or `error`
- `New()`
  - Returns `*Assister`

Additionally, functions are provided for use with ephemeral DB connections (open a connection to the DB, execute an operation, close the connection to the DB).
The functions:
- `EphmrlUpdateSingleRow()`
  - Updates a single record or returns `error`
- `EphmrlSingleRowScanner()`
  - Returns `*sql.Row` or `error`
- `EphmrlSingleRowScannerWithArgs()`
  - Requires at least a single argument to be passed with the query & returns `*sql.Row` or `error`
- `EphmrlMultipleRowScannerWithArgs()`
  - Requires at least a single argument to be passed with the query & returns `*sql.Rows` or `error`

### Assister & StatementAssister Interface
sqlAssister provides methods exposed through an interface that expect a persistent connection to the DB
Example:
```
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

    row, err := statementAssister.SingleRowScannerWithArgs(statement, bookId)
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
```

### Ephemeral function example
Use when you expect to open & close a connection to a DB during each operation execution

```
db, err := sql.Open("postgres", "DB info placeholder")
if err != nil {
    log.Fatalln(err)
}
defer db.Close()

yourStruct := &YourStruct{}
row, err := sqlAssister.EphmrlSingleRowScannerWithArgs(db, statement, args)
if err != nil {
    return nil, err
}

err = row.Scan(&yourStruct)
if err != nil {
    return nil, err
}

```
