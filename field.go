package assert

import "reflect"

//Field contains methods for verifying the field
type Field struct {
	assert      *StructAssert
	name        string
	structField *reflect.StructField
}

//Assert checks the tag (name) with the specified value
func (f *Field) Assert(name, value string) *Field {
	f.assert.t.Helper()
	t := f.ExpectTag(name)
	if t.Field == nil {
		return f
	}

	if !t.HasValue(value) {
		f.assert.t.Errorf("%s: Tag <%s> does not have a value of <%s>,but actual <%s>", f.getFullName(), t.Name, value, t.Value)
	}
	return f
}

func (f *Field) getFullName() string {
	if f.assert == nil || f.assert.failed {
		return f.name
	}
	structName := f.assert.vtype.Name()
	if structName == "" {
		structName = "<Unnamed>"
	}
	return structName + "." + f.name
}

//ExpectTag waiting for a tag with name to verify assert
func (f *Field) ExpectTag(name string) *Tag {
	f.assert.t.Helper()
	if f.structField == nil {
		f.assert.t.Errorf("%s: Tag <%s> not found", f.getFullName(), name)
		return &Tag{Name: name}
	}

	value, ok := f.structField.Tag.Lookup(name)
	if !ok {
		f.assert.t.Errorf("%s: Tag <%s> not found", f.getFullName(), name)
		return &Tag{Name: name}
	}
	return &Tag{
		Field: f,
		Name:  name,
		Value: value,
	}
}

//HasTag checks the existence of a tag in the field
func (f *Field) HasTag(name string) *Field {
	f.assert.t.Helper()
	if f.structField == nil {
		f.assert.t.Errorf("%s: Tag <%s> not found", f.getFullName(), name)
		return f
	}
	_, ok := f.structField.Tag.Lookup(name)
	if !ok {
		f.assert.t.Errorf("%s: Tag <%s> not found", f.getFullName(), name)
	}
	return f
}

//HasTags checks the existence of tags in the field
func (f *Field) HasTags(names ...string) *Field {
	for _, name := range names {
		f.HasTag(name)
	}
	return f
}

//Empty verifies that the tag is empty
func (f *Field) Empty() *Field {
	if string(f.structField.Tag) != "" {
		f.assert.t.Errorf("%s: Not empty", f.getFullName())
	}
	return f
}
