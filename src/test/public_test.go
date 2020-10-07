package test

import (
	"context"
	"fmt"
	"github.com/libmonsoon-dev/gomut/src/testutil"
	"github.com/libmonsoon-dev/jsontest"
	"golang.org/x/tools/cover"
	"reflect"
	"testing"
)

func TestGetCoverage(t *testing.T) {
	type Args struct {
		ctx  context.Context
		path string
	}
	tests := []struct {
		args    Args
		want    Profiles
		wantErr bool
	}{
		{
			args: Args{
				nil,
				arithmeticV1,
			},
			want: NewProfiles([]*cover.Profile{
				{
					FileName: "github.com/libmonsoon-dev/gomut/testdata/arithmetic/v1/sum.go",
					Mode:     "set",
					Blocks:   []cover.ProfileBlock{{StartLine: 3, StartCol: 24, EndLine: 5, EndCol: 2, NumStmt: 1, Count: 1}},
				},
			}),
		},
		{
			args: Args{
				nil,
				arithmeticV2,
			},
			want: NewProfiles([]*cover.Profile{{
				FileName: "github.com/libmonsoon-dev/gomut/testdata/arithmetic/v2/sum.go",
				Mode:     "set",
				Blocks:   []cover.ProfileBlock{{StartLine: 3, StartCol: 24, EndLine: 5, EndCol: 2, NumStmt: 1, Count: 1}},
			}}),
		},
		{
			args: Args{
				nil,
				arithmeticV3,
			},
			want: NewProfiles([]*cover.Profile{{
				FileName: "github.com/libmonsoon-dev/gomut/testdata/arithmetic/v3/sum.go",
				Mode:     "set",
				Blocks:   []cover.ProfileBlock{{StartLine: 3, StartCol: 24, EndLine: 5, EndCol: 2, NumStmt: 1, Count: 1}},
			}}),
		},
	}
	for i, test := range tests {
		test := test

		t.Run(fmt.Sprintf("test#%v", i+1), func(t *testing.T) {
			t.Parallel()

			got, err := GetCoverage(test.args.ctx, test.args.path)
			if (err != nil) != test.wantErr {
				t.Errorf("GetCoverage() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if len(got) != len(test.want) {
				t.Errorf("Coverage profiles length not equal: got: %v want: %v", len(got), len(test.want))
			}
			for i := 0; i < len(got); i++ {
				if !reflect.DeepEqual(got, test.want) {
					t.Errorf("GetCoverage(): \ngot = %#v\nwant %#v", got, test.want)
				}
			}
		})
	}
}

func TestRunTest(t *testing.T) {
	tests := []struct {
		arg  Config
		want []Event
	}{
		{
			arg: Config{Paths: []string{arithmeticV1}},
			want: []Event{
				{
					Event: jsontest.Event{
						Action:  jsontest.Run,
						Package: arithmeticV1,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV1,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV1,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Pass,
						Package: arithmeticV1,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV1,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV1,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Pass,
						Package: arithmeticV1,
						Test:    "TestAdd",
					},
				},
			},
		},
		{
			arg: Config{Paths: []string{arithmeticV2}},
			want: []Event{
				{
					Event: jsontest.Event{
						Action:  jsontest.Run,
						Package: arithmeticV2,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV2,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV2,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Pass,
						Package: arithmeticV2,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV2,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV2,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Pass,
						Package: arithmeticV2,
						Test:    "TestAdd",
					},
				},
			},
		},
		{
			arg: Config{Paths: []string{arithmeticV3}},
			want: []Event{
				{
					Event: jsontest.Event{
						Action:  jsontest.Run,
						Package: arithmeticV3,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV3,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV3,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Pass,
						Package: arithmeticV3,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV3,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Output,
						Package: arithmeticV3,
						Test:    "TestAdd",
					},
				},
				{
					Event: jsontest.Event{
						Action:  jsontest.Pass,
						Package: arithmeticV3,
						Test:    "TestAdd",
					},
				},
			},
		},
		{
			arg: Config{
				Paths: []string{"./notExist"},
			},
			want: []Event{{
				Err: fmt.Errorf(
					"test command error: exit status 1: stat %v/src/test/notExist: directory not found",
					testutil.ProjectPath(),
				),
			}},
		},
	}

	for i, test := range tests {
		test := test

		t.Run(fmt.Sprintf("test#%v", i+1), func(t *testing.T) {
			t.Parallel()

			var i int
			for msg := range Run(test.arg) {
				if msg.Err != nil && test.want[i].Err != nil {
					if msg.Err.Error() != test.want[i].Err.Error() {
						t.Errorf("Run(): %v\nwantErr: %v", msg.Err, test.want[i].Err)
					}

				} else if msg.Err != nil {
					t.Errorf("Unexpected error in event #%v: %v", i, msg.Err)
				} else if !eventEqual(msg.Event, test.want[i].Event) {
					t.Errorf("Message #%v\nevent = %#v, \nexpect = %#v", i+1, msg.Event, test.want[i])
				}
				i++
			}
			if i != len(test.want) {
				t.Errorf("Expect %v messages, but got %v", len(test.want), i)
			}

		})
	}
}

func eventEqual(a, b jsontest.Event) bool {
	if a.Action != b.Action {
		return false
	}
	if a.Package != b.Package {
		return false
	}
	if a.Test != b.Test {
		return false
	}

	return true
}
