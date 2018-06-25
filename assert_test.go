package assert

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

type testAssert struct {
	mockT      *MockTB
	t          tb
	controller *gomock.Controller
}

func setUp(t *testing.T) *testAssert {
	t.Helper()
	mockT := NewMockTB(gomock.NewController(t))

	test := &testAssert{
		controller: mockT.ctrl,
		mockT:      mockT,
		t:          mockT,
	}
	return test
}
func (t *testAssert) tearDown() {
	t.controller.Finish()
}

type SubStruct struct {
	Name string
}

//nolint
type TestStruct struct {
	private     string
	Public      string `tag1:"pub" tag2:"public,options"`
	WithoutTags string
	SubStruct
}

//TODO: when failed not checks fields

func TestExpect(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()
	v := TestStruct{}
	assert := Expect(test.t, &v)
	if assert == nil {
		t.Fatal("Unexpected nil")
	}
	if assert.t != test.t {
		t.Errorf("Expected %v,got %v", test.t, assert.t)
	}
	if reflect.ValueOf(&v).Pointer() != reflect.ValueOf(assert.value).Pointer() {
		t.Errorf("Unexpected %p = %p", &v, assert.value)

	}

	assert2 := assert.Expect(&v)
	if assert2 == nil {
		t.Fatal("Unexpected nil")
	}
	if reflect.ValueOf(assert).Pointer() == reflect.ValueOf(assert2).Pointer() {
		t.Errorf("Unexpected %p = %p", assert, assert2)

	}

	if reflect.ValueOf(&v).Pointer() != reflect.ValueOf(assert2.value).Pointer() {
		t.Errorf("Unexpected %p = %p", &v, assert2.value)

	}

}

func TestAreStructs(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	type TypeStruct struct{}
	var nilTypeStruct *TypeStruct
	cases := []struct {
		Name       string
		Value      interface{}
		MustFailed bool
		Expected   func()
	}{
		{
			Name:       "Nil",
			Value:      nil,
			MustFailed: true,
			Expected: func() {
				test.mockT.EXPECT().Helper().AnyTimes()
				test.mockT.EXPECT().Fatal(ErrUnxpectedNil)
			},
		},
		{
			Name:       "NotStruct",
			Value:      0,
			MustFailed: true,
			Expected: func() {
				test.mockT.EXPECT().Helper().AnyTimes()
				test.mockT.EXPECT().Fatal(ErrNotStruct)
			},
		},

		{
			Name:       "AnonymousStruct",
			Value:      struct{}{},
			MustFailed: false,
			Expected: func() {
				test.mockT.EXPECT().Helper().AnyTimes()
			},
		},

		{
			Name:       "StructNil",
			Value:      nilTypeStruct,
			MustFailed: false,
			Expected: func() {
				test.mockT.EXPECT().Helper().AnyTimes()

			},
		},
		{
			Name:       "Struct",
			Value:      TypeStruct{},
			MustFailed: false,
			Expected: func() {
				test.mockT.EXPECT().Helper().AnyTimes()

			},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			c.Expected()
			assertion := Expect(test.t, c.Value)
			if assertion.failed != c.MustFailed {
				t.Error("Expected failed")
			}
		})

	}
}

func TestHasValue(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	assert := Expect(test.t, TestStruct{})

	has := assert.ExpectField("Public").ExpectTag("tag1").HasValue("12")
	if has {
		t.Error("Unexpected has")
	}
	has = assert.ExpectField("Public").ExpectTag("tag1").HasValue("pub")
	if !has {
		t.Error("Expected has")
	}

	test.mockT.EXPECT().Errorf("%s: Tag <%s> not found", "TestStruct.Public", "Unknown")

	has =
		assert.ExpectField("Public").
			ExpectTag("Unknown").
			HasValue("pub")
	if has {
		t.Error("Unexpected has")
	}
}
