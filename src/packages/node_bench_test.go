package packages

import (
	"testing"
)

var out string

func BenchmarkNode_String(t *testing.B) {
	t.ReportAllocs()

	pkg := mustLoad(arithmeticV1)[0]
	file := pkg.Syntax[0]
	node := Node{
		Package:  pkg,
		File:     file,
		FileName: pkg.CompiledGoFiles[0],
		Node:     file,
	}

	for i := 0; i < t.N; i++ {
		out = node.String()
	}
}
