package cli

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/params"
)

func checkFlags(t *testing.T, f *flag.FlagSet, argv interface{}) {
	t.Helper()

	reflectVal := reflect.Indirect(reflect.ValueOf(argv))
	reflectTyp := reflectVal.Type()

	for i := 0; i < reflectTyp.NumField(); i++ {
		tfield := reflectTyp.Field(i)
		vfield := reflectVal.Field(i)

		name, isFlag := tfield.Tag.Lookup("flag")
		if !isFlag {
			continue
		}

		fl := f.Lookup(name)
		if fl == nil {
			t.Fatalf("Expected there to be the flag %s\n", name)
		}

		if fl.Name != name {
			t.Errorf("Expected flag %s to have name %s but it has %s\n", name, name, fl.Name)
		}

		if fl.Usage != tfield.Tag.Get("descr") {
			t.Errorf("Expected flag %s to have usage %s but it has %s\n", name, tfield.Tag.Get("descr"), fl.Usage)
		}

		if tfield.Type.Kind() != reflect.Func { // Callbacks do not have a "default"
			if expdef := fmt.Sprintf("%v", vfield); fl.DefValue != expdef {
				t.Errorf("Expected flag %s to have default value %s but it has %s\n", name, expdef, fl.DefValue)
			}
		}
	}
}

func checkParams(t *testing.T, p *params.ParamSet, argv interface{}) {
	t.Helper()

	var paramsSeen int // keeps track of how many parameters we have seen so far

	reflectVal := reflect.Indirect(reflect.ValueOf(argv))
	reflectTyp := reflectVal.Type()

	for i := 0; i < reflectTyp.NumField(); i++ {
		tfield := reflectTyp.Field(i)
		// later I will need vfield := reflectVal.Field(i)

		name, isArg := tfield.Tag.Lookup("arg")
		if !isArg {
			continue
		}

		if tfield.Type.Kind() == reflect.Slice {
			// variadic
			param := p.Variadic()
			if param.Name != name {
				t.Errorf("Expected the parameter's name to be %s but it is %s", name, param.Name)
			}

			if param.Desc != tfield.Tag.Get("descr") {
				t.Errorf("Expected parameter %s's usage to be %s but it is %s", name, tfield.Tag.Get("descr"), param.Desc)
			}

			// If we are here, there couldn't have been a parsing error on min
			if min, _ := strconv.Atoi(tfield.Tag.Get("min")); param.Min != min {
				t.Errorf("Unexpected min on variadic variable, got %d but expected %d", param.Min, min)
			}
		} else {
			// non-variadic
			if paramsSeen > p.NParams() {
				t.Fatal("We have now seen more parameters in the spec than in the set")
			}

			param := p.Param(paramsSeen)
			paramsSeen++

			if param.Name != name {
				t.Errorf("Expected the parameter's name to be %s but it is %s", name, param.Name)
			}

			if param.Desc != tfield.Tag.Get("descr") {
				t.Errorf("Expected parameter %s's usage to be %s but it is %s", name, tfield.Tag.Get("descr"), param.Desc)
			}
		}
	}
}

func checkFlagsParams(t *testing.T, f *flag.FlagSet, p *params.ParamSet, argv interface{}) {
	t.Helper()
	checkFlags(t, f, argv)
	checkParams(t, p, argv)
}

