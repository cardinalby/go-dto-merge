package tests

import (
	"reflect"
	"strings"
	"testing"

	dtomerge "github.com/cardinalby/go-dto-merge"
	"github.com/stretchr/testify/require"
)

type CustomStr string

func (c CustomStr) Merge(override CustomStr) (CustomStr, error) {
	if strings.TrimSpace(string(override)) != "" {
		return override, nil
	}
	return c, nil
}

type TestNestedStruct struct {
	prInt     int
	Int       int
	Str       string
	CustomStr CustomStr
}

type TestNestedStruct2 struct {
	IntPtr  *int
	BoolPtr *bool
	Slice   []string
}

func (t TestNestedStruct2) GetMergeOptions() dtomerge.Options {
	return dtomerge.Options{
		DeRefPointers: false,
	}
}

type TestStruct struct {
	N1   TestNestedStruct
	N2   *TestNestedStruct2
	Bool bool
}

func ptr[T any](v T) *T {
	return &v
}

func checkPtrEquals[T comparable](t *testing.T, exp, act *T) {
	t.Helper()
	require.Equal(t, exp == nil, act == nil)
	if act != nil {
		require.Equal(t, *exp, *act)
	}
}

func checkEquals(t *testing.T, exp, act TestStruct) {
	t.Helper()
	require.Equal(t, exp.N1.prInt, act.N1.prInt)
	require.Equal(t, exp.N1.Int, act.N1.Int)
	require.Equal(t, exp.N1.Str, act.N1.Str)
	require.Equal(t, exp.N1.CustomStr, act.N1.CustomStr)

	checkPtrEquals(t, exp.N2.IntPtr, act.N2.IntPtr)
	checkPtrEquals(t, exp.N2.BoolPtr, act.N2.BoolPtr)
	require.EqualValues(t, exp.N2.Slice, act.N2.Slice)
	require.Equal(t, exp.Bool, act.Bool)
}

