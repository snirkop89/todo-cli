package main

import (
	"fmt"
	"strings"

	"github.com/snirkop89/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

func addCmd(tm *todo.TodoManager) *cli.Command {
	cmd := &cli.Command{
		Name:  "add",
		Usage: "add task @context +project [due:XXX pri:XXX]",
		Description: `
Add works with todo.txt format.
Use '@' to set the context of the task.
Use '+' to set the project of the task.
Use 'due:time' to set a due date. See docs.
Use 'pri:XXX' to set a priority for the task. See docs. 
		`,
		Aliases: []string{"a"},

		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				return fmt.Errorf("too little information. provide a task to add")
			}
			return addAction(tm, ctx.Args().Slice())
		},
	}

	return cmd
}

func addAction(tm *todo.TodoManager, args []string) error {
	if err := tm.Add(strings.Join(args, " ")); err != nil {
		return err
	}
	fmt.Println("Added task")
	return nil
}
