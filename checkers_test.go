package checkers_test

import (
	"testing"

	"github.com/howbazaar/checkers"
)

func TestIsNil(t *testing.T) {
	err := checkers.IsNil.Check(nil)
	if err != nil {
		t.Fatalf("IsNil(nil) returned error: %v", err)
	}
	// Anything else is not nil.
	type anything struct{}
	err = checkers.IsNil.Check(anything{})
	if err == nil {
		t.Fatal("IsNil(anything{}) should return an error")
	}
	err = checkers.IsNil.Check(&anything{})
	if err == nil {
		t.Fatal("IsNil(&anything{}) should return an error")
	}
}

func TestEquals(t *testing.T) {
	for _, test := range []struct {
		description string
		obtained    interface{}
		expected    interface{}
		err         string
	}{
		{
			description: "bool, both true",
			obtained:    true,
			expected:    true,
		}, {
			description: "bool, both false",
			obtained:    false,
			expected:    false,
		}, {
			description: "bool, unequal",
			obtained:    true,
			expected:    false,
			err:         "expected bool value false, got true",
		}, {
			description: "bool, different type",
			obtained:    "true",
			expected:    false,
			err:         `obtained type string does not match expected type bool`,
		}, {
			description: "string, both empty",
			obtained:    "",
			expected:    "",
		}, {
			description: "string, both same",
			obtained:    "something",
			expected:    "something",
		}, {
			description: "string, different",
			obtained:    "something",
			expected:    "different",
			err:         "expected string value different, got something",
		}, {
			description: "int, same",
			obtained:    1234,
			expected:    1234,
		}, {
			description: "int, different",
			obtained:    1234,
			expected:    4321,
			err:         "expected int value 4321, got 1234",
		}, {
			description: "int, different types",
			obtained:    int32(1234),
			expected:    int64(1234),
			err:         "obtained type int32 does not match expected type int64",
		},
	} {
		err := checkers.Equals.Check(test.obtained, test.expected)
		if err == nil {
			if test.err != "" {
				t.Errorf("%s: expected error: %q", test.description, test.err)
			}
		} else {
			if test.err == "" {
				t.Errorf("%s: unexpected error: %v", test.description, err)
			} else {
				if err.Error() != test.err {
					t.Errorf("%s: error mismatch: \n\tobtained: %q\n\texpected: %q", test.description, err.Error(), test.err)
				}
			}
		}
	}
}

func TestDeepEquals(t *testing.T) {
	err := checkers.DeepEquals.Check(nil)
	if err.Error() != "missing 'expected' value" {
		t.Errorf("incorrect error for missing expected value: %v", err)
	}
	obtained := map[string]interface{}{
		"foo": 1234,
		"bar": "result",
	}
	expected := map[string]interface{}{
		"foo": 1234,
		"bar": "something",
	}
	// The rest of the deep equals checks are done in deepequal_test.go
	err = checkers.DeepEquals.Check(obtained, expected)
	if err.Error() != `mismatch at ["bar"]: unequal; obtained "result"; expected "something"` {
		t.Errorf("incorrect error response: %v", err)
	}
}
