package params_test

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/mailund/cli/internal/failure"
	"github.com/mailund/cli/internal/params"
	"github.com/mailund/cli/internal/vals"
)

func TestMakeParamSet(t *testing.T) {
	p := params.NewParamSet("test", failure.ExitOnError)
	if p.Name != "test" {
		t.Fatalf("New params set got the wrong name")
	}
}

func TestShortUsage(t *testing.T) {
	p := params.NewParamSet("test", failure.ExitOnError)
	if p.ShortUsage() != "" {
		t.Errorf(`Short usage of empty paramset should be ""`)
	}

	s := "foo"
	p.Var((*vals.StringValue)(&s), "foo", "")

	if p.ShortUsage() != "foo" {
		t.Errorf(`Short usage of with one "foo" parameter should be "foo"`)
	}

	p.Var((*vals.StringValue)(&s), "bar", "")

	if p.ShortUsage() != "foo bar" {
		t.Errorf(`Short usage of with parameters "foo" and "bar" should be "foo bar"`)
	}

	// just testing that I can get the vals.VariadicStringValue this way...
	x := vals.AsVariadicValue(reflect.ValueOf(&([]string{})))
	p.VariadicVar(x, "...", "", 0)

	if p.ShortUsage() != "foo bar ..." {
		t.Errorf(`Short usage of with parameters "foo", "bar" and moreshould be "foo bar ..."`)
	}
}

func TestStringParam(t *testing.T) {
	var (
		x string
		y string
	)

	p := params.NewParamSet("test", failure.ExitOnError)
	p.Var((*vals.StringValue)(&x), "x", "")
	p.Var((*vals.StringValue)(&y), "y", "")

	_ = p.Parse([]string{"foo", "bar"})

	if x != "foo" {
		t.Errorf(`Expected var x to hold "foo"`)
	}

	if y != "bar" {
		t.Errorf(`Expected var y to hold "bar"`)
	}
}

func TestPrintDefault(t *testing.T) {
	p := params.NewParamSet("test", failure.ExitOnError)

	// without any arguments, it shouldn't print anything
	builder := new(strings.Builder)
	p.SetOutput(builder)
	p.PrintDefaults()

	res := builder.String()

	if res != "" {
		t.Errorf("Unexpected usage output: %s\n", res)
	}

	builder = new(strings.Builder)
	p.SetOutput(builder)

	var x, y string

	p.Var((*vals.StringValue)(&x), "x", "an x")
	p.Var((*vals.StringValue)(&y), "y", "a y")

	p.PrintDefaults()

	res = builder.String()
	expected := `Arguments:
  x
	an x
  y
	a y
`

	if res != expected {
		t.Fatalf("PrintDefault gave us \n----\n%s\n----\n\nbut we exected\n\n----\n%s\n----",
			res, expected)
	}
}

func TestPrintDefaultVariadic(t *testing.T) {
	builder := new(strings.Builder)
	p := params.NewParamSet("test", failure.ExitOnError)
	p.SetOutput(builder)

	var (
		x, y vals.StringValue
		z    = []string{}
	)

	p.Var(&x, "x", "an x")
	p.Var(&y, "y", "a y")
	p.VariadicVar((*vals.VariadicStringValue)(&z), "...", "arguments for ...", 0)
	p.PrintDefaults()

	expected := `Arguments:
  x
	an x
  y
	a y
  ...
	arguments for ...
`

	if res := builder.String(); res != expected {
		t.Fatalf("PrintDefault gave us \n----\n%s\n----\n\nbut we exected\n\n----\n%s\n----",
			res, expected)
	}
}

func TestParseVariadic(t *testing.T) {
	p := params.NewParamSet("test", failure.ExitOnError)

	var (
		x    vals.StringValue
		rest = []string{}
	)

	p.Var(&x, "x", "")
	p.VariadicVar((*vals.VariadicStringValue)(&rest), "...", "arguments for ...", 0)

	args := []string{"for x", "y", "z"}
	_ = p.Parse(args)

	if !reflect.DeepEqual(rest, []string{"y", "z"}) {
		t.Fatalf("The parser ate more than it should!")
	}
}

