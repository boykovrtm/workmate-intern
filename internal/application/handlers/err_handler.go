package handlers

import (
	"context"
	"errors"
)

type ErrHandler struct {
}

func (e ErrHandler) Handle(ctx context.Context, payload string) (string, error) {
	return "", errors.New("oops")
}

func (e ErrHandler) Name() string {
	return "ErrHandler"
}
