package packages

import (
	"context"
	"fmt"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"testing"
	"time"
)

const (
	arithmeticV1 = "github.com/libmonsoon-dev/gomut/testdata/arithmetic/v1"
	arithmeticV2 = "github.com/libmonsoon-dev/gomut/testdata/arithmetic/v2"
	arithmeticV3 = "github.com/libmonsoon-dev/gomut/testdata/arithmetic/v3"

	arithmeticV1sum = "gomut/testdata/arithmetic/v1/sum.go"
	arithmeticV2sum = "gomut/testdata/arithmetic/v2/sum.go"
	arithmeticV3sum = "gomut/testdata/arithmetic/v3/sum.go"
)

func TestLoad(t *testing.T) {
	type Test struct {
		args          []string
		expectedIds   []string
		expectedError error
	}

	tests := []Test{
		{
			[]string{"./notExist"},
			nil,
			ErrMatchedNoPackages,
		},
		{
			[]string{"../../testdata/arithmetic/v1"},
			[]string{arithmeticV1},
			nil,
		},
		{
			[]string{arithmeticV1},
			[]string{arithmeticV1},
			nil,
		},
		{
			[]string{"../../testdata/arithmetic/v2"},
			[]string{arithmeticV2},
			nil,
		},
		{
			[]string{arithmeticV2},
			[]string{arithmeticV2},
			nil,
		},
		{
			[]string{"../../testdata/arithmetic/v3"},
			[]string{arithmeticV3},
			nil,
		},
		{
			[]string{arithmeticV3},
			[]string{arithmeticV3},
			nil,
		},
		{
			[]string{
				"../../testdata/arithmetic/v1",
				"../../testdata/arithmetic/v2",
				"../../testdata/arithmetic/v3",
			},
			[]string{arithmeticV1, arithmeticV2, arithmeticV3},
			nil,
		},
		{
			[]string{arithmeticV1, arithmeticV2, arithmeticV3},
			[]string{arithmeticV1, arithmeticV2, arithmeticV3},
			nil,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test#%v", i), func(t *testing.T) {
			got, err := Load(test.args...)
			if err != test.expectedError {
				t.Errorf("Load() error: %v, expectedError: %v", err, test.expectedError)
				return
			}
			if len(got) != len(test.expectedIds) {
				t.Errorf("Load()\ngot: %v\nexpected: %v", got, test.expectedIds)
			}
			for i := range got {
				if got[i].ID != test.expectedIds[i] {
					t.Errorf("Load()\ngot: %#v\nexpected: %v", got[i], test.expectedIds[i])
				}
			}
		})
	}
}

func TestWalk(t *testing.T) {
	type Test struct {
		pkgs     []*packages.Package
		messages []Message
	}

	tests := []Test{
		{
			[]*packages.Package{
				{
					Fset: token.NewFileSet(),
					ID:   "inMemoryPkg",
					Syntax: []*ast.File{
						{
							Name: &ast.Ident{
								Name: "inMemoryPkg",
							},
						},
					},
					CompiledGoFiles: []string{"inMemoryFileName"},
				},
			},
			[]Message{
				{
					"inMemoryPkg",
					"inMemoryFileName",
					"package inMemoryPkg\n",
				},
				{
					"inMemoryPkg",
					"inMemoryFileName",
					"inMemoryPkg",
				},
			},
		},
		{
			mustLoad(arithmeticV1),
			[]Message{
				{
					arithmeticV1,
					arithmeticV1sum,
					"package arithmetic\n\nfunc Add(a, b int) int {\n\treturn a + b\n}\n",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"arithmetic",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"func Add(a, b int) int {\n\treturn a + b\n}",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"Add",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"func(a, b int) int",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"a",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"b",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"int",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"int",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"{\n\treturn a + b\n}",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"return a + b",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"a + b",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"a",
				},
				{
					arithmeticV1,
					arithmeticV1sum,
					"b",
				},
			},
		},
		{
			mustLoad(arithmeticV2),
			[]Message{
				{
					arithmeticV2,
					arithmeticV2sum,
					"package arithmetic\n\nfunc Add(a, b int) int {\n\treturn a + b\n}\n",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"arithmetic",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"func Add(a, b int) int {\n\treturn a + b\n}",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"Add",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"func(a, b int) int",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"a",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"b",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"int",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"int",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"{\n\treturn a + b\n}",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"return a + b",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"a + b",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"a",
				},
				{
					arithmeticV2,
					arithmeticV2sum,
					"b",
				},
			},
		},
		{
			mustLoad(arithmeticV3),
			[]Message{
				{
					arithmeticV3,
					arithmeticV3sum,
					"package arithmetic\n\nfunc Add(a, b int) int {\n\treturn a + b\n}\n",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"arithmetic",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"func Add(a, b int) int {\n\treturn a + b\n}",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"Add",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"func(a, b int) int",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"a",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"b",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"int",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"int",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"{\n\treturn a + b\n}",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"return a + b",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"a + b",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"a",
				},
				{
					arithmeticV3,
					arithmeticV3sum,
					"b",
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test#%v", i), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
			defer cancel()

			ch := Walk(ctx, test.pkgs)

			for i, message := range test.messages {
				actual, ok := <-ch

				if !ok {
					t.Errorf("(index %v) Channel already closed, but expect %#v", i, message)
					break
				}

				if err := message.AssertMatch(actual); err != nil {
					t.Errorf("(index %v) %s", i, err)
				}
			}

			if actual, ok := <-ch; ok {
				t.Errorf("Expect: channel closed but received \n\"%s\"", actual)
			}
		})
	}
}

func mustLoad(paths ...string) []*packages.Package {
	pkgs, err := Load(paths...)
	if err != nil {
		panic(err)
	}
	return pkgs
}
