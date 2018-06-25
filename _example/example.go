package example

type ExampleStruct struct {
	Name       string `xml:"Name" json:"name,omitempty"`
	ID         int    `xml:"ID" json:"rn"`
	private    string `xml:"private"`
	WithoutTag string
}

func main() {
	//see: example1_test.go
}
