package assert

import (
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
	private string
	Public  string
	SubStruct
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

func TestExpectField(t *testing.T) {
	test := setUp(t)
	defer test.tearDown()

	test.mockT.EXPECT().Helper().AnyTimes()

	assert := Expect(test.t, TestStruct{})
	field := assert.ExpectField("Public")
	if field == nil {
		t.Error("Unexpected nil")
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
	//assert.ExpectField

}
