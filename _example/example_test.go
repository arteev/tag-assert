package example

import (
	"testing"

	"github.com/arteev/tag-assert"
)

func TestExampleStructTagsValueSuccess(t *testing.T) {
	v := ExampleStruct{}
	assert.Expect(t, v).ExpectField("Name").
		Assert("xml", "Name").
		Assert("json", "name,omitempty")

	assert.Expect(t, v).ExpectField("WithoutTag").Empty()
}

func TestExampleStructTagsValueFailed(t *testing.T) {
	v := ExampleStruct{}
	assert.Expect(t, v).ExpectField("ID").
		Assert("xml", "ID").
		Assert("json", "id"). // this error
		HasTag("bson").
		ExpectTag("json").Equal("rn").Equal("id") //this error

	assert.Expect(t, v).ExpectField("SN").
		Assert("xml", "SN").
		Assert("json", "social_number")

	assert.Expect(t, v).ExpectField("private").
		Assert("xml", "private")

}
