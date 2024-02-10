package dtomerge

// OptDeRefPointers Sets Options.DeRefPointers
func OptDeRefPointers(deRefPointers bool) Option {
	return func(o *Options) {
		o.DeRefPointers = deRefPointers
	}
}

// OptRespectMergers sets Options.RespectMergers
func OptRespectMergers(respectMergers bool) Option {
	return func(o *Options) {
		o.RespectMergers = respectMergers
	}
}

// OptAtomicTypes sets Options.AtomicTypes
func OptAtomicTypes(atomicTypes AtomicTypes) Option {
	return func(o *Options) {
		o.AtomicTypes = atomicTypes
	}
}

// OptCustomMergeFuncs sets Options.CustomMergeFuncs
func OptCustomMergeFuncs(customMergeFuncs CustomMergeFuncs) Option {
	return func(o *Options) {
		o.CustomMergeFuncs = customMergeFuncs
	}
}

// OptRespectMergeOptionsProviders sets Options.RespectMergeOptionsProviders
func OptRespectMergeOptionsProviders(respectMergeOptionsProviders bool) Option {
	return func(o *Options) {
		o.RespectMergeOptionsProviders = respectMergeOptionsProviders
	}
}

// OptCustomMergeOptions sets Options.CustomMergeOptions
func OptCustomMergeOptions(customMergeOptions CustomMergeOptions) Option {
	return func(o *Options) {
		o.CustomMergeOptions = customMergeOptions
	}
}

// OptIterateMaps sets Options.IterateMaps
func OptIterateMaps(iterateMaps bool) Option {
	return func(o *Options) {
		o.IterateMaps = iterateMaps
	}
}

// OptMergeSlices sets Options.SlicesMerge
func OptMergeSlices(iterateSlices SlicesMergeStrategy) Option {
	return func(o *Options) {
		o.SlicesMerge = iterateSlices
	}
}
