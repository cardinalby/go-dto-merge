package dtomerge

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergePointers(t *testing.T) {
	t.Parallel()

	t.Run("both filled", func(t *testing.T) {
		t.Parallel()
		src := 1
		patch := 2
		res, err := mergePointers(reflect.ValueOf(&src), reflect.ValueOf(&patch), Options{})
		require.NoError(t, err)
		require.Equal(t, 2, res.Elem().Interface())
	})

	t.Run("nil patch", func(t *testing.T) {
		t.Parallel()
		src := 1
		res, err := mergePointers(reflect.ValueOf(&src), reflect.ValueOf((*int)(nil)), Options{})
		require.NoError(t, err)
		require.Equal(t, 1, res.Elem().Interface())
	})

	t.Run("zero filled patch", func(t *testing.T) {
		t.Parallel()
		src := 1
		patch := 0
		res, err := mergePointers(reflect.ValueOf(&src), reflect.ValueOf(&patch), Options{})
		require.NoError(t, err)
		require.Equal(t, 1, res.Elem().Interface())
	})

	t.Run("zero filled patch with nil src", func(t *testing.T) {
		t.Parallel()
		patch := 0
		res, err := mergePointers(reflect.ValueOf((*int)(nil)), reflect.ValueOf(&patch), Options{})
		require.NoError(t, err)
		require.Equal(t, 0, res.Elem().Interface())
	})
}
