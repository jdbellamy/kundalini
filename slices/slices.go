package slices

import (
	"fmt"
	"reflect"
)

var TypeMismatchError = fmt.Errorf("type mismatch between wrapped value and operand")
var ExportTargetIsNotPointerError = fmt.Errorf("buf must be a pointer to a slice of the correct type")

// Types returns a mapping of the `Type`s of the elements of slice `s`
func Types(s reflect.Value) []reflect.Type {
	types := make([]reflect.Type, s.Len())
	for i := 0; i < s.Len(); i++ {
		types[i] = s.Index(i).Type()
	}
	return types
}

// Map applys `fn` over each element encoiled by `k`
func Map(s reflect.Value, fn func(interface{}) interface{}) interface{} {
	r := reflect.MakeSlice(s.Type(), s.Len(), s.Len())

	for i := 0; i < s.Len(); i++ {
		v := s.Index(i)
		mapped := fn(v.Interface())
		if mapped == nil {
			r.Index(i).Set(v)
		} else {
			r.Index(i).Set(reflect.ValueOf(mapped))
		}
	}

	return r.Interface()
}

// Filter keeps the elements of `k` that predicate `p` is true for
func Filter(s reflect.Value, p func(interface{}) bool) interface{} {
	if s.Len() == 0 {
		return s.Interface()
	}

	tmp := make([]interface{}, 0)
	for i := 0; i < s.Len(); i++ {
		v := s.Index(i).Interface()
		if p(v) {
			tmp = append(tmp, v)
		}
	}

	r := reflect.MakeSlice(s.Type(), len(tmp), len(tmp))
	for i, v := range tmp {
		r.Index(i).Set(reflect.ValueOf(v))
	}

	return r.Interface()
}

// Reduce applys 'fn' over the elements of `k` and accumulates the results
func Reduce(s reflect.Value, acc interface{}, fn func(interface{}, interface{}) interface{}) interface{} {

	if s.Len() == 0 {
		return s.Interface()
	}

	for i := 0; i < s.Len(); i++ {
		v := s.Index(i).Interface()
		acc = fn(acc, v)
	}

	var r reflect.Value

	accT := reflect.TypeOf(acc)
	accV := reflect.ValueOf(acc)
	if accT.Kind() != reflect.Slice {
		rT := reflect.SliceOf(accT)
		r = reflect.MakeSlice(rT, 1, 1)
		r.Index(0).Set(accV)
	} else {
		r = accV
	}

	return r.Interface()
}

// Concat appends the elements of `s` to the elements of `k`
func Concat(s reflect.Value, e interface{}) (interface{}, error) {
	eT := reflect.TypeOf(e)
	eV := reflect.ValueOf(e)

	if eT != s.Type() {
		return nil, TypeMismatchError
	}

	rLen := s.Len() + eV.Len()
	r := reflect.MakeSlice(s.Type(), rLen, rLen)

	for i := 0; i < s.Len(); i++ {
		r.Index(i).Set(s.Index(i))
	}

	for i := s.Len(); i < r.Len(); i++ {

		r.Index(i).Set(eV.Index(i - s.Len()))
	}

	return r.Interface(), nil
}

// Export attempts to copy the current elements of `k` to the provided target
// a new slice with len and cap based on `k`s elements is generated at `ptr`
func Export(s reflect.Value, ptr reflect.Value) interface{} {
	if ptr.Type().Kind() != reflect.Ptr {
		panic(ExportTargetIsNotPointerError)
	}

	kLen := s.Len()
	slice := reflect.MakeSlice(ptr.Type().Elem(), kLen, kLen)
	ptr.Elem().Set(slice)
	ptr.Elem().SetLen(kLen)
	reflect.Copy(ptr.Elem(), s)

	return s.Interface()
}
