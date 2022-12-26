package todo

import (
	"strings"
)

func (t *TodoManager) ParseFilters(args []string) []keepFunc {
	var filters []keepFunc

	for _, w := range args {
		switch {
		case strings.HasPrefix(w, "+"):
			filters = append(filters, WithProject(strings.TrimPrefix(w, "+")))
		case strings.HasPrefix(w, "@"):
			filters = append(filters, WithContext(strings.TrimPrefix(w, "@")))
		case strings.HasPrefix(w, "pri:"):
			filters = append(filters, WithPriority(strings.TrimPrefix(w, "pri:")))
		case strings.EqualFold(w, "done"):
			filters = append(filters, Completed)
		case strings.EqualFold(w, "todo"):
			filters = append(filters, NotCompleted)
		}
	}

	return filters
}

type keepFunc func(t *Todo) bool

func Filter(todos []*Todo, f keepFunc) []*Todo {
	result := []*Todo{}

	for _, t := range todos {
		if f(t) {
			result = append(result, t)
		}
	}
	return result
}

func Completed(t *Todo) bool {
	return t.Completed
}

func NotCompleted(t *Todo) bool {
	return !t.Completed
}

func WithProject(project string) func(t *Todo) bool {
	return func(t *Todo) bool {
		return t.Project == project
	}
}

func WithContext(context string) func(t *Todo) bool {
	return func(t *Todo) bool {
		return t.Context == context
	}
}

func WithPriority(pri string) func(t *Todo) bool {
	return func(t *Todo) bool {
		return strings.EqualFold(t.Priority.String(), pri)
	}
}
