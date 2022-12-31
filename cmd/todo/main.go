package main

import (
	"fmt"
	"os"

	"github.com/snirkop89/todo-cli/pkg/todo"
	"github.com/snirkop89/todo-cli/pkg/todo/repository/filestore"
	"github.com/urfave/cli/v2"
)

const (
	todoFile   = ".todos.json"
	timeFormat = "2006-01-02"
	checkMark  = "\u2713"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fileStore := filestore.New(todoFile)

	tm, err := todo.NewManager(fileStore)
	if err != nil {
		return err
	}

	// Initalize cli application and run
	app := cli.NewApp()
	app.Name = "todo"
	app.Usage = "Manage tasks in your CLI"
	app.Commands = []*cli.Command{
		addCmd(tm),
		completeCmd(tm),
		deleteCmd(tm),
		editCmd(tm),
		listCmd(tm),
		statsCmd(tm),
	}

	return app.Run(os.Args)
}
