package main

import (
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreatPost(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()
	mock.ExpectBegin()

	tags := []*tag{
		{id: 1},
		{id: 2},
	}

	post := &post{
		title:       "test title",
		description: "test description",
		tags:        tags,
	}

	query := `
	INSERT INTO 
		posts
    `

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(1)

	mock.ExpectQuery(query).
		WithArgs("test title", "test description").WillReturnRows(rows)

	query = `
		 INSERT INTO posts_tags
		`

	//tag_id -> 1, post_id, tag_id -> 2, post_id -> 1
	mock.ExpectExec(query).WithArgs(1, 1, 2, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err = createPost(mockDB, post); err != nil {
		t.Errorf("Error inserting order: %s", err.Error())
	}

}
