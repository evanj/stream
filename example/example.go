// Example of converting a Writer into a Reader
package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/evanj/stream"
)

func newExampleReader() io.Reader {
	return strings.NewReader("hello some example data")
}

func main() {
	// Read an entire stream and write it to os.Stdout
	fmt.Println("Original source stream:")
	source := newExampleReader()
	_, err := io.Copy(os.Stdout, source)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println()

	// Wrap os.Stdout in a Base64 encoder
	fmt.Println("Base64 encoding with a Writer:")
	source = newExampleReader()
	dest := base64.NewEncoder(base64.StdEncoding, os.Stdout)
	_, err = io.Copy(dest, source)
	if err != nil {
		panic(err)
	}
	// Must close to flush buffered state
	err = dest.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println()

	// Use the processor to wrap the encoder
	fmt.Println("Base64 encoding with a *Reader*:")
	source = stream.ProcessReaderWithWriter(
		newExampleReader(),
		func(w io.Writer) io.WriteCloser { return base64.NewEncoder(base64.StdEncoding, w) })
	_, err = io.Copy(os.Stdout, source)
	if err != nil {
		panic(err)
	}
	fmt.Println()
	fmt.Println()
}
