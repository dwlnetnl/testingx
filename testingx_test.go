package testingx

import (
	"errors"
	"fmt"
	"regexp"
	"testing"
)

func TestEqualErrors(t *testing.T) {
	err := errors.New("error")

	cases := []struct {
		equal bool
		lhs   error
		rhs   error
	}{
		{true, err, err},
		{true, errors.New("error"), errors.New("error")},
		{false, errors.New("error"), fmt.Errorf("error rhs")},
	}

	for _, c := range cases {
		got := EqualErrors(c.lhs, c.rhs)
		if got != c.equal {
			t.Errorf("EqualErrors(%#v, %#v) = %v, want: %v", c.lhs, c.rhs, got, c.equal)
		}
	}
}

func TestEqualError(t *testing.T) {
	cases := []struct {
		in   error
		want bool
	}{
		{errors.New("error"), true},
		{nil, false},
	}

	for _, c := range cases {
		const str = "error"
		got := EqualError(c.in, str)

		if got != c.want {
			t.Errorf("EqualError(%#v, %q) = %v, want: %v", c.in, str, got, c.want)
		}
	}
}

func TestMatchError(t *testing.T) {
	cases := []struct {
		in   error
		want bool
	}{
		{errors.New("error"), true},
		{nil, false},
	}

	for _, c := range cases {
		const re = "^error$"
		got := MatchError(c.in, re)

		if got != c.want {
			t.Errorf("MatchError(%#v, %q) = %v, want: %v", c.in, re, got, c.want)
		}
	}
}

func TestMatchError_panic(t *testing.T) {
	paniced := Panics(func() {
		MatchError(errors.New("error"), "\\")
	})

	if !paniced {
		t.Error("MatchError did not panic when regular expression cannot be compiled")
	}
}

func TestMatchErrorRegexp(t *testing.T) {
	cases := []struct {
		in   error
		want bool
	}{
		{errors.New("error"), true},
		{nil, false},
	}

	re := regexp.MustCompile("^error$")
	for _, c := range cases {
		got := MatchErrorRegexp(c.in, re)

		if got != c.want {
			t.Errorf("MatchErrorRegexp(%#v, %v) = %v, want: %v", c.in, re, got, c.want)
		}
	}
}

func TestMatchErrorRegexp_panic(t *testing.T) {
	paniced := Panics(func() {
		MatchErrorRegexp(errors.New("error"), nil)
	})

	if !paniced {
		t.Error("MatchErrorRegexp did not panic when regular expression is nil")
	}
}

func TestPanics(t *testing.T) {
	paniced := Panics(func() {
		panic("trigger panic")
	})

	if !paniced {
		t.Error("Panics did not panic")
	}
}
