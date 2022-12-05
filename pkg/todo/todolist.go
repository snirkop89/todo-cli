package todo

import (
	"fmt"
	"io"

	"github.com/gofrs/uuid"
)

// TodoManager manages a list of todos in a file.
type TodoManager struct {
	repo  repository
	Todos []*Todo
	Out   io.Writer
}

type repository interface {
	Load() ([]*Todo, error)
	Save(todos []*Todo) error
}

// New manager returns a new TodoManager.
func NewManager(repo repository) (*TodoManager, error) {
	tm := &TodoManager{
		repo: repo,
	}
	todos, err := repo.Load()
	if err != nil {
		return nil, err
	}
	tm.Todos = todos
	return tm, nil
}

// Add a new task to the list.
func (t *TodoManager) Add(task string) error {
	guid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	todo, err := parseTask(task)
	if err != nil {
		return err
	}

	todo.ID = t.MaxID() + 1
	todo.UUID = guid.String()

	t.Todos = append(t.Todos, &todo)
	if err := t.repo.Save(t.Todos); err != nil {
		return err
	}
	return nil
}

// Completed marks a task as completed.
func (t *TodoManager) Complete(id int) error {
	if id <= 0 || id > len(t.Todos)-1 {
		return fmt.Errorf("invalid id: %d", id)
	}

	for _, todo := range t.Todos {
		if todo.ID == id {
			todo.Done()
			break
		}
	}

	return nil
}

// Edit allows to edit any property of a task.
func (t *TodoManager) Edit(id int, subject string) error {
	todo, err := t.findByID(id)
	if err != nil {
		return err
	}

	// TODO: parse props of format prop:value and apply
	todo.Task = subject
	return t.repo.Save(t.Todos)
}

// TODO differentiate between completed and not completed. delete not completed
// by default unless provided all
// TODO allow deleting range
// Delete removes a task from the list completely.
func (t *TodoManager) Delete(id int) error {
	if id <= 0 || id > t.MaxID() {
		return fmt.Errorf("invalid id: %d", id)
	}

	idx := t.indexByID(id)
	t.Todos = append(t.Todos[:idx], t.Todos[idx+1:]...)

	return t.repo.Save(t.Todos)
}

// List outputs the existing todos based on provided filters
func (t *TodoManager) List(filters ...string) ([]*Todo, error) {
	todos := t.Todos

	for _, filter := range filters {
		if f, ok := filtersFuncs[filter]; ok {
			todos = f(todos)
		} else {
			return nil, fmt.Errorf("unknonwn filter: %s", filter)
		}
	}

	return todos, nil
}

// Uncomplete sets a completed task as not-completed.
func (t *TodoManager) Uncomplete(id int) error {
	if id <= 0 || id > len(t.Todos)-1 {
		return fmt.Errorf("invalid id: %d", id)
	}

	for _, todo := range t.Todos {
		if todo.ID == id {
			todo.Undone()
			break
		}
	}

	return nil
}

func (t *TodoManager) findByID(id int) (*Todo, error) {
	if id <= 0 || id > t.MaxID() {
		return nil, fmt.Errorf("invalid id: %d", id)
	}

	for _, todo := range t.Todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	return nil, fmt.Errorf("todo not found")
}

func (t *TodoManager) indexByID(id int) int {
	for i, todo := range t.Todos {
		if todo.ID == id {
			return i
		}
	}
	return -1
}

func (t *TodoManager) MaxID() int {
	var max int
	for _, todo := range t.Todos {
		if todo.ID > max {
			max = todo.ID
		}
	}
	return int(max)
}

func (t *TodoManager) countPending() int {
	var pending int

	for _, todo := range t.Todos {
		if !todo.Completed {
			pending++
		}
	}

	return pending
}

// TODO show statistics:
// total todos, done, pending, overdue
func (t *TodoManager) Stats() {}
