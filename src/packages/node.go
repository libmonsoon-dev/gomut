package packages

import (
	"fmt"
	"go/ast"
	"go/format"
	"golang.org/x/tools/go/packages"
	"io"
)

var errPackageIsNil = fmt.Errorf("node.Package is nil")
var errPackageFsetIsNil = fmt.Errorf("node.Package.Fset is nil")
var errNodeIsNil = fmt.Errorf("node.Node is nil")

// Node is the structure that Walk generates. Is ast.Node with context.
type Node struct {
	Package  *packages.Package
	File     *ast.File
	FileName string
	ast.Node
}

func (n Node) format(dst io.Writer) error {
	if n.Package == nil {
		return errPackageIsNil
	}
	if n.Package.Fset == nil {
		return errPackageFsetIsNil
	}
	if n.Node == nil {
		return errNodeIsNil
	}

	return format.Node(dst, n.Package.Fset, n.Node)
}

func (n Node) String() string {
	if !n.isSafeForFormatting() {
		return ""
	}

	buf := bufPool.Get()
	defer bufPool.Put(buf)

	if err := n.format(buf); err != nil {
		panic(fmt.Errorf("could not format %#v: %w", n.Node, err))
	}
	return buf.String()
}

func (n Node) isSafeForFormatting() bool {
	switch n.Node.(type) {
	case ast.Expr, ast.Stmt, ast.Decl, ast.Spec, *ast.File:
		return true
	default:
		return false
	}
}
