package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ErrUnexpectedStatus is returned by PostJSON when the server responds with a
// non-2xx status code. Callers can type-assert to access the status and body.
type ErrUnexpectedStatus struct {
	StatusCode int
	Body       []byte
}

func (unexpectedStatusErr ErrUnexpectedStatus) Error() string {
	return fmt.Sprintf("unexpected status %d", unexpectedStatusErr.StatusCode)
}

// PostJSON is intended for external API calls that return JSON responses.
// It marshals body as JSON, POSTs it to url, and decodes the response into out,
// which must be a non-nil pointer.
// If the server responds with Content-Encoding: gzip the body is decompressed
// transparently; otherwise the raw body is decoded directly.
// A non-2xx status code is returned as ErrUnexpectedStatus, which carries the
// raw response body so callers can apply their own error parsing.
// opts are applied to the request before it is sent (e.g. setting headers).
func (client *Client) PostJSON(ctx context.Context, url string, body any, out any, opts ...RequestOption) error {
	encoded, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(encoded))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := client.Do(req, opts...)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	reader, err := DecompressResponse(resp)
	if err != nil {
		return fmt.Errorf("decompress response: %w", err)
	}
	defer func() { _ = reader.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		rawBody, _ := io.ReadAll(reader)
		return ErrUnexpectedStatus{StatusCode: resp.StatusCode, Body: rawBody}
	}

	if err := json.NewDecoder(reader).Decode(out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}
