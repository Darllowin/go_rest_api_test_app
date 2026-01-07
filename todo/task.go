package todo

import "time"

type Task struct {
	title       string
	Description string
	Completed   bool

	CreatedAt   time.Time
	CompletedAt *time.Time
}

func NewTask(title string, description string) Task {
	return Task{
		title:       title,
		Description: description,
		Completed:   false,

		CreatedAt:   time.Now(),
		CompletedAt: nil,
	}
}

func (t *Task) Complete() {
	completeTime := time.Now()

	t.Completed = true
	t.CompletedAt = &completeTime
}

func (t *Task) UnComplete() {
	t.Completed = false
	t.CompletedAt = nil
}
