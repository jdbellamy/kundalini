package kundalini

import (
	"fmt"
	"reflect"
)

// Concat appends the elements of `op` to the elements of `k`
func (k *K) Concat(op interface{}) Kundalini {
	if k.err != nil {
		return k
	}

	kT := reflect.TypeOf(k.wrapped)
	opT := reflect.TypeOf(op)
	opV := reflect.ValueOf(op)

	// handle when `op` is a scalar value
	if opT.Kind() != reflect.Slice {
		sliceOfOpT := reflect.SliceOf(opT)
		sliceOfOp := reflect.MakeSlice(sliceOfOpT, 1, 1)
		sliceOfOp.Index(0).Set(opV)
		op = sliceOfOp.Interface()
		opT = reflect.TypeOf(op)
		opV = reflect.ValueOf(op)
	}

	if opT != kT {
		k.err = fmt.Errorf("type mismatch between wrapped value and operand")
		return k
	}

	if opV.Len() == 0 {
		return k
	}

	kV := reflect.ValueOf(k.wrapped)

	rLen := kV.Len() + opV.Len()
	r := reflect.MakeSlice(kT, rLen, rLen)

	for i := 0; i < kV.Len(); i++ {
		r.Index(i).Set(kV.Index(i))
	}

	for i := kV.Len(); i < r.Len(); i++ {
		r.Index(i).Set(opV.Index(i - kV.Len()))
	}

	k.wrapped = r.Interface()
	return k
}

// Map applys `fn` over each element of `k`
func (k *K) Map(fn func(interface{}) interface{}) Kundalini {
	if k.err != nil {
		return k
	}

	kV := reflect.ValueOf(k.wrapped)

	if kV.Len() == 0 {
		return k
	}

	r := reflect.MakeSlice(kV.Type(), kV.Len(), kV.Len())

	for i := 0; i < kV.Len(); i++ {
		v := kV.Index(i).Interface()
		mapped := fn(v)
		r.Index(i).Set(reflect.ValueOf(mapped))
	}

	k.wrapped = r.Interface()

	return k
}

// Filter keeps the elements of `k` that predicate `p` is true for
func (k *K) Filter(p func(interface{}) bool) Kundalini {
	if k.err != nil {
		return k
	}

	kV := reflect.ValueOf(k.wrapped)

	if kV.Len() == 0 {
		return k
	}

	tmp := make([]interface{}, 0)
	for i := 0; i < kV.Len(); i++ {
		v := kV.Index(i).Interface()
		if p(v) {
			tmp = append(tmp, v)
		}
	}

	r := reflect.MakeSlice(kV.Type(), len(tmp), len(tmp))
	for i, v := range tmp {
		r.Index(i).Set(reflect.ValueOf(v))
	}

	k.wrapped = r.Interface()

	return k
}

// Reduce applys 'fn' over the elements of `k` and accumulates the results
func (k *K) Reduce(acc interface{}, fn func(interface{}, interface{}) interface{}) Kundalini {
	if k.err != nil {
		return k
	}

	kV := reflect.ValueOf(k.wrapped)

	if kV.Len() == 0 {
		return k
	}

	for i := 0; i < kV.Len(); i++ {
		v := kV.Index(i).Interface()
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

	k.wrapped = r.Interface()

	return k
}
