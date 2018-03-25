package assert

import (
	"errors"
	"reflect"
	"testing"
	"unicode"

	goerror "github.com/pkg/errors"
)

//Errors
var (
	ErrNotStruct     = errors.New("Must be struct")
	ErrUnxpectedNil  = errors.New("Unexpected nil")
	ErrFieldNotFound = errors.New("Field not found")
	ErrUnexported    = errors.New("Field unexported")
)

type TB interface {
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
type internalTest struct {
	*testing.T
}

type structAssert struct {
	t      TB
	value  interface{}
	failed bool
	fields map[string]*field
}

type field struct {
	Assert      *structAssert
	name        string
	structField *reflect.StructField
}

func Expect(t TB, v interface{}) *structAssert {
	t.Helper()

	check := &structAssert{
		t:      t,
		value:  v,
		fields: make(map[string]*field),
	}
	return check.assertStruct()
}

func (a *structAssert) assertStruct() *structAssert {
	value := reflect.ValueOf(a.value)

	if value.Kind() == reflect.Invalid {
		a.failed = true
		a.t.Fatal(ErrUnxpectedNil)
		return a
	}

	vtype := value.Type()
	if vtype.Kind() == reflect.Ptr {
		vtype = vtype.Elem()
	}

	if vtype.Kind() != reflect.Struct {
		//	log.Println(a.value, vtype.Kind())
		a.failed = true
		a.t.Fatal(ErrNotStruct)
	}
	return a
}

func (a *structAssert) mustStructField(name string) (*field, bool) {

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
		a.t.Error(goerror.Wrap(ErrFieldNotFound, name))
		//	a.t.Error(ErrFieldNotFound)
		return nil, false
	}
	if unicode.IsLower(rune(name[0])) {
		a.t.Error(goerror.Wrap(ErrUnexported, name))
		return nil, false
	}
	//ErrUnexported
	a.fields[name] = &field{
		name:        name,
		structField: &structField,
		Assert:      a,
	}
	return a.fields[name], true

}

func (a *structAssert) HasField(name string) *structAssert {
	a.mustStructField(name)
	return a
}

func (a *structAssert) ExpectField(name string) *field {
	structField, ok := a.mustStructField(name)
	if ok {
		return structField
	}
	//TODO: mark not found
	structField = &field{
		name:        name,
		structField: nil,
		Assert:      a,
	}
	return structField
}
