package main

import (
	"encoding/json"
	//"github.com/nu7hatch/gouuid"
	"github.com/dpritchett/go-analyst/connection"
	"github.com/hoisie/web"
	"github.com/joho/godotenv"
	"html/template"
	"log"
)

var templates = template.Must(template.ParseFiles(
  "templates/index.html",
  "templates/rowset.html",
  "templates/builder.html"))

type Rowset struct {
	Columns []string
	Rows    [][]string
}

type QueryResult struct {
  Rowset *Rowset
  Query  string
  Error  error
}

func renderAsTable(rs *Rowset, ctx *web.Context) {
	templates.ExecuteTemplate(ctx, "rowset.html", rs)
}

func hisqlite() *Rowset {
	queryString := "SELECT * FROM Queries"

	db, err := connection.Connect("sqlite3", "db/development.sqlite3.db")
	if err != nil {
		log.Fatal(err)
	}

	columns, rows, err := connection.Query(db, queryString)

	if err != nil {
		log.Fatal(err)
	}

	return &Rowset{Columns: columns, Rows: rows}
}

func execQuery(queryString string) (results *QueryResult) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return
	}

	connString := myEnv["CONN_STRING"]

	db, err := connection.Connect("postgres", connString)
	if err != nil {
		return
	}

    results = &QueryResult{Query: queryString}

	columns, rows, err := connection.Query(db, queryString)

	if err != nil {
        results.Error = err
		return
	}

	results.Rowset = &Rowset{Columns: columns, Rows: rows}

	return
}

func hiPg(ctx *web.Context) (err error) {
	result := execQuery("select * from spree_states")
	renderAsTable(result.Rowset, ctx)
	return
}

func helloSQLite(ctx *web.Context) (body []byte, err error) {
	ctx.ContentType("json")
	body, err = json.Marshal(hisqlite())
	return
}

func handleQuery(ctx *web.Context) (body []byte, err error) {
	log.Print(ctx.Params)

	results := execQuery(ctx.Params["query"])

	if err != nil {
		log.Printf("Error: %v", err)
		body = []byte("Error executing query!")
	} else {
      body, err = json.Marshal(results)
	}
	return
}

func helloWorld(ctx *web.Context) (err error) {
	templates.ExecuteTemplate(ctx, "index.html", nil)
	return
}

func heavySQLite(ctx *web.Context) {
	rs := hisqlite()
	renderAsTable(rs, ctx)
}

func buildQuery(ctx *web.Context) {
  templates.ExecuteTemplate(ctx, "builder.html", &QueryResult{Query: "select 1 as one;"})
}

func serve() {
	web.Post("/sql-query", handleQuery)
	web.Get("/", buildQuery)
	web.Get("/lite", helloSQLite)
	web.Get("/heavy", heavySQLite)
	web.Get("/pg", hiPg)

	web.Run("0.0.0.0:9999")
}

func main() {
	serve()
}
