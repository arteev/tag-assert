package assert

import (
	"testing"
)

func TestTagHasValue(t *testing.T) {
	tag := &Tag{
		Field: &Field{
			assert: &StructAssert{
				t: t,
			},
		},
		Value: "test",
	}
	if !tag.HasValue("test") {
		t.Error("Expected HasValue: test")
	}
}

func TestTagNotEmpty(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()
	assert := Expect(t, &struct {
		Test         string `json:""`
		TestNotEmpty string `json:"test"`
	}{})
	tag := assert.ExpectField("Test").ExpectTag("json")
	assert.t = test.mockT
	test.mockT.EXPECT().Helper().AnyTimes()
	test.mockT.EXPECT().Errorf("%s: Tag <%s> is empty", "<Unnamed>.Test", "json")
	tag.NotEmpty()
	assert.ExpectField("TestNotEmpty").ExpectTag("json").NotEmpty()
}

func TestTagEqual(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()
	assert := Expect(t, &struct {
		Test string `json:"test"`
	}{})
	tag := assert.ExpectField("Test").ExpectTag("json")
	assert.t = test.mockT
	test.mockT.EXPECT().Helper().AnyTimes()
	tag.Equal("test")
	test.mockT.EXPECT().Errorf("%s: Tag <%s> does not have a value of <%s>,but actual <%s>", "<Unnamed>.Test", "json", "unknown", "test")
	tag.Equal("unknown")
}

func TestHasTag(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	assert := Expect(test.t, TestStruct{})

	test.mockT.EXPECT().Errorf("%s: Field <%s> is private", "TestStruct", "private")
	test.mockT.EXPECT().Errorf("%s: Tag <%s> not found", "TestStruct.private", "Unknown")
	assert.ExpectField("private").HasTag("Unknown")

	test.mockT.EXPECT().Errorf("%s: Tag <%s> not found", "TestStruct.WithoutTags", "Unknown")
	assert.ExpectField("WithoutTags").HasTag("Unknown")

	test.mockT.EXPECT().Errorf("%s: Tag <%s> not found", "TestStruct.Public", "Unknown")
	assert.ExpectField("Public").HasTag("Unknown")

	//check no call Errorf
	assert.ExpectField("Public").HasTag("tag1").HasTag("tag2")
}

func TestHasTags(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	assert := Expect(test.t, TestStruct{})
	test.mockT.EXPECT().Errorf("%s: Tag <%s> not found", "TestStruct.Public", "unknown")
	assert.ExpectField("Public").HasTags("tag1", "tag2", "unknown")
}

func TestAssertTag(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	assert := Expect(test.t, TestStruct{})

	field := assert.ExpectField("Public")

	f := field.Assert("tag1", "pub")

	if f != field {
		t.Errorf("Expected %p, got %p", field, f)
	}

	test.mockT.EXPECT().Errorf("%s: Tag <%s> does not have a value of <%s>,but actual <%s>", "TestStruct.Public", "tag1", "unknown", "pub")
	field.Assert("tag1", "unknown")

	test.mockT.EXPECT().Errorf("%s: Tag <%s> not found", "TestStruct.Public", "unknown")
	field.Assert("unknown", "unknown")
}

func TestExpectTag(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	test.mockT.EXPECT().Errorf("%s: Field <%s> is private", "TestStruct", "private")
	test.mockT.EXPECT().Errorf("%s: Tag <%s> not found", "TestStruct.private", "Unknown")

	assert := Expect(test.t, TestStruct{})

	assert.ExpectField("private").ExpectTag("Unknown")

	test.mockT.EXPECT().Errorf("%s: Tag <%s> not found", "TestStruct.WithoutTags", "Unknown")
	assert.ExpectField("WithoutTags").ExpectTag("Unknown")
	test.mockT.EXPECT().Errorf("%s: Tag <%s> not found", "TestStruct.Public", "Unknown")
	assert.ExpectField("Public").ExpectTag("Unknown")

	assert.ExpectField("Public").ExpectTag("tag1")
	assert.ExpectField("Public").ExpectTag("tag2")

}
