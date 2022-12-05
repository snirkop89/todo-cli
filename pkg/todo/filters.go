package todo

import "strings"

func (t *TodoManager) ParseFilters(arg string) []string {
	var filters []string

	opts := strings.Fields(arg)

	// Expect status:!completed
	for _, opt := range opts {
		parts := strings.Split(opt, ":")
		if len(parts) != 2 {
			continue
		}
		filterType, value := parts[0], parts[1]
		switch filterType {
		case "status":
			if strings.HasPrefix(value, "-") {
				filters = append(filters, "notcompleted")
			} else {
				filters = append(filters, value)
			}
		}
	}

	return filters
}

type filterFunc func([]*Todo) []*Todo

var filtersFuncs = map[string]filterFunc{
	"completed":    filterCompleted,
	"notcompleted": filterNotCompleted,
}

var filterCompleted = func(todos []*Todo) []*Todo {
	var filteredTodos []*Todo
	for _, t := range todos {
		if t.Completed {
			filteredTodos = append(filteredTodos, t)
		}
	}
	return filteredTodos
}

var filterNotCompleted = func(todos []*Todo) []*Todo {
	var filteredTodos []*Todo
	for _, t := range todos {
		if !t.Completed {
			filteredTodos = append(filteredTodos, t)
		}
	}
	return filteredTodos
}
