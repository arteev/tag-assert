package main

type ExampleStruct struct {
	Name string `xml:"Name" json:"name,omitempty"`
	ID   int    `xml:"ID" json:"rn"`
}

func main() {
	//see: example1_test.go
}
