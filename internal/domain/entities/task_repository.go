package entities

import (
	"github.com/google/uuid"
)

type TaskRepository interface {
	Save(task Task) error
	Get(id uuid.UUID) (Task, error)
	Delete(id uuid.UUID) error
	Take() (*Task, error)
}
