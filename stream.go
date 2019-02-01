// Package stream converts Readers to Writers and vice-versa.
package stream

import "io"

// ProcessReaderWithWriter returns an io.Reader that will process bytes read from it using a Writer
// created by newWriter. The function newWriter must return a WriteCloser that writes to its
// argument writer. Close will be called on it to flush any buffered state.
func ProcessReaderWithWriter(source io.Reader, newWriter func(writer io.Writer) io.WriteCloser) io.Reader {
	reader, writer := io.Pipe()
	processor := newWriter(writer)
	go func() {
		// Copy the source into the processor. Any errors get returned on reader.Read().
		_, err := io.Copy(processor, source)
		if err != nil {
			writer.CloseWithError(err)
			return
		}
		err = processor.Close()
		if err != nil {
			writer.CloseWithError(err)
			return
		}

		err = writer.Close()
		if err != nil {
			// This should never happen: a PipeWriter Close() should not return an error
			panic(err)
		}
	}()

	return reader
}

// Wraps an io.PipeWriter so .Close() blocks until the read side has finished. Avoids a race
// condition where the caller calls Close, then expects the output to be complete. See the unit
// test for a case that fails if we just use io.PipeWriter directly.
type blockingPipeWriter struct {
	writer *io.PipeWriter
	done   chan struct{}
}

func (b *blockingPipeWriter) Write(source []byte) (int, error) {
	return b.writer.Write(source)
}
func (b *blockingPipeWriter) Close() error {
	err := b.writer.Close()
	if err != nil {
		return err
	}
	// wait for the reading side of the pipe to finish flushing/writing to the destination
	<-b.done
	return nil
}

// ProcessWriterWithReader returns an io.WriteCloser that will process bytes written to it using the
// reader created by newReader. The function newReader must return an io.Reader that reads from its
// argument reader. You must close the WriteCloser returned by this function to flush any buffered
// state to the destination.
func ProcessWriterWithReader(destination io.Writer, newReader func(reader io.Reader) io.Reader) io.WriteCloser {
	reader, writer := io.Pipe()
	processor := newReader(reader)
	blockingWriter := &blockingPipeWriter{writer, make(chan struct{})}

	go func() {
		// Copy from the processor into the destination
		_, err := io.Copy(destination, processor)
		if err != nil {
			reader.CloseWithError(err)
			return
		}
		err = reader.Close()
		if err != nil {
			// io.PipeReader should never return an error
			panic(err)
		}

		// indicate we are actually done
		close(blockingWriter.done)
	}()

	return blockingWriter
}
