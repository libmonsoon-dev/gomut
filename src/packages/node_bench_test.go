package packages

import (
	"testing"
)

var out string

func BenchmarkNode_String(b *testing.B) {
	type Benchmark struct {
		pkgPattern string
		fileIndex  int
	}

	benchmarks := []Benchmark{
		{arithmeticV1, 0},
	}

	for _, bench := range benchmarks {
		b.Run(bench.pkgPattern, func(b *testing.B) {
			pkg := mustLoad(bench.pkgPattern)[0]
			file := pkg.Syntax[bench.fileIndex]
			node := Node{
				Package:  pkg,
				File:     file,
				FileName: pkg.CompiledGoFiles[bench.fileIndex],
				Node:     file,
			}

			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				out = node.String()
			}
		})
	}
}
