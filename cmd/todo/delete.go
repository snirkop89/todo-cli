package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/snirkop89/todo-cli/pkg/todo"
	"github.com/urfave/cli/v2"
)

func deleteCmd(tm *todo.TodoManager) *cli.Command {
	cmd := &cli.Command{
		Name:    "delete",
		Aliases: []string{"del"},
		Usage:   "Delete task(s)",
		Action: func(ctx *cli.Context) error {
			if ctx.NArg() == 0 {
				return errors.New("expected at least 1 id, got none")
			}
			return deleteAction(tm, ctx.Args().Slice())
		},
	}

	return cmd
}

func deleteAction(tm *todo.TodoManager, args []string) error {
	var errs []error

	for _, a := range args {

		// Provided range to ids.
		if strings.Contains(a, "-") {
			nums := strings.Split(a, "-")
			start, err := strconv.Atoi(nums[0])
			if err != nil {
				return err
			}
			end, err := strconv.Atoi(nums[1])
			if err != nil {
				return err
			}
			ids := []int{}
			for i := start; i <= end; i++ {
				ids = append(ids, i)
			}
			if err := tm.Delete(ids...); err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "Deleted items %s\n", a)

		} else {
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
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
