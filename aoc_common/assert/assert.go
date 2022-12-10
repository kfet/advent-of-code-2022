package assert

import "reflect"

func Equals(expected, actual interface{}, message string) {
	if !reflect.DeepEqual(expected, actual) {
		panic(message)
	}
}
