package dtomerge

import (
	"reflect"
	"unsafe"
)

type mergeContext struct {
	options        Options
	mergedPointers map[unsafe.Pointer]reflect.Value
}

func newMergeContext(options Options) mergeContext {
	return mergeContext{
		options:        options,
		mergedPointers: make(map[unsafe.Pointer]reflect.Value),
	}
}
