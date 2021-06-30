package flags_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mailund/cli/interfaces"
	"github.com/mailund/cli/internal/failure"
	"github.com/mailund/cli/internal/flags"
	"github.com/mailund/cli/internal/vals"
)

func TestBasicInt(t *testing.T) { //nolint:funlen,gocyclo // test functions are just sometimes long
	f := flags.NewFlagSet("test", failure.ExitOnError)

	builder := new(strings.Builder)
	f.SetOutput(builder)

	if f.Output() != builder {
		t.Error("somehow the output wasn't set")
	}

	f.Usage()

	expected := ""
	if usageMsg := builder.String(); usageMsg != expected {
		t.Errorf("Unexpected usage: %s", usageMsg)
	}

	var i vals.IntValue

	if err := f.Var(&i, "int", "i", "integer"); err != nil {
		t.Fatal("error inserting an integer")
	}

	if f.Lookup("i").Value != &i {
		t.Error("Expected i to be i")
	}

	if f.Lookup("int").Value != &i {
		t.Error("Expected int to be i as well")
	}

	if f.NFlags() != 1 {
		t.Error("there should be one flag now")
	}

	if f.Flag(0).Value != &i {
		t.Error("expected the first flag to be i")
	}

	builder = new(strings.Builder)
	f.SetOutput(builder)

	if f.Output() != builder {
		t.Error("somehow the output wasn't set")
	}

	f.Usage()

	expected = "Flags:\n  -i,--int value\n\tinteger (default 0)\n"
	if usageMsg := builder.String(); usageMsg != expected {
		t.Errorf("Unexpected usage: %s", usageMsg)
	}

	if err := f.Parse([]string{"-i", "42", "foo"}); err != nil {
		t.Errorf("unexpected parse error: %s", err)
	}

	if i != 42 {
		t.Errorf("i wasn't set correctly: %d", i)
	}

	if f.Args()[0] != "foo" {
		t.Error("Args are in an unexpected state")
	}

	if err := f.Parse([]string{"-i", "42"}); err != nil {
		t.Errorf("unexpected parse error: %s", err)
	}

	if i != 42 {
		t.Errorf("i wasn't set correctly: %d", i)
	}

	if len(f.Args()) != 0 {
		t.Error("there should not be Args left")
	}

	if err := f.Parse([]string{"--int=24", "foo"}); err != nil {
		t.Errorf("unexpected parse error: %s", err)
	}

	if i != 24 {
		t.Errorf("i wasn't set correctly: %d", i)
	}

	if f.Args()[0] != "foo" {
		t.Error("Args are in an unexpected state")
	}

	f.SetErrFlag(failure.ContinueOnError)

	if err := f.Parse([]string{"-i", "foo"}); err == nil {
		t.Error("This should have been a failure")
	}

	var failed = false

	f.SetErrFlag(failure.ExitOnError)

	failure.Failure = func() { failed = true }
	builder = new(strings.Builder)

	f.SetOutput(builder)

	if err := f.Parse([]string{"-i"}); err == nil {
		t.Error("We expected a failure here")
	} else if err.Error() != "flag -i needs an argument" {
		t.Errorf("unexpected error message: %s", err)
	}

	if !failed {
		t.Error("We expected a failure")
	}

	expected = "Error parsing flag:  flag -i needs an argument."
	if errMsg := builder.String(); errMsg != expected {
		t.Errorf("Unexpected error: %s", errMsg)
	}

	builder = new(strings.Builder)

	f.SetOutput(builder)

	if err := f.Parse([]string{"--int"}); err == nil {
		t.Error("We expected a failure here")
	} else if err.Error() != "flag --int needs an argument" {
		t.Errorf("unexpected error message: %s", err)
	}

	if !failed {
		t.Error("We expected a failure")
	}

	expected = "Error parsing flag:  flag --int needs an argument."
	if errMsg := builder.String(); errMsg != expected {
		t.Errorf("Unexpected error: %s", errMsg)
	}

	if err := f.Parse([]string{"--int", "13"}); err != nil {
		t.Error("This is a valid argument")
	}

	if i != 13 {
		t.Error("We expected i to be 13 now")
	}

	if err := f.Parse([]string{"--int", "foo"}); err == nil {
		t.Error("This should still fail")
	}
}

