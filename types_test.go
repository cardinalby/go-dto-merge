package dtomerge

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type mergerTestType1 struct {
	f string
}

func (m mergerTestType1) String() string {
	return m.f
}

func (m mergerTestType1) Merge(override mergerTestType1) (mergerTestType1, error) {
	if override.f != "" {
		m.f = override.f
	}
	return m, nil
}

type mergerTestType1Ptr struct {
	f string
}

func (m *mergerTestType1Ptr) Merge(override *mergerTestType1Ptr) (*mergerTestType1Ptr, error) {
	if override == nil {
		return m, nil
	}
	if m == nil {
		return override, nil
	}
	mCopy := *m
	if override.f != "" {
		mCopy.f = override.f
	}
	return &mCopy, nil
}

type mergerTestType1PtrWrong struct {
}

func (m *mergerTestType1PtrWrong) Merge(_ mergerTestType1PtrWrong) (mergerTestType1PtrWrong, error) {
	return mergerTestType1PtrWrong{}, nil
}

type mergerTestType2 struct {
}

func (m mergerTestType2) String() string {
	return ""
}

func (m mergerTestType2) Merge(_ mergerTestType1) (mergerTestType1, error) {
	return mergerTestType1{}, nil
}

func TestAsMergerTestType1(t *testing.T) {
	t.Parallel()

	t.Run("mergerTestType1", func(t *testing.T) {
		t.Parallel()
		val := mergerTestType1{}
		fn, ok := AsMerger(reflect.ValueOf(val))
		require.True(t, ok)
		resultRValue, err := fn(reflect.ValueOf(mergerTestType1{f: "test"}))
		require.NoError(t, err)
		require.Equal(t, "test", resultRValue.Interface().(mergerTestType1).f)
	})

	t.Run("mergerTestType1Ptr", func(t *testing.T) {
		t.Parallel()
		var anyM any = &mergerTestType1Ptr{}
		fn, ok := AsMerger(reflect.ValueOf(anyM))
		require.True(t, ok)
		resultRValue, err := fn(reflect.ValueOf(&mergerTestType1Ptr{f: "test"}))
		require.NoError(t, err)
		require.Equal(t, "test", resultRValue.Interface().(*mergerTestType1Ptr).f)
	})

	t.Run("mergerTestType1_any", func(t *testing.T) {
		t.Parallel()
		var anyM any = mergerTestType1{}
		fn, ok := AsMerger(reflect.ValueOf(anyM))
		require.True(t, ok)
		resultRValue, err := fn(reflect.ValueOf(mergerTestType1{f: "test"}))
		require.NoError(t, err)
		require.Equal(t, "test", resultRValue.Interface().(mergerTestType1).f)
	})

	t.Run("mergerTestType1PtrWrong", func(t *testing.T) {
		t.Parallel()
		var val any = &mergerTestType1PtrWrong{}
		_, ok := AsMerger(reflect.ValueOf(val))
		require.False(t, ok)
	})

	t.Run("mergerTestType2", func(t *testing.T) {
		t.Parallel()
		var val any = mergerTestType2{}
		_, ok := AsMerger(reflect.ValueOf(val))
		require.False(t, ok)
	})
}

func TestAsCustomReflectMergeFunc(t *testing.T) {
	t.Parallel()
	_, _, err := asCustomReflectMergeFunc(nil)
	require.Error(t, err)

	fn, rType, err := asCustomReflectMergeFunc(func(src int, override int) (int, error) {
		return src + override, nil
	})
	require.NoError(t, err)
	require.Equal(t, reflect.TypeOf(0), rType)
	fnRes, fnErr := fn(reflect.ValueOf(1), reflect.ValueOf(2))
	require.NoError(t, fnErr)
	require.Equal(t, 3, fnRes.Interface().(int))
}

func TestAsCustomReflectMergeFunc2(t *testing.T) {

}
