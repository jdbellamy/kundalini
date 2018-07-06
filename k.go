package kundalini

import (
	"reflect"
)

// Kundalini is a chaining map/filter/reduce library
type Kundalini interface {
	Concat(slice interface{}) Kundalini
	Map(fn func(interface{}) interface{}) Kundalini
	Filter(p func(interface{}) bool) Kundalini
	Reduce(acc interface{}, fn func(interface{}, interface{}) interface{}) Kundalini
	Release() (interface{}, error)
}

// K holds the elements that kundalini operates on
type K struct {
	wrapped interface{}
	err     error
}

// Coil wraps a slice or scaler in an instance of `k`
func Coil(e interface{}) Kundalini {
	eT := reflect.TypeOf(e)
	if eT.Kind() == reflect.Slice {
		return &K{
			wrapped: e,
		}
	}
	eV := reflect.ValueOf(e)
	sliceT := reflect.SliceOf(eT)
	sliceOfE := reflect.MakeSlice(sliceT, 1, 1)
	sliceOfE.Index(0).Set(eV)
	return &K{
		wrapped: sliceOfE.Interface(),
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
