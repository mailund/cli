package cli

import (
	"flag"
	"fmt"
	"reflect"

	"github.com/mailund/cli/params"
)

// TODO: add variadic and Func arguments.
// TODO: connect specs with commands

func prepareSpecs(f *flag.FlagSet, p *params.ParamSet, argv interface{}) {
	reflectVal := reflect.Indirect(reflect.ValueOf(argv))
	reflectTyp := reflectVal.Type()

	for i := 0; i < reflectTyp.NumField(); i++ {
		tfield := reflectTyp.Field(i)
		vfield := reflectVal.Field(i)

		if name, isFlag := tfield.Tag.Lookup("flag"); isFlag {
			switch tfield.Type.Kind() {
			case reflect.Bool:
				f.BoolVar(vfield.Addr().Interface().(*bool), name, vfield.Bool(), tfield.Tag.Get("descr"))

			case reflect.Int:
				f.IntVar(vfield.Addr().Interface().(*int), name, int(vfield.Int()), tfield.Tag.Get("descr"))

			case reflect.Float64:
				f.Float64Var(vfield.Addr().Interface().(*float64), name, vfield.Float(), tfield.Tag.Get("descr"))

			case reflect.String:
				f.StringVar(vfield.Addr().Interface().(*string), name, vfield.String(), tfield.Tag.Get("descr"))

			default:
				fmt.Printf("Unknown type %q", tfield.Type.Kind())
			}
		}

		if name, isArg := tfield.Tag.Lookup("arg"); isArg {
			switch tfield.Type.Kind() {
			case reflect.Bool:
				p.BoolVar(vfield.Addr().Interface().(*bool), name, tfield.Tag.Get("descr"))

			case reflect.Int:
				p.IntVar(vfield.Addr().Interface().(*int), name, tfield.Tag.Get("descr"))

			case reflect.Float64:
				p.FloatVar(vfield.Addr().Interface().(*float64), name, tfield.Tag.Get("descr"))

			case reflect.String:
				p.StringVar(vfield.Addr().Interface().(*string), name, tfield.Tag.Get("descr"))

			case reflect.Slice:
				switch tfield.Type.Elem().Kind() {
				case reflect.Int:
					fmt.Println("Slice of int")
				case reflect.String:
					fmt.Println("slice of strings")

				default:
					fmt.Printf("Unknown slice type %q", tfield.Type.Elem().Kind())
				}

			default:
				fmt.Printf("Unknown type %q", tfield.Type.Kind())
			}
		}
	}
}
