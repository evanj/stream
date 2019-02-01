package stream

import (
	"bytes"
	"encoding/base64"
	"io"
	"strings"
	"testing"
)

// needs a pading byte when base64 encoded
const decoded = "hello input"

var encoded = base64.StdEncoding.EncodeToString([]byte(decoded))

func TestProcessReader(t *testing.T) {
	r := ProcessReaderWithWriter(
		strings.NewReader(decoded),
		func(w io.Writer) io.WriteCloser { return base64.NewEncoder(base64.StdEncoding, w) })

	output := &bytes.Buffer{}
	_, err := io.Copy(output, r)
	if err != nil {
		t.Fatal(err)
	}

	out := output.String()
	if out != encoded {
		t.Errorf("expected %s got %s", encoded, out)
	}
}

func TestProcessWriter(t *testing.T) {
	output := &bytes.Buffer{}
	w := ProcessWriterWithReader(
		output,
		func(r io.Reader) io.Reader { return base64.NewDecoder(base64.StdEncoding, r) })

	_, err := io.Copy(w, strings.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
	}
	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}

	// fails with race error without blockingPipeWriter and calling Close
	out := output.String()
	if out != decoded {
		t.Errorf("expected %s got %s", decoded, out)
	}
}
