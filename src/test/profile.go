package test

import (
	"github.com/libmonsoon-dev/gomut/src/packages"
	"golang.org/x/tools/cover"
)

// NewProfiles create Profiles from []*cover.Profile
func NewProfiles(input []*cover.Profile) Profiles {
	output := make(Profiles, len(input))

	for _, profile := range input {
		if profile == nil {
			continue
		}
		output[profile.FileName] = profile.Blocks
	}

	return output
}

// Profiles represents the profiling data
type Profiles map[string]ProfileBlocks

// IsCover check if node covered by tests
func (p Profiles) IsCover(node packages.Node) bool {
	return p[node.ProjectFilePath].IsCover(node)
}
