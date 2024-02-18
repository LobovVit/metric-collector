package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("write data to compress temporary buffer: %w", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("compress data: %w", err)
	}
	return b.Bytes(), nil
}
