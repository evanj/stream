# Converting Byte Streams: Readers into Writers and Writers into Readers

A common use of Go's `Reader` and `Writer` interface is to transform a byte stream into a different byte stream. For example, compression is a `Writer`, and decryption is a `Reader`. This package contains utilities to convert a Reader into a Writer using io.Pipe. They correctly forward any errors and syncronize on Close to prevent lost data.

For details, see my blog post (TODO: coming soon)


## io.PipeWriter race on on Close

It is easy to screw up using io.PipeWriter. Let's say we want to wrap an existing Writer, but transform the bytes using a Reader. We do something that looks like this:

    Tranform Goroutine: PipeReader --> Transform Reader --> io.Copy --> Output io.Writer

This transform Goroutine calls Transform.Read, which calls PipeReader.Read. This transforms bytes, which it then writes to the final output. This Goroutine runs until the Pipe is closed.

We use this from a "main" Goroutine, by writing into the pipe:

    Main Goroutine --> PipeWriter

When we have written all our data, we needs to call PipeWriter.Close(), to indicate the stream is done. We then will want to do something with the output. This is a race!

* Transform Goroutine: On close, it may need to transform and flush one last chunk to the output.
* Main Goroutine: After close, it will try to read the output.

To solve it, we need to block PipeWriter.Close() until the transform Goroutine is actually done. The `blockingPipeWriter` type in the `stream` package exists for this purpose. See `examplewrong` for a standalone program that triggers this problem.

### Example

The `examplewrong` program demonstrates this problem. Try running:

```shell
echo "hello" | base64 | go run -race examplewrong.go
```

The unit test in stream_test.go also covers this case.
