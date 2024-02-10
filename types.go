package dtomerge

import (
	"fmt"
	"reflect"

	"github.com/cardinalby/go-dto-merge/types"
)

type CustomMergeFunc[T any] func(src, override T) (T, error)

// CustomMergeFuncs is a map of types to custom merge functions.
// Values should be of CustomMergeFunc[T] type where T is the same type as the key points to.
type CustomMergeFuncs map[reflect.Type]any
type CustomMergeOptions map[reflect.Type]Options
type AtomicTypes []reflect.Type

type MergeOptionsProvider interface {
	GetMergeOptions() Options
}

type Merger[T any] interface {
	Merge(override T) (T, error)
}

func AsMerger(
	value reflect.Value,
) (
	mergeFn func(override reflect.Value) (reflect.Value, error),
	isMerger bool,
) {
	vType := value.Type()
	for i := 0; i < vType.NumMethod(); i++ {
		method := vType.Method(i)
		if method.Name != "Merge" {
			continue
		}
		methodType := method.Type
		if methodType.NumIn() != 2 {
			return nil, false
		}
		if methodType.NumOut() != 2 {
			return nil, false
		}
		if methodType.In(1) != vType {
			return nil, false
		}
		if methodType.Out(0) != vType {
			return nil, false
		}
		if methodType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			return nil, false
		}
		return func(override reflect.Value) (merged reflect.Value, err error) {
			// call method
			results := method.Func.Call([]reflect.Value{value, override})
			if results[0].IsValid() && !results[1].IsZero() && !results[1].IsNil() {
				err = results[1].Interface().(error)
			}
			return results[0], err
		}, true
	}
	return nil, false
}

type customReflectMergeFunc func(src reflect.Value, override reflect.Value) (reflect.Value, error)
type customReflectMergeFuncs map[reflect.Type]customReflectMergeFunc

func asCustomReflectMergeFunc(f any) (customReflectMergeFunc, reflect.Type, error) {
	fValue := reflect.ValueOf(f)
	if !fValue.IsValid() || fValue.IsZero() || fValue.IsNil() {
		return nil, nil, fmt.Errorf("%w: custom merge function is nil", types.ErrInvalidTypes)
	}
	fType := fValue.Type()
	if fType.Kind() != reflect.Func {
		return nil, nil, fmt.Errorf(
			"%w: element of custom merge functions is not a function", types.ErrInvalidTypes,
		)
	}
	if fType.NumIn() != 2 {
		return nil, nil, fmt.Errorf(
			"%w: custom merge function must have 2 arguments", types.ErrInvalidTypes,
		)
	}
	if fType.In(0) != fType.In(1) {
		return nil, nil, fmt.Errorf(
			"%w: both arguments of custom merge function must be of the same type", types.ErrInvalidTypes,
		)
	}
	if fType.NumOut() != 2 {
		return nil, nil, fmt.Errorf(
			"%w: custom merge function must have 2 return values", types.ErrInvalidTypes,
		)
	}
	if fType.Out(0) != fType.In(0) {
		return nil, nil, fmt.Errorf(
			"%w: first return value of custom merge function must be of the same type as the first argument",
			types.ErrInvalidTypes,
		)
	}
	if fType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
		return nil, nil, fmt.Errorf(
			"%w: second return value of custom merge function must be of type error",
			types.ErrInvalidTypes,
		)
	}

	return func(src reflect.Value, override reflect.Value) (res reflect.Value, err error) {
		// call method
		results := fValue.Call([]reflect.Value{src, override})
		if !results[1].IsZero() && !results[1].IsNil() {
			err = results[1].Interface().(error)
		}
		return results[0], err
	}, fType.Out(0), nil
}

func (mf CustomMergeFuncs) toReflectMergeFuncs() (customReflectMergeFuncs, error) {
	result := make(customReflectMergeFuncs, len(mf))
	for k, v := range mf {
		mergeFn, argType, err := asCustomReflectMergeFunc(v)
		if err != nil {
			return nil, fmt.Errorf("invalid custom merge func for type '%s': %w", k.Name(), err)
		}
		if argType != k {
			return nil, fmt.Errorf(
				"custom merge func for type '%s' has different args and return type '%s'",
				k.Name(), argType.Name(),
			)
		}
		result[k] = mergeFn
	}
	return result, nil
}
