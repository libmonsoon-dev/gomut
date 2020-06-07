package test

import (
	"context"
	"fmt"
	"github.com/libmonsoon-dev/gomut/src/packages"
	"testing"
	"time"
)

func TestProfiles_IsCover(t *testing.T) {
	tests := []struct {
		path    string
		timeout time.Duration
		expect  []bool
	}{
		{
			arithmeticV1,
			time.Second,
			[]bool{false, false, false, false, false, false, false, false, false, true, true, true, true, true},
		},
		{
			arithmeticV2,
			time.Second,
			[]bool{false, false, false, false, false, false, false, false, false, true, true, true, true, true},
		},
		{
			arithmeticV3,
			time.Second,
			[]bool{false, false, false, false, false, false, false, false, false, true, true, true, true, true},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test#%v", i+1), func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), test.timeout)
			defer cancel()

			coverage, err := GetCoverage(ctx, test.path)
			if err != nil {
				t.Errorf("GetCoverage(%#v): %s", test.path, err)
				return
			}

			pkg, err := packages.Load(test.path)
			if err != nil {
				t.Errorf("packages.Load(%#v): %s", test.path, err)
				return
			}

			var i int
			for node := range packages.Walk(ctx, pkg) {
				if i >= len(test.expect) {
					t.Errorf("expect: %v nodes, got: %v+", len(test.expect), i+1)
					return
				}

				if actual, expect := coverage.IsCover(node), test.expect[i]; actual != expect {
					t.Errorf("node #%v: IsCover(): %t, expect %t, source \"%s\"", i+1, actual, expect, node)
				}

				i++
			}

			if i != len(test.expect) {
				t.Errorf("expect: %v nodes, got: %v", len(test.expect), i)
			}
		})
	}
}
