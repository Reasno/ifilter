package ifilter filters a collection of unknown values by a known interface

## Quick Start
```go
package main

import (
	"bytes"
	"fmt"
	"github.com/Reasno/ifilter"
	"io"
	"os"
)

func main() {
	var collection = ifilter.Collection{&os.File{}, &bytes.Buffer{}, struct{}{}, nil, 42}
	collection.Filter(func(readers []io.Reader) {
		for _, reader := range readers {
			fmt.Printf("%T\n", reader)
		}
	})
	/*
		*os.File
		*bytes.Buffer
	*/
}
```


 
