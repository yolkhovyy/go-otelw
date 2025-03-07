package domain

import (
	"context"
)

type Controller struct{}

//nolint:ireturn
func New(
	_ context.Context,
) (*Controller, error) {
	controller := Controller{}

	return &controller, nil
}

func (u Controller) Close() error {
	return nil
}
