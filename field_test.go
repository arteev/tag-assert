package assert

import (
	"testing"
)

func TestHasField(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	test.mockT.EXPECT().Errorf("%s: Field <%s> is private", "TestStruct", "private")
	test.mockT.EXPECT().Errorf("%s: Field <%s> not found", "TestStruct", "Unknown")

	Expect(test.t, &TestStruct{}).
		HasField("Public").
		HasField("private").
		HasField("SubStruct").
		HasField("Unknown")
}

func TestExpectField(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	assert := Expect(test.t, TestStruct{})
	field := assert.ExpectField("Public")
	if field == nil {
		t.Error("Unexpected nil")
	}

	if assert != field.assert {
		t.Errorf("Expected %p, got %p", assert, field.Assert)
	}
	if field.name != "Public" {
		t.Errorf("Expected %q, got %q", "Public", field.name)
	}
	if field.structField == nil {
		t.Error("Unexpected field.structField=nil")
	}

	assert.ExpectField("Public")

	test.mockT.EXPECT().Errorf("%s: Field <%s> is private", "TestStruct", "private")
	field = assert.ExpectField("private")
	if field == nil {
		t.Error("Unexpected nil")
	}

	if field.structField != nil {
		t.Errorf("Expected field.structField nil,got %v", field.structField)
	}
}
