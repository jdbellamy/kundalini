package kundalini_test

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	. "gitlab.com/jdbellamy/kundalini"
)

func TestWrap(t *testing.T) {

	t.Run("should release correct values", func(t *testing.T) {
		type Test struct {
			input    interface{}
			expected interface{}
			err      error
		}

		tests := []Test{{
			input:    []string{},
			expected: []string{},
		}, {
			input:    []string{"0", "a", "b", "C", "∆"},
			expected: []string{"0", "a", "b", "C", "∆"},
		}, {
			input:    []int{},
			expected: []int{},
		}, {
			input:    []int{0, 0, 1, 1000, 21345},
			expected: []int{0, 0, 1, 1000, 21345},
		}, {
			input:    []int{math.MaxUint32},
			expected: []int{math.MaxUint32},
		}, {
			input:    []string{strings.Repeat("∆", 9999)},
			expected: []string{strings.Repeat("∆", 9999)},
		}}

		for _, tt := range tests {
			var v interface{}
			v = tt.input

			actual, err := Wrap(v).Release()

			expected := tt.expected

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})
}

func TestTypes(t *testing.T) {

	t.Run("should return correct types", func(t *testing.T) {
		type Test struct {
			input    interface{}
			expected []reflect.Type
			err      error
		}

		var strT = reflect.TypeOf("")
		var intT = reflect.TypeOf(0)

		tests := []Test{{
			input:    []string{},
			expected: []reflect.Type{},
		}, {
			input:    []string{""},
			expected: []reflect.Type{strT},
		}, {
			input:    []string{"", "", "a", "b", ""},
			expected: []reflect.Type{strT, strT, strT, strT, strT},
		}, {
			input:    []int{},
			expected: []reflect.Type{},
		}, {
			input:    []int{0},
			expected: []reflect.Type{intT},
		}, {
			input:    []int{0, 0, 1, 1000, 21345},
			expected: []reflect.Type{intT, intT, intT, intT, intT},
		}, {
			input:    []int{math.MaxInt8},
			expected: []reflect.Type{intT},
		}, {
			input:    []int{math.MaxUint32},
			expected: []reflect.Type{intT},
		}, {
			input:    []string{strings.Repeat("", 9999)},
			expected: []reflect.Type{strT},
		}, {
			input:    []string{strings.Repeat("∆", 9999)},
			expected: []reflect.Type{strT},
		}}

		for _, tt := range tests {
			var v interface{}
			v = tt.input

			actual, err := Wrap(v).Types().Release()

			expected := tt.expected

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("should raise error when input type is not supported", func(t *testing.T) {
		type Test struct {
			input  interface{}
			actual []reflect.Type
			err    error
		}

		tests := []Test{{
			input: 0,
			err:   UnsupportedWrappedTypeError,
		}, {
			input: "a",
			err:   UnsupportedWrappedTypeError,
		}, {
			input: map[string]string{},
			err:   UnsupportedWrappedTypeError,
		}}

		for _, tt := range tests {
			actual, err := Wrap(tt.input).Types().Release()

			assert.EqualError(t, err, tt.err.Error())
			assert.Nil(t, actual)
		}
	})

	t.Run("should forward received error", func(t *testing.T) {
		var noop Fn = func(x interface{}) interface{} { return nil }

		actual, err := Wrap(0).Map(noop).Types().Release()

		assert.EqualError(t, err, UnsupportedWrappedTypeError.Error())
		assert.Nil(t, actual)
	})
}

func TestMap(t *testing.T) {

	t.Run("should apply Fn correctly", func(t *testing.T) {
		type Test struct {
			input    interface{}
			expected interface{}
			fn       Fn
			err      error
		}

		var noop Fn = func(x interface{}) interface{} { return nil }

		var double Fn = func(x interface{}) interface{} {
			return x.(int) * 2
		}

		var incr Fn = func(x interface{}) interface{} {
			return x.(int) + 1
		}

		tests := []Test{{
			input:    []string{},
			expected: []string{},
			fn:       noop,
		}, {
			input:    []int{1, 2, 3},
			expected: []int{2, 4, 6},
			fn:       double,
		}, {
			input:    []int{1, 2, 3},
			expected: []int{2, 3, 4},
			fn:       incr,
		}}

		for _, tt := range tests {
			var v interface{}
			v = tt.input

			actual, err := Wrap(v).Map(tt.fn).Release()

			expected := tt.expected

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("should raise error when input type is not supported", func(t *testing.T) {
		type Test struct {
			input  interface{}
			actual interface{}
			fn     Fn
			err    error
		}

		var noop Fn = func(x interface{}) interface{} { return nil }

		tests := []Test{{
			input: 0,
			err:   UnsupportedWrappedTypeError,
			fn:    noop,
		}, {
			input: "a",
			err:   UnsupportedWrappedTypeError,
			fn:    noop,
		}, {
			input: map[string]string{},
			err:   UnsupportedWrappedTypeError,
			fn:    noop,
		}}
		for _, tt := range tests {
			actual, err := Wrap(tt.input).Map(tt.fn).Release()
			assert.EqualError(t, err, tt.err.Error())
			assert.Nil(t, actual)
		}
	})

	t.Run("should forward received error", func(t *testing.T) {
		var noop Fn = func(x interface{}) interface{} { return nil }

		actual, err := Wrap(0).Map(noop).Map(noop).Release()

		assert.EqualError(t, err, UnsupportedWrappedTypeError.Error())
		assert.Nil(t, actual)
	})
}

func TestFilter(t *testing.T) {

	t.Run("should filter Predicate correctly", func(t *testing.T) {
		type Test struct {
			input    interface{}
			expected interface{}
			fn       Predicate
			err      error
		}

		var none Predicate = func(x interface{}) bool { return false }

		tests := []Test{{
			input:    []string{},
			expected: []string{},
			fn:       none,
		}}

		for _, tt := range tests {
			var v interface{}
			v = tt.input

			actual, err := Wrap(v).Filter(tt.fn).Release()

			expected := tt.expected

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("should raise error when input type is not supported", func(t *testing.T) {
		type Test struct {
			input  interface{}
			actual interface{}
			fn     Predicate
			err    error
		}

		var none Predicate = func(x interface{}) bool { return false }

		tests := []Test{{
			input: 0,
			err:   UnsupportedWrappedTypeError,
			fn:    none,
		}, {
			input: "a",
			err:   UnsupportedWrappedTypeError,
			fn:    none,
		}, {
			input: map[string]string{},
			err:   UnsupportedWrappedTypeError,
			fn:    none,
		}}

		for _, tt := range tests {
			actual, err := Wrap(tt.input).Filter(tt.fn).Release()

			assert.EqualError(t, err, tt.err.Error())
			assert.Nil(t, actual)
		}
	})

	t.Run("should forward received error", func(t *testing.T) {
		var noop Fn = func(x interface{}) interface{} { return nil }
		var none Predicate = func(x interface{}) bool { return false }

		actual, err := Wrap(0).Map(noop).Filter(none).Release()

		assert.EqualError(t, err, UnsupportedWrappedTypeError.Error())
		assert.Nil(t, actual)
	})
}

func TestReduce(t *testing.T) {

	t.Run("should apply Transform correctly", func(t *testing.T) {
		type Test struct {
			input    interface{}
			expected interface{}
			fn       Transform
			acc      interface{}
			err      error
		}

		var noacc Transform = func(acc interface{}, x interface{}) interface{} { return acc }

		tests := []Test{{
			input:    []string{},
			expected: []string{},
			fn:       noacc,
			acc:      []string{},
		}}

		for _, tt := range tests {
			var v interface{}
			v = tt.input

			actual, err := Wrap(v).Reduce(tt.acc, tt.fn).Release()

			expected := tt.expected

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("should raise error when input type is not supported", func(t *testing.T) {
		type Test struct {
			input  interface{}
			actual interface{}
			fn     Transform
			acc    interface{}
			err    error
		}

		var noacc Transform = func(acc interface{}, x interface{}) interface{} { return acc }

		tests := []Test{{
			input: 0,
			err:   UnsupportedWrappedTypeError,
			fn:    noacc,
			acc:   []int{},
		}, {
			input: "a",
			err:   UnsupportedWrappedTypeError,
			fn:    noacc,
			acc:   []string{},
		}, {
			input: map[string]string{},
			err:   UnsupportedWrappedTypeError,
			fn:    noacc,
			acc:   []map[string]string{},
		}}

		for _, tt := range tests {
			actual, err := Wrap(tt.input).Reduce(tt.acc, tt.fn).Release()

			assert.EqualError(t, err, tt.err.Error())
			assert.Nil(t, actual)
		}
	})

	t.Run("should forward received error", func(t *testing.T) {
		var noop Fn = func(x interface{}) interface{} { return nil }
		var noacc Transform = func(acc interface{}, x interface{}) interface{} { return acc }
		acc := []map[string]string{}

		actual, err := Wrap(0).Map(noop).Reduce(acc, noacc).Release()

		assert.EqualError(t, err, UnsupportedWrappedTypeError.Error())
		assert.Nil(t, actual)
	})
}

func TestConcat(t *testing.T) {

	t.Run("should concat the elements of k with the given operand correctly", func(t *testing.T) {
		type Test struct {
			input    interface{}
			expected interface{}
			op       interface{}
			err      error
		}
		tests := []Test{{
			input:    []string{},
			expected: []string{},
			op:       []string{},
		}}
		for _, tt := range tests {
			actual, err := Wrap(tt.input).Concat(tt.op).Release()

			expected := tt.expected

			assert.NoError(t, err)
			assert.Equal(t, expected, actual)
		}
	})

	t.Run("should raise error when input type is not supported", func(t *testing.T) {
		type Test struct {
			input  interface{}
			actual interface{}
			op     interface{}
			err    error
		}

		tests := []Test{{
			input: 0,
			err:   UnsupportedWrappedTypeError,
			op:    []int{},
		}, {
			input: "a",
			err:   UnsupportedWrappedTypeError,
			op:    []string{},
		}, {
			input: map[string]string{},
			err:   UnsupportedWrappedTypeError,
			op:    []map[string]string{},
		}}
		for _, tt := range tests {
			actual, err := Wrap(tt.input).
				Concat(tt.op).
				Release()

			assert.EqualError(t, err, tt.err.Error())
			assert.Nil(t, actual)
		}
	})

	t.Run("should forward received error", func(t *testing.T) {
		var noop Fn = func(x interface{}) interface{} { return nil }
		var op = []int{0, 1, 3}
		v := 0

		actual, err := Wrap(v).Map(noop).Concat(op).Release()

		assert.EqualError(t, err, UnsupportedWrappedTypeError.Error())
		assert.Nil(t, actual)
	})

	t.Run("should raise an error when given an incorrectly typed operand", func(t *testing.T) {
		v := []string{}

		actual, err := Wrap(v).
			Concat([]int{}).
			Release()

		errString := OperandTypeMismatchError.Error()

		assert.Error(t, err)
		assert.EqualError(t, err, errString)
		assert.Nil(t, actual)
	})
}

func TestPushPop(t *testing.T) {

	t.Run("push forwards received error", func(t *testing.T) {
		var noop Fn = func(x interface{}) interface{} { return nil }

		actual, err := Wrap(0).
			Map(noop).
			Push().
			Release()

		assert.EqualError(t, err, UnsupportedWrappedTypeError.Error())
		assert.Nil(t, actual)
	})

	t.Run("pop forwards received error", func(t *testing.T) {
		var noop Fn = func(x interface{}) interface{} { return nil }

		actual, err := Wrap(0).
			Map(noop).
			Pop().
			Release()

		assert.EqualError(t, err, UnsupportedWrappedTypeError.Error())
		assert.Nil(t, actual)
	})

	t.Run("stack state is correctly managed for happy path", func(t *testing.T) {
		var even Predicate = func(x interface{}) bool { return x.(int)%2 == 0 }
		var none Predicate = func(x interface{}) bool { return false }

		v := []int{1, 2, 3}
		popped := make([][]int, 3)
		Wrap(v).
			Push().
			Concat([]int{4}).
			Pop().
			Export(reflect.ValueOf(&popped[0])).
			Concat([]int{7}).
			Push().
			Concat([]int{8}).
			Pop().
			Export(reflect.ValueOf(&popped[1])).
			Filter(even).
			Push().
			Filter(none).
			Pop().
			Export(reflect.ValueOf(&popped[2])).
			ReleaseOrPanic()

		assert.Equal(t, []int{1, 2, 3}, popped[0])
		assert.Equal(t, []int{1, 2, 3, 7}, popped[1])
		assert.Equal(t, []int{2}, popped[2])
	})
}

func Test_Slices_Concat(t *testing.T) {
	v := []string{"a", "b"}
	actual, err := Wrap(v).Concat([]string{"c"}).Release()
	expected := []string{"a", "b", "c"}
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_Export_EncoiledError(t *testing.T) {
	var noop Fn = func(x interface{}) interface{} { return nil }
	v := 0
	exp := reflect.ValueOf(&[]int{})
	actual, err := Wrap(v).Map(noop).Export(exp).Release()
	assert.EqualError(t, err, UnsupportedWrappedTypeError.Error())
	assert.Nil(t, actual)
}

func Test_Export_UnsupportedEncoiledType(t *testing.T) {
	v := 0
	exp := reflect.ValueOf(&[]int{})
	actual, err := Wrap(v).Export(exp).Release()
	assert.EqualError(t, err, UnsupportedWrappedTypeError.Error())
	assert.Nil(t, actual)
}

func Test_ReleaseOrPanic_NoPanic(t *testing.T) {
	v := []int{1, 2, 3}
	actual := Wrap(v).
		Concat([]int{4, 5, 6}).
		ReleaseOrPanic()
	expected := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, expected, actual)
}

func Test_ReleaseOrPanic_PanicsOnError(t *testing.T) {
	v := []string{}
	panics := func() {
		Wrap(v).
			Concat([]int{4, 5, 6}).
			ReleaseOrPanic()
	}
	assert.Panics(t, panics)
}

func Test_Export(t *testing.T) {
	v := []int{1, 2, 3}
	actual := []int{}
	Wrap(v).
		Concat([]int{4, 5, 6}).
		Export(reflect.ValueOf(&actual)).
		ReleaseOrPanic()
	expected := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, expected, actual)
}
