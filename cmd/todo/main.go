package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/snirkop89/todo-cli/pkg/todo"
	"github.com/snirkop89/todo-cli/pkg/todo/repository/filestore"
	"github.com/urfave/cli/v2"
)

const (
	todoFile   = ".todos.json"
	timeFormat = "2006-01-02 @ 15:04"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fileStore := filestore.New(todoFile)
	app := &cli.App{
		Name:  "todo",
		Usage: "Manage tasks in your CLI",
		Commands: []*cli.Command{
			{
				Name:    "add",
				Usage:   "add a new task",
				Aliases: []string{"a"},
				Action: func(ctx *cli.Context) error {
					if ctx.NArg() == 0 {
						return fmt.Errorf("too little information. provide a task to add")
					}

					tm, err := todo.NewManager(fileStore)
					if err != nil {
						return err
					}

					return addAction(tm, ctx.Args().Slice())
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "list all your tasks",
				Action: func(ctx *cli.Context) error {
					tm, err := todo.NewManager(fileStore)
					if err != nil {
						return err
					}
					return listAction(tm, ctx.Args().Slice())
				},
			},
			{
				Name:    "complete",
				Aliases: []string{"c"},
				Usage:   "Complete task(s)",
				Action: func(ctx *cli.Context) error {
					tm, err := todo.NewManager(fileStore)
					if err != nil {
						return err
					}
					if ctx.Args().Len() == 0 {
						return errors.New("expected at least 1 id, got none")
					}
					return completeAction(tm, ctx.Args().Slice())
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"del"},
				Usage:   "Delete task(s)",
				Action: func(ctx *cli.Context) error {
					tm, err := todo.NewManager(fileStore)
					if err != nil {
						return err
					}
					if ctx.Args().Len() == 0 {
						return errors.New("expected at least 1 id, got none")
					}
					return deleteAction(tm, ctx.Args().Slice())
				},
			},
			{
				Name:    "edit",
				Aliases: []string{"e"},
				Usage:   "Edit task",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "subject",
						Required: true,
						Usage:    "Change todo main subject",
					},
				},
				// TODO: Add option for due date and tags
				Action: func(ctx *cli.Context) error {
					tm, err := todo.NewManager(fileStore)
					if err != nil {
						return err
					}
					if ctx.Args().Len() == 0 {
						return errors.New("expected at least one id, got none")
					}
					if len(ctx.String("subject")) < 5 {
						return errors.New("subject is too short")
					}
					return editAction(tm, ctx.Args().First(), ctx.String("subject"))
				},
			},
		},
	}

	return app.Run(os.Args)
}
