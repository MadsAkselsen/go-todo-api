package tododb

import (
	"database/sql"
	"log"

	"github.com/mattn/go-sqlite3"
)

type TodoEntry struct {
	Id				int64
	Done			bool
	Description		string
}

func TryCreate(db *sql.DB) {
	_, err := db.Exec(`
		CREATE TABLE todo (
			id				INTEGER PRIMARY KEY,
			done			BOOLEAN NOT NULL CHECK (done IN (0, 1)),
			description		TEXT
		)
	`)

	if err != nil {
		// 'err.(sqlite3.Error)' means the error is being casted to a sqlite3 error
		if sqlError, ok := err.(sqlite3.Error); ok {
			// code 1 == "table already exists"
			if sqlError.Code != 1 {
				log.Fatal(sqlError)
			}
		} else {
			log.Fatal(err)
		}
	}
}

func CreateTodoEntry(db *sql.DB, description string) error {
	_, err := db.Exec(`INSERT INTO
		todo(done, description)
		VALUES(?, ?)`, false, description)
	
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func UpdateTodoEntry(db *sql.DB, id int, done bool) error {
	// convertedDone := 0
	// if done == true {
	// 	convertedDone = 1
	// }

	_, err := db.Exec(`
			UPDATE todo
			SET done=?
			WHERE id=?
	`, done, id)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func DeleteTodoEntry(db *sql.DB, id int) error {
	_, err := db.Exec(`
			DELETE FROM todo
			WHERE id = ?
	`, id)

	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetAllTodoEntries(db *sql.DB) ([]TodoEntry, error) {
	var empty []TodoEntry

	rows, err := db.Query(`
		SELECT id, done, description
		FROM todo`) // the last line 'LIMIT ? OFFSET ?' is what enables pagination
		
	if err != nil {
		log.Println(err)
		return empty, err
	}

	defer rows.Close()

	todoEntries := make([]TodoEntry, 0, 10)

	for rows.Next() {
		todoEntry, err := todoEntryFromRow(rows)
		if err != nil {
			return nil, err
		}
		todoEntries = append(todoEntries, *todoEntry)
	}

	return todoEntries, nil
}

// creates an email data structure from a database row
func todoEntryFromRow(row *sql.Rows) (*TodoEntry, error) {
	var id int64
	var done bool
	var description string

	err := row.Scan(&id, &done, &description)
	
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &TodoEntry{Id: id, Done: done, Description: description}, nil
}