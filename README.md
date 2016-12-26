# genval
Generates Validate() methods for all structs in pkg by tags
- no reflection in generated code - it means fast  
- generator not needed on runtime
- possibilities to override generated behavior for local purposes

Installation
------------
    go get github.com/l1va/genval

Usage and documentation
------
    ./genval packageWithStructsForGeneration
##### Addition flags
    outputFile - output file name (default: validators_generated.go)
    needValidatableCheck - check struct on Validatable before calling Validate() (default: true)

##### Supported tags:
    String: min_len, max_len
    Number: min, max
    Array:  min_items, max_items, item
    Pointer: nullable, not_null
    Interface: func
    Struct: method, func
    Map: min_items, max_items, key, value

##### Default validation:
    String: min_len=1
    Pointer: not_null
    Other: no validation

##### Some tips:
1. not use interface{} if you can
2. try flag needValidatableCheck=false 
3. commit generated code under source control
4. **read generated code** if needed, do not afraid it

##### Ways to override generated behavior(see examples or try yourself): 
1. custom `(r Struct)Validate() error` method
2. `func` or `method` tags for fields

##### Examples:
[Generated code](https://github.com/l1va/genval/blob/master/examples/simple/validators_generated.go) for next structs:
```go
type User struct {
    Name string `validate:"min_len=3,max_len=64"`
    Age  uint   `validate:"min=18,max=95"`
    Dog  Dog
    Emails map[int]string `validate:"min_items=1,key=[max=3],value=[min_len=5]"`
}
type Dog struct {
    Name string `validate:"min_len=1,max_len=64"`
}
```

- [Simple](https://github.com/l1va/genval/tree/master/examples/simple)
- [Complicated](https://github.com/l1va/genval/tree/master/examples/complicated)
- [Overriding generated validators](https://github.com/l1va/genval/tree/master/examples/overriding)
