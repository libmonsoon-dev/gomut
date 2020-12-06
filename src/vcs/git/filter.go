package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
)

type Config struct {
	Path string

	// TODO: rename
	Skip int
}

type NodeFilter struct {
	patch []diff.FilePatch
}

func NewNodeFilter(config Config) (*NodeFilter, error) {
	patch, err := getDiff(config)
	if err != nil {
		return nil, fmt.Errorf("getDiff: %w", err)
	}

	return &NodeFilter{patch}, nil
}

func getDiff(config Config) (patches []diff.FilePatch, err error) {
	repo, err := git.PlainOpen(config.Path)
	if err != nil {
		return nil, fmt.Errorf("open repository: %w", err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("get HEAD reference: %w", err)
	}

	head, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, fmt.Errorf("get HEAD commit: %w", err)
	}

	cIter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return nil, fmt.Errorf("get commit iterator: %w", err)
	}
	defer cIter.Close()

	for i := 0; i < config.Skip; i++ {
		_, err := cIter.Next()
		if err != nil {
			return nil, fmt.Errorf("skip commit (%v): %w", i, err)
		}
	}

	commit, err := cIter.Next()
	if err != nil {
		return nil, fmt.Errorf("get target commit: %w", err)
	}

	patch, err := commit.Patch(head)
	if err != nil {
		return nil, fmt.Errorf("generating patch: %w", err)
	}

	for _, patch := range patch.FilePatches() {
		if patch.IsBinary() {
			continue
		}

		from, to := patch.Files()
		var path string
		if from != nil && to != nil {
			if from.Path() != to.Path() {
				path = fmt.Sprintf("%v -> %v", from.Path(), to.Path())
			} else {
				path = to.Path()
			}
		} else if from != nil {
			path = from.Path()
		} else if to != nil {
			path = to.Path()
		} else {
			panic("Imposable case")
		}

		if !strings.HasSuffix(path, ".go") {
			continue
		}

		patches = append(patches, patch)
	}
	return
}
