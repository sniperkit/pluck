package json

import (
	"github.com/mspark777/go/fmt"
)

func GetArrayValue(array Array, index int) interface{} {
	var length = len(array)
	if index < 0 {
		return nil
	} else if index >= length {
		return nil
	}

	var value = array[index]
	return value
}

func GetArrayInt(array Array, index int) int {
	var value = GetArrayValue(array, index)
	var i = fmt.ToInt(value)

	return i
}

func GetArrayInt64(array Array, index int) int64 {
	var value = GetArrayValue(array, index)
	var i64 = fmt.ToInt64(value)

	return i64
}

func GetArrayString(array Array, index int) string {
	var value = GetArrayValue(array, index)
	var str = fmt.ToString(value)

	return str
}