func TestFailure(t *testing.T) { //nolint:funlen // A test function can have as many statements as it likes
	p := params.NewParamSet("test", failure.ExitOnError)
	x := vals.StringValue("")

	p.Var(&x, "x", "")

	var failed = false

	failure.Failure = func() { failed = true }
	builder := new(strings.Builder)

	p.SetOutput(builder)

	err := p.Parse([]string{})

	if !failed {
		t.Errorf("We expected a parse failure")
	}

	if errmsg := builder.String(); !strings.HasPrefix(errmsg, "Too few arguments") {
		t.Errorf("We got an unexpected error message: %s", errmsg)
	}

	if err == nil {
		t.Fatal("expected an error from the parse error")
	}

	if err.Error() != "too few arguments" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	failed = false
	builder = new(strings.Builder)
	p.SetOutput(builder)
	err = p.Parse([]string{"x", "y"})

	if !failed {
		t.Errorf("We expected a parse failure")
	}

	if errmsg := builder.String(); !strings.HasPrefix(errmsg, "Too many arguments") {
		t.Errorf("We got an unexpected error message: %s", errmsg)
	}

	if err == nil {
		t.Fatal("expected an error from the parse error")
	}

	if err.Error() != "too many arguments" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	// Add a variadic that wants at least one argument
	y := []string{}

	p.VariadicVar((*vals.VariadicStringValue)(&y), "y", "", 1)

	failed = false
	builder = new(strings.Builder)

	p.SetOutput(builder)

	err = p.Parse([]string{"x"}) // too few

	if !failed {
		t.Errorf("We expected a parse failure")
	}

	if errmsg := builder.String(); !strings.HasPrefix(errmsg, "Too few arguments") {
		t.Errorf("We got an unexpected error message: %s", errmsg)
	}

	if err == nil {
		t.Fatal("expected an error from the parse error")
	}

	if err.Error() != "too few arguments" {
		t.Errorf("Unexpected error message: %s", err.Error())
	}

	failed = false
	err = p.Parse([]string{"x", "y"}) // now we have the minimum

	if failed {
		t.Errorf("We shouldn't fail")
	}

	if err != nil {
		t.Fatal("We shouldn't get an error now")
	}

	err = p.Parse([]string{"x", "y", "z", "w"}) // and we can have more

	if failed {
		t.Errorf("We shouldn't fail")
	}

	if err != nil {
		t.Fatal("We shouldn't get an error now")
	}
}

func TestFailureContinue(t *testing.T) {
	var failed = false

	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", failure.ContinueOnError)
	x := vals.StringValue("")

	builder := new(strings.Builder)
	p.SetOutput(builder)
	p.Var(&x, "x", "")

	if err := p.Parse([]string{}); err == nil {
		t.Fatalf("expected an error from the parse error")
	}

	if errmsg := builder.String(); errmsg != "" {
		t.Fatalf("We shouldn't have an error message (but have %s)", errmsg)
	}

	if failed {
		t.Fatal("We shouldn't terminate now")
	}
}

func TestFailureSetFlag(t *testing.T) {
	var failed = false

	failure.Failure = func() { failed = true }

	// create a paramset that will crash on errors
	p := params.NewParamSet("test", failure.ExitOnError)
	x := vals.StringValue("")

	// but then change the flag
	p.SetErrFlag(failure.ContinueOnError)

	builder := new(strings.Builder)
	p.SetOutput(builder)
	p.Var(&x, "x", "")

	if err := p.Parse([]string{}); err == nil {
		t.Fatalf("expected an error from the parse error")
	}

	if errmsg := builder.String(); errmsg != "" {
		t.Fatalf("We shouldn't have an error message (but have %s)", errmsg)
	}

	if failed {
		t.Fatal("We shouldn't terminate now")
	}
}

func TestFuncCallback(t *testing.T) {
	var called = false

	p := params.NewParamSet("test", failure.ExitOnError)
	f := func(arg string) error {
		if arg != "arg" {
			t.Errorf("Got unexpected argument for callback")
		}

		called = true

		return nil
	}

	p.Var(vals.FuncValue(f), "foo", "")

	_ = p.Parse([]string{"arg"})

	if !called {
		t.Fatalf("Callback was not called")
	}
}

func TestFuncCallbackError(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", failure.ExitOnError)
	builder := new(strings.Builder)
	p.SetOutput(builder)

	f := vals.FuncValue(
		func(arg string) error {
			return errors.New("foo failed to bar") //nolint:goerr113 // Testing error
		})

	p.Var(f, "foo", "")

	if err := p.Parse([]string{"arg"}); err == nil {
		t.Fatalf("Expected an error")
	}

	if !failed {
		t.Errorf("Expected parse error to trigger a failure")
	}

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, "foo failed to bar.\n") {
		t.Errorf("Unexpected error message %s", errmsg)
	}
}

