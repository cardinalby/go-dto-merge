package dtomerge

import (
	"fmt"
	"reflect"

	"github.com/cardinalby/go-dto-merge/types"
)

func mergeAny(src reflect.Value, patch reflect.Value, mCtx mergeContext) (reflect.Value, error) {
	srcType := src.Type()
	patchType := patch.Type()
	if srcType != patchType {
		return reflect.Value{},
			fmt.Errorf("%w: src type '%s' != patch type '%s'",
				types.ErrInvalidTypes, srcType.Name(), patchType.Name())
	}

	customMergeFuncs, err := mCtx.options.getCustomReflectMergeFuncs()
	if err != nil {
		return reflect.Value{}, err
	}

	if mergeFn, has := customMergeFuncs[srcType]; has {
		return mergeFn(src, patch)
	}
	if mCtx.options.RespectMergers {
		mergeFn, isMerger := AsMerger(src)
		if isMerger {
			return mergeFn(patch)
		}
	}
	if customOptions, ok := mCtx.options.CustomMergeOptions[srcType]; ok {
		mCtx.options = customOptions
	} else if mCtx.options.RespectMergeOptionsProviders {
		if optionsProvider, ok := src.Interface().(MergeOptionsProvider); ok {
			mCtx.options = optionsProvider.GetMergeOptions()
		}
	}

	if mCtx.options.DeRefPointers && srcType.Kind() == reflect.Ptr {
		return mergePointers(src, patch, mCtx)
	}

	if srcType.Kind() == reflect.Struct {
		if _, isAtomic := mCtx.options.getAtomicTypesMap()[srcType]; !isAtomic {
			return mergeStructs(src, patch, mCtx)
		}
	}

	if mCtx.options.IterateMaps && srcType.Kind() == reflect.Map {
		return mergeMaps(src, patch, mCtx)
	}

	if srcType.Kind() == reflect.Slice {
		if srcType.Elem().Comparable() && mCtx.options.SlicesMerge == SlicesMergeStrategyUnique {
			return mergeComparableSlicesUnique(src, patch)
		}
		if mCtx.options.SlicesMerge == SlicesMergeStrategyByIndex {
			return mergeSlicesByIndex(src, patch, mCtx)
		}
	}

	return mergeSimpleValue(src, patch)
}
