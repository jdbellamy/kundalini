# kundalini

Chaining map/filter/reduce in go  


## Example
```go
v := []int{0, 1, 2, 3, 4, 5}

even := func(x interface{}) bool {
  return x.(int)%2 == 0
}

double := func(x interface{}) interface{} {
  return x.(int) * 2
}

sum := func(acc interface{}, x interface{}) interface{} {
  return acc.(int) + x.(int)
}

k, _ := Coil(v).   // []int{0, 1, 2, 3, 4, 5}
  Concat(6).       // []int{0, 1, 2, 3, 4, 5, 6}
  Filter(even).    // []int{0, 2, 4, 6}
  Map(double).     // []int{0, 4, 8, 12}
  Reduce(0, sum).  // []int{24}
  Release()
```

## Docs
This project was just an excuse for me to learn more about reflection in go.  I wouldn't suggest using it for any actual things.


#### Kundalini
``` go
type Kundalini interface {
    Concat(slice interface{}) Kundalini
    Map(fn func(interface{}) interface{}) Kundalini
    Filter(p func(interface{}) bool) Kundalini
    Reduce(acc interface{}, fn func(interface{}, interface{}) interface{}) Kundalini
    Release() (interface{}, error)
}
```


#### Coil
``` go
func Coil(e interface{}) Kundalini
```
Coil wraps a slice or scaler in an instance of `k`


#### Concat
``` go
func (k *K) Concat(op interface{}) Kundalini
```
Concat appends the elements of `op` to the elements of `k`


#### Filter
``` go
func (k *K) Filter(p func(interface{}) bool) Kundalini
```
Filter keeps the elements of `k` that predicate `p` is true for


#### Map
``` go
func (k *K) Map(fn func(interface{}) interface{}) Kundalini
```
Map applys `fn` over each element of `k`


#### Reduce
``` go
func (k *K) Reduce(acc interface{}, fn func(interface{}, interface{}) interface{}) Kundalini
```
Reduce applys 'fn' over the elements of `k` and accumulates the results


#### Release
``` go
func (k *K) Release() (val interface{}, err error)
```
Release returns the elements wrapped by `k`
`val` is always nil when `err` is populated and vice-versa
