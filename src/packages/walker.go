package packages

import (
	"context"
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"os"
	"path"
)

type walker struct {
	ctx              context.Context
	out              chan Node
	currentPkg       *packages.Package
	currentFile      *ast.File
	projectFilePath  string
	absoluteFilePath string
}

func newWalker(ctx context.Context) *walker {
	return &walker{ctx: ctx, out: make(chan Node)}
}

func (w walker) pkgWalk(pkg *packages.Package) bool {
	w.setCurrentPkg(pkg)

	if syntaxLen, filesLen := len(pkg.Syntax), len(pkg.CompiledGoFiles); syntaxLen != filesLen {
		panic(fmt.Errorf("len(pkg.Syntax) != len(pkg.CompiledGoFiles) (%v != %v)", syntaxLen, filesLen))
	}

	for i := range pkg.Syntax {
		astFile := pkg.Syntax[i]
		projectFilePath := fmt.Sprintf("%s%c%s", pkg.ID, os.PathSeparator, path.Base(pkg.CompiledGoFiles[i]))

		w.setCurrentFile(astFile)
		w.setCurrentProjectFilePath(projectFilePath)
		w.setCurrentAbsoluteFilePath(pkg.CompiledGoFiles[i])

		ast.Inspect(astFile, w.astWalk)

		select {
		case <-w.ctx.Done():
			return false
		default:
			// continue
		}
	}

	return true
}

func (w walker) astWalk(node ast.Node) bool {
	msg := w.newMessage(node)

	if node == nil || !msg.isSafeForFormatting() {
		return true
	}

	select {
	case <-w.ctx.Done():
		return false
	case w.out <- msg:
		return true
	}
}

func (w *walker) newMessage(node ast.Node) Node {
	return Node{w.currentPkg, w.currentFile, w.projectFilePath, w.absoluteFilePath, node}
}

func (w *walker) setCurrentPkg(pkg *packages.Package) {
	w.currentPkg = pkg
}

func (w *walker) setCurrentFile(file *ast.File) {
	w.currentFile = file
}

func (w *walker) setCurrentProjectFilePath(projectFilePath string) {
	w.projectFilePath = projectFilePath
}

func (w *walker) setCurrentAbsoluteFilePath(absoluteFilePath string) {
	w.absoluteFilePath = absoluteFilePath
}
