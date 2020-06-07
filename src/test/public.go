package test

import (
	"context"
	"fmt"
	"github.com/libmonsoon-dev/gomut/src/rand"
	"golang.org/x/tools/cover"
	"os"
)

// GetCoverage exec go test with -coverprofile flag.
// After successful exit of test command parse end return coverage profile
func GetCoverage(ctx context.Context, path string) (Profiles, error) {
	coveragePath := fmt.Sprintf("%v/%v.coverage", os.TempDir(), rand.String(10))
	coverProfileFlag := fmt.Sprintf("-coverprofile=%s", coveragePath)

	cfg := Config{
		Ctx:            ctx,
		Path:           path,
		BuildTestFlags: []string{coverProfileFlag},
	}

	for event := range Run(cfg) {
		if err := event.Err; err != nil {
			return nil, fmt.Errorf("got error from Run(%#v): %w", path, err)
		}
	}

	profiles, err := cover.ParseProfiles(coveragePath)
	if err != nil {
		return nil, fmt.Errorf("cover.ParseProfiles(%v): %w", coveragePath, err)
	}
	return NewProfiles(profiles), nil
}

// Run exec go test command
func Run(config Config) <-chan Event {
	ch := make(chan Event)
	go run(config, ch)
	return ch
}
