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

// Assert expects to succeed, and if not, causes the test to fail immediately.
func (t *Test) Assert(obtained interface{}, checker Checker, extras ...interface{}) {
	if ok := checker.Check(t.T, obtained, extras...); !ok {
		t.FailNow()
	}
}

// Check will mark the test as a failure if the checker fails. The test continues.
func (t *Test) Check(obtained interface{}, checker Checker, extras ...interface{}) bool {
	return checker.Check(t.T, obtained, extras...)
}
