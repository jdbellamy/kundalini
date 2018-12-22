package kundalini

import (
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
	"gitlab.com/jdbellamy/kundalini/slices"
)

// Kundalini is a chaining map/filter/reduce library
type Kundalini interface {
	Concat(slice interface{}) Kundalini
	Map(fn Fn) Kundalini
	Filter(p func(interface{}) bool) Kundalini
	Reduce(acc interface{}, fn Transform) Kundalini
	Release() (interface{}, error)
	ReleaseOrPanic() interface{}
	Types() Kundalini
	Export(reflect.Value) Kundalini
	Push() Kundalini
	Pop() Kundalini
}

// K holds the elements that kundalini operates on
type K struct {
	wrapped interface{}
	err     error
	stack   []interface{}
}

type Fn func(interface{}) interface{}
type Predicate func(interface{}) bool
type Transform func(interface{}, interface{}) interface{}

var UnsupportedWrappedTypeError = fmt.Errorf("Unsupported encoiled type")
var OperandTypeMismatchError = fmt.Errorf("type mismatch between wrapped value and operand")

// Wrap wraps an element in an instance of `k`
func Wrap(e interface{}) Kundalini {
	logrus.Debug("  wrap: ", e)
	return &K{
		wrapped: e,
		stack:   make([]interface{}, 0),
	}
}

// Release returns the elements wrapped by `k`
// `val` is always nil when `err` is populated and vice-versa
func (k *K) Release() (val interface{}, err error) {
	if k.err != nil {
		return nil, k.err
	}
	return k.wrapped, nil
}

// ReleaseOrPanic either returns the elements wrapped by `k` or panics
func (k *K) ReleaseOrPanic() interface{} {
	if k.err != nil {
		panic(k.err)
	}
	return k.wrapped
}

// Types returns a mapping of the types of each element encoiled by `k`
func (k *K) Types() Kundalini {
	if k.err != nil {
		return k
	}
	switch reflect.TypeOf(k.wrapped).Kind() {
	case reflect.Slice:
		v := slices.Types(reflect.ValueOf(k.wrapped))
		logrus.Debug(" types: ", v)
		return &K{
			wrapped: v,
			stack:   k.stack,
		}
	}
	return &K{err: UnsupportedWrappedTypeError}
}

// Export attempts to copy the current elements of `k` to the provided target
// a new slice with len and cap based on `k`s elements is generated at `ptr`
func (k *K) Export(ptr reflect.Value) Kundalini {
	if k.err != nil {
		return k
	}
	switch reflect.TypeOf(k.wrapped).Kind() {
	case reflect.Slice:
		v := slices.Export(reflect.ValueOf(k.wrapped), ptr)
		logrus.Debug("export: ", v)
		return &K{
			wrapped: v,
			stack:   k.stack,
		}
	}
	return &K{err: UnsupportedWrappedTypeError}
}

// Map applys `fn` over each element encoiled by `k`
func (k *K) Map(fn Fn) Kundalini {
	if k.err != nil {
		return k
	}
	switch reflect.TypeOf(k.wrapped).Kind() {
	case reflect.Slice:
		v := slices.Map(reflect.ValueOf(k.wrapped), fn)
		logrus.Debug("   map: ", v)
		return &K{
			wrapped: v,
			stack:   k.stack,
		}
	}
	return &K{err: UnsupportedWrappedTypeError}
}

// Filter keeps the elements of `k` that predicate `p` is true for
func (k *K) Filter(p func(interface{}) bool) Kundalini {
	if k.err != nil {
		return k
	}
	switch reflect.TypeOf(k.wrapped).Kind() {
	case reflect.Slice:
		v := slices.Filter(reflect.ValueOf(k.wrapped), p)
		logrus.Debug("filter: ", v)
		return &K{
			wrapped: v,
			stack:   k.stack,
		}
	}
	return &K{err: UnsupportedWrappedTypeError}
}

// Reduce applys 'fn' over the elements of `k` and accumulates the results
func (k *K) Reduce(acc interface{}, fn Transform) Kundalini {
	if k.err != nil {
		return k
	}
	switch reflect.TypeOf(k.wrapped).Kind() {
	case reflect.Slice:
		v := slices.Reduce(reflect.ValueOf(k.wrapped), acc, fn)
		logrus.Debug("reduce: ", v)
		return &K{
			wrapped: v,
			stack:   k.stack,
		}
	}
	return &K{err: UnsupportedWrappedTypeError}
}

// Concat appends the elements of `e` to the elements wrapped by `k`
func (k *K) Concat(e interface{}) Kundalini {
	if k.err != nil {
		return k
	}
	switch reflect.TypeOf(k.wrapped).Kind() {
	case reflect.Slice:
		v, err := slices.Concat(reflect.ValueOf(k.wrapped), e)
		logrus.Debug("concat: ", v)
		if err != nil {
			return &K{err: err}
		} else {
			return &K{
				wrapped: v,
				stack:   k.stack,
			}
		}
	}
	return &K{err: UnsupportedWrappedTypeError}
}

// Push appends the elements wrapped by `k` to an internal stack
func (k *K) Push() Kundalini {
	if k.err != nil {
		return k
	}
	logrus.Debug("pushed: ", k.wrapped)
	stack := append(k.stack, k.wrapped)
	return &K{
		wrapped: k.wrapped,
		stack:   stack,
	}
}

// Pop sets the value wrapped by `k` to the tail of the internal stack
func (k *K) Pop() Kundalini {
	if k.err != nil {
		return k
	}
	logrus.Debug("popped: ", k.wrapped)
	idx := len(k.stack) - 1
	tail := k.stack[idx]
	left := k.stack[:idx]
	return &K{
		wrapped: tail,
		stack:   left,
	}
}
