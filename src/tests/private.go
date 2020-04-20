package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/libmonsoon-dev/jsontest"
	"io"
)

func run(config Config, ch chan<- Event) {
	defer close(ch)
	cmd := config.GetCommand()
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		ch <- Event{Err: fmt.Errorf("could not pipe test cmd stdout: %w", err)}
		return
	}

	stderrBuf := new(bytes.Buffer)
	cmd.Stderr = stderrBuf

	if err := cmd.Start(); err != nil {
		ch <- Event{Err: fmt.Errorf("could not start test cmd: %w", err)}
		return
	}

	decoder := json.NewDecoder(stdout)
	var e jsontest.Event

	for {
		if err := decoder.Decode(&e); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			ch <- Event{Err: fmt.Errorf("could not unmarshal test Event: %w", err)}
			return
		}

		ch <- Event{Event: e}
	}

	if err := cmd.Wait(); err != nil {
		stdout := stderrBuf.Bytes()
		if len(stdout) > 0 && stdout[len(stdout)-1] == '\n' {
			stdout = stdout[:len(stdout)-1]
		}
		ch <- Event{Err: fmt.Errorf("test command error: %w: %s", err, stdout)}
	}
}
