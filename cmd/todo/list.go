package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/snirkop89/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

func listCmd(tm *todo.TodoManager) *cli.Command {
	cmd := &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List all your tasks",
		Action: func(ctx *cli.Context) error {

			return listAction(tm, ctx.Args().Slice())
		},
	}

	return cmd
}

func listAction(tm *todo.TodoManager, args []string) error {
	filters := tm.ParseFilters(args)

	todos, err := tm.List(filters...)
	if err != nil {
		return err
	}

	if len(todos) == 0 {
		fmt.Println("No todos yet, how about adding one?")
		return nil
	}

	tw := table.NewWriter()
	tw.SetOutputMirror(os.Stdout)
	tw.AppendHeader(table.Row{"ID", "Done?", "Task", "Tags", "Due At", "Created At", "Completed At", "Priority"})

	tw.SetColumnConfigs([]table.ColumnConfig{
		{Number: 2, Align: text.AlignCenter},
		{Number: 6, Align: text.AlignCenter},
	})

	for _, t := range todos {
		prefix := ""
		if t.Completed {
			prefix = checkMark
		}
		tags := strings.Join(t.Tags, " ")
		completedAt := "TDB"
		if !t.CompletedAt.IsZero() {
			completedAt = t.CreatedAt.Local().Format(timeFormat)
		}
		dueDate := "No due date"
		if !t.DueDate.IsZero() {
			dueDate = t.DueDate.Local().Format(timeFormat)
		}
		tw.AppendRow(table.Row{fmt.Sprintf("%d", t.ID), prefix, t.Task, tags, dueDate, t.CreatedAt.Local().Format(timeFormat), completedAt, t.Priority})
	}
	tw.Render()
	return nil
}
