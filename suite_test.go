// Add a copyright
// Add a licence

package checkers

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSetTestingT(t *testing.T) {
	t.Run("embed suite", func(t *testing.T) {
		type Embed struct {
			Test
		}
		aT := &testing.T{}
		s := &Embed{}
		ok := setTestingT(aT, reflect.ValueOf(s))
		if !ok {
			t.Fatalf("unable to set the testing.T")
		}
		if fmt.Sprintf("%p", s.T) != fmt.Sprintf("%p", aT) {
			t.Fatalf("nested testing.T not set")
		}
	})
	t.Run("embed suite pointer", func(t *testing.T) {
		type Embed struct {
			*Test
		}
		aT := &testing.T{}
		s := &Embed{&Test{}}
		ok := setTestingT(aT, reflect.ValueOf(s))
		if !ok {
			t.Fatalf("unable to set the testing.T")
		}
		if fmt.Sprintf("%p", s.T) != fmt.Sprintf("%p", aT) {
			t.Fatalf("nested testing.T not set")
		}
	})
	t.Run("embed suite nil pointer", func(t *testing.T) {
		type Embed struct {
			*Test
		}
		aT := &testing.T{}
		s := &Embed{}
		ok := setTestingT(aT, reflect.ValueOf(s))
		if ok {
			t.Fatalf("unexpected setting of the testing.T")
		}
	})

}
