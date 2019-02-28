// Example of an io.Pipe race. See: https://github.com/evanj/stream
package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// Returns an io.WriteClose that base64 decodes data written to it, and writes the plain bytes
// to output. This is an example of a race with io.Pipe writers.
func base64DecodedWriter(output io.Writer) io.WriteCloser {
	readPipe, writePipe := io.Pipe()
	go func() {
		decoder := base64.NewDecoder(base64.StdEncoding, readPipe)
		_, err := io.Copy(output, decoder)
		if err != nil {
			readPipe.CloseWithError(err)
			return
		}
		err = readPipe.Close()
		if err != nil {
			panic("readPipe must not return error on close")
		}
	}()

	return writePipe
}

func main() {
	fmt.Println("reading base64 encoded data from stdin ...")
	// take plain bytes from stdin, then base64-decode it using a Writer
	output := &bytes.Buffer{}
	decodedWriter := base64DecodedWriter(output)

	// we now have a Writer that base64 decodes data: copy into it
	_, err := io.Copy(decodedWriter, os.Stdin)
	if err != nil {
		panic(err)
	}
	// close the writer to indicate we are done
	err = decodedWriter.Close()
	if err != nil {
		panic(err)
	}

	// fails with race error without blockingPipeWriter and calling Close
	out := output.String()
	fmt.Println("encoded output:")
	os.Stdout.WriteString(out)
	fmt.Println("\n\nencoded output:")

	fmt.Println("decoded:")
	b, err := base64.StdEncoding.DecodeString(out)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(b)
	fmt.Println()
}
