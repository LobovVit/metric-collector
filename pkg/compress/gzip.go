// Package compress - contains methods for compressing and decompressing data
package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// Compress - methods for compressing data
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

// UnCompress - methods for decompressing data
func UnCompress(data []byte) ([]byte, error) {
	b := bytes.NewBuffer(data)
	var r io.Reader
	r, err := gzip.NewReader(b)
	if err != nil {
		return nil, fmt.Errorf("gzip reader: %w", err)
	}
	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("read data to uncompress temporary buffer: %w", err)
	}
	resData := resB.Bytes()
	return resData, nil
}
