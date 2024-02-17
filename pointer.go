package dtomerge

import (
	"reflect"
)

func mergePointers(src reflect.Value, patch reflect.Value, mCtx mergeContext) (reflect.Value, error) {
	if patch.IsZero() || patch.IsNil() {
		return src, nil
	}
	if src.IsZero() || src.IsNil() {
		return patch, nil
	}
	result := reflect.New(src.Type().Elem())

	if mergeResult, has := mCtx.mergedPointers[src.UnsafePointer()]; has {
		return mergeResult, nil
	}
	srcElem := src.Elem()
	patchElem := patch.Elem()
	mCtx.mergedPointers[src.UnsafePointer()] = result

	mergedElem, err := mergeAny(srcElem, patchElem, mCtx)
	if err != nil {
		return reflect.Value{}, err
	}
	result.Elem().Set(mergedElem)

	return result, nil
}
