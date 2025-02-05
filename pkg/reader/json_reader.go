package reader

import (
	"bytes"
	"io"
)

func GetBody(r io.ReadCloser) ([]byte, error) {
	var buf bytes.Buffer

	n, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	if err = r.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes()[:n], nil
}
