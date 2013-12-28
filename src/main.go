package main

import (
	"bitbucket.org/dpritchett/analyst"
	"github.com/hoisie/web"
    "github.com/joho/godotenv"
    "encoding/json"
	"log"
    "database/sql"
)

var db *sql.DB

func report() (results [][]string) {
    myEnv, err := godotenv.Read()

	connString := myEnv["CONN_STRING"]
	queryString := "SELECT * FROM spree_states order by name asc"

    if db == nil {
      db, err = analyst.Connect(connString)
      if err != nil {
          log.Fatal(err)
      }
    }


	columns, rows, err := analyst.Query(db, queryString)

	if err != nil {
		log.Fatal(err)
	}

	results = append(results, columns)

	for _, row := range rows {
		results = append(results, row)
	}

	return
}

func helloSQL(ctx *web.Context) (body []byte, err error) {
  ctx.ContentType("json")
  body, err = json.Marshal(report())
  return
}

func helloWorld(ctx *web.Context) (body []byte, err error) {
  body = []byte{'h', 'i'}
  return
}

func serve() {
  web.Get("/sql", helloSQL)
  web.Get("/", helloWorld)
  web.Run("0.0.0.0:9999")
}

func main() {
  //analyst.Hello()
  serve()
}
