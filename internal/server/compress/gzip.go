package compress

import (
	"compress/gzip"
	"io"
	"net/http"
)

// compressWriter реализует интерфейс http.ResponseWriter и позволяет прозрачно для сервера
// сжимать передаваемые данные и выставлять правильные HTTP-заголовки
type СompressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewCompressWriter(w http.ResponseWriter) *СompressWriter {
	return &СompressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *СompressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *СompressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *СompressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *СompressWriter) Close() error {
	return c.zw.Close()
}

type СompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewCompressReader(r io.ReadCloser) (*СompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &СompressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c *СompressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *СompressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
