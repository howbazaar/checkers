// Add a copyright
// Add a licence

package checkers

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Checker defines the interface for any specific checker.
type Checker interface {
	// Check determins if the obtained value is sufficient.
	// If the checker requires an expected value, it should be the first
	// extra value.
	Check(obtained interface{}, extras ...interface{}) error
}

type isNil struct{}

// IsNil checker will return an error if the obtained value is not nil.
var IsNil Checker = isNil{}

func (isNil) Check(obtained interface{}, extras ...interface{}) error {
	if obtained == nil {
		return nil
	}
	return errors.New("obtained value is non-nil")
}

type equals struct{}

// Equals checker tests for equality.
var Equals Checker = equals{}

// TODO: add describer interface, and pass failing values to the describers
// in the checkers.

func (equals) Check(obtained interface{}, extras ...interface{}) error {
	if len(extras) == 0 {
		return errors.New("missing 'expected' value")
	}
	expected, extras := extras[0], extras[1:]
	exValue := reflect.ValueOf(expected)
	value := reflect.ValueOf(obtained)
	if value.Kind() != exValue.Kind() {
		return fmt.Errorf("obtained type %T does not match expected type %T", obtained, expected)
	}
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() == exValue.Bool() {
			return nil
		}
	case reflect.String:
		if value.String() == exValue.String() {
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() == exValue.Int() {
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value.Uint() == exValue.Uint() {
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if value.Float() == exValue.Float() {
			return nil
		}
	default:
		return fmt.Errorf("Equals checker does not support type %T", obtained)
	}
	return fmt.Errorf("expected %T value %v, got %v", expected, expected, obtained)
}

type deepEquals struct{}

// DeepEquals checker tests for equality of complex types.
var DeepEquals Checker = deepEquals{}

func (deepEquals) Check(obtained interface{}, extras ...interface{}) error {
	if len(extras) == 0 {
		return errors.New("missing 'expected' value")
	}
	expected, extras := extras[0], extras[1:]

	if ok, err := DeepEqual(obtained, expected); !ok {
		return err
	}
	return nil
}

type isFalse struct{}

// IsFalse checker will return an error if the obtained value is not a bool
// or if the bool value is true.
var IsFalse Checker = isFalse{}

func (isFalse) Check(obtained interface{}, extras ...interface{}) error {
	value := reflect.ValueOf(obtained)
	switch value.Kind() {
	case reflect.Bool:
		if !value.Bool() {
			return nil
		}
	default:
		return fmt.Errorf("IsFalse checker expected bool, obtained was type %T", obtained)
	}
	return errors.New("obtained value is true")
}

type isTrue struct{}

// IsTrue checker will return an error if the obtained value is not a bool
// of if the bool value is false.
var IsTrue Checker = isTrue{}

func (isTrue) Check(obtained interface{}, extras ...interface{}) error {
	value := reflect.ValueOf(obtained)
	switch value.Kind() {
	case reflect.Bool:
		if value.Bool() {
			return nil
		}
	default:
		return fmt.Errorf("IsTrue checker expected bool, obtained was type %T", obtained)
	}
	return errors.New("obtained value is false")
}

type hasLen struct{}

// HasLen checker will return an error of the type does not support the
// getting the length using the `len` function, or if the length does not
// match the specified value.
var HasLen Checker = hasLen{}

func (hasLen) Check(obtained interface{}, extras ...interface{}) error {
	if len(extras) == 0 {
		return errors.New("missing 'expected' value")
	}
	expected, extras := extras[0], extras[1:]
	// TODO: deal with panics of wrong types of unexpected types
	size := reflect.ValueOf(expected).Int()

	value := reflect.ValueOf(obtained)
	var length int
	switch value.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		length = value.Len()
	default:
		return fmt.Errorf("HasLen checker expected array, channel, map, slice or string, obtained was type %T", obtained)
	}

	if int64(length) != size {
		return fmt.Errorf("expected length %d, obtained %d", size, length)
	}

	return nil
}

type matches struct{}

// Matches checker will use regex to match against a string, or Stringer.
var Matches Checker = matches{}

func (matches) Check(obtained interface{}, extras ...interface{}) error {
	if len(extras) == 0 {
		return errors.New("missing 'expected' value")
	}
	expected, extras := extras[0], extras[1:]
	pattern, ok := expected.(string)
	if !ok {
		return errors.New("expected value must be a string containing a regexp pattern")
	}
	var value string
	if s, ok := obtained.(string); ok {
		value = s
	} else if s, ok := obtained.(fmt.Stringer); ok {
		value = s.String()
	} else {
		return fmt.Errorf("%T(%#v) is neither a string nor has a 'String() string' method", obtained, obtained)
	}

	return checkMatch(value, pattern)
}

func checkMatch(obtained, pattern string) error {
	if !strings.HasPrefix(pattern, "^") {
		pattern = "^" + pattern
	}
	if !strings.HasSuffix(pattern, "$") {
		pattern = pattern + "$"
	}
	ok, err := regexp.MatchString(pattern, obtained)
	if err != nil {
		return fmt.Errorf("unable to compile regexp: %v", err)
	}
	if ok {
		return nil
	}
	return fmt.Errorf("%q did not match pattern %q", obtained, pattern)
}

type panicMatches struct{}

// PanicMatches checker will match an error or string panic result against
// the specified pattern.
var PanicMatches Checker = panicMatches{}

func (panicMatches) Check(obtained interface{}, extras ...interface{}) (err error) {
	if len(extras) == 0 {
		return errors.New("missing 'expected' value")
	}
	expected, extras := extras[0], extras[1:]
	pattern, ok := expected.(string)
	if !ok {
		return errors.New("expected value must be a string containing a regexp pattern")
	}
	// First arg must be a function with no args.
	f := reflect.ValueOf(obtained)
	if f.Kind() != reflect.Func || f.Type().NumIn() != 0 {
		return errors.New("first arg must be a function that takes no args")
	}

	defer func() {
		v := recover()
		if v == nil {
			// No panic, that's bad, but the error already says that.
		} else if e, ok := v.(error); ok {
			err = checkMatch(e.Error(), pattern)
		} else if s, ok := v.(string); ok {
			err = checkMatch(s, pattern)
		} else {
			err = fmt.Errorf("recovered panic value %T(%#v) is not a string nor an error", v, v)
		}
	}()
	f.Call(nil)
	return errors.New("no panic")
}
