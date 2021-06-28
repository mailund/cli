package vals_test

import (
	"flag"
	"reflect"
	"strings"
	"testing"

	"github.com/mailund/cli/inter"
	"github.com/mailund/cli/internal/vals"
)

func simplifyUsage(x string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(x)), " ")
}

func TestValsFlags(t *testing.T) {
	var (
		i vals.IntValue    = 42
		s vals.StringValue = "foo"
		f                  = flag.NewFlagSet("test", flag.ContinueOnError)
	)

	f.Var(&i, "i", "an integer")
	f.Var(&s, "s", "a string")

	builder := new(strings.Builder)
	f.SetOutput(builder)
	f.PrintDefaults()

	usage := simplifyUsage(builder.String())
	expected := strings.Join(
		[]string{
			"-i value an integer (default 42)",
			"-s value a string (default foo)",
		}, " ")

	if usage != expected {
		t.Errorf("Unexpected usage string:\n'%s'\n", usage)
	}

	err := f.Parse([]string{"-i=13", "-s", "bar"})

	if err != nil {
		t.Fatal("Didn't expect parsing to return an error")
	}

	if i != 13 {
		t.Errorf("Expected i to have the value 13, it has %d", i)
	}

	if s != "bar" {
		t.Errorf("Expected s to have the value bar, it has %s", s)
	}
}

type TestValue struct{}

func (t TestValue) String() string   { return "" }
func (t TestValue) Set(string) error { return nil }

func TestAsValue(t *testing.T) {
	var (
		// Implements the interface (if we use its address)
		iv vals.IntValue

		// doesn't implement the interface, so we will wrap it
		i = 42

		// implements the interface, so we should get it back unchanged
		tv TestValue

		val inter.FlagValue
	)

	val = vals.AsValue(reflect.ValueOf(&i))
	if val == nil {
		t.Fatal("We should have wrapped int")
	}

	if cast, ok := val.(*vals.IntValue); !ok {
		t.Error("Unexpected type for wrapped integer")
	} else if i != int(*cast) {
		t.Error("Expected value and i to hold the same value")
	}

	val = vals.AsValue(reflect.ValueOf(&iv))
	if val == nil || !reflect.DeepEqual(val, &iv) {
		t.Fatal("As pointer receiver, iv should implement the interface")
	}

	val = vals.AsValue(reflect.ValueOf(tv))
	if val == nil || !reflect.DeepEqual(val, tv) {
		t.Fatal("We should have gotten tv back unchanged")
	}
}
