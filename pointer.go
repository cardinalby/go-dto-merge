package dtomerge

import (
	"reflect"
)

func mergePointers(src reflect.Value, patch reflect.Value, options Options) (reflect.Value, error) {
	if patch.IsZero() || patch.IsNil() {
		return src, nil
	}
	if src.IsZero() || src.IsNil() {
		return patch, nil
	}
	srcElem := src.Elem()
	patchElem := patch.Elem()
	mergedElem, err := mergeAny(srcElem, patchElem, options)
	if err != nil {
		return reflect.Value{}, err
	}
	result := reflect.New(src.Type().Elem())
	result.Elem().Set(mergedElem)
	return result, nil
}
