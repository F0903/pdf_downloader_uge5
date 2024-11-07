package utils

import (
	"context"
	"io"
)

// Based on https://ixday.github.io/post/golang-cancel-copy/

// here is some syntaxic sugar inspired by the Tomas Senart's video,
// it allows me to inline the Reader interface
type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

func CancelableCopy(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	written, err := io.Copy(dst, readerFunc(func(p []byte) (int, error) {
		select {

		// if context has been canceled
		case <-ctx.Done():
			// stop process and propagate "context canceled" error
			return 0, ctx.Err()
		default:
			// otherwise just run default io.Reader implementation
			return src.Read(p)
		}
	}))
	return written, err
}
