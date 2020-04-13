package tests

import (
	"context"
	"os/exec"
)

const testCommand = "test"
const jsonFlag = "-json"

type config struct {
	Ctx             context.Context
	GoBin           string
	BuildTestFlags  []string
	Path            string
	TestBinaryFlags []string
}

func (c config) GetCommand() *exec.Cmd {
	goBin := c.GoBin
	if goBin == "" {
		goBin = "go"
	}

	ctx := c.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	args := []string{testCommand, jsonFlag}

	args = append(args, c.BuildTestFlags...)
	args = append(args, c.Path)
	args = append(args, c.TestBinaryFlags...)

	return exec.CommandContext(ctx, goBin, args...)
}