func TestDoubleInsertion(t *testing.T) {
	f := flags.NewFlagSet("test", failure.ExitOnError)

	var i vals.IntValue

	f.SetErrFlag(failure.ContinueOnError)

	if err := f.Var(&i, "", "i", "integer"); err != nil {
		t.Fatal("error inserting an integer")
	}

	if err := f.Var(&i, "", "i", "integer"); err == nil {
		t.Fatal("no error inserting a flag twice")
	}

	if err := f.Var(&i, "int", "", "integer"); err != nil {
		t.Fatal("error inserting an integer")
	}

	if err := f.Var(&i, "int", "", "integer"); err == nil {
		t.Fatal("no error inserting a flag twice")
	}
}

func TestDashDashTerminate(t *testing.T) {
	f := flags.NewFlagSet("test", failure.ExitOnError)

	// '--' stops parsing, so the missing flags foo and bar shouldn't
	// be a problem.
	if err := f.Parse([]string{"--", "--foo", "--bar"}); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if !reflect.DeepEqual(f.Args(), []string{"--foo", "--bar"}) {
		t.Errorf("Args is not in is correct state")
	}
}

func TestBool(t *testing.T) {
	f := flags.NewFlagSet("test", failure.ExitOnError)

	var b vals.BoolValue

	if err := f.Var(&b, "bool", "b", "boolean"); err != nil {
		t.Fatal("It should not be an error to insert a bool")
	}

	if err := f.Parse([]string{"-b"}); err != nil {
		t.Error("We should be able to set b without values")
	}

	if !b {
		t.Error("b was not set")
	}

	b = false

	if err := f.Parse([]string{"--bool"}); err != nil {
		t.Error("We should be able to set bool without values")
	}

	if !b {
		t.Error("bool was not set")
	}
}

func TestCallback(t *testing.T) { //nolint:funlen // test functions are just sometimes long
	called := false
	cb := vals.FuncNoValue(func() error { called = true; return nil })

	f := flags.NewFlagSet("test", failure.ContinueOnError)

	if err := f.Var(&cb, "foo", "f", "callback"); err != nil {
		t.Error("didn't expect an error inserting a function value")
	}

	flag := f.Lookup("f")
	if flag == nil {
		t.Fatal("we should be able to get cb back")
	}

	if flag.Value != &cb {
		t.Error("what happened to our callback?")
	}

	nv, ok := flag.Value.(interfaces.NoValueFlag)

	if !ok {
		t.Error("the callback should implement this!")
	}

	if nv == nil {
		t.Error("nv should not be nil")
	} else if !nv.NoValueFlag() {
		t.Error("it should return true for no value")
	}

	if err := f.Parse([]string{"-f", "foo", "bar"}); err != nil {
		t.Errorf("didn't expect a parser error here: %s", err)
	}

	if !called {
		t.Error("We expected the function to be called")
	}

	if !reflect.DeepEqual(f.Args(), []string{"foo", "bar"}) {
		t.Errorf("unexpected Args: %v", f.Args())
	}

	called = false

	if err := f.Parse([]string{"-f=foo", "bar"}); err == nil {
		t.Error("this should be an error, -f doesn't take arguments")
	}

	if called {
		t.Error("We didn't expect the function to be called")
	}

	if err := f.Parse([]string{"--foo=foo", "bar"}); err == nil {
		t.Error("this should be an error, --foo doesn't take arguments")
	}

	if called {
		t.Error("We didn't expect the function to be called")
	}

	if err := f.Parse([]string{"--foo", "bar"}); err != nil {
		t.Error("this should be okay")
	}

	if !called {
		t.Error("We expect the function to be called")
	}
}

func TestErrors(t *testing.T) {
	f := flags.NewFlagSet("test", failure.ContinueOnError)

	if err := f.Var(nil, "foo", "bar", ""); err == nil {
		t.Error("The short name is too long")
	} else if err.Error() != "short arguments can only be one character: bar" {
		t.Errorf("unexpected error: %s", err)
	}

	if err := f.Parse([]string{"-="}); err == nil {
		t.Error("-= is an invalid flag")
	}

	if err := f.Parse([]string{"-aab=c"}); err == nil {
		t.Error("-aab=c is an invalid flag")
	}

	if err := f.Parse([]string{"---"}); err == nil {
		t.Error("--- is an invalid flag")
	}

	if err := f.Parse([]string{"--foo"}); err == nil {
		t.Error("--foo isn't defined")
	}

	if err := f.Parse([]string{"-x"}); err == nil {
		t.Error("-x isn't defined")
	}
}

