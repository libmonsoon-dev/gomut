package packages

import (
	"context"
	"fmt"
	"github.com/go-toolsmith/pkgload"
	"golang.org/x/tools/go/packages"
)

// ErrMatchedNoPackages - the error that Load returns in case the underlying loader does not return an error,
// but returns an empty slice of packages
var ErrMatchedNoPackages = fmt.Errorf("matched no packages")

// Load parse and returns the Go packages named by the given patterns.
// Uses packages.Load under the hood with an application-specific config.
// In case of a successful load, returns the first error of the package, if any
func Load(paths ...string) ([]*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedImports |
			packages.NeedTypes |
			packages.NeedCompiledGoFiles |
			packages.NeedSyntax,
	}

	patterns := make([]string, len(paths))
	for i := range paths {
		patterns[i] = "pattern=" + paths[i]
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("could not load package: %w", err)
	}

	findFirstError := func(pkg *packages.Package) bool {
		for _, pkgErr := range pkg.Errors {
			err = fmt.Errorf("package %s: %w", pkg, pkgErr)
			return false
		}

		return true
	}

	packages.Visit(pkgs, findFirstError, nil)
	if err != nil {
		return nil, err
	}

	if len(pkgs) == 0 {
		return nil, ErrMatchedNoPackages
	}

	return pkgload.Deduplicate(pkgs), nil
}

// Walk is a Node generator that bypasses the packages specified in the arguments.
// The first argument, ctx, can safely be nil
func Walk(ctx context.Context, pkgs []*packages.Package) <-chan Node {
	if ctx == nil {
		ctx = context.Background()
	}
	w := newWalker(ctx)

	go func() {
		defer close(w.out)
		packages.Visit(pkgs, w.pkgWalk, nil)
	}()
	return w.out
}
