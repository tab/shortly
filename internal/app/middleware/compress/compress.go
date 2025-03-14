package compress

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// Middleware is a middleware for compressing requests and responses
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalWriter := w

		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			compressedReader, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = compressedReader
			defer compressedReader.Close()
		}

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			compressedWriter := newCompressWriter(w)
			originalWriter = compressedWriter

			next.ServeHTTP(originalWriter, r)
			compressedWriter.Close()
			return
		}

		next.ServeHTTP(originalWriter, r)
	})
}

type compressReader struct {
	reader     io.ReadCloser
	gzipReader *gzip.Reader
}

// Read reads the compressed data
func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.gzipReader.Read(p)
}

// Close closes the reader
func (c *compressReader) Close() error {
	if err := c.reader.Close(); err != nil {
		return err
	}
	return c.gzipReader.Close()
}

func newCompressReader(reader io.ReadCloser) (*compressReader, error) {
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		reader:     reader,
		gzipReader: gzipReader,
	}, nil
}

type compressWriter struct {
	writer     http.ResponseWriter
	gzipWriter *gzip.Writer
}

// Header returns the writer header
func (c *compressWriter) Header() http.Header {
	return c.writer.Header()
}

// Write writes the compressed data
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.gzipWriter.Write(p)
}

// WriteHeader writes the header
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode >= http.StatusContinue && statusCode <= http.StatusIMUsed {
		c.writer.Header().Set("Content-Encoding", "gzip")
	}

	c.writer.WriteHeader(statusCode)
}

// Close closes the writer
func (c *compressWriter) Close() error {
	return c.gzipWriter.Close()
}

func newCompressWriter(writer http.ResponseWriter) *compressWriter {
	return &compressWriter{
		writer:     writer,
		gzipWriter: gzip.NewWriter(writer),
	}
}
