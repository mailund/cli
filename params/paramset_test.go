package params_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/mailund/cli/internal/failure"
	"github.com/mailund/cli/params"
)

func TestMakeParamSet(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	if p.Name != "test" {
		t.Fatalf("New params set got the wrong name")
	}
}

func TestShortUsage(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	if p.ShortUsage() != "" {
		t.Errorf(`Short usage of empty paramset should be ""`)
	}

	p.String("foo", "")

	if p.ShortUsage() != "foo" {
		t.Errorf(`Short usage of with one "foo" parameter should be "foo"`)
	}

	p.String("bar", "")

	if p.ShortUsage() != "foo bar" {
		t.Errorf(`Short usage of with parameters "foo" and "bar" should be "foo bar"`)
	}

	p.VariadicString("...", "", 0)

	if p.ShortUsage() != "foo bar ..." {
		t.Errorf(`Short usage of with parameters "foo", "bar" and moreshould be "foo bar ..."`)
	}
}

func TestStringParam(t *testing.T) {
	var (
		x string
		y *string
	)

	p := params.NewParamSet("test", params.ExitOnError)
	p.StringVar(&x, "x", "")
	y = p.String("y", "")

	_ = p.Parse([]string{"foo", "bar"})

	if x != "foo" {
		t.Errorf(`Expected var x to hold "foo"`)
	}

	if *y != "bar" {
		t.Errorf(`Expected var *y to hold "bar"`)
	}
}

func TestPrintDefault(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)

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

	var x string

	p.StringVar(&x, "x", "an x")
	p.String("y", "a y")

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

func TestPrintDefaultMore(t *testing.T) {
	builder := new(strings.Builder)
	p := params.NewParamSet("test", params.ExitOnError)
	p.SetOutput(builder)

	var x string

	p.StringVar(&x, "x", "an x")
	p.String("y", "a y")
	p.VariadicString("...", "arguments for ...", 0)
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

func TestParseMore(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)

	p.String("x", "")
	rest := p.VariadicString("...", "arguments for ...", 0)

	args := []string{"for x", "y", "z"}
	_ = p.Parse(args)

	if !reflect.DeepEqual(*rest, []string{"y", "z"}) {
		t.Fatalf("The parser ate more than it should!")
	}
}

func TestFailure(t *testing.T) { //nolint:funlen // A test function can have as many statements as it likes
	p := params.NewParamSet("test", params.ExitOnError)

	p.String("x", "")

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
	p.VariadicString("y", "", 1)

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

	p := params.NewParamSet("test", params.ContinueOnError)

	builder := new(strings.Builder)
	p.SetOutput(builder)
	p.String("x", "")

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

	p := params.NewParamSet("test", params.ExitOnError)
	f := func(arg string) error {
		if arg != "arg" {
			t.Errorf("Got unexpected argument for callback")
		}

		called = true

		return nil
	}

	p.Func("foo", "", f)

	_ = p.Parse([]string{"arg"})

	if !called {
		t.Fatalf("Callback was not called")
	}
}

func TestFuncCallbackError(t *testing.T) {
	failed := false
	failure.Failure = func() { failed = true }

	p := params.NewParamSet("test", params.ExitOnError)
	builder := new(strings.Builder)
	p.SetOutput(builder)

	p.Func("foo", "",
		func(arg string) error {
			return errors.New("foo failed to bar") //nolint:goerr113 // Testing error
		})

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

	p := params.NewParamSet("test", params.ExitOnError)
	builder := new(strings.Builder)
	p.SetOutput(builder)

	p.VariadicFunc("foo", "", 0,
		func(args []string) error {
			return errors.New("foo failed to bar") //nolint:goerr113 // Testing error
		})

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

	p := params.NewParamSet("test", params.ExitOnError)
	ip := p.Int("i", "int")

	_ = p.Parse([]string{"42"})

	if *ip != 42 {
		t.Errorf("Parse error, ip is %d", *ip)
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

	p := params.NewParamSet("test", params.ExitOnError)
	bp := p.Bool("var", "")

	_ = p.Parse([]string{"1"})

	if !*bp {
		t.Errorf("Parse error, val is %t", *bp)
	}

	_ = p.Parse([]string{"0"})

	if *bp {
		t.Errorf("Parse error, val is %t", *bp)
	}

	_ = p.Parse([]string{"false"})

	if *bp {
		t.Errorf("Parse error, val is %t", *bp)
	}

	_ = p.Parse([]string{"true"})

	if !*bp {
		t.Errorf("Parse error, val is %t", *bp)
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

	p := params.NewParamSet("test", params.ExitOnError)
	x := p.Float("var", "")

	_ = p.Parse([]string{"3.14"})

	if *x != 3.14 {
		t.Errorf("Parse error, x is %f", *x)
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
	p := params.NewParamSet("test", params.ExitOnError)
	x := p.String("x", "")
	res := p.VariadicString("x [x...]", "", 0)
	args := []string{"x", "y", "z"}

	_ = p.Parse(args)

	if *x != "x" {
		t.Errorf("Argument x should be x, is %s", *x)
	}

	if !reflect.DeepEqual(args[1:], *res) {
		t.Errorf("Variadic argument got %v, expected [y, z]", *res)
	}
}

func TestVariadicBools(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	res := p.VariadicBool("x [x...]", "", 0)
	args := []string{"1", "true", "0", "false", "t", "FALSE"}
	expected := []bool{true, true, false, false, true, false}

	_ = p.Parse(args)

	if !reflect.DeepEqual(*res, expected) {
		t.Errorf("Variadic argument got %v, expected %v", *res, expected)
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

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, "cannot parse 'foo' as boolean.\n") {
		t.Errorf("Unexpected error message: '%s'\n", errmsg)
	}
}

func TestVariadicInts(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	res := p.VariadicInt("x [x...]", "", 0)
	args := []string{"1", "2", "3", "-4", "0x05", "6"}
	expected := []int{1, 2, 3, -4, 5, 6}

	_ = p.Parse(args)

	if !reflect.DeepEqual(*res, expected) {
		t.Errorf("Variadic argument got %v, expected %v", *res, expected)
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

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, "cannot parse 'foo' as integer.\n") {
		t.Errorf("Unexpected error message: '%s'\n", errmsg)
	}
}

func TestVariadicFloats(t *testing.T) {
	p := params.NewParamSet("test", params.ExitOnError)
	res := p.VariadicFloat("x [x...]", "", 0)
	args := []string{"1", "0.2", "3e4", "-4.1", "3.14"}
	expected := []float64{1.0, 0.2, 3e4, -4.1, 3.14}

	_ = p.Parse(args)

	if !reflect.DeepEqual(*res, expected) {
		t.Errorf("Variadic argument got %v, expected %v", *res, expected)
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

	if errmsg := builder.String(); !strings.HasSuffix(errmsg, "cannot parse 'foo' as float.\n") {
		t.Errorf("Unexpected error message: '%s'\n", errmsg)
	}
}
