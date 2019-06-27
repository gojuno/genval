# genval [![GoDoc](https://godoc.org/github.com/gojuno/genval?status.svg)](http://godoc.org/github.com/gojuno/genval) [![Build Status](https://travis-ci.org/gojuno/genval.svg?branch=master)](https://travis-ci.org/gojuno/genval)

Generates Validate() methods for all structs in package by tags
- no reflection in generated code - it means fast  
- possibilities to override generated behavior
- can be used as `//go:generate genval pkg` 
- Enum support 

## Installation

    go get github.com/gojuno/genval


Usage
------

```go
type User struct {
	Name   string   `validate:"max_len=64"`
	Age    uint     `validate:"min=18"`
	Emails []string `validate:"min_items=1,item=[min_len=5]"`
}

//generates:
func (r User) Validate() error {
	var errs errlist.List
	if utf8.RuneCountInString(string(r.Name)) > 64 {
		errs.AddFieldf("Name", "longer than 64 chars")
	}
	if r.Age < 18 {
		errs.AddFieldf("Age", "less than 18")
	}
	if len(r.Emails) < 1 {
		errs.AddFieldf("Emails", "less items than 1")
	}
	for kEmails, vEmails := range r.Emails {
		if utf8.RuneCountInString(string(vEmails)) < 5 {
			errs.AddFieldf(fmt.Sprintf("Emails"+".%v", kEmails), "shorter than 5 chars")
		}
	}
	return errs.ErrorOrNil()
}
```

##### Some other examples:
- [Simple](https://github.com/gojuno/genval/tree/master/examples/simple)
- [Complicated](https://github.com/gojuno/genval/tree/master/examples/complicated)
- [Overriding generated validators](https://github.com/gojuno/genval/tree/master/examples/overriding)

#### How to generate?

    genval mypkg

or you can use it as `go:generate` directive  

    //go:generate genval mypkg

### Supported tags
- *String*: **min_len**, **max_len** - min and max valid lenghth 
- *Number*: **min**, **max** - min and max valid value (can be float)
- *Array*:  **min_items**, **max_items** - min and max count of items in array  
    **item** - scope tag, contains validation tags for each item
- *Pointer*: **nullable**, **not_null** - it's clear
- *Interface*: **func** - the same as for struct (`func NameOfTheFunc(i interface{})error{..}`)
- *Struct*: **func** - name of the method of this struct (`func(s Struct) MethodName()error{..}`)  
    or name of the func that will be used for validation (`func nameOfTheFunc(s Struct)error{..}`)      
    *Methods should starts from '.'*  
    *Can be used not once:* `func=.MethodName,func=nameOfTheFunc` or even `func=.MethodName;nameOfTheFunc`    
- *Map*: **min_items**, **max_items** - min and max count of items in map  
    **key**, **value** - scope tags, contains validation tags for key or value 

### Enum support
Go doesn`t support enums, but you can create some custom type and add few constants with required values.

```go
type State int

const (
	StateOk    State = 200
	StateError State = 400
)

//generates:
func (r State) Validate() error {
	switch r {
	case StateOk:
	case StateError:
	default:
		return fmt.Errorf("invalid value for enum State: %v", r)
	}
	return nil
}
```

#### Some tips
1. don`t use interface{} if you can
2. commit generated code under source control
3. **read generated code** if needed, do not afraid it

### Custom validation
#### Additional validation
In some cases it\`s required to add some custom validation. You can just add unexported `validate` method.
```go
type User struct {
	Name  string `validate:"max_len=64"`
	Age   uint   `validate:"min=16"`
	Email string
}

func (u User) validate() error {
	if u.Age < 18 && u.Email == "" {
		return errors.New("email is required for people younger than 18")
	}
	return nil
}

// generates:
func (r User) Validate() error {
	var errs errlist.List
	if utf8.RuneCountInString(string(r.Name)) > 64 {
		errs.AddFieldf("Name", "longer than 64 chars")
	}
	if r.Age < 16 {
		errs.AddFieldf("Age", "less than 16")
	}
	errs.Add(r.validate())
	return errs.ErrorOrNil()
}
```
#### Override validation
If you don\`t want to use genval for some structs or use just custom validation then you can **override exported `Validate`** method. In this case genval **will generate nothing** for this struct.
