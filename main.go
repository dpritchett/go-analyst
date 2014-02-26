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

var templates = template.Must(template.ParseFiles("index.html", "rowset.html"))

type Rowset struct {
	Columns []string
	Rows    [][]string
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

func execQuery(queryString string) (rs *Rowset, err error) {
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

	rs = &Rowset{Columns: columns, Rows: rows}

	return
}

func hiPg(ctx *web.Context) (err error) {
	rs, err := execQuery("select * from spree_states")
	renderAsTable(rs, ctx)
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

func helloWorld(ctx *web.Context) (err error) {
	templates.ExecuteTemplate(ctx, "index.html", nil)
	return
}

func heavySQLite(ctx *web.Context) {
	rs := hisqlite()
	renderAsTable(rs, ctx)
}

func serve() {
	web.Post("/sql-query", handleQuery)
	web.Get("/", helloWorld)
	web.Get("/lite", helloSQLite)
	web.Get("/heavy", heavySQLite)
	web.Get("/pg", hiPg)

	web.Run("0.0.0.0:9999")
}

func main() {
	serve()
}
