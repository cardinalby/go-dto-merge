package dtomerge

import (
	"reflect"
)

// SlicesMergeStrategy defines how to merge slices
type SlicesMergeStrategy string

const (
	// SlicesMergeStrategyAtomic considers slices as atomic values. Patch will entirely replace src if patch is not nil.
	SlicesMergeStrategyAtomic SlicesMergeStrategy = "atomic"

	// SlicesMergeStrategyUnique considers slices as sets of unique elements.
	// Patch element will be appended to src slice if it's not present in src.
	// If slice contains not comparable elements, it roll backs SlicesMergeStrategyAtomic
	SlicesMergeStrategyUnique SlicesMergeStrategy = "unique"

	// SlicesMergeStrategyByIndex merges elements with the same indexes.
	// Patch slice elements will replace the src element if the patch slice element is not a zero value
	// If patch slice is longer than src slice, remaining patch slice elements will be appended to src slice
	SlicesMergeStrategyByIndex SlicesMergeStrategy = "by_index"
)

type Options struct {
	// internal caches
	customReflectMergeFuncs customReflectMergeFuncs
	atomicTypesMap          map[reflect.Type]struct{}

	// RespectMergers sets whether to respect Merger interface for fields. If true, fields that implement
	// Merger will be merged using their Merge method.
	// Default: true
	// Use OptRespectMergers to set this option
	RespectMergers bool

	// CustomMergeFuncs sets custom merge functions for specific types. If merge function is set for a type,
	// it will be used instead of default merge logic.
	// See CustomMergeFuncs for more details about functions signature.
	// Use OptCustomMergeFuncs to set this option
	CustomMergeFuncs CustomMergeFuncs

	// RespectMergeOptionsProviders sets whether to respect MergeOptionsProvider interface for fields. If true,
	// fields that implement MergeOptionsProvider will be merged with options returned from their GetMergeOptions method.
	// If CustomMergeOptions is set for a type, it will be used instead.
	// Default: true
	// Use OptRespectMergeOptionsProviders to set this option
	RespectMergeOptionsProviders bool

	// CustomMergeOptions sets custom merge options for specific types. If merge options are set for a type,
	// they will be used instead of merge options passed to Merge function.
	// Use OptCustomMergeOptions to set this option
	CustomMergeOptions CustomMergeOptions

	// DeRefPointers sets comparison mode for fields containing pointers (only if they are not Merger or
	// have defined custom merge function)
	//   - If true, pointers will be de-referenced (if not nil) and their values will be compared.
	//     if `patch` pointer is nil, `src` pointer will be used.
	//   - If false, pointers will be compared directly.
	//
	// Default: true
	// Use OptDeRefPointers to set this option
	DeRefPointers bool

	// sets types that should be treated as atomic. For a struct it means fields will not be iterated
	// and two structs will be compared with reflect.DeepEqual. If they are not equal, patch value will be used.
	// Use OptAtomicTypes to set this option
	AtomicTypes AtomicTypes

	// IterateMaps sets whether to iterate maps considering their keys as individual values.
	// If true, maps will be iterated and merged by keys (calling merge func with options recursively for values).
	// If false, maps will be compared with reflect.DeepEqual and patch value will be used if they are not equal.
	// Default: false
	// Use OptIterateMaps to set this option
	IterateMaps bool

	// MergeSlices sets how to merge slices.
	// See SlicesMergeStrategy for more details.
	// Default: SlicesMergeStrategyAtomic
	// Use OptMergeSlices to set this option
	SlicesMerge SlicesMergeStrategy
}

type Option func(*Options)

func NewOptions(opts ...Option) Options {
	options := Options{
		DeRefPointers:                true,
		RespectMergers:               true,
		RespectMergeOptionsProviders: true,
		SlicesMerge:                  SlicesMergeStrategyAtomic,
	}
	for _, option := range opts {
		option(&options)
	}
	return options
}

func (o Options) getCustomReflectMergeFuncs() (customReflectMergeFuncs, error) {
	if o.customReflectMergeFuncs == nil {
		var err error
		o.customReflectMergeFuncs, err = o.CustomMergeFuncs.toReflectMergeFuncs()
		if err != nil {
			return nil, err
		}
	}
	return o.customReflectMergeFuncs, nil
}

func (o Options) getAtomicTypesMap() map[reflect.Type]struct{} {
	if o.atomicTypesMap == nil {
		o.atomicTypesMap = make(map[reflect.Type]struct{})
		for _, sType := range o.AtomicTypes {
			o.atomicTypesMap[sType] = struct{}{}
		}
	}
	return o.atomicTypesMap
}
