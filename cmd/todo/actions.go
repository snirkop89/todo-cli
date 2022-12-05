package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/snirkop89/todo-cli/pkg/todo"
)

const checkMark = "\u2713"

func addAction(tm *todo.TodoManager, args []string) error {
	return tm.Add(strings.Join(args, " "))
}

func listAction(tm *todo.TodoManager, args []string) error {
	filters := tm.ParseFilters(strings.Join(args, " "))

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
	tw.AppendHeader(table.Row{"#", "Done?", "Task", "Tags", "Due At", "Created At", "Completed At", "Priority"})

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
			dueDate = t.CreatedAt.Local().Format(timeFormat)
		}
		tw.AppendRow(table.Row{fmt.Sprintf("%d", t.ID), prefix, t.Task, tags, dueDate, t.CreatedAt.Local().Format(timeFormat), completedAt, t.Priority})
	}
	tw.Render()
	return nil
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

func deleteAction(tm *todo.TodoManager, args []string) error {
	var errs []error
	for _, a := range args {
		id, err := strconv.Atoi(a)
		if err != nil {
			return errors.New("id must be numeric")
		}
		if err := tm.Delete(id); err != nil {
			errs = append(errs, err)
			continue
		}
		fmt.Fprintf(os.Stdout, "Deleted item %d\n", id)
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func editAction(tm *todo.TodoManager, id string, subject string) error {
	todoID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	if err := tm.Edit(todoID, subject); err != nil {
		return err
	}

	return nil
}
