package go_util

import (
	"fmt"
	"reflect"
)

// unique Stringer list
// eg: Unique(oldList, &newList)
func Unique(ary interface{}, out interface{}) {
	rv := reflect.ValueOf(ary)
	outRv := reflect.ValueOf(out).Elem()

	m := make(map[string]interface{})

	for i := 0; i < rv.Len(); i++ {
		item := rv.Index(i).Interface()
		key := item.(fmt.Stringer).String()
		_, found := m[key]
		if found {
			continue
		}
		m[key] = true

		outRv = reflect.Append(outRv, rv.Index(i))
	}

	reflect.ValueOf(out).Elem().Set(outRv)
}
