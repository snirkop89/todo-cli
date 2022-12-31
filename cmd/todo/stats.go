package main

import (
	"fmt"
	"io"
	"os"

	"github.com/snirkop89/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

func statsCmd(tm *todo.TodoManager) *cli.Command {
	cmd := &cli.Command{
		Name:    "stats",
		Aliases: []string{"st"},
		Usage:   "Show stats about your todos",
		Action: func(ctx *cli.Context) error {
			return statsAction(os.Stdout, tm, ctx.Args().Slice())
		},
	}

	return cmd
}

func statsAction(out io.Writer, tm *todo.TodoManager, args []string) error {
	stats := tm.Stats()

	fmt.Fprintf(out, "Done: %d\n", stats.Total)
	fmt.Fprintf(out, "Pending: %d\n", stats.Pending)
	fmt.Fprintf(out, "Done: %d\n", stats.Done)
	fmt.Fprintf(out, "Overdue: %d\n", stats.Overdue)
	fmt.Fprintln(out)

	return nil
}
