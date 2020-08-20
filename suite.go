// Add a copyright
// Add a licence

package checkers

import (
	"reflect"
	"testing"
)

// Suite extends the standard library T.
type Suite struct {
	t          *testing.T
	setupFuncs []func(t *testing.T)
}

func NewSuite(t *testing.T, setupFuncs ...func(*testing.T)) *Suite {
	return &Suite{
		t:          t,
		setupFuncs: setupFuncs,
	}
}

// Run will run the test as a subtest after calling each of the setup functions.
func (s *Suite) Run(name string, test func(*testing.T)) bool {
	for _, f := range s.setupFuncs {
		f(s.t)
	}
	return s.t.Run(name, test)
}

// Assert expects to succeed, and if not, causes the test to fail immediately.
func (s *Suite) Assert(obtained interface{}, checker Checker, extras ...interface{}) {
	if ok := checker.Check(s.t, obtained, extras...); !ok {
		s.t.FailNow()
	}
}

// Check will mark the test as a failure if the checker fails. The test continues.
func (s *Suite) Check(obtained interface{}, checker Checker, extras ...interface{}) bool {
	return checker.Check(s.t, obtained, extras...)
}

// Checker defines the interface for any specific checker.
type Checker interface {
	// Check determins if the obtained value is sufficient.
	// If the checker requires an expected value, it should be the first
	// extra value.
	Check(t *testing.T, obtained interface{}, extras ...interface{}) bool
}

type isNil struct{}

var IsNil Checker = isNil{}

func (isNil) Check(t *testing.T, obtained interface{}, extras ...interface{}) bool {
	if obtained == nil {
		return true
	}
	t.Error("obtained value is non-nil")
	return false
}

type equals struct{}

// Equals checker tests for equality.
var Equals Checker = equals{}

// TODO: add describer interface, and pass failing values to the describers
// in the checkers.

func (equals) Check(t *testing.T, obtained interface{}, extras ...interface{}) bool {
	if len(extras) == 0 {
		t.Errorf("missing 'expected' value")
		return false
	}
	expected, extras := extras[0], extras[1:]
	exValue := reflect.ValueOf(expected)
	value := reflect.ValueOf(obtained)
	if value.Kind() != exValue.Kind() {
		t.Errorf("obtained type %T, does not match expected %T", obtained, expected)
		return false
	}
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() == exValue.Bool() {
			return true
		}
	case reflect.String:
		if value.String() == exValue.String() {
			return true
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() == exValue.Int() {
			return true
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value.Uint() == exValue.Uint() {
			return true
		}
	case reflect.Float32, reflect.Float64:
		if value.Float() == exValue.Float() {
			return true
		}
	default:
		t.Errorf("Equals checker does not support type %T", obtained)
		return false
	}
	t.Errorf("expected %T value %v, got %v", expected, expected, obtained)
	return false
}

type isFalse struct{}

var IsFalse Checker = isFalse{}

func (isFalse) Check(t *testing.T, obtained interface{}, extras ...interface{}) bool {
	value := reflect.ValueOf(obtained)
	switch value.Kind() {
	case reflect.Bool:
		if !value.Bool() {
			return true
		}
	default:
		t.Errorf("IsFalse checker expected bool, obtained was type %T", obtained)
		return false
	}
	t.Error("obtained value is true")
	return false
}

type isTrue struct{}

var IsTrue Checker = isTrue{}

func (isTrue) Check(t *testing.T, obtained interface{}, extras ...interface{}) bool {
	value := reflect.ValueOf(obtained)
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() {
			return true
		}
	default:
		t.Errorf("IsTrue checker expected bool, obtained was type %T", obtained)
		return false
	}
	t.Error("obtained value is false")
	return false
}
