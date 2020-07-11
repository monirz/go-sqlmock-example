package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password"
	dbname   = "testdb"
)

type tag struct {
	id   int64
	name string
}

type post struct {
	title       string
	description string
	tags        []*tag
}

//createPost creates new post
func createPost(db *sql.DB, post *post) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	query := `
	 INSERT INTO posts(title, description) VALUES($1,$2) returning id 
	`

	var lastInsertID int64
	err = tx.QueryRow(query, post.title, post.description).Scan(&lastInsertID)
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Error inserting into posts")
	}

	query, args := queryInsertArray(post.tags, lastInsertID)
	log.Println(query, args)

	//bulk insert from array
	_, err = tx.Exec(query, args...)

	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Error inserting authors_posts")
	}

	return err
}

//queryInsertArray makes query string from bulk insert
func queryInsertArray(arr []*tag, postID int64) (string, []interface{}) {
	query := `INSERT INTO posts_tags (tag_id, post_id) values `

	values := []interface{}{}

	for i, t := range arr {

		values = append(values, t.id, postID)

		numFields := 2
		n := i * numFields

		query += `(`
		for j := 0; j < numFields; j++ {
			query += `$` + strconv.Itoa(n+j+1) + `,`
		}
		query = query[:len(query)-1] + `),`
	}
	query = query[:len(query)-1]
	return query, values
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	pst := &post{
		title:       "This is title",
		description: "Some description",
	}

	err = createPost(db, pst)
	if err != nil {
		panic(err)
	}

}
