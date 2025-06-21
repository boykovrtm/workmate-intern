package main

import (
	"context"
	"github.com/boykovrtm/workmate-intern/internal/application/handlers"
	"github.com/boykovrtm/workmate-intern/internal/application/interfaces"
	"github.com/boykovrtm/workmate-intern/internal/facade"
	"github.com/boykovrtm/workmate-intern/internal/infrastructure/storage/in_memory"
	"github.com/boykovrtm/workmate-intern/internal/infrastructure/workers"
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
