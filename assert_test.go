package assert

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	goerror "github.com/pkg/errors"
)

type testAssert struct {
	mockT      *MockTB
	t          TB
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
				test.mockT.EXPECT().Helper()
				test.mockT.EXPECT().Fatal(ErrUnxpectedNil)
			},
		},
		{
			Name:       "NotStruct",
			Value:      0,
			MustFailed: true,
			Expected: func() {
				test.mockT.EXPECT().Helper()
				test.mockT.EXPECT().Fatal(ErrNotStruct)
			},
		},

		{
			Name:       "AnonymousStruct",
			Value:      struct{}{},
			MustFailed: false,
			Expected: func() {
				test.mockT.EXPECT().Helper()
			},
		},

		{
			Name:       "StructNil",
			Value:      nilTypeStruct,
			MustFailed: false,
			Expected: func() {
				test.mockT.EXPECT().Helper()

			},
		},
		{
			Name:       "Struct",
			Value:      TypeStruct{},
			MustFailed: false,
			Expected: func() {
				test.mockT.EXPECT().Helper()

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

func TestHasField(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	test.mockT.EXPECT().Error(gomock.Any()).Do(func(args ...interface{}) {
		if goerror.Cause(args[0].(error)) != ErrUnexported {
			t.Errorf("Expected %v", ErrUnexported)
		}
	})

	test.mockT.EXPECT().Error(gomock.Any()).Do(func(args ...interface{}) {
		if goerror.Cause(args[0].(error)) != ErrFieldNotFound {
			t.Errorf("Expected %v", ErrFieldNotFound)
		}
	})

	Expect(test.t, &TestStruct{}).
		HasField("Public").
		HasField("private").
		HasField("SubStruct").
		HasField("Unknown")

}

func TestHasTag(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	assert := Expect(test.t, TestStruct{})

	test.mockT.EXPECT().Error(gomock.Any()).Do(func(args ...interface{}) {
		if goerror.Cause(args[0].(error)) != ErrUnexported {
			t.Errorf("Expected %v, got %v", ErrUnexported, args[0])
		}
	})

	test.mockT.EXPECT().Error(gomock.Any()).Times(3).Do(func(args ...interface{}) {
		if goerror.Cause(args[0].(error)) != ErrTagNotFound {
			t.Errorf("Expected %v, got %v", ErrTagNotFound, args[0])
		}
	})
	assert.ExpectField("private").HasTag("Unknown")
	assert.ExpectField("WithoutTags").HasTag("Unknown")
	assert.ExpectField("Public").HasTag("Unknown")

	//check for absence of call Error
	assert.ExpectField("Public").HasTag("tag1").HasTag("tag2")

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

	field2 := assert.ExpectField("Public")

	if reflect.ValueOf(field).Pointer() != reflect.ValueOf(field2).Pointer() {
		t.Errorf("Unexpected %p = %p", &field, field2)
	}

	test.mockT.EXPECT().Error(gomock.Any()).Do(func(args ...interface{}) {
		if goerror.Cause(args[0].(error)) != ErrUnexported {
			t.Errorf("Expected %v", ErrUnexported)
		}
	})

	field = assert.ExpectField("private")
	if field == nil {
		t.Error("Unexpected nil")
	}

	if field.structField != nil {
		t.Errorf("Expected field.structField nil,got %v", field.structField)
	}
}
