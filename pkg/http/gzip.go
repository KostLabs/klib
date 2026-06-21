package http

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// DecompressResponse returns a reader over the response body, transparently
// decompressing gzip when the Content-Encoding header says so.
// The caller is responsible for closing the returned ReadCloser.
func DecompressResponse(resp *http.Response) (io.ReadCloser, error) {
	if resp.Header.Get("Content-Encoding") != "gzip" {
		return resp.Body, nil
	}

	gzipReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gzip reader: %w", err)
	}

	return gzipReader, nil
}
