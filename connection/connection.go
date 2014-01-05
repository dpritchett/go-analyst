package connection

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	"log"
)


func isDbAcceptable(dbType string) bool {
  acceptableDBs := []string{"sqlite3", "postgres"}

  for _, t := range(acceptableDBs) {
    if t == dbType {
      return true
    }
  }
  return false
}

func Connect(dbType string, fileName string) (db *sql.DB, err error) {
    if !isDbAcceptable(dbType) {
      log.Fatal("Invalid DB type: ", dbType)
    }

	db, err = sql.Open(dbType, fileName)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// Borrowed from http://stackoverflow.com/a/14500756 [by user ANisus]
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
