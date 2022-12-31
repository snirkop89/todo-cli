package todo

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/gofrs/uuid"
)

// TodoManager manages a list of todos in a file.
type TodoManager struct {
	repo  repository
	Todos []*Todo
	eow   time.Weekday
	Out   io.Writer
}

type repository interface {
	Load() ([]*Todo, error)
	Save(todos []*Todo) error
}

// New manager returns a new TodoManager.
func NewManager(repo repository) (*TodoManager, error) {
	tm := &TodoManager{
		eow:  time.Saturday,
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

	todo := t.parse(task)

	todo.ID = t.MaxID() + 1
	todo.UUID = guid.String()
	todo.CreatedAt = time.Now()

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

	return t.repo.Save(t.Todos)
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
func (t *TodoManager) Delete(ids ...int) error {
	for _, id := range ids {
		if id <= 0 || id > t.MaxID() {
			return fmt.Errorf("invalid id: %d", id)
		}

		idx := t.indexByID(id)
		t.Todos = append(t.Todos[:idx], t.Todos[idx+1:]...)
	}

	return t.repo.Save(t.Todos)
}

// List outputs the existing todos based on provided filters
func (t *TodoManager) List(sortBy string, filters ...keepFunc) ([]*Todo, error) {
	todos := t.Todos

	for _, f := range filters {
		todos = Filter(todos, f)
	}

	switch sortBy {
	case "priority":
		sort.Sort(ByPriority(todos))
	case "created":
		sort.Sort(ByCreated(todos))
	case "due":
		sort.Sort(ByDue(todos))

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

type Stats struct {
	Total   int
	Done    int
	Pending int
	Overdue int
}

// Stats returns the number of total todos, done, pending, overdue
func (t *TodoManager) Stats() Stats {
	var stats Stats

	stats.Total = len(t.Todos)

	for _, t := range t.Todos {
		if t.Completed {
			stats.Done++
		} else {
			stats.Pending++
		}

		if !t.DueDate.IsZero() && t.DueDate.Before(time.Now()) {
			stats.Overdue++
		}
	}

	return stats
}

type ByPriority []*Todo

func (bp ByPriority) Len() int      { return len(bp) }
func (bp ByPriority) Swap(i, j int) { bp[i], bp[j] = bp[j], bp[i] }
func (bp ByPriority) Less(i, j int) bool {
	return bp[i].Priority < bp[j].Priority
}

type ByCreated []*Todo

func (bc ByCreated) Len() int      { return len(bc) }
func (bc ByCreated) Swap(i, j int) { bc[i], bc[j] = bc[j], bc[i] }
func (bc ByCreated) Less(i, j int) bool {
	return bc[i].CreatedAt.Before(bc[j].CreatedAt)
}

type ByDue []*Todo

func (bd ByDue) Len() int      { return len(bd) }
func (bd ByDue) Swap(i, j int) { bd[i], bd[j] = bd[j], bd[i] }
func (bd ByDue) Less(i, j int) bool {
	return bd[i].DueDate.Before(bd[j].DueDate)
}
