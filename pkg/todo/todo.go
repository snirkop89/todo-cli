package todo

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

type priority int

const (
	none priority = iota
	high
	urgent
)

func (p priority) String() string {
	switch p {
	case none:
		return "None"
	case high:
		return "High"
	case urgent:
		return "Urgent"
	default:
		return "Unknown"
	}
}

func (p *priority) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		*p = none
		return nil
	}

	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*p = priority(i)
	return nil
}

func (p priority) MarshalJSON() ([]byte, error) {
	pri := int(p)
	data, err := json.Marshal(pri)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type Todo struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	Task        string    `json:"task"`
	Completed   bool      `json:"done"`
	Tags        []string  `json:"tags"`
	Notes       []string  `json:"notes"`
	DueDate     time.Time `json:"due_date"`
	Priority    priority  `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"updated_at"`
}

func (t *Todo) Done() {
	t.Completed = true
	t.CompletedAt = time.Now().UTC()
}

func (t *Todo) Undone() {
	t.Completed = false
	t.CompletedAt = time.Time{}
}

func (t *Todo) AddNote(note string) {
	t.Notes = append(t.Notes, note)
}

func (t *Todo) Prioritize() {
	if t.Priority < urgent {
		t.Priority++
	}
}

func (t *Todo) Unprioritze() {
	if t.Priority > none {
		t.Priority--
	}
}

func (t *Todo) SetPriority(pri string) error {
	switch strings.ToLower(pri) {
	case "n", "none":
		t.Priority = none
	case "h", "high":
		t.Priority = high
	case "u", "urgent":
		t.Priority = urgent
	default:
		return errors.New("unknown priority")
	}
	return nil
}
