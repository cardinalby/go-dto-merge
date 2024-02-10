package dtomerge

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeComparableSlicesUnique(t *testing.T) {
	src := []int{1, 2, 3}
	patch := []int{3, 4, 5}
	merged, err := mergeComparableSlicesUnique(reflect.ValueOf(src), reflect.ValueOf(patch), Options{})
	require.NoError(t, err)
	require.Equal(t, []int{1, 2, 3, 4, 5}, merged.Interface().([]int))
	require.Equal(t, []int{1, 2, 3}, src)
}

func TestMergeSlicesByIndex(t *testing.T) {
	t.Run("ints", func(t *testing.T) {
		src := []int{1, 2, 3}
		patch := []int{3, 0, 5, 6}
		merged, err := mergeSlicesByIndex(reflect.ValueOf(src), reflect.ValueOf(patch), Options{})
		require.NoError(t, err)
		require.Equal(t, []int{3, 2, 5, 6}, merged.Interface().([]int))
		require.Equal(t, []int{1, 2, 3}, src)
	})

	t.Run("funcs", func(t *testing.T) {
		src := []func() int{
			func() int { return 1 },
			func() int { return 2 },
		}
		patch := []func() int{
			func() int { return 3 },
			nil,
			func() int { return 4 },
		}
		merged, err := mergeSlicesByIndex(reflect.ValueOf(src), reflect.ValueOf(patch), Options{})
		require.NoError(t, err)
		require.Len(t, merged.Interface().([]func() int), 3)
		require.Equal(t, 3, merged.Index(0).Interface().(func() int)())
		require.Equal(t, 2, merged.Index(1).Interface().(func() int)())
		require.Equal(t, 4, merged.Index(2).Interface().(func() int)())
	})
}
