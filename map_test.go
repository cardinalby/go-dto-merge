package dtomerge

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMergeMaps(t *testing.T) {
	t.Parallel()

	t.Run("comparable", func(t *testing.T) {
		src := map[string]string{
			"a": "b",
			"c": "d",
			"d": "e",
		}
		patch := map[string]string{
			"a": "bb",
			"c": "",
			"f": "g",
			"h": "i",
		}
		res, err := mergeMaps(reflect.ValueOf(src), reflect.ValueOf(patch), NewOptions())
		require.NoError(t, err)
		require.Equal(t, map[string]string{
			"a": "bb",
			"c": "d",
			"d": "e",
			"f": "g",
			"h": "i",
		}, res.Interface())
		require.Equal(t, map[string]string{
			"a": "b",
			"c": "d",
			"d": "e",
		}, src)
	})

	t.Run("non-comparable", func(t *testing.T) {
		src := map[string]func() string{
			"a": func() string { return "src" },
			"b": nil,
			"c": func() string { return "src" },
			"g": func() string { return "src" },
		}
		patch := map[string]func() string{
			"a": func() string { return "patch" },
			"b": func() string { return "patch" },
			"c": nil,
			"f": func() string { return "patch" },
		}
		res, err := mergeMaps(reflect.ValueOf(src), reflect.ValueOf(patch), NewOptions())
		require.NoError(t, err)
		require.Len(t, res.Interface(), 5)
		resMap := res.Interface().(map[string]func() string)
		require.Equal(t, "patch", resMap["a"]())
		require.Equal(t, "patch", resMap["b"]())
		require.Equal(t, "src", resMap["c"]())
		require.Equal(t, "patch", resMap["f"]())
		require.Equal(t, "src", resMap["g"]())
		require.Equal(t, "src", src["a"]())
	})
}
