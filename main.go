package main

import (
	"encoding/json"
	"github.com/dpritchett/go-analyst/pg"
	"github.com/dpritchett/go-analyst/sqlite"
	"github.com/hoisie/web"
	"github.com/joho/godotenv"
	"log"
)

func hisqlite() (results [][]string) {

	queryString := "SELECT * FROM Queries"

	db, err := sqlite.Connect("db/development.sqlite3.db")
	if err != nil {
		log.Fatal(err)
	}

	columns, rows, err := sqlite.Query(db, queryString)

	if err != nil {
		log.Fatal(err)
	}

	results = append(results, columns)

	for _, row := range rows {
		results = append(results, row)
	}

	return
}

func report() (results [][]string) {
	myEnv, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}

	connString := myEnv["CONN_STRING"]
	queryString := "SELECT * FROM spree_states order by name asc"

	db, err := pg.Connect(connString)
	if err != nil {
		log.Fatal(err)
	}

	columns, rows, err := pg.Query(db, queryString)

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

func helloSQLite(ctx *web.Context) (body []byte, err error) {
	ctx.ContentType("json")
	body, err = json.Marshal(hisqlite())
	return
}

func helloWorld(ctx *web.Context) (body []byte, err error) {
	body = []byte{'h', 'i'}
	return
}

func serve() {
	web.Get("/sql", helloSQL)
	web.Get("/", helloWorld)
	web.Get("/lite", helloSQLite)
	web.Run("0.0.0.0:9999")
}

func main() {
	//analyst.Hello()
	serve()
}
