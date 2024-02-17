package dtomerge

import (
	"fmt"
	"reflect"
)

func mergeMaps(src, patch reflect.Value, mCtx mergeContext) (reflect.Value, error) {
	result := reflect.MakeMap(src.Type())
	for _, key := range src.MapKeys() {
		result.SetMapIndex(key, src.MapIndex(key))
	}
	for _, key := range patch.MapKeys() {
		patchValue := patch.MapIndex(key)
		srcValue := result.MapIndex(key)
		resultValue := patchValue
		if srcValue.IsValid() {
			var err error
			resultValue, err = mergeAny(srcValue, patchValue, mCtx)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("error merging map key '%v': %w", key, err)
			}
		}
		result.SetMapIndex(key, resultValue)
	}
	return result, nil
}
