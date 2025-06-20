package main

import (
	"awesomeProject2/internal/application/handlers"
	"awesomeProject2/internal/application/interfaces"
	"awesomeProject2/internal/facade"
	"awesomeProject2/internal/infrastructure/storage/in_memory"
	"awesomeProject2/internal/infrastructure/workers"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func main() {
	webApp := fiber.New()
	logger := logrus.New()
	ctx := context.Background()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	// Здесь регистрируются все хендлеры
	handlerCollection := interfaces.HandlerCollection{}
	handlerCollection.Add(handlers.NewTestHandler(logger))
	handlerCollection.Add(&handlers.ErrHandler{})

	storage := in_memory.NewTasksStorage()

	facade.NewTaskController(webApp, storage, logger, handlerCollection)
	worker := workers.NewWorker(storage, logger)
	runWorkers(ctx, worker, 10)

	logrus.Fatal(webApp.Listen(":8080"))
}

func runWorkers(ctx context.Context, worker *workers.Worker, count int) {
	for i := 0; i < count; i++ {
		go func() {
			err := worker.Run(ctx)
			if err != nil {
				panic(err)
			}
		}()
	}
}
