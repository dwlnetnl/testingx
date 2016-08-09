package testingx

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/exec"
	"regexp"
	"strings"
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

func TestInDelta(t *testing.T) {
	cases := []struct {
		lhs, rhs, delta float64
	}{
		{math.Sqrt(3), 1.732, 1e-3},
		{math.Pow(3, 1.2), 3.737192, 1e-6},
	}

	for _, c := range cases {
		if !InDelta(c.lhs, c.rhs, c.delta) {
			t.Errorf("%f != %f ± %f", c.lhs, c.rhs, c.delta)
		}
	}
}

func TestSkipIfShort(t *testing.T) {
	if testing.Short() {
		defer func() {
			if !t.Skipped() {
				t.Error("test is not skipped with -short")
			}
		}()

		SkipIfShort(t)
		return
	}

	if os.Getenv("SKIPIFSHORT") != "1" {
		cmd := exec.Command("go", "test", "-short", "-run=TestSkipIfShort")
		cmd.Env = append(cmd.Env, "GOPATH="+os.Getenv("GOPATH"))
		cmd.Env = append(cmd.Env, "SKIPIFSHORT=1")

		if testing.Verbose() {
			cmd.Args = append(cmd.Args, "-v")
		}

		output, err := cmd.Output()
		if err != nil {
			err, ok := err.(*exec.ExitError)
			// exit status 2: invocation failed
			// or not an *exec.ExitError
			if !ok || ok && len(err.Stderr) > 0 {
				t.Errorf("go test -short returned error: %v", err)
				return
			}

			// exit status 1: test failed
			for _, line := range strings.Split(string(output), "\n") {
				t.Log(line)
			}

			t.Fail()
		}
	}
}
