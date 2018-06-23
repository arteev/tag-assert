package example

type ExampleStruct struct {
	Name    string `xml:"Name" json:"name,omitempty"`
	ID      int    `xml:"ID" json:"rn"`
	private string `xml:"private"`
}

func main() {
	//see: example1_test.go
}
