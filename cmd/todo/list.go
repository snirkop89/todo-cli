package main

import (
	"fmt"
	"strings"

	"github.com/cheynewallace/tabby"
	"github.com/snirkop89/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

func listCmd(tm *todo.TodoManager) *cli.Command {
	cmd := &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "List all your tasks",
		Description: `Show your todos and apply filters, if any.
todo ls [filter name ...filter]

Available Filters:
+project_name     Show tasks belong to a project
@context          Show tasks belong to a context
pri:priority      Show tasks of specific priority
done              Show completed tasks
todo              Show yet to be completed tasks`,
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

	tw := tabby.New()
	tw.AddHeader("ID", "Done?", "Task", "Tags", "Due At", "Created At", "Completed At", "Priority")

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

		tw.AddLine(fmt.Sprintf("%d", t.ID), prefix, t.Task, tags, dueDate, t.CreatedAt.Local().Format(timeFormat), completedAt, t.Priority)
	}

	tw.Print()
	return nil
}
