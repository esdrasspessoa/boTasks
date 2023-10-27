package task

import (
	"time"

	"github.com/rs/xid"
)

type Task struct {
	ID          string
	Title       string
	Description string
	Deadline    time.Time
}

func NewTask(title, description string, deadline time.Time) *Task {
	if title == "" {
		return nil
	} else if deadline.Before(time.Now()) {
		return nil
	}

	return &Task{
		ID:          xid.New().String(),
		Title:       title,
		Description: description,
		Deadline:    deadline,
	}
}

func (t *Task) UpdateTask(title, description string, deadline time.Time) {
	if title != "" {
		t.Title = title
	} else if description != "" {
		t.Description = description
	} else if !deadline.IsZero() {
		t.Deadline = deadline
	}
}

/*
func (t *Task) DeleteTask() {

}
*/
