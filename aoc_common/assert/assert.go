package assert

import (
	"fmt"
	"reflect"
	"testing"
)

func Equals(expected, actual interface{}, message string) {
	if !reflect.DeepEqual(expected, actual) {
		msg := fmt.Sprintf("%s: expected: %v, actual %v", message, expected, actual)
		panic(msg)
	}
}

func NoErr(err error) {
	if err != nil {
		panic(err)
	}
}

func True(c bool) {
	if !c {
		panic("condition should be true")
	}
}

func False(c bool) {
	if c {
		panic("condition should be false")
	}
}

func EqualsT(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		msg := fmt.Sprintf("expected: %v, actual %v", expected, actual)
		t.Error(msg)
	}
}

func NoErrT(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}
