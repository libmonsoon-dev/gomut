package test

import (
	"context"
	"os/exec"
)

const (
	testCommand = "test"
	jsonFlag    = "-json"
)

// Config is a structure with which you can configure a test command.
// All fields can safely have zero value.
// For more info see "go help test".
type Config struct {
	Ctx             context.Context
	GoBin           string
	BuildTestFlags  []string
	Paths           []string
	TestBinaryFlags []string
}

// GetCommand returns not started test *exec.Cmd
func (c Config) GetCommand() *exec.Cmd {
	if c.Ctx == nil {
		c.Ctx = context.Background()
	}

	if c.GoBin == "" {
		c.GoBin = "go"
	}

	if len(c.Paths) == 0 {
		c.Paths = []string{"./..."}
	}

	args := []string{testCommand, jsonFlag}

	args = append(args, c.BuildTestFlags...)
	args = append(args, c.Paths...)
	args = append(args, c.TestBinaryFlags...)

	return exec.CommandContext(c.Ctx, c.GoBin, args...)
}
