package rand

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	type Test struct {
		Length int
	}
	tests := []Test{
		{},
		{1},
		{10},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test#%v", i+1), func(t *testing.T) {
			if got := String(test.Length); len(got) != test.Length {
				t.Errorf("String() = %v (len = %v), expect: len %v", got, len(got), test.Length)
			}
		})
	}
}
