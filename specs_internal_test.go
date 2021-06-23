package cli

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/mailund/cli/params"
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

		if expdef := fmt.Sprintf("%v", vfield); fl.DefValue != expdef {
			t.Errorf("Expected flag %s to have default value %s but it has %s\n", name, expdef, fl.DefValue)
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

func checkFlagsParams(t *testing.T, f *flag.FlagSet, p *params.ParamSet, argv interface{}) {
	t.Helper()
	checkFlags(t, f, argv)
	checkParams(t, p, argv)
}

func Test_prepareSpecs(t *testing.T) { //nolint:funlen // Test functions can be long...
	type args struct {
		f    *flag.FlagSet
		p    *params.ParamSet
		argv interface{}
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "No flags or arguments",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				new(struct{}),
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
			},
		},
		{
			name: "Int and string params",
			args: args{
				flag.NewFlagSet("test", flag.ExitOnError),
				params.NewParamSet("test", flag.ExitOnError),
				func() interface{} {
					var x struct {
						Foo int    `arg:"foo" descr:"foo"`
						Bar string `arg:"bar" descr:"bar"`
					}

					x.Foo = 42
					x.Bar = "qux"

					return &x
				}(),
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

					A1 int     `arg:"f1"`
					A2 bool    `arg:"f2"`
					A3 float64 `arg:"f3"`
					A4 string  `arg:"f4"`
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepareSpecs(tt.args.f, tt.args.p, tt.args.argv)
			checkFlagsParams(t, tt.args.f, tt.args.p, tt.args.argv)
		})
	}
}
