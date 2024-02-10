package dtomerge

import (
	"fmt"
	"reflect"

	"github.com/cardinalby/go-dto-merge/types"
)

func mergeAny(src reflect.Value, patch reflect.Value, opts Options) (reflect.Value, error) {
	srcType := src.Type()
	patchType := patch.Type()
	if srcType != patchType {
		return reflect.Value{},
			fmt.Errorf("%w: src type '%s' != patch type '%s'",
				types.ErrInvalidTypes, srcType.Name(), patchType.Name())
	}

	customMergeFuncs, err := opts.getCustomReflectMergeFuncs()
	if err != nil {
		return reflect.Value{}, err
	}

	if mergeFn, has := customMergeFuncs[srcType]; has {
		return mergeFn(src, patch)
	}
	if opts.RespectMergers {
		mergeFn, isMerger := AsMerger(src)
		if isMerger {
			return mergeFn(patch)
		}
	}
	if customOptions, ok := opts.CustomMergeOptions[srcType]; ok {
		opts = customOptions
	} else if opts.RespectMergeOptionsProviders {
		if optionsProvider, ok := src.Interface().(MergeOptionsProvider); ok {
			opts = optionsProvider.GetMergeOptions()
		}
	}

	if opts.DeRefPointers && srcType.Kind() == reflect.Ptr {
		return mergePointers(src, patch, opts)
	}

	if srcType.Kind() == reflect.Struct {
		if _, isAtomic := opts.getAtomicTypesMap()[srcType]; !isAtomic {
			return mergeStructs(src, patch, opts)
		}
	}

	if opts.IterateMaps && srcType.Kind() == reflect.Map {
		return mergeMaps(src, patch, opts)
	}

	if srcType.Kind() == reflect.Slice {
		if srcType.Elem().Comparable() && opts.SlicesMerge == SlicesMergeStrategyUnique {
			return mergeComparableSlicesUnique(src, patch, opts)
		}
		if opts.SlicesMerge == SlicesMergeStrategyByIndex {
			return mergeSlicesByIndex(src, patch, opts)
		}
	}

	return mergeSimpleValue(src, patch, opts)
}
