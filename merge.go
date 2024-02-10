package dtomerge

import (
	"reflect"
)

func Merge[T any](src T, patch T, opts ...Option) (T, error) {
	res, err := mergeAny(reflect.ValueOf(src), reflect.ValueOf(patch), NewOptions(opts...))
	if err != nil {
		var empty T
		return empty, err
	}
	return res.Interface().(T), nil
}
