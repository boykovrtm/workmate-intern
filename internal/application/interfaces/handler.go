package interfaces

import (
	"context"
)

type Handler interface {
	Handle(ctx context.Context, payload string) (string, error)
	Name() string
}
