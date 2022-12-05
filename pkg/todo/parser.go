package todo

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func parseTask(task string) (Todo, error) {
	parts := strings.Fields(task)

	var todo Todo
	for _, p := range parts {
		switch {
		case strings.HasPrefix(p, "due:"):
			fmt.Println("found due date, go parse")
			parts := strings.Split(p, ":")
			if len(parts) != 2 {
				return Todo{}, errors.New("invalid due date format. use due:xxx")
			}
			todo.DueDate = parseDueDate(parts[1])
		case strings.HasPrefix(p, "#"):
			fmt.Println("found tag")
		case strings.HasPrefix(p, "pri"):
			fmt.Println("found priority")
		case strings.HasPrefix(p, "@"):
			fmt.Println("found context")
		}
	}
	todo.Task = task
	return todo, nil
}

func parseDueDate(d string) time.Time {
	switch d {
	case "tom", "tomorrow":
		return time.Now().AddDate(0, 0, 1).UTC()
	case "sun", "sunday":
		return nextAfter(time.Sunday)
	case "mon", "monday":
		return nextAfter(time.Monday)
	case "thu", "thursday":
		return nextAfter(time.Thursday)
	case "wed", "wednesday":
		return nextAfter(time.Wednesday)
	case "tue", "tuesday":
		return nextAfter(time.Tuesday)
	case "fri", "friday":
		return nextAfter(time.Friday)
	case "sat", "saturday":
		return nextAfter(time.Saturday)
	}

	// If user specific a number, add this to the current date
	toAdd, err := strconv.Atoi(d)
	if err != nil {
		return time.Time{}
	}
	return time.Now().AddDate(0, 0, toAdd)
}

func nextAfter(weekday time.Weekday) time.Time {
	diff := int(weekday) - int(time.Now().Weekday())
	if diff <= 0 {
		diff += 7
	}
	return time.Now().AddDate(0, 0, diff)
}
