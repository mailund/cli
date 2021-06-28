package vals_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mailund/cli/internal/vals"
)

func TestCallbacks(t *testing.T) {
	var x, y, z, w string

	f := reflect.ValueOf(
		func(s string) error { x = s; return nil })
	g := reflect.ValueOf(
		func(s string, i interface{}) error { y = s + i.(string); return nil })
	h := reflect.ValueOf(
		func(s []string) error { z = strings.Join(s, ""); return nil })
	i := reflect.ValueOf(
		func(s []string, i interface{}) error { w = strings.Join(s, "") + i.(string); return nil })

	fv := vals.AsCallback(&f, nil)
	gv := vals.AsCallback(&g, "bar")
	hv := vals.AsVariadicCallback(&h, nil)
	iv := vals.AsVariadicCallback(&i, "baz")

	if err := fv.Set("foo"); err != nil {
		t.Fatal("Error calling f")
	}

	if err := gv.Set("foo"); err != nil {
		t.Fatal("Error calling g")
	}

	if err := hv.Set([]string{"foo", "bar"}); err != nil {
		t.Fatal("Error calling h")
	}

	if err := iv.Set([]string{"foo", "bar"}); err != nil {
		t.Fatal("Error calling i")
	}

	if x != "foo" {
		t.Error("f didn't work")
	}

	if y != "foobar" {
		t.Error("g didn't work")
	}

	if z != "foobar" {
		t.Error("h didn't work")
	}

	if w != "foobarbaz" {
		t.Error("i didn't work")
	}

	if fv.String() != "" || gv.String() != "" || hv.String() != "" || iv.String() != "" {
		t.Error("Functions should have empty default strings")
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
