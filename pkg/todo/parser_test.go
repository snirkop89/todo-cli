package todo

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	task := "speak with @john about +project xxx due: tomorrow pri: high"

	tm := TodoManager{}
	todo := tm.parse(task)

	expTodo := Todo{
		Task:        task,
		Completed:   false,
		CreatedAt:   time.Now().UTC(),
		CompletedAt: time.Time{},
		Project:     "project",
		Context:     "john",
		DueDate:     time.Now().UTC().AddDate(0, 0, 1),
		Priority:    high,
	}

	if todo.Task != expTodo.Task {
		t.Errorf("expected task %q, got %q", expTodo.Task, task)
	}

	if todo.Project != expTodo.Project {
		t.Errorf("expected project %q, got %q", expTodo.Project, todo.Project)
	}

	if todo.Completed != expTodo.Completed {
		t.Errorf("expected completed %t got %t", expTodo.Completed, todo.Completed)
	}

	if todo.CreatedAt.Format("2006-01-02") != expTodo.CreatedAt.Format("2006-01-02") {
		t.Errorf("expected createdAt %q, got %q",
			expTodo.CreatedAt.Format("2006-01-02"),
			todo.CreatedAt.Format("2006-01-02"),
		)
	}

	if !todo.CompletedAt.IsZero() {
		t.Errorf("expected completed at to be zero time")
	}

	if todo.DueDate.Format("2006-01-02") != expTodo.DueDate.Format("2006-01-02") {
		t.Errorf("expected due date %q, got %q", expTodo.DueDate.Format("2006-01-02"), todo.DueDate.Format("2006-01-02"))
	}

}

func TestEndOfWeek(t *testing.T) {

	tests := []struct {
		name    string
		dow     time.Weekday
		expDate time.Time
	}{
		{name: "Sunday", dow: time.Sunday},
		{name: "Monday", dow: time.Monday},
		{name: "Tuesday", dow: time.Tuesday},
		{name: "Wednesday", dow: time.Wednesday},
		{name: "Thursday", dow: time.Thursday},
		{name: "Friday", dow: time.Friday},
		{name: "Saturday", dow: time.Saturday},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := TodoManager{
				eow: tt.dow,
			}
			eow := endOfWeek(tm.eow)

			if eow.Weekday() != tt.dow {
				t.Errorf("expected day %d, got %d", time.Saturday, eow.Weekday())
			}
		})
	}
}

func BenchmarkParse(b *testing.B) {
	task := "speak with @john about +project xxx due: tomorrow pri: high"

	tm := TodoManager{}
	for i := 0; i < b.N; i++ {
		_ = tm.parse(task)
	}
}
