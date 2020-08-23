// Add a copyright
// Add a licence

package checkers

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

var testMethodMatch = regexp.MustCompile(`^Test([A-Z]\w*)$`)

// RunSuite runs a collection of methods as subtests.
func RunSuite(t *testing.T, suite interface{}) {
	v := reflect.ValueOf(suite)
	if v.Kind() != reflect.Ptr {
		t.Fatalf("suite must be passed in with pointer, not value")
	}
	// Find the *testing.T in the suite, and set it.
	if ok := setTestingT(t, v); !ok {
		t.Fatal("unable to initialize the suite *testing.T")
		return
	}
	// See if there is a method called SetUpTest, and if there is
	// save it in the suite.
	setup := v.MethodByName("SetUpTest")
	if !setup.IsZero() {
		fmt.Println("found SetUpTest")
		// There is a setup method, ensure it takes no args.
		methodType := setup.Type()
		if methodType.NumIn() != 0 {
			fmt.Println("doesn't take no args")
			t.Fatal("SetUpTest should take no arguments")
		}
	}

	testMethods := findTestMethods(v)

	for _, method := range testMethods {
		// We know that the method name starts with Test,
		// so remove that for the subtest name.
		short := method.Name[4:]
		testFunc := v.MethodByName(method.Name)
		t.Run(short, func(*testing.T) {
			if !setup.IsZero() {
				setup.Call(nil)
			}
			funcType := testFunc.Type()
			if count := funcType.NumIn(); count != 0 {
				t.Fatalf("Test method %q takes %d args, should take none", method.Name, count)
			}
			if count := funcType.NumOut(); count != 0 {
				t.Fatalf("Test method %q returns %d values, should return none", method.Name, count)
			}
			testFunc.Call(nil)
		})
	}
}

func findTestMethods(v reflect.Value) []reflect.Method {
	result := []reflect.Method{}

	t := v.Type()
	numMethods := t.NumMethod()
	fmt.Printf("suite type %v has %d method\n", t.Name(), numMethods)
	for i := 0; i < numMethods; i++ {
		method := t.Method(i)
		fmt.Printf("looking at method %q\n", method.Name)
		match := testMethodMatch.FindStringSubmatch(method.Name)
		if len(match) > 0 {
			fmt.Printf("it's a match\n")
			result = append(result, method)
		}
	}
	return result
}

func setTestingT(t *testing.T, v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr:
		return setTestingT(t, v.Elem())
	case reflect.Struct:
		// Falls through to the rest of the function.
	default:
		return false
	}

	tValue := reflect.ValueOf(t)
	tType := tValue.Type()
	fieldCount := v.NumField()

	for i := 0; i < fieldCount; i++ {
		field := v.Field(i)
		if field.Type() == tType {
			if field.CanSet() {
				field.Set(tValue)
				return true
			}
		}
		switch field.Kind() {
		case reflect.Struct, reflect.Ptr:
			if ok := setTestingT(t, field); ok {
				return true
			}
		}
	}
	return false
}
