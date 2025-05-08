package transport

import (
	"io"
	"os"

	"go.lsp.dev/jsonrpc2"
	"go.uber.org/multierr"
)

type readWriteCloser struct {
	readCloser  io.ReadCloser
	writeCloser io.WriteCloser
}

func (r *readWriteCloser) Read(b []byte) (int, error) {
	return r.readCloser.Read(b)
}

func (r *readWriteCloser) Write(b []byte) (int, error) {
	return r.writeCloser.Write(b)
}

func (r *readWriteCloser) Close() error {
	return multierr.Append(r.readCloser.Close(), r.writeCloser.Close())
}

func CreateTransport() jsonrpc2.Conn {
	return jsonrpc2.NewConn(
		jsonrpc2.NewStream(
			&readWriteCloser{
				readCloser:  os.Stdin,
				writeCloser: os.Stdout,
			},
		),
	)
}
