# Converting Byte Streams: Readers into Writers and Writers into Readers

A common use of Go's `Reader` and `Writer` interface is to transform a byte stream into a different byte stream. For example, compression is a `Writer`, and decryption is a `Reader`. This package contains utilities to convert a Reader into a Writer using io.Pipe. They correctly forward any errors and syncronize on Close to prevent lost data.

For details, see my blog post (TODO: coming soon)

You can also see the example.
