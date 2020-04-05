package arithmetic

import (
	"testing"
)

type testCase struct {
	a,
	b,
	want int
}

func TestAdd(t *testing.T) {
	test := testCase{3, 4, 7}

	if got := Add(test.a, test.b); got != test.want {
		t.Errorf("Add(%v, %v) = %v, want %v", test.a, test.b, got, test.want)
	}
}
