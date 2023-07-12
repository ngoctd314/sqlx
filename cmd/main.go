package main

import (
	"database/sql"
	"fmt"
	sqlx "learn-sqlx"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func main() {
	db, _ = sqlx.Open("sqlite3", ":memory:")
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	createSchema()
	insertData()
	queryx()
}

func createSchema() {
	schema := `CREATE TABLE place (
        country text,
        city text NULL,
        telcode integer);`

	// execute a query on the server
	result, err := db.Exec(schema)
	if err != nil {
		log.Println(err)
	}
	log.Println("create schema result", result)
}

func insertData() {
	cityState := `INSERT INTO place (country, telcode) VALUES (?, ?)`
	countryCity := `INSERT INTO place (country, city, telcode) VALUES (?, ?, ?)`
	db.MustExec(cityState, "Hong Kong", 852)
	db.MustExec(cityState, "Singapore", 65)
	db.MustExec(countryCity, "South Africa", "Johannesburg", 27)
}

// query is the primary way to run queries with database/sql that return row results. Query returns an sql.Rows object and an error
func query() {
	rows, err := db.Query("SELECT country, city, telcode FROM place")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var country string
		var city sql.NullString
		var telcode int
		err = rows.Scan(&country, &city, &telcode)
		fmt.Printf("country: %s, city: %s, telcode: %d\n", country, city.String, telcode)
	}
	// check the error from rows
	err = rows.Err()
}

type Place struct {
	Country       string
	City          sql.NullString
	TelephoneCode int `db:"telcode"`
}

func queryx() {
	rows, err := db.Queryx("SELECT telcode FROM place")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var p Place
		err := rows.StructScan(&p)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(p)
	}
}
