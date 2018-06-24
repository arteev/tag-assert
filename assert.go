package assert

import (
	"errors"
	"reflect"
	"unicode"
)

//Errors
var (
	ErrNotStruct    = errors.New("Must be struct")
	ErrUnxpectedNil = errors.New("Unexpected nil")
)

type tb interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	Helper()
}

//StructAssert contains methods for verifying the structure
type StructAssert struct {
	t      tb
	vtype  reflect.Type
	value  interface{}
	failed bool
	fields map[string]*Field
}

//Expect waiting for a structure to verify assert
func Expect(t tb, v interface{}) *StructAssert {
	t.Helper()

	check := &StructAssert{
		t:      t,
		value:  v,
		fields: make(map[string]*Field),
	}
	return check.assertStruct()
}

//Expect waiting for a structure to verify assert
func (a *StructAssert) Expect(v interface{}) *StructAssert {
	a.t.Helper()
	return Expect(a.t, v)
}

func (a *StructAssert) assertStruct() *StructAssert {
	a.t.Helper()
	value := reflect.ValueOf(a.value)

	if value.Kind() == reflect.Invalid {
		a.failed = true
		a.t.Fatal(ErrUnxpectedNil)
		return a
	}

	vtype := value.Type()
	a.vtype = vtype
	if vtype.Kind() == reflect.Ptr {
		vtype = vtype.Elem()
	}

	if vtype.Kind() != reflect.Struct {
		a.failed = true
		a.t.Fatal(ErrNotStruct)
	}

	return a
}

func (a *StructAssert) mustStructField(name string) (*Field, bool) {
	a.t.Helper()
	if field, ok := a.fields[name]; ok {
		return field, true
	}

	value := reflect.ValueOf(a.value)
	vtype := value.Type()
	if vtype.Kind() == reflect.Ptr {
		vtype = vtype.Elem()
	}
	structField, ok := vtype.FieldByName(name)
	if !ok {
		a.t.Errorf("%s: Field <%s> not found", vtype.Name(), name)
		return nil, false
	}
	if unicode.IsLower(rune(name[0])) {
		a.t.Errorf("%s: Field <%s> is private", vtype.Name(), name)
		return nil, false
	}
	a.fields[name] = &Field{
		name:        name,
		structField: &structField,
		assert:      a,
	}
	return a.fields[name], true
}

//HasField checks the existence of a field in the structure
func (a *StructAssert) HasField(name string) *StructAssert {
	a.mustStructField(name)
	return a
}

//ExpectField waiting for a field with name to verify assert
func (a *StructAssert) ExpectField(name string) *Field {
	a.t.Helper()
	structField, ok := a.mustStructField(name)
	if ok {
		return structField
	}
	//TODO: mark not found
	structField = &Field{
		name:        name,
		structField: nil,
		assert:      a,
	}
	return structField
}
