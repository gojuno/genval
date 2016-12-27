# genval
Generates Validate() methods for all structs in pkg by tags
- no reflection in generated code - it means fast  
- generator not needed on runtime
- possibilities to override generated behavior for local purposes
- can be used as `//go:generate genval pkg` 
- if type has constants, validation will be by all found constants 

Installation
------------
    go get github.com/l1va/genval

Usage and documentation
------
    ./genval packageWithStructsForGeneration 
or as go:generate directive
    //go:generate genval packageWithStructsForGeneration
    
##### Additional flags
    outputFile - output file name (default: validators.go)
    needValidatableCheck - check struct on Validatable before calling Validate() (default: true)

##### Supported tags:
String: *min_len*, *max_len* - min and max valid lenghth 
Number: *min*, *max* - min and max valid value (can be float)
Array:  *min_items*, *max_items* - min and max count of items in array
    *item* - scope tag, contains validation tags for each item
Pointer: *nullable*, *not_null* - it's clear
Interface: *func* - name of func that will be used for validation,
    function should be like `func NameOfFunc(i interface{})error{..}` 
Struct: *method* - name of the method of this struct, `func(s Struct) MethodName()error{..}`
    *func* - the same as for interface, but param can be as struct type.
Map: *min_items*, *max_items* - min and max count of items in map
    *key*, *value* - scope tags, contains validation tags for key or value 

##### Default validation - no Validation    

##### Some tips:
1. not use interface{} if you can
2. commit generated code under source control
3. **read generated code** if needed, do not afraid it

##### Ways to override generated behavior(see examples or try yourself): 
1. custom `(r Struct)Validate() error` method
2. `func` or `method` tags for fields

##### Examples:
[Generated code](https://github.com/l1va/genval/blob/master/examples/simple/validators.go) for next structs:
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

Validation by constants :
```go
type State int

const (
    StateOk    State = 200
    StateError State = 400
)
//generated:
func (r State) validate() error {
    switch r {
    case StateOk:
    case StateError:
    default:
        return fmt.Errorf("invalid value for enum State: %v", r)
    }
    return nil
}
```