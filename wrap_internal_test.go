package cli // white box test

import "testing"

func TestWrapSimplifySpace(t *testing.T) {
	x := "  foo\nbar\t\n\t\tbaz\t"

	const linewidth = 70

	if wrapped, expected := wordWrap(x, linewidth), "foo bar baz"; wrapped != expected {
		t.Errorf(`Expected "%s" but got "%s"`, expected, wrapped)
	}
}

func TestWrapNewlines(t *testing.T) {
	x := "  foo\nbar\t\n\t\tbaz\t"

	if wrapped, expected := wordWrap(x, 8), "foo bar\nbaz"; wrapped != expected {
		t.Errorf(`Expected "%s" but got "%s"`, expected, wrapped)
	}
}
