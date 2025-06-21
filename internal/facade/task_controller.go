package facade

import (
	"github.com/boykovrtm/workmate-intern/internal/application/interfaces"
	"github.com/boykovrtm/workmate-intern/internal/domain/entities"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type TaskController struct {
	repository entities.TaskRepository
	logger     *logrus.Logger
	collection interfaces.HandlerCollection
}

func NewTaskController(webApp *fiber.App, repository entities.TaskRepository, logger *logrus.Logger, collection interfaces.HandlerCollection) {
	controller := &TaskController{
		repository: repository,
		logger:     logger,
		collection: collection,
	}

	group := webApp.Group("api/v1/tasks")
	group.Post("/", func(c *fiber.Ctx) error {
		return controller.CreateTask(c)
	})

	group.Get("/:id", func(c *fiber.Ctx) error {
		return controller.GetTask(c)
	})

	group.Patch("/:id/retry", func(c *fiber.Ctx) error {
		return controller.TryAgain(c)
	})

	group.Delete("/:id", controller.DeleteTask)
}

func (tc *TaskController) CreateTask(c *fiber.Ctx) error {
	var req CreateTaskRequest

	err := c.BodyParser(&req)
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"error":   err,
			"handler": "CreateTask",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("error parsing request req")
		return c.SendStatus(fiber.StatusBadRequest)
	}

	taskHandler, ok := tc.collection[req.Name]
	if !ok {
		tc.logger.WithFields(logrus.Fields{
			"error":   "handler not found",
			"handler": "CreateTask",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("handler not found")
	}

	task, err := entities.NewTask(req.Payload, taskHandler)
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "CreateTask",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("Failed to create task")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	err = tc.repository.Save(task)
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "CreateTask",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("Failed to save task")

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusCreated).JSON(mapTask(task))
}

func (tc *TaskController) GetTask(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "GetTask",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("Failed parse id")

		return c.SendStatus(fiber.StatusBadRequest)
	}
	task, err := tc.repository.Get(id)
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "GetTask",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("Failed to get task")

		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.Status(fiber.StatusOK).JSON(mapTask(task))
}

func (tc *TaskController) TryAgain(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "TryAgain",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("Failed parse id")

		return c.SendStatus(fiber.StatusBadRequest)
	}
	task, err := tc.repository.Get(id)
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "TryAgain",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("Failed to get task")

		return c.SendStatus(fiber.StatusNotFound)
	}

	task.Retry()
	return c.Status(fiber.StatusOK).JSON(mapTask(task))
}

func (tc *TaskController) DeleteTask(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "DeleteTask",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("Failed parse id")

		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = tc.repository.Delete(id)
	if err != nil {
		tc.logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"handler": "DeleteTask",
			"method":  c.Method(),
			"path":    c.Path(),
		}).Error("Failed to delete task")
		return c.SendStatus(fiber.StatusNotFound)
	}

	return c.SendStatus(fiber.StatusOK)
}

type CreateTaskRequest struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

type TaskView struct {
	ID                 uuid.UUID           `json:"id"`
	CreatedAt          time.Time           `json:"created"`
	TakenAt            time.Time           `json:"taken_at"`
	CompletedAt        time.Time           `json:"completed_at"`
	ProcessingDuration string              `json:"processing-duration"`
	Payload            string              `json:"payload"`
	Status             entities.TaskStatus `json:"status"`
	Result             string              `json:"result"`
}

func mapTask(task entities.Task) TaskView {
	return TaskView{
		ID:                 task.ID,
		CreatedAt:          task.CreatedAt,
		TakenAt:            task.TakenAt,
		CompletedAt:        task.CompletedAt,
		ProcessingDuration: task.ProcessingDuration().String(),
		Payload:            task.Payload,
		Status:             task.Status,
		Result:             task.Result,
	}
}
