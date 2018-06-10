package json

import (
	"github.com/mspark777/go/encoding"
	"github.com/mspark777/go/fmt"
)

func NewObject() Object {
	var object = make(Object)
	return object
}

func SetObjectValue(object Object, key string, value interface{}) {
	object[key] = value
}

func DeleteObjectValue(object Object, keys ...string) {
	for _, key := range keys {
		delete(object, key)
	}
}

func GetObjectValue(object Object, key string) interface{} {
	var value, _ = object[key]
	return value
}

func GetObjectInt(object Object, key string) int {
	var value = GetObjectValue(object, key)
	var i = fmt.ToInt(value)

	return i
}

func GetObjectInt64(object Object, key string) int64 {
	var value = GetObjectValue(object, key)
	var i64 = fmt.ToInt64(value)

	return i64
}

func GetObjectString(object Object, key string) string {
	var value = GetObjectValue(object, key)
	var str = fmt.ToString(value)

	return str
}

func EncodeObject(object Object) []byte {
	var data = encoding.EncodeJSON(object)
	return data
}

func DecodeObject(data []byte) Object {
	var object = NewObject()
	encoding.DecodeJSON(data, &object)

	return object
}

func MergeObject(objects ...Object) Object {
	var result = NewObject()
	for _, obj := range objects {
		for key, value := range obj {
			result[key] = value
		}
	}

	return result
}

func MergeData(datas ...[]byte) []byte {
	var objects = make([]Object, len(datas))
	for i, data := range datas {
		var object = DecodeObject(data)
		objects[i] = object
	}

	var object = MergeObject(objects...)
	var data = EncodeObject(object)
	return data
}
