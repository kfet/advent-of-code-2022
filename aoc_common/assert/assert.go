package assert

import (
	"fmt"
	"reflect"
)

func Equals(expected, actual interface{}, message string) {
	if !reflect.DeepEqual(expected, actual) {
		msg := fmt.Sprintf("%s: expected: %v, actual %v", message, expected, actual)
		panic(msg)
	}
}