func TestVariadicFuncError(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", failure.ExitOnError)
	builder := new(strings.Builder)
	p.SetOutput(builder)

	var (
		f = func(args []string) error {
			return errors.New("foo failed to bar") //nolint:goerr113 // Testing error
		}
		vf = vals.VariadicFuncValue(f)
	)

	p.VariadicVar(vf, "foo", "", 0)

	if err := p.Parse([]string{"arg"}); err == nil {
		t.Fatalf("Expected an error")
	}

	if !failed {
		t.Errorf("Expected parse error to trigger a failure")
	}

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, "foo failed to bar.\n") {
		t.Errorf("Unexpected error message %s", errmsg)
	}
}

func TestInt(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", failure.ExitOnError)
	i := vals.IntValue(0)
	p.Var(&i, "i", "int")

	_ = p.Parse([]string{"42"})

	if i != 42 {
		t.Errorf("Parse error, i is %d", i)
	}

	builder := new(strings.Builder)
	p.SetOutput(builder)

	if err := p.Parse([]string{"foo"}); err == nil {
		t.Error("Expected an error")
	}

	if !failed {
		t.Error("Expected to fail")
	}

	if errmsg := builder.String(); !strings.HasPrefix(errmsg, `Error parsing parameter i='foo'`) {
		t.Errorf("Unexpected error message %s", errmsg)
	}
}

func TestBool(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", failure.ExitOnError)
	b := vals.BoolValue(false)
	p.Var(&b, "var", "")

	_ = p.Parse([]string{"1"})

	if !b {
		t.Errorf("Parse error, val is %t", b)
	}

	_ = p.Parse([]string{"0"})

	if b {
		t.Errorf("Parse error, val is %t", b)
	}

	_ = p.Parse([]string{"false"})

	if b {
		t.Errorf("Parse error, val is %t", b)
	}

	_ = p.Parse([]string{"true"})

	if !b {
		t.Errorf("Parse error, val is %t", b)
	}

	builder := new(strings.Builder)
	p.SetOutput(builder)

	if err := p.Parse([]string{"foo"}); err == nil {
		t.Error("Expected an error")
	}

	if !failed {
		t.Error("Expected to fail")
	}

	if errmsg := builder.String(); !strings.HasPrefix(errmsg, `Error parsing parameter var='foo'`) {
		t.Errorf("Unexpected error message %s", errmsg)
	}
}

func TestFloat(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", failure.ExitOnError)
	x := vals.Float64Value(0.0)
	p.Var(&x, "var", "")

	_ = p.Parse([]string{"3.14"})

	if x != 3.14 {
		t.Errorf("Parse error, x is %f", x)
	}

	builder := new(strings.Builder)
	p.SetOutput(builder)

	if err := p.Parse([]string{"foo"}); err == nil {
		t.Error("Expected an error")
	}

	if !failed {
		t.Error("Expected to fail")
	}

	if errmsg := builder.String(); !strings.HasPrefix(errmsg, `Error parsing parameter var='foo'`) {
		t.Errorf("Unexpected error message %s", errmsg)
	}
}

func TestVariadicStrings(t *testing.T) {
	p := params.NewParamSet("test", failure.ExitOnError)
	x := vals.StringValue("")
	args := []string{"x", "y", "z"}
	res := []string{}

	p.Var(&x, "x", "")
	p.VariadicVar((*vals.VariadicStringValue)(&res), "x [x...]", "", 0)

	_ = p.Parse(args)

	if x != "x" {
		t.Errorf("Argument x should be x, is %s", x)
	}

	if !reflect.DeepEqual(args[1:], res) {
		t.Errorf("Variadic argument got %v, expected [y, z]", res)
	}
}

func TestVariadicBools(t *testing.T) {
	p := params.NewParamSet("test", failure.ExitOnError)
	res := []bool{}

	p.VariadicVar((*vals.VariadicBoolValue)(&res), "x [x...]", "", 0)

	args := []string{"1", "true", "0", "false", "t", "FALSE"}
	expected := []bool{true, true, false, false, true, false}

	_ = p.Parse(args)

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Variadic argument got %v, expected %v", res, expected)
	}

	// Testing errors
	failed := false
	failure.Failure = func() { failed = true }
	builder := new(strings.Builder)
	p.SetOutput(builder)
	_ = p.Parse([]string{"foo", "bar"})

	if !failed {
		t.Error("Expected a parser failure")
	}

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, "cannot parse 'foo' as bool.\n") {
		t.Errorf("Unexpected error message: '%s'\n", errmsg)
	}
}

