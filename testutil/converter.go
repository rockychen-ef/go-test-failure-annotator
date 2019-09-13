package testutil

import (
	"bytes"
	"io"
)

func ToBytes(readCloser io.ReadCloser) []byte {
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(readCloser)

	return buffer.Bytes()
}
