package entities

import (
	"github.com/google/uuid"
)

type TaskRepository interface {
	Save(task Task) error
	Get(id uuid.UUID) (Task, error)
	Take() (*Task, error)
}
