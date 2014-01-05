package main

import (
	"encoding/json"
	"github.com/dpritchett/go-analyst/connection"
	"github.com/hoisie/web"
	"github.com/joho/godotenv"
	"log"
)

func hisqlite() (results [][]string) {

	queryString := "SELECT * FROM Queries"

	db, err := connection.Connect("sqlite3", "db/development.sqlite3.db")
	if err != nil {
		log.Fatal(err)
	}

	columns, rows, err := connection.Query(db, queryString)

	if err != nil {
		log.Fatal(err)
	}

	results = append(results, columns)

	for _, row := range rows {
		results = append(results, row)
	}

	return
}

func report() (results [][]string, err error) {
	queryString := "SELECT * FROM spree_states order by name asc"

	log.Printf("Executing [%s]", queryString)
	return execQuery(queryString)
}

func execQuery(queryString string) (results [][]string, err error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return
	}

	connString := myEnv["CONN_STRING"]

	db, err := connection.Connect("postgres", connString)
	if err != nil {
		return
	}

	columns, rows, err := connection.Query(db, queryString)

	if err != nil {
		return
	}

	results = append(results, columns)

	for _, row := range rows {
		results = append(results, row)
	}

	return
}

func helloSQLite(ctx *web.Context) (body []byte, err error) {
	ctx.ContentType("json")
	body, err = json.Marshal(hisqlite())
	return
}

func handleQuery(ctx *web.Context) (body []byte, err error) {
	ctx.ContentType("json")
	log.Print(ctx.Params)

	results, err := execQuery(ctx.Params["query"])

	if err != nil {
		log.Printf("Error: %v", err)
		body = []byte("Error executing query!")
	} else {
		body, err = json.Marshal(results)
	}
	return
}

func helloWorld(ctx *web.Context) (body []byte, err error) {
	body = []byte("Hello world")
	return
}

func serve() {
	web.Post("/sql-query", handleQuery)
	web.Get("/", helloWorld)
	web.Get("/lite", helloSQLite)
	web.Run("0.0.0.0:9999")
}

func main() {
	//analyst.Hello()
	serve()
}