func TestMerge(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		t.Parallel()
		src := TestStruct{
			N1: TestNestedStruct{
				prInt:     23,
				Int:       1,
				Str:       "a",
				CustomStr: "b",
			},
			N2: &TestNestedStruct2{
				IntPtr:  ptr(2),
				BoolPtr: ptr(true),
				Slice:   []string{"a", "b"},
			},
			Bool: true,
		}
		patch := TestStruct{
			N1: TestNestedStruct{
				prInt:     24,
				Int:       0,
				Str:       "b",
				CustomStr: "  ",
			},
			N2: &TestNestedStruct2{
				IntPtr:  ptr(0),
				BoolPtr: ptr(false),
				Slice:   []string{"c", "d"},
			},
			Bool: false,
		}

		t.Run("default options", func(t *testing.T) {
			t.Parallel()
			res, err := dtomerge.Merge(src, patch)
			require.NoError(t, err)
			checkEquals(t, TestStruct{
				N1: TestNestedStruct{
					prInt:     23,
					Int:       1,
					Str:       "b",
					CustomStr: "b",
				},
				N2: &TestNestedStruct2{
					IntPtr:  ptr(0),
					BoolPtr: ptr(false),
					Slice:   []string{"c", "d"},
				},
				Bool: true,
			}, res)
		})

		t.Run("do not respect interfaces", func(t *testing.T) {
			t.Parallel()
			res, err := dtomerge.Merge(src, patch,
				dtomerge.OptRespectMergeOptionsProviders(false),
				dtomerge.OptRespectMergers(false),
			)
			require.NoError(t, err)
			checkEquals(t, TestStruct{
				N1: TestNestedStruct{
					prInt:     23,
					Int:       1,
					Str:       "b",
					CustomStr: "  ",
				},
				N2: &TestNestedStruct2{
					IntPtr:  ptr(2),
					BoolPtr: ptr(true),
					Slice:   []string{"c", "d"},
				},
				Bool: true,
			}, res)
		})

		t.Run("atomic types, dont deref pointers", func(t *testing.T) {
			t.Parallel()
			res, err := dtomerge.Merge(src, patch,
				dtomerge.OptRespectMergeOptionsProviders(false),
				dtomerge.OptRespectMergers(false),
				dtomerge.OptDeRefPointers(false),
				dtomerge.OptAtomicTypes(dtomerge.AtomicTypes{reflect.TypeOf(TestNestedStruct{})}),
			)
			require.NoError(t, err)
			checkEquals(t, TestStruct{
				N1: TestNestedStruct{
					prInt:     24,
					Int:       0,
					Str:       "b",
					CustomStr: "  ",
				},
				N2: &TestNestedStruct2{
					IntPtr:  ptr(0),
					BoolPtr: ptr(false),
					Slice:   []string{"c", "d"},
				},
				Bool: true,
			}, res)
		})

		t.Run("with custom merge funcs", func(t *testing.T) {
			t.Parallel()

			src := TestStruct{
				N1: TestNestedStruct{
					prInt:     23,
					Int:       1,
					Str:       "a",
					CustomStr: "b",
				},
				N2: &TestNestedStruct2{
					IntPtr:  ptr(2),
					BoolPtr: ptr(true),
					Slice:   []string{"a", "b"},
				},
				Bool: true,
			}
			patch := TestStruct{
				N1: TestNestedStruct{
					prInt:     24,
					Int:       0,
					Str:       "b",
					CustomStr: "--",
				},
				N2: &TestNestedStruct2{
					IntPtr:  ptr(0),
					BoolPtr: ptr(false),
					Slice:   []string{"c", "d"},
				},
				Bool: false,
			}

			res, err := dtomerge.Merge(src, patch, dtomerge.OptCustomMergeFuncs(
				dtomerge.CustomMergeFuncs{
					reflect.TypeOf(CustomStr("")): func(src, override CustomStr) (CustomStr, error) {
						if strings.Trim(string(override), "-") != "" {
							return override, nil
						}
						return src, nil
					},
				}),
				dtomerge.OptRespectMergeOptionsProviders(false),
				dtomerge.OptCustomMergeOptions(dtomerge.CustomMergeOptions{
					reflect.TypeOf(TestNestedStruct2{}): dtomerge.Options{
						DeRefPointers: false,
					},
				}),
			)
			require.NoError(t, err)
			checkEquals(t, TestStruct{
				N1: TestNestedStruct{
					prInt:     23,
					Int:       1,
					Str:       "b",
					CustomStr: "b",
				},
				N2: &TestNestedStruct2{
					IntPtr:  ptr(0),
					BoolPtr: ptr(false),
					Slice:   []string{"c", "d"},
				},
				Bool: true,
			}, res)
		})
	})

	t.Run("slices", func(t *testing.T) {
		t.Parallel()
		src := TestNestedStruct2{
			Slice: []string{"a", "b"},
		}
		patch := TestNestedStruct2{
			Slice: []string{"", "b", "d"},
		}

		t.Run("default options", func(t *testing.T) {
			t.Parallel()
			res, err := dtomerge.Merge(src, patch,
				dtomerge.OptRespectMergeOptionsProviders(false),
			)
			require.NoError(t, err)
			require.Equal(t, []string{"", "b", "d"}, res.Slice)
			require.Equal(t, &patch.Slice[0], &res.Slice[0])
		})

		t.Run("merge unique", func(t *testing.T) {
			t.Parallel()
			res, err := dtomerge.Merge(src, patch,
				dtomerge.OptRespectMergeOptionsProviders(false),
				dtomerge.OptMergeSlices(dtomerge.SlicesMergeStrategyUnique),
			)
			require.NoError(t, err)
			require.Equal(t, []string{"a", "b", "", "d"}, res.Slice)
			res.Slice[0] = "X"
		})

		t.Run("merge by index", func(t *testing.T) {
			t.Parallel()
			res, err := dtomerge.Merge(src, patch,
				dtomerge.OptRespectMergeOptionsProviders(false),
				dtomerge.OptMergeSlices(dtomerge.SlicesMergeStrategyByIndex),
			)
			require.NoError(t, err)
			require.Equal(t, []string{"a", "b", "d"}, res.Slice)
			res.Slice[0] = "X"
		})

		require.Equal(t, []string{"a", "b"}, src.Slice)
		require.Equal(t, []string{"", "b", "d"}, patch.Slice)
	})

	type mTestNestedStruct struct {
		CustomStr CustomStr
	}

	type mTestStruct struct {
		M map[string]mTestNestedStruct
	}

	t.Run("maps", func(t *testing.T) {
		t.Parallel()
		src := mTestStruct{
			M: map[string]mTestNestedStruct{
				"a": {CustomStr: "src_a"},
				"b": {CustomStr: "src_b"},
				"c": {CustomStr: "src_c"},
			},
		}
		patch := mTestStruct{
			M: map[string]mTestNestedStruct{
				"a": {CustomStr: "patch_a"},
				"b": {CustomStr: "  "},
				"d": {CustomStr: "patch_d"},
			},
		}

		t.Run("default options", func(t *testing.T) {
			res, err := dtomerge.Merge(src, patch)
			require.NoError(t, err)
			require.Equal(t, map[string]mTestNestedStruct{
				"a": {CustomStr: "patch_a"},
				"b": {CustomStr: "  "},
				"d": {CustomStr: "patch_d"},
			}, res.M)
			require.Equal(t, &patch.M, &res.M)
		})

		t.Run("iterate maps", func(t *testing.T) {
			res, err := dtomerge.Merge(src, patch, dtomerge.OptIterateMaps(true))
			require.NoError(t, err)
			require.Equal(t, map[string]mTestNestedStruct{
				"a": {CustomStr: "patch_a"},
				"b": {CustomStr: "src_b"},
				"c": {CustomStr: "src_c"},
				"d": {CustomStr: "patch_d"},
			}, res.M)
			res.M["a"] = mTestNestedStruct{CustomStr: "X"}
		})

		require.Equal(t, map[string]mTestNestedStruct{
			"a": {CustomStr: "src_a"},
			"b": {CustomStr: "src_b"},
			"c": {CustomStr: "src_c"},
		}, src.M)
		require.Equal(t, map[string]mTestNestedStruct{
			"a": {CustomStr: "patch_a"},
			"b": {CustomStr: "  "},
			"d": {CustomStr: "patch_d"},
		}, patch.M)
	})
}
