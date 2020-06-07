package test

import (
	"github.com/libmonsoon-dev/gomut/src/packages"
	"golang.org/x/tools/cover"
)

// ProfileBlocks represents the profiling data for specific file.
type ProfileBlocks []cover.ProfileBlock

// IsCover check if node covered by tests
func (pb ProfileBlocks) IsCover(node packages.Node) bool {
	nodeStart, nodeEnd := node.Position()

	if !nodeStart.IsValid() || !nodeEnd.IsValid() {
		return false
	}

	for _, block := range pb {
		if block.StartLine > nodeStart.Line ||
			block.EndLine < nodeEnd.Line ||
			(block.StartLine == nodeStart.Line && block.StartCol > nodeStart.Column) ||
			(block.EndLine == nodeEnd.Line && block.EndCol < nodeEnd.Column) {
			continue
		}

		return true
	}

	return false
}
