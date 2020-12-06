package git

import (
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/format/diff"
)

// Config stores configuration information for NodeFilter
type Config struct {
	Path string

	// TODO: rename
	N int
}

// NodeFilter check if a node was changed in the previous N commits
type NodeFilter struct {
	patch []diff.FilePatch
}

// NewNodeFilter return new *NodeFilter
func NewNodeFilter(config Config) (*NodeFilter, error) {
	diffObj, err := getDiff(config)
	if err != nil {
		return nil, fmt.Errorf("getDiff: %w", err)
	}

	return &NodeFilter{filterDiff(diffObj)}, nil
}
