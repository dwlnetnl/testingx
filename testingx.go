// Package testingx implements utilities for tests.
package testingx

import (
	"math"
	"regexp"
	"testing"
)

// EqualErrors compares errors lhs and rhs for equality.
func EqualErrors(lhs, rhs error) bool {
	if lhs == rhs {
		return true
	}

	if lhs != nil && rhs != nil {
		if lhs.Error() == rhs.Error() {
			return true
		}
	}

	return false
}

// EqualError returns true if error err is non-nil and the error string matches
// string str.
func EqualError(err error, str string) bool {
	return err != nil && err.Error() == str
}

// MatchError returns true if error err is non-nil and the error string matches
// regular expression re. It will panic if re cannot compiled.
func MatchError(err error, re string) bool {
	return MatchErrorRegexp(err, regexp.MustCompile(re))
}

// MatchErrorRegexp returns true if error err is non-nil and the error string
// matches regular expression re. It will panic if re is nil.
func MatchErrorRegexp(err error, re *regexp.Regexp) bool {
	if err == nil {
		return false
	}

	return re.MatchString(err.Error())
}

// Panics returns true if function fn panics.
func Panics(fn func()) bool { return Recover(fn) != nil }

// Recover calls fn and returns the error value passed to panic.
func Recover(fn func()) (v interface{}) {
	defer func() {
		v = recover()
	}()

	fn()
	return
}

// InDelta returns true if floats lhs and rhs are equal within delta.
func InDelta(lhs, rhs, delta float64) bool {
	return math.Abs(lhs-rhs) < delta
}

// SkipIfShort causes the test to be skipped when running with -short.
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test with -short")
	}
}
