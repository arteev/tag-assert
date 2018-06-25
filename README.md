
[![Build Status](https://travis-ci.org/arteev/tag-assert.svg?branch=master)](https://travis-ci.org/arteev/tag-assert)
[![Coverage Status](https://coveralls.io/repos/arteev/tag-assert/badge.svg?branch=master&service=github)](https://coveralls.io/github/arteev/tag-assert?branch=master)
[![GoDoc](https://godoc.org/github.com/arteev/tag-assert?status.png)](https://godoc.org/github.com/arteev/tag-assert)


# tag-assert

Checking tags of Golang structures

## Install

``` 
go get github.com/arteev/tag-assert
```

## Usage

```go

//example.go
package example
type ExampleStruct struct {
	Name string `xml:"Name" json:"name,omitempty"`
	ID   int    `xml:"ID" json:"rn"`
}


//example_test.go
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
}

func TestExampleStructTagsValueFailed(t *testing.T) {
	v := ExampleStruct{}
	assert.Expect(t, v).ExpectField("ID").
		Assert("xml", "ID").
		Assert("json", "id") // this error

	assert.Expect(t, v).ExpectField("SN").
		Assert("xml", "SN").
		Assert("json", "social_number")

}

```

```bash
 ~: go test

--- FAIL: TestExampleStructTagsValueFailed (0.00s)
	example_test.go:22: ExampleStruct.ID: Tag <json> does not have a value of <id>,but actual <rn>
	example_test.go:23: ExampleStruct.ID: Tag <bson> not found
	example_test.go:24: ExampleStruct.ID: Tag <json> does not have a value of <id>,but actual <rn>
	example_test.go:26: ExampleStruct: Field <SN> not found
	example_test.go:27: ExampleStruct.SN: Tag <xml> not found
	example_test.go:28: ExampleStruct.SN: Tag <json> not found
	example_test.go:30: ExampleStruct: Field <private> is private
	example_test.go:31: ExampleStruct.private: Tag <xml> not found
FAIL
exit status 1
FAIL	github.com/arteev/tag-assert/_example	0.001s

```