package vals_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mailund/cli/internal/vals"
)

func TestCallbacks(t *testing.T) {
	var x, y string

	f := reflect.ValueOf(
		func(s string) error { x = s; return nil })
	g := reflect.ValueOf(
		func(s string, i interface{}) error { y = s + i.(string); return nil })

	fv := vals.AsCallback(&f, nil)
	gv := vals.AsCallback(&g, "bar")

	if err := fv.Set("foo"); err != nil {
		t.Fatal("Error calling f")
	}

	if err := gv.Set("foo"); err != nil {
		t.Fatal("Error calling g")
	}

	if x != "foo" { //nolint:goconst // no, foo should not be a constant
		t.Error("f didn't work")
	}

	if y != "foobar" {
		t.Error("g didn't work")
	}

	if fv.String() != "" || gv.String() != "" {
		t.Error("Functions should have empty default strings")
	}
}

func TestBoolCallbacks(t *testing.T) {
	var x, y string

	f := reflect.ValueOf(func() { x = "foo" })
	g := reflect.ValueOf(func(i interface{}) { y = "foo" + i.(string) })

	fv := vals.AsCallback(&f, nil)
	gv := vals.AsCallback(&g, "bar")

	if h, ok := fv.(vals.FuncBoolValue); !ok {
		t.Error("We should have a bool function here")
	} else if !h.IsBoolFlag() {
		t.Error("We should have IsBoolFlag()")
	}

	if h, ok := gv.(vals.FuncBoolValue); !ok {
		t.Error("We should have a bool function here")
	} else if !h.IsBoolFlag() {
		t.Error("We should have IsBool()")
	}

	if err := fv.Set("foo"); err != nil {
		t.Fatal("Error calling f")
	}

	if err := gv.Set("foo"); err != nil {
		t.Fatal("Error calling g")
	}

	if x != "foo" {
		t.Error("f didn't work")
	}

	if y != "foobar" {
		t.Error("g didn't work")
	}

	if fv.String() != "" || gv.String() != "" {
		t.Error("Functions should have empty default strings")
	}
}

func TestBoolECallbacks(t *testing.T) {
	var x, y string

	f := reflect.ValueOf(
		func() error { x = "foo"; return nil })
	g := reflect.ValueOf(
		func(i interface{}) error { y = "foo" + i.(string); return nil })

	fv := vals.AsCallback(&f, nil)
	gv := vals.AsCallback(&g, "bar")

	if h, ok := fv.(vals.FuncBoolValue); !ok {
		t.Error("We should have a bool function here")
	} else if !h.IsBoolFlag() {
		t.Error("We should have IsBoolFlag()")
	}

	if h, ok := gv.(vals.FuncBoolValue); !ok {
		t.Error("We should have a bool function here")
	} else if !h.IsBoolFlag() {
		t.Error("We should have IsBool()")
	}

	if err := fv.Set("foo"); err != nil {
		t.Fatal("Error calling f")
	}

	if err := gv.Set("foo"); err != nil {
		t.Fatal("Error calling g")
	}

	if x != "foo" {
		t.Error("f didn't work")
	}

	if y != "foobar" {
		t.Error("g didn't work")
	}

	if fv.String() != "" || gv.String() != "" {
		t.Error("Functions should have empty default strings")
	}
}

func TestVariadicCallbacks(t *testing.T) {
	var x, y string

	f := reflect.ValueOf(
		func(s []string) error { x = strings.Join(s, ""); return nil })
	g := reflect.ValueOf(
		func(s []string, i interface{}) error { y = strings.Join(s, "") + i.(string); return nil })

	fv := vals.AsVariadicCallback(&f, nil)
	gv := vals.AsVariadicCallback(&g, "baz")

	if err := fv.Set([]string{"foo", "bar"}); err != nil {
		t.Fatal("Error calling f")
	}

	if err := gv.Set([]string{"foo", "bar"}); err != nil {
		t.Fatal("Error calling g")
	}

	if x != "foobar" {
		t.Error("h didn't work")
	}

	if y != "foobarbaz" {
		t.Error("i didn't work")
	}
}

func TestWrongSignature(t *testing.T) {
	f := reflect.ValueOf(func(int) error { return nil })

	if vals.AsCallback(&f, nil) != nil {
		t.Error("Unexpected callback")
	}

	if vals.AsVariadicCallback(&f, nil) != nil {
		t.Error("Unexpected callback")
	}
}
