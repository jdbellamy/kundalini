package slices_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	. "gitlab.com/jdbellamy/kundalini"
	"gitlab.com/jdbellamy/kundalini/slices"
)

var noop Fn = func(x interface{}) interface{} { return nil }

func TestTypes_Strings(t *testing.T) {
	v := []string{"a", "b"}
	actual, err := Wrap(v).Types().Release()
	expected := []reflect.Type{reflect.TypeOf(""), reflect.TypeOf("")}
	assert.NoError(t, err)
	assert.IsType(t, []reflect.Type{}, actual)
	assert.Equal(t, expected, actual)
}

func TestConcat_OpIsSliceOfSingleValue(t *testing.T) {
	v := []string{}
	actual, err := Wrap(v).
		Concat([]string{"a"}).
		Release()
	expected := []string{"a"}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []string{}, actual)
}

func TestConcat_OpIsEmptySlice(t *testing.T) {
	v := []string{}
	actual, err := Wrap(v).
		Concat([]string{}).
		Release()
	expected := []string{}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []string{}, actual)
}

func TestConcat_ChainConcat(t *testing.T) {
	v := []string{"1"}
	actual, err := Wrap(v).
		Concat([]string{"2"}).
		Concat([]string{"3", "4"}).
		Concat([]string{"5"}).
		Concat([]string{"6"}).
		Concat([]string{"7"}).
		Release()
	expected := []string{"1", "2", "3", "4", "5", "6", "7"}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []string{}, actual)
}

func TestConcat_ChainConcat_TypeMismatchError(t *testing.T) {
	v := []string{"1"}
	actual, err := Wrap(v).
		Concat([]string{"2"}).
		Concat([]int{3}).
		Concat([]string{"4"}).
		Release()
	assert.EqualError(t, err, slices.TypeMismatchError.Error())
	assert.Nil(t, actual)
}

func TestMap_ChainConcat_TypeMismatchError(t *testing.T) {
	v := []int{1}
	double := func(x interface{}) interface{} {
		return x.(int) * 2
	}
	actual, err := Wrap(v).
		Concat("2").
		Map(double).
		Release()
	assert.EqualError(t, err, slices.TypeMismatchError.Error())
	assert.Nil(t, actual)
}

func TestMap_ChainConcat_DoubeValues(t *testing.T) {
	v := []int{1}
	double := func(x interface{}) interface{} {
		return x.(int) * 2
	}
	actual, err := Wrap(v).
		Concat([]int{2}).
		Concat([]int{3, 4}).
		Map(double).
		Concat([]int{9, 10, 11}).
		Release()
	expected := []int{2, 4, 6, 8, 9, 10, 11}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []int{}, actual)
}

func TestMap_KIsEmpty(t *testing.T) {
	v := []int{}
	double := func(x interface{}) interface{} {
		return x.(int) * 2
	}
	actual, err := Wrap(v).
		Map(double).
		Release()
	expected := []int{}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []int{}, actual)
}

func TestFilter_InitiallyEmpty(t *testing.T) {
	v := []int{}
	p := func(x interface{}) bool {
		return true
	}
	actual, err := Wrap(v).
		Filter(p).
		Release()
	expected := []int{}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []int{}, actual)
}

func TestFilter_RetainsWhenValueIs100(t *testing.T) {
	v := []int{0, 1, 2, 3, 100, 4, 5}
	isOneHundred := func(x interface{}) bool {
		if x == 100 {
			return true
		}
		return false
	}
	actual, err := Wrap(v).
		Concat([]int{6, 100, 7, 100}).
		Filter(isOneHundred).
		Release()
	expected := []int{100, 100, 100}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []int{}, actual)
}

func TestFilter_ChainConcat_TypeMismatchError(t *testing.T) {
	v := []int{0, 1}
	p := func(x interface{}) bool {
		return false
	}
	actual, err := Wrap(v).
		Concat("2").
		Filter(p).
		Release()
	errMsg := "type mismatch between wrapped value and operand"
	assert.EqualError(t, err, errMsg)
	assert.Nil(t, actual)
}

func TestReduce_AccIsScalar(t *testing.T) {
	v := []int{0, 1, 2}
	sum := func(acc interface{}, x interface{}) interface{} {
		return acc.(int) + x.(int)
	}
	actual, err := Wrap(v).
		Reduce(0, sum).
		Release()
	expected := []int{3}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []int{}, actual)
}

func TestReduce_AccIsSlice(t *testing.T) {
	v := []int{0, 1, 2, 3, 4, 5}
	progressiveSum := func(acc interface{}, x interface{}) interface{} {
		sliceAcc := acc.([]int)
		tailInx := len(sliceAcc) - 1
		n := sliceAcc[tailInx] + x.(int)
		return append(sliceAcc, n)
	}
	actual, err := Wrap(v).
		Reduce([]int{0}, progressiveSum).
		Release()
	expected := []int{0, 0, 1, 3, 6, 10, 15}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []int{}, actual)
}

func TestReduce_InitiallyEmpty(t *testing.T) {
	v := []int{}
	sum := func(acc interface{}, x interface{}) interface{} {
		return acc.(int) + x.(int)
	}
	actual, err := Wrap(v).
		Reduce(0, sum).
		Release()
	expected := []int{}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []int{}, actual)
}

func TestReduce_ChainConcat_TypeMismatchError(t *testing.T) {
	v := []int{}
	sum := func(acc interface{}, x interface{}) interface{} {
		return acc.(int) + x.(int)
	}
	actual, err := Wrap(v).
		Concat("1").
		Reduce(0, sum).
		Release()
	errMsg := "type mismatch between wrapped value and operand"
	assert.EqualError(t, err, errMsg)
	assert.Nil(t, actual)
}

func TestReduce_All(t *testing.T) {
	v := []int{0, 1, 2, 3, 4, 5}
	progressiveSum := func(acc interface{}, x interface{}) interface{} {
		sliceAcc := acc.([]int)
		tailInx := len(sliceAcc) - 1
		n := sliceAcc[tailInx] + x.(int)
		return append(sliceAcc, n)
	}
	even := func(x interface{}) bool {
		return x.(int)%2 == 0
	}
	double := func(x interface{}) interface{} {
		return x.(int) * 2
	}
	sum := func(acc interface{}, x interface{}) interface{} {
		return acc.(int) + x.(int)
	}
	actual, err := Wrap(v).
		Concat([]int{6}).
		Concat([]int{7, 8}).
		Reduce([]int{0}, progressiveSum).
		Filter(even).
		Map(double).
		Reduce(0, sum).
		Release()
	expected := []int{160}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []int{}, actual)
}

func TestMap_Types_FnReturnsNil(t *testing.T) {
	v := []int{1}
	actual, err := Wrap(v).Map(noop).Release()
	expected := []int{1}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExport_WritesToPointer(t *testing.T) {
	v := []int{1, 2, 3}
	actual := []int{}
	Wrap(v).
		Concat([]int{4, 5, 6}).
		Export(reflect.ValueOf(&actual)).
		ReleaseOrPanic()
	expected := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, expected, actual)
}

func TestExport_ErrorsOnTypeMismatch(t *testing.T) {
	v := []int{}
	panics := func() {
		Wrap([]int{}).
			Concat([]int{}).
			Export(reflect.ValueOf(v)).
			Release()
	}
	assert.Panics(t, panics)
}
