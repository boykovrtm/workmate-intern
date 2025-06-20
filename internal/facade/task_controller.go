package facade

import (
	"awesomeProject2/internal/application/interfaces"
	"awesomeProject2/internal/domain/entities"
	"awesomeProject2/internal/infrastructure/storage/in_memory"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type TaskController struct {
	storage    *in_memory.InMemoryTasksStorage
	logger     *logrus.Logger
	collection interfaces.HandlerCollection
}

func NewTaskController(webApp *fiber.App, storage *in_memory.InMemoryTasksStorage, logger *logrus.Logger, collection interfaces.HandlerCollection) {
	controller := &TaskController{
		storage:    storage,
		logger:     logger,
		collection: collection,
	}

	webApp.Post("/task", func(c *fiber.Ctx) error {
		return controller.CreateTask(c)
	})

	webApp.Get("/task/:id", func(c *fiber.Ctx) error {
		return controller.GetTask(c)
	})

	webApp.Patch("/task/:id", func(c *fiber.Ctx) error {
		return controller.TryAgain(c)
	})

	webApp.Delete("/task/:id", func(c *fiber.Ctx) error {
		return controller.DeleteTask(c)
	})
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

	err = tc.storage.Save(task)
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
	task, err := tc.storage.Get(uuid.MustParse(c.Params("id")))
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
	task, err := tc.storage.Get(uuid.MustParse(c.Params("id")))
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
	err := tc.storage.Delete(uuid.MustParse(c.Params("id")))
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
	ProcessingDuration string              `json:"completingTime"`
	Status             entities.TaskStatus `json:"status"`
	Result             string              `json:"result"`
}

func mapTask(task entities.Task) TaskView {
	return TaskView{
		ID:                 task.ID,
		CreatedAt:          task.CreatedAt,
		ProcessingDuration: task.ProcessingDuration().String(),
		Status:             task.Status,
		Result:             task.Result,
	}
}
