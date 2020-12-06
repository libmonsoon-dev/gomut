package app

import "context"

func New() (*App, error) {
	return &App{}, nil
}

type App struct {
}

func (a App) Run(ctx context.Context) error {
	// TODO
	return nil
}
