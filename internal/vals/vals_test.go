package vals_test

import (
	"flag"
	"reflect"
	"strings"
	"testing"

	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/vals"
)

func simplifyUsage(x string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(x)), " ")
}

func TestBoolFlag(t *testing.T) {
	var b vals.BoolValue

	if !b.IsBoolFlag() {
		t.Fatal("boolean values should work as boolean flags")
	}
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

		// cannot be translated into a value
		m map[int]int

		pval interfaces.PosValue
		fval interfaces.FlagValue
	)

	pval = vals.AsPosValue(reflect.ValueOf(&i))
	if pval == nil {
		t.Fatal("We should have wrapped int")
	}

	fval = vals.AsFlagValue(reflect.ValueOf(&i))
	if fval == nil {
		t.Fatal("We should have wrapped int")
	}

	if cast, ok := fval.(*vals.IntValue); !ok {
		t.Error("Unexpected type for wrapped integer")
	} else if i != int(*cast) {
		t.Error("Expected value and i to hold the same value")
	}

	fval = vals.AsFlagValue(reflect.ValueOf(&iv))
	if fval == nil || !reflect.DeepEqual(fval, &iv) {
		t.Fatal("As pointer receiver, iv should implement the interface")
	}

	fval = vals.AsFlagValue(reflect.ValueOf(tv))
	if fval == nil || !reflect.DeepEqual(fval, tv) {
		t.Fatal("We should have gotten tv back unchanged")
	}

	pval = vals.AsPosValue(reflect.ValueOf(tv))
	if pval == nil || !reflect.DeepEqual(pval, tv) {
		t.Fatal("We should have gotten tv back unchanged")
	}

	fval = vals.AsFlagValue(reflect.ValueOf(&m))
	if fval != nil {
		t.Fatal("The map should not be wrapped")
	}

	pval = vals.AsPosValue(reflect.ValueOf(&m))
	if pval != nil {
		t.Fatal("The map should not be wrapped")
	}
}

func TestVariadic(t *testing.T) {
	ints := []int{}
	vi := (*vals.VariadicIntValue)(&ints)

	err := vi.Set([]string{"1", "2", "3"})
	if err != nil {
		t.Fatalf("Did not expect errors, but got %s", err.Error())
	}

	if len(*vi) != 3 {
		t.Fatal("Expected *vi to have three values now")
	}

	for i := 0; i < 3; i++ {
		if (*vi)[i] != i+1 {
			t.Errorf("Unexpected (*vi)[%d] = %d", i+1, (*vi)[i])
		}
	}
}

func TestVariadic2(t *testing.T) {
	ints := []int{}
	vi := vals.VariadicIntValueConstructor(reflect.ValueOf(&ints))

	if vi == nil {
		t.Fatal("Wrapping created nil")
	}

	vii, ok := vi.(*vals.VariadicIntValue)
	if !ok {
		t.Fatal("We should be able to cast to *vals.VariadicIntValue")
	}

	err := vi.Set([]string{"1", "2", "3"})
	if err != nil {
		t.Fatalf("Did not expect errors, but got %s", err.Error())
	}

	if len(*vii) != 3 {
		t.Fatal("Expected *vi to have three values now")
	}

	for i := 0; i < 3; i++ {
		if (*vii)[i] != i+1 {
			t.Errorf("Unexpected (*vi)[%d] = %d", i+1, (*vii)[i])
		}
	}
}

func TestVariadic3(t *testing.T) {
	ints := []int{}
	vi := vals.AsVariadicValue(reflect.ValueOf(&ints))

	if vi == nil {
		t.Fatal("Wrapping created nil")
	}

	vii, ok := vi.(*vals.VariadicIntValue)
	if !ok {
		t.Fatal("We should be able to cast to *vals.VariadicIntValue")
	}

	err := vi.Set([]string{"1", "2", "3"})
	if err != nil {
		t.Fatalf("Did not expect errors, but got %s", err.Error())
	}

	if len(*vii) != 3 {
		t.Fatal("Expected *vi to have three values now")
	}

	for i := 0; i < 3; i++ {
		if (*vii)[i] != i+1 {
			t.Errorf("Unexpected (*vi)[%d] = %d", i+1, (*vii)[i])
		}
	}
}

type VTestValue struct{}

func (t VTestValue) Set([]string) error { return nil }

func TestAsVariadicValue(t *testing.T) {
	var (
		// Implements the interface (if we use its address)
		iv vals.VariadicIntValue

		// doesn't implement the interface, so we will wrap it
		i = []int{}

		// implements the interface, so we should get it back unchanged
		tv VTestValue

		// cannot be translated into a value
		m map[int]int

		val interfaces.VariadicValue
	)

	val = vals.AsVariadicValue(reflect.ValueOf(&i))
	if val == nil {
		t.Fatal("We should have wrapped []int")
	}

	if cast, ok := val.(*vals.VariadicIntValue); !ok {
		t.Error("Unexpected type for wrapped integer")
	} else {
		j := (*[]int)(cast)
		if !reflect.DeepEqual(i, *j) {
			t.Error("Expected j and i to be equal")
		}
	}

	val = vals.AsVariadicValue(reflect.ValueOf(&iv))
	if val == nil || !reflect.DeepEqual(val, &iv) {
		t.Fatal("As pointer receiver, iv should implement the interface")
	}

	val = vals.AsVariadicValue(reflect.ValueOf(tv))
	if val == nil || !reflect.DeepEqual(val, tv) {
		t.Fatal("We should have gotten tv back unchanged")
	}

	val = vals.AsVariadicValue(reflect.ValueOf(&m))
	if val != nil {
		t.Fatal("The map should not be wrapped")
	}
}
