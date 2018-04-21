
[![Build Status](https://travis-ci.org/arteev/tag-assert.svg?branch=master)](https://travis-ci.org/arteev/tag-assert)
[![Coverage Status](https://coveralls.io/repos/arteev/tag-assert/badge.svg?branch=master&service=github)](https://coveralls.io/github/arteev/tag-assert?branch=master)
[![GoDoc](https://godoc.org/github.com/arteev/tag-assert?status.png)](https://godoc.org/github.com/arteev/go-assert)


# tag-assert

Checking tags of Golang structures

## Install

``` 
go get github.com/arteev/tag-assert
```

## Using

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