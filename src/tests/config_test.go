package tests

import (
	"context"
	"fmt"
	"path"
	"reflect"
	"testing"
	"time"
)

func TestGetCommand(t *testing.T) {
	type Key int

	const (
		KeyA Key = iota
		KeyB
	)

	timeoutCtx, cancel1 := context.WithTimeout(context.Background(), time.Second)
	defer cancel1()
	cancelCtx, cancel2 := context.WithCancel(context.Background())
	defer cancel2()
	deadlineCtx, cancel3 := context.WithDeadline(context.Background(), time.Now().Add(time.Second))
	defer cancel3()

	tests := []Config{
		{
			Path: "fmt",
		},
		{
			Path:           "fmt",
			BuildTestFlags: []string{"-coverprofile=cover.out"},
		},
		{
			Path:            "fmt",
			TestBinaryFlags: []string{"-test.short"},
		},
		{
			Path:            "fmt",
			BuildTestFlags:  []string{"-coverprofile=cover.out"},
			TestBinaryFlags: []string{"-test.short"},
		},
		{
			Path:  "fmt",
			GoBin: "go1.14.1",
		},
		{
			Path:  "fmt",
			GoBin: "go",
		},
		{
			Path: "fmt",
			Ctx:  context.Background(),
		},
		{
			Path: "fmt",
			Ctx:  context.TODO(),
		},
		{
			Path: "fmt",
			Ctx:  context.WithValue(context.Background(), KeyA, "value"),
		},
		{
			Path: "fmt",
			Ctx:  timeoutCtx,
		},
		{
			Path: "fmt",
			Ctx:  cancelCtx,
		},
		{
			Path: "fmt",
			Ctx:  deadlineCtx,
		},
		{
			Path:            "fmt",
			GoBin:           "go1.14.1",
			BuildTestFlags:  []string{"-coverprofile=cover.out"},
			TestBinaryFlags: []string{"-test.short"},
			Ctx:             context.WithValue(context.Background(), KeyB, "value"),
		},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("test#%v", i+1), func(t *testing.T) {
			expectedCtx := test.Ctx
			cmd := test.GetCommand()

			ctxReflectValue := reflect.ValueOf(cmd).Elem().FieldByName("ctx")

			if test.Ctx == nil {
				if ctxReflectValue.IsNil() {
					t.Error("cmd.ctx is nil, expect: ctx.Background()")
					return
				}
				expectedCtx = context.Background()
			}
			actualCtxType := ctxReflectValue.Elem().Elem().Type()
			expectType := reflect.ValueOf(expectedCtx).Elem().Type()
			if actualCtxType != expectType {
				t.Errorf("cmd.ctx type not equal ctx.Background() (%s != %s)", actualCtxType, expectType)
				return
			}

			goBin := test.GoBin
			if goBin == "" {
				goBin = "go"
			}

			if path.Base(cmd.Path) != goBin {
				t.Errorf("cmd.Path not match goBin, (%#v, %#v)", cmd.Path, goBin)
				return
			}

			if path.Base(cmd.Args[0]) != goBin {
				t.Errorf("cmd.Args[0] not match goBin, (%#v, %#v)", cmd.Args[0], goBin)
				return
			}

			if cmd.Args[1] != "test" {
				t.Errorf("cmd.Args[1] \"test\" (cmd.Args[1] %#v)", cmd.Args[1])
				return
			}

			if cmdJsonFlag := cmd.Args[2]; cmdJsonFlag != jsonFlag {
				t.Errorf(" cmdJsonFlag != jsonFlag (%#v != %#v)", cmdJsonFlag, jsonFlag)
				return
			}

			buildTestFlagsStart := 3
			buildTestFlagsEnd := len(test.BuildTestFlags) + buildTestFlagsStart
			if actual := cmd.Args[buildTestFlagsStart:buildTestFlagsEnd]; !stringsEqual(actual, test.BuildTestFlags) {
				t.Errorf(
					"cmd.Args[1:buildTestFlagsLen+1] != config.BuildTestFlags (%v != %v)",
					actual,
					test.BuildTestFlags,
				)
				return
			}

			cmdPathIndex := buildTestFlagsEnd
			if cmdPath := cmd.Args[cmdPathIndex]; cmdPath != test.Path {
				t.Errorf("cmdPath != test.Path (%#v != %#v)", cmdPath, test.Path)
				return
			}

			cmdTestBinaryFlagsStart := cmdPathIndex + 1
			cmdTestBinaryFlags := cmd.Args[cmdTestBinaryFlagsStart:]
			if !stringsEqual(cmdTestBinaryFlags, test.TestBinaryFlags) {
				t.Errorf(
					"cmd TestBinaryFlags != config.TestBinaryFlags (%#v != %#v)",
					cmdTestBinaryFlags,
					test.TestBinaryFlags,
				)
			}

		})
	}
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
