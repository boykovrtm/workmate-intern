package entities

import (
	"awesomeProject2/internal/application/interfaces"
	"context"
	"github.com/google/uuid"
	"time"
)

type TaskStatus int

const (
	TaskStatusUnknown TaskStatus = iota
	TaskStatusCreated
	TaskStatusInWork
	TaskStatusCompleted
	TaskStatusFailed
)

type Task struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	TakenAt     time.Time
	CompletedAt time.Time
	Payload     string
	Status      TaskStatus
	Result      string
	handler     interfaces.Handler
}

func NewTask(payload string, taskHandler interfaces.Handler) (Task, error) {
	return Task{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		Status:    TaskStatusCreated,
		Payload:   payload,
		handler:   taskHandler,
	}, nil
}

func (t *Task) ProcessingDuration() time.Duration {
	if t.Status == TaskStatusCompleted {
		return t.CompletedAt.Sub(t.TakenAt)
	}

	if t.Status == TaskStatusInWork {
		return time.Now().Sub(t.TakenAt)
	}

	return 0
}

func (t *Task) MarkTaken() {
	if t.Status == TaskStatusCreated {
		t.Status = TaskStatusInWork
		t.TakenAt = time.Now()
	}
}

func (t *Task) Retry() {
	if t.Status == TaskStatusFailed {
		t.Status = TaskStatusCreated
		t.TakenAt = time.Time{}
	}
}

func (t *Task) Complete(ctx context.Context) error {
	result, err := t.handler.Handle(ctx, t.Payload)
	if err != nil {
		t.Status = TaskStatusFailed
		return err
	}
	t.Result = result
	t.Status = TaskStatusCompleted
	t.CompletedAt = time.Now()
	return nil
}
