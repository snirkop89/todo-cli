package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/snirkop89/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

func completeCmd(tm *todo.TodoManager) *cli.Command {
	cmd := &cli.Command{
		Name:    "complete",
		Aliases: []string{"c"},
		Usage:   "Complete task(s)",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() == 0 {
				return errors.New("expected at least 1 id, got none")
			}
			return completeAction(tm, ctx.Args().Slice())
		},
	}

	return cmd
}

func completeAction(tm *todo.TodoManager, args []string) error {
	var errs []error
	for _, a := range args {
		id, err := strconv.Atoi(a)
		if err != nil {
			return errors.New("id must be numerics")
		}
		if err := tm.Complete(id); err != nil {
			errs = append(errs, err)
			continue
		}
		fmt.Fprintf(os.Stdout, "Completed item %d\n", id)
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
