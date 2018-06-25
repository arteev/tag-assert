package assert

//Tag of field
type Tag struct {
	Field *Field
	Name  string
	Value string
}

//HasValue checks the tag for the specified value
func (t *Tag) HasValue(value string) bool {
	if t.Field == nil {
		return false
	}
	t.Field.assert.t.Helper()
	return t.Value == value
}

//Equal checks the tag for the specified value
func (t *Tag) Equal(value string) *Tag {
	t.Field.assert.t.Helper()
	if t.Value != value {
		t.Field.assert.t.Errorf("%s: Tag <%s> does not have a value of <%s>,but actual <%s>", t.Field.getFullName(), t.Name, value, t.Value)

	}
	return t
}

//NotEmpty check for empty value
func (t *Tag) NotEmpty() *Tag {
	t.Field.assert.t.Helper()
	if t.Value == "" {
		t.Field.assert.t.Errorf("%s: Tag <%s> is empty", t.Field.getFullName(), t.Name)
	}
	return t
}
