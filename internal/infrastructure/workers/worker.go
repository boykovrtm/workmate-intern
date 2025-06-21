package workers

import (
	"context"
	"github.com/boykovrtm/workmate-intern/internal/domain/entities"
	"github.com/sirupsen/logrus"
	"time"
)

type Worker struct {
	taskRepository entities.TaskRepository
	logger         *logrus.Logger
}

func NewWorker(taskRepository entities.TaskRepository, logger *logrus.Logger) *Worker {
	return &Worker{
		taskRepository: taskRepository,
		logger:         logger,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Millisecond * 50):
		}

		task, err := w.taskRepository.Take()
		if err != nil {
			w.logger.WithFields(logrus.Fields{
				"error": err,
				"task":  task,
			}).Error("failed to take task")
		}

		if task == nil {
			continue
		}

		err = task.Complete(ctx)
		if err != nil {
			w.logger.WithFields(logrus.Fields{
				"error": err,
				"task":  task,
			}).Error("failed to complete task")
		}

		err = w.taskRepository.Save(*task)
		if err != nil {
			w.logger.WithFields(logrus.Fields{
				"error": err,
				"task":  task,
			}).Error("failed to save task")
		}
	}
}
