package kundalini_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "gitlab.com/jdbellamy/kundalini"
)

func TestConcat_OpIsSingleValue(t *testing.T) {

	v := []int{}

	actual, err := Coil(v).
		Concat(0).
		Release()

	expected := []int{0}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []int{}, actual)
}

func TestConcat_OpIsSliceOfSingleValue(t *testing.T) {

	v := []string{}

	actual, err := Coil(v).
		Concat([]string{"a"}).
		Release()

	expected := []string{"a"}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []string{}, actual)
}

func TestConcat_OpIsEmptySlice(t *testing.T) {

	v := []string{}

	actual, err := Coil(v).
		Concat([]string{}).
		Release()

	expected := []string{}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []string{}, actual)
}

func TestConcat_OpTypeMismatch(t *testing.T) {

	v := []string{}

	actual, err := Coil(v).
		Concat([]int{}).
		Release()

	errString := "type mismatch between wrapped value and operand"

	assert.Error(t, err)
	assert.EqualError(t, err, errString)
	assert.Nil(t, actual)
}

func TestConcat_ChainConcat(t *testing.T) {

	v := []string{"1"}

	actual, err := Coil(v).
		Concat("2").
		Concat([]string{"3", "4"}).
		Concat([]string{"5"}).
		Concat("6").
		Concat([]string{"7"}).
		Release()

	expected := []string{"1", "2", "3", "4", "5", "6", "7"}

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
	assert.IsType(t, []string{}, actual)
}

func TestConcat_ChainConcat_TypeMismatchError(t *testing.T) {

	v := []string{"1"}

	actual, err := Coil(v).
		Concat([]string{"2"}).
		Concat([]int{3}).
		Concat([]string{"4"}).
		Release()

	errMsg := "type mismatch between wrapped value and operand"

	assert.EqualError(t, err, errMsg)
	assert.Nil(t, actual)
}

func TestMap_ChainConcat_TypeMismatchError(t *testing.T) {

	v := []int{1}

	double := func(x interface{}) interface{} {
		return x.(int) * 2
	}

	actual, err := Coil(v).
		Concat("2").
		Map(double).
		Release()

	errMsg := "type mismatch between wrapped value and operand"

	assert.EqualError(t, err, errMsg)
	assert.Nil(t, actual)
}

func TestMap_ChainConcat_DoubeValues(t *testing.T) {

	v := []int{1}

	double := func(x interface{}) interface{} {
		return x.(int) * 2
	}

	actual, err := Coil(v).
		Concat(2).
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

	actual, err := Coil(v).
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

	actual, err := Coil(v).
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

	actual, err := Coil(v).
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

	actual, err := Coil(v).
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

	actual, err := Coil(v).
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

	actual, err := Coil(v).
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

	actual, err := Coil(v).
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

	actual, err := Coil(v).
		Concat("1").
		Reduce(0, sum).
		Release()

	errMsg := "type mismatch between wrapped value and operand"

	assert.EqualError(t, err, errMsg)
	assert.Nil(t, actual)
}

func TestReduce_All(t *testing.T) {

	v := []int{0, 1, 2, 3, 4, 5}

	sum := func(acc interface{}, x interface{}) interface{} {
		return acc.(int) + x.(int)
	}

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

	actual, err := Coil(v).
		Concat(6).
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
