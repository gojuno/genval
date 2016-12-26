# genval
Generates Validate() methods for all structs in pkg by tags
- no reflection in generated code - it means fast  
- generator not needed on runtime

Installation
------------
    go get github.com/l1va/genval

Usage and documentation
------
    ./genval pkg
##### Flags
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
    Other: empty value (no validation)

##### Some tips:
    1. not use interface{} if you can
    2. try flag needValidatableCheck=false 
    3. commit generated code under source control
    4. **read generated code** if needed, do not afraid it

##### Examples:
Next code :
```go
type User struct {
    Name string `validate:"min_len=3,max_len=64"`
    Age  uint   `validate:"min=18,max=95"`
    Dog  Dog
    Emails []string `validate:"min_items=1,item=[min_len=5]"`
}
type Dog struct {
    Name string `validate:"min_len=1,max_len=64"`
}
```
will generate : [Generated code](https://github.com/l1va/genval/examples/simple/validators_generated.go)

- [Simple](https://github.com/l1va/genval/tree/master/examples/simple)
- [Complicated](https://github.com/l1va/genval/tree/master/examples/complicated)
- [Overriding generated validators](https://github.com/l1va/genval/tree/master/examples/overriding)
