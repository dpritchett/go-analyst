package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Connect(fileName string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", fileName)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func Query(db *sql.DB, queryString string) (columns []string, rowset [][]string, err error) {
	rows, err := db.Query(queryString)
	if err != nil {
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		return
	}

	columns = cols

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		result := make([]string, len(cols))
		err = rows.Scan(dest...)
		if err != nil {
			return
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}

		rowset = append(rowset, result)
	}

	return
}
