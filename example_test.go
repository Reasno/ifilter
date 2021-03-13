package ifilter_test

import (
	"bytes"
	"fmt"
	"github.com/Reasno/ifilter"
	"io"
	"os"
)

func Example() {
	var collection = ifilter.Collection{&os.File{}, &bytes.Buffer{}, struct{}{}, nil, 42, (io.Reader)(nil)}
	collection.Filter(func(reader io.Reader) {
		fmt.Printf("%T\n", reader)
	})
	// Output:
	// *os.File
	// *bytes.Buffer
}
