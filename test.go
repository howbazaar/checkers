// Add a copyright
// Add a licence

package checkers

import (
	"testing"
)

// Test is a simple wrapper around a testing.T to add Assert and Check methods.
type Test struct {
	*testing.T
}

// Check will mark the test as a failure if the checker fails. The test continues.
func (t *Test) Check(obtained interface{}, checker Checker, extras ...interface{}) bool {
	if err := checker.Check(obtained, extras...); err != nil {
		t.Error(err.Error())
		return false
	}
	return true
}

// Assert expects to succeed, and if not, causes the test to fail immediately.
func (t *Test) Assert(obtained interface{}, checker Checker, extras ...interface{}) {
	if ok := t.Check(obtained, checker, extras...); !ok {
		t.FailNow()
	}
}
