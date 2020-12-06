package app

import "context"

// New create App object
func New() (*App, error) {
	return &App{}, nil
}

// App is a main object
type App struct {
}

// Run start application
func (a App) Run(ctx context.Context) error {
	// TODO
	return nil
}