func Test_prepareSpecs(t *testing.T) { //nolint:funlen // Test functions can be long...
	type args struct {
		f             *flag.FlagSet
		p             *params.ParamSet
		argv          interface{}
		allowVariadic bool
	}

	tests := []struct {
		name string
		args args
		err  error
		hook func(*testing.T, args)
	}{
		{
			name: "No flags or arguments",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct{}),
				true,
			},
		},
		{
			name: "String flag",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					var x struct {
						Foo string `flag:"foo" descr:"foobar"`
					}

					x.Foo = "qux"

					return &x
				}(),
				true,
			},
		},
		{
			name: "Int flag",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					var x struct {
						Foo int `flag:"foo" descr:"foobar"`
					}

					x.Foo = 42

					return &x
				}(),
				true,
			},
		},
		{
			name: "Int and string params",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					var x struct {
						Foo int    `pos:"foo" descr:"foo"`
						Bar string `pos:"bar" descr:"bar"`
					}

					x.Foo = 42
					x.Bar = "qux"

					return &x
				}(),
				true,
			},
		},
		{
			name: "All currently supported types",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					F1 int     `flag:"f1"`
					F2 bool    `flag:"f2"`
					F3 float64 `flag:"f3"`
					F4 string  `flag:"f4"`

					A1 int     `pos:"f1"`
					A2 bool    `pos:"f2"`
					A3 float64 `pos:"f3"`
					A4 string  `pos:"f4"`
				}),
				true,
			},
		},

		{
			name: "Unsupported flag type",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					B []bool `flag:"b"`
				}),
				true,
			},
			err: interfaces.SpecErrorf(`unsupported type for flag b: "slice"`),
		},
		{
			name: "Unsupported parameter type",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					B uintptr `pos:"b"`
				}),
				true,
			},
			err: interfaces.SpecErrorf(`unsupported type for parameter b: "uintptr"`),
		},

		{
			name: "Variadic bool",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					X []bool `pos:"x"`
				}),
				true,
			},
		},
		{
			name: "Variadic int",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					X []int `pos:"x" min:"2"`
				}),
				true,
			},
		},
		{
			name: "Variadic float",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					X []float64 `pos:"x"`
				}),
				true,
			},
		},
		{
			name: "Variadic string",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					X []string `pos:"x" descr:"foo"`
				}),
				true,
			},
		},
		{
			name: "Variadic with invalid min",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					X []string `pos:"x" descr:"foo" min:"not an int"`
				}),
				true,
			},
			err: interfaces.SpecErrorf(`unexpected min value for variadic parameter x: not an int`),
		},

		{
			name: "More than one variadic",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					A []int `pos:"a"`
					B []int `pos:"b"`
				}),
				true,
			},
			err: interfaces.SpecErrorf("a command spec cannot contain more than one variadic parameter"),
		},

		{
			name: "Flag callback nil",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					A func(string) error `flag:"a"`
				}),
				true,
			},
			err: interfaces.SpecErrorf("callbacks cannot be nil"),
		},
		{
			name: "Flag callback wrong signature 1",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					x := new(struct {
						A func(int) error `flag:"a"`
					})
					x.A = func(int) error { return nil }
					return x
				}(),
				true,
			},
			err: interfaces.SpecErrorf(`incorrect signature for callbacks: "func(int) error"`),
		},
		{
			name: "Flag callback wrong signature 2",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					x := new(struct {
						A func(string) `flag:"a"`
					})
					x.A = func(string) {}
					return x
				}(),
				true,
			},
			err: interfaces.SpecErrorf(`incorrect signature for callbacks: "func(string)"`),
		},
		{
			name: "Flag callbacks non-nil",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					var f = func(x string) error { return nil }
					var x = struct {
						A func(string) error `flag:"a"`
					}{A: f}
					return &x
				}(),
				true,
			},
			hook: func(t *testing.T, a args) {
				t.Helper()

				funcA := reflect.ValueOf(a).FieldByName("A")
				funcB := reflect.ValueOf(a).FieldByName("B")
				if funcA != funcB {
					t.Errorf("The callback function is no longer its default")
				}
			},
		},

		{
			name: "Param callback nil",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct {
					A func(string) error `pos:"a"`
				}),
				true,
			},
			err: interfaces.SpecErrorf("callbacks cannot be nil"),
		},
		{
			name: "Params callback wrong signature 1",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					x := new(struct {
						A func(int) error `pos:"a"`
					})
					x.A = func(int) error { return nil }
					return x
				}(),
				true,
			},
			err: interfaces.SpecErrorf(`incorrect signature for callbacks: "func(int) error"`),
		},
		{
			name: "Params callback wrong signature 2",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					x := new(struct {
						A func(string) `pos:"a"`
					})
					x.A = func(string) {}
					return x
				}(),
				true,
			},
			err: interfaces.SpecErrorf(`incorrect signature for callbacks: "func(string)"`),
		},
		{
			name: "Param callbacks non-nil",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					var f = func(x string) error { return nil }
					var x = struct {
						A func(string) error `pos:"a"`
					}{A: f}
					return &x
				}(),
				true,
			},
			hook: func(t *testing.T, a args) {
				t.Helper()

				funcA := reflect.ValueOf(a).FieldByName("A")
				funcB := reflect.ValueOf(a).FieldByName("B")
				if funcA != funcB {
					t.Errorf("The callback function is no longer its default")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &Command{flags: tt.args.f, params: tt.args.p}
			if err := connectSpecsFlagsAndParams(cmd, tt.args.argv); err != nil {
				if tt.err == nil {
					t.Fatalf("Got an error, but did not expect one. Got: %s", err.Error())
				}
				// FIXME: This is a bit vulnerable. I'm checking the string in the errors. I should
				// add parameters to the error type so I could check without expecting error messages
				// never to change.
				if err.Error() != tt.err.Error() {
					t.Fatalf("Unexpected error, expected %s but got %s", tt.err.Error(), err.Error())
				}
			} else {
				// No preparation error
				if tt.err != nil {
					t.Fatalf("Expected error %s here, but got nothing", tt.err.Error())
				}
				checkFlagsParams(t, tt.args.f, tt.args.p, tt.args.argv)
				if tt.hook != nil {
					tt.hook(t, tt.args)
				}
			}
		})
	}
}
