package assert

import (
	"errors"
	"reflect"
	"unicode"

	goerror "github.com/pkg/errors"
)

//Errors
var (
	ErrNotStruct     = errors.New("Must be struct")
	ErrUnxpectedNil  = errors.New("Unexpected nil")
	ErrFieldNotFound = errors.New("Field not found")
	ErrUnexported    = errors.New("Field unexported")
	ErrTagNotFound   = errors.New("Tag not found")
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
	value  interface{}
	failed bool
	fields map[string]*Field
}

//Field contains methods for verifying the field
type Field struct {
	assert      *StructAssert
	name        string
	structField *reflect.StructField
}

//Tag of field
type Tag struct {
	Field *Field
	Name  string
	Value string
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
		a.t.Error(goerror.Wrap(ErrFieldNotFound, name))
		//	a.t.Error(ErrFieldNotFound)
		return nil, false
	}
	if unicode.IsLower(rune(name[0])) {
		a.t.Error(goerror.Wrap(ErrUnexported, name))
		return nil, false
	}
	//ErrUnexported
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

//TODO: check Tag

//HasTag checks the existence of a tag in the field
func (f *Field) HasTag(name string) *Field {
	if f.structField == nil {
		f.assert.t.Error(goerror.Wrap(ErrTagNotFound, name))
		return f
	}
	_, ok := f.structField.Tag.Lookup(name)
	if !ok {
		f.assert.t.Error(goerror.Wrap(ErrTagNotFound, name))
	}
	return f
}

//ExpectTag waiting for a tag with name to verify assert
func (f *Field) ExpectTag(name string) *Tag {
	f.assert.t.Helper()
	if f.structField == nil {
		f.assert.t.Error(goerror.Wrap(ErrTagNotFound, name))
		return &Tag{Name: name}
	}

	value, ok := f.structField.Tag.Lookup(name)
	if !ok {
		f.assert.t.Error(goerror.Wrap(ErrTagNotFound, name))
	}
	return &Tag{
		Field: f,
		Name:  name,
		Value: value,
	}
}

//HasValue checks the tag for the specified value
func (t *Tag) HasValue(value string) bool {
	if t.Field == nil {
		return false
	}
	return t.Value == value
}

//Assert checks the tag (name) with the specified value
func (f *Field) Assert(name, value string) *Field {
	f.assert.t.Helper()
	t := f.ExpectTag(name)
	if t.Field == nil {
		return f
	}

	if !t.HasValue(value) {
		f.assert.t.Errorf("%s: value is not %s, actual %s", t.Name, value, t.Value)
	}
	return f
}
