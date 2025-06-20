package handlers

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"
)

type TestHandler struct {
	logger *logrus.Logger
}

func NewTestHandler(logger *logrus.Logger) *TestHandler {
	return &TestHandler{
		logger: logger,
	}
}

func (h *TestHandler) Handle(ctx context.Context, payload string) (string, error) {
	time.Sleep(10 * time.Second)
	h.logger.WithFields(logrus.Fields{
		"payload":     payload,
		"handlerName": h.Name(),
	}).Info("Handled request of payload")

	return "Done", nil
}

func (h *TestHandler) Name() string {
	return "TestHandler"
}