func TestRunOfShort(t *testing.T) {
	var (
		f      = flags.NewFlagSet("test", failure.ContinueOnError)
		a      vals.BoolValue
		b      vals.BoolValue
		c      vals.IntValue
		called = false
		cb     = vals.FuncNoValue(func() error { called = true; return nil })
	)

	if err := f.Var(&a, "", "a", ""); err != nil {
		t.Errorf("we should be able to set a: %s", err)
	}

	if err := f.Var(&b, "", "b", ""); err != nil {
		t.Errorf("we should be able to set b: %s", err)
	}

	if err := f.Var(&c, "", "c", ""); err != nil {
		t.Errorf("we should be able to set c: %s", err)
	}

	if err := f.Var(&cb, "", "f", ""); err != nil {
		t.Errorf("we should be able to set f: %s", err)
	}

	if err := f.Parse([]string{"-fabc", "42"}); err != nil {
		t.Error("this is a valid flag!")
	}

	if !called {
		t.Error("f should be called")
	}

	if !a {
		t.Error("a should be set")
	}

	if !b {
		t.Error("b should also be set")
	}

	if c != 42 {
		t.Error("c should have gotten a value")
	}

	if err := f.Parse([]string{"-afcb"}); err == nil {
		t.Error("this should be an error, because c didn't get a value")
	}

	if err := f.Parse([]string{"-abfc"}); err == nil {
		t.Error("this should be an error, because c didn't get a value")
	}

	if err := f.Parse([]string{"-xab"}); err == nil {
		t.Error("this should be an error, because x doesn't exit")
	}

	if err := f.Parse([]string{"-abx"}); err == nil {
		t.Error("this should be an error, because x doesn't exit")
	}
}

// these are a bit stilly to get full coverage...
func checkProtocol(t *testing.T, f *flags.FlagSet) {
	t.Helper()

	if err := f.Parse([]string{"--int"}); err == nil {
		t.Error("This should fail when MyIntDef.Set() raises an error")
	} else if err.Error() != "error parsing flag --int: myint" {
		t.Errorf("unexpected error: '%s'", err)
	}

	if err := f.Parse([]string{"-ai"}); err == nil {
		t.Error("This should fail when MyIntDef.Set() raises an error")
	} else if err.Error() != "error parsing flag -i: myint" {
		t.Errorf("unexpected error: '%s'", err)
	}

	if err := f.Parse([]string{"-ia"}); err == nil {
		t.Error("This should fail when MyIntDef.Set() raises an error")
	} else if err.Error() != "error evaluating flag -i: myint" {
		t.Errorf("unexpected error: '%s'", err)
	}
}

type MyIntDef int // a value with a default that raises an error

func (val *MyIntDef) String() string  { return "foo" }
func (val *MyIntDef) Default() string { return "42" }
func (val *MyIntDef) Set(x string) error {
	return interfaces.ParseErrorf("myint")
}

func TestDefProtocolErrors(t *testing.T) {
	f := flags.NewFlagSet("test", failure.ContinueOnError)
	a := vals.BoolValue(false)
	i := MyIntDef(0)

	if err := f.Var(&a, "", "a", ""); err != nil {
		t.Error("inserting the myint should be okay")
	}

	if err := f.Var(&i, "int", "i", ""); err != nil {
		t.Error("inserting the myint should be okay")
	}

	checkProtocol(t, f)
}

type MyIntNoVal int // a flag without values that raises an error

func (val *MyIntNoVal) String() string    { return "foo" }
func (val *MyIntNoVal) NoValueFlag() bool { return true }
func (val *MyIntNoVal) Set(x string) error {
	return interfaces.ParseErrorf("myint")
}

func TestNoValProtocolErrors(t *testing.T) {
	f := flags.NewFlagSet("test", failure.ContinueOnError)
	a := vals.BoolValue(false)
	i := MyIntNoVal(0)

	if err := f.Var(&a, "", "a", ""); err != nil {
		t.Error("inserting the myint should be okay")
	}

	if err := f.Var(&i, "int", "i", ""); err != nil {
		t.Error("inserting the myint should be okay")
	}

	checkProtocol(t, f)
}
