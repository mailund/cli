package cli // white box test

import "testing"

func TestWrapSimplifySpace(t *testing.T) {
	x := "  foo\nbar\t\n\t\tbaz\t"
	expected := "foo bar baz"
	wrapped := wordWrap(x, 70)
	if wrapped != expected {
		t.Errorf(`Expected "%s" but got "%s"`, expected, wrapped)
	}
}

func TestWrapNewlines(t *testing.T) {
	x := "  foo\nbar\t\n\t\tbaz\t"
	expected := "foo bar\nbaz"
	wrapped := wordWrap(x, 8)
	if wrapped != expected {
		t.Errorf(`Expected "%s" but got "%s"`, expected, wrapped)
	}
}
