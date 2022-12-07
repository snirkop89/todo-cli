package main

import (
	"errors"
	"strconv"

	"github.com/snirkop89/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

func editCmd(tm *todo.TodoManager) *cli.Command {
	cmd := &cli.Command{
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
			if ctx.Args().Len() == 0 {
				return errors.New("expected at least one id, got none")
			}
			if len(ctx.String("subject")) < 5 {
				return errors.New("subject is too short")
			}
			return editAction(tm, ctx.Args().First(), ctx.String("subject"))
		},
	}

	return cmd
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
