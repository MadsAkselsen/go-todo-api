package main

import (
	"database/sql"
	"go-todo-api/jsonapi"
	"go-todo-api/tododb"
	"log"
	"sync"

	"github.com/alexflint/go-arg"
)

var args struct {
	// We can specify the todo_db in the commandline, otherwise it is set to a default
	DbPath string `arg:"env:TODO_DB`
	BindJson string `arg:"env:TODO_BIND_JSON"`
}

func main() {
	arg.MustParse(&args)

	if args.DbPath == "" {
		args.DbPath = "todo.db" // default DB location
	}
	if args.BindJson == "" {
		args.BindJson = ":8080" // default port
	}

	log.Printf("using database '%v'\n", args.DbPath)
	db, err := sql.Open("sqlite3", args.DbPath)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tododb.TryCreate(db)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Printf("starting JSON API server...\n")
		jsonapi.Serve(db, args.BindJson)
		wg.Done()
	}()

	wg.Wait()
}

