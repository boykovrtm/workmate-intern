package in_memory

import (
	"awesomeProject2/internal/domain/entities"
	"errors"
	"github.com/google/uuid"
	"sync"
)

type InMemoryTasksStorage struct {
	tasks map[uuid.UUID]entities.Task
	mutex *sync.Mutex
}

func NewTasksStorage() *InMemoryTasksStorage {
	return &InMemoryTasksStorage{
		tasks: make(map[uuid.UUID]entities.Task),
		mutex: new(sync.Mutex),
	}
}

func (s *InMemoryTasksStorage) Save(task entities.Task) error {
	s.tasks[task.ID] = task
	return nil
}

func (s *InMemoryTasksStorage) Get(id uuid.UUID) (entities.Task, error) {
	task, ok := s.tasks[id]
	if !ok {
		return entities.Task{}, errors.New("task not found")
	}

	return task, nil
}

func (s *InMemoryTasksStorage) Delete(id uuid.UUID) error {
	_, ok := s.tasks[id]
	if !ok {
		return errors.New("task not found")
	}

	delete(s.tasks, id)

	return nil
}

func (s *InMemoryTasksStorage) Take() (*entities.Task, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, task := range s.tasks {
		if task.Status == entities.TaskStatusCreated {
			task.MarkTaken()

			err := s.Save(task)
			if err != nil {
				return nil, err
			}

			return &task, nil
		}
	}

	return nil, nil
}