func TestVariadicInts(t *testing.T) {
	p := params.NewParamSet("test", failure.ExitOnError)
	res := []int{}

	p.VariadicVar((*vals.VariadicIntValue)(&res), "x [x...]", "", 0)

	args := []string{"1", "2", "3", "-4", "0x05", "6"}
	expected := []int{1, 2, 3, -4, 5, 6}

	_ = p.Parse(args)

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Variadic argument got %v, expected %v", res, expected)
	}

	// Testing errors
	failed := false
	failure.Failure = func() { failed = true }
	builder := new(strings.Builder)
	p.SetOutput(builder)

	_ = p.Parse([]string{"foo", "bar"})

	if !failed {
		t.Error("Expected a parser failure")
	}

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, "cannot parse 'foo' as int.\n") {
		t.Errorf("Unexpected error message: '%s'\n", errmsg)
	}
}

func TestVariadicFloats(t *testing.T) {
	p := params.NewParamSet("test", failure.ExitOnError)
	res := []float64{}

	p.VariadicVar((*vals.VariadicFloat64Value)(&res), "x [x...]", "", 0)

	args := []string{"1", "0.2", "3e4", "-4.1", "3.14"}
	expected := []float64{1.0, 0.2, 3e4, -4.1, 3.14}

	_ = p.Parse(args)

	if !reflect.DeepEqual(res, expected) {
		t.Errorf("Variadic argument got %v, expected %v", res, expected)
	}

	// Testing errors
	failed := false
	failure.Failure = func() { failed = true }
	builder := new(strings.Builder)
	p.SetOutput(builder)
	_ = p.Parse([]string{"foo", "bar"})

	if !failed {
		t.Error("Expected a parser failure")
	}

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, "cannot parse 'foo' as float64.\n") {
		t.Errorf("Unexpected error message: '%s'\n", errmsg)
	}
}

func TestParamVariadic(t *testing.T) {
	p := params.NewParamSet("p", failure.ExitOnError)

	var (
		i vals.IntValue
		b vals.BoolValue
	)

	p.Var(&i, "i", "int")
	p.Var(&b, "b", "bool")

	res := []string{}

	p.VariadicVar((*vals.VariadicStringValue)(&res), "s", "strings", 0)

	if p.NParams() != 2 {
		t.Fatalf("Expected two paramteres, but paramset says there are %d", p.NParams())
	}

	first := p.Param(0)
	second := p.Param(1)

	if first.Name != "i" {
		t.Error("first parameter name should be 'i'")
	}

	if first.Desc != "int" {
		t.Error("first param descr should be 'int'")
	}

	if second.Name != "b" {
		t.Error("second parameter name should be 'b'")
	}

	if second.Desc != "bool" {
		t.Error("second param descr should be 'bool'")
	}

	if p.Variadic() == nil {
		t.Fatal("Expected there to be a variadic argument")
	}

	variadic := p.Variadic()

	if variadic.Name != "s" {
		t.Error("variadic argument should be 's'")
	}

	if variadic.Desc != "strings" {
		t.Error("variadic description should be 'strings'")
	}
}

type Choice struct {
	Choice  string
	Options []string
}

// We don't use Set() for actual parsing in the text...
func (c *Choice) Set(x string) error { return nil }
func (c *Choice) String() string     { return c.Choice }

func (c *Choice) ArgumentDescription(flag bool, descr string) string {
	return descr + " (choose from " + "{" + strings.Join(c.Options, ",") + "})"
}

func TestDescriptionHooks(t *testing.T) {
	var (
		a = Choice{"A", []string{"A", "B", "C"}}
		p = params.NewParamSet("test", failure.ExitOnError)
	)

	p.Var(&a, "foo", "a choice")

	builder := new(strings.Builder)
	p.SetOutput(builder)

	p.Usage()

	expected := `Arguments:
	foo a choice (choose from {A,B,C})
	`
	msg := builder.String()

	space := regexp.MustCompile(`\s+`)
	msg = space.ReplaceAllString(msg, " ")
	expected = space.ReplaceAllString(expected, " ")

	if msg != expected {
		t.Errorf("unexpected: %s", msg)
		t.Errorf("unexpected: %s", expected)
	}
}
