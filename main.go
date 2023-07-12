package sqlx

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

func main() {
	db, err := sqlx.Connect("postgres", "host=192.168.49.2 port=30101 user=admin password=secret dbname=db sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	// rs := db.MustExec("INSERT INTO persons (first_name, last_name, email) VALUES ($1, $2, $3)", "JaSon", "Moiron", "jmoiron@jmoiron.net")
	// fmt.Println(rs.LastInsertId())
	// rs = db.MustExec("INSERT INTO persons (first_name, last_name, email) VALUES ($1, $2, $3)", "John", "Doe", "johndoeDNE@gmail.net")
	// fmt.Println(rs.LastInsertId())

	people := []Person{}
	db.Select(&people, "SELECT * FROM persons ORDER BY first_name ASC")
	fmt.Println(people)

	jason := Person{}
	err = db.Get(&jason, "SELECT * FROM persons where first_name = 'JaSon'")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(jason)

	rows, err := db.NamedQuery(`SELECT * FROM persons WHERE first_name=:first_name`, jason)
	if err != nil {
		log.Println(err)
	}
	p := Person{}
	for rows.Next() {
		err := rows.StructScan(&p)
		if err != nil {
			log.Println(err)
		}
	}
}
