# kundalini

Chaining map/filter/reduce in Go  

## Example
```go
a := []int{0, 1, 2, 3, 4}

v, err := Wrap(a).       // {0, 1, 2, 3, 4}
  Filter(even).          // {0, 2, 4}
  Map(double).           // {0, 4, 8}
  Export(ptr).           // {0, 4, 8} copied into `buf`
  Filter(firstN(1)).     // {4, 8}
  Concat(Wrap(buf).      // - {0, 4, 8}
      Filter(firstN(1)). // - {4, 8}
      ReleaseOrPanic()). // {4, 8, 4, 8}
  Reduce(8, sum).        // {32}
  Types().               // {int}
  Release()        
```

## Docs
Just an excuse to learn about Go reflection - not intended for any actual things
