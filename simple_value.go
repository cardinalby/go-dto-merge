package dtomerge

import (
	"reflect"
)

func mergeSimpleValue(src reflect.Value, patch reflect.Value, _ Options) (reflect.Value, error) {
	if patch.IsZero() || reflect.DeepEqual(src, patch) {
		return src, nil
	}
	result := reflect.New(src.Type()).Elem()
	result.Set(patch)
	return result, nil
}
