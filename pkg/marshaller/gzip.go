package marshaller

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
)

// MarshalGzip serialises value to JSON then gzip-compresses the result.
func MarshalGzip(value any) ([]byte, error) { //goverifier:ignore:any-type
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(jsonBytes); err != nil {
		return nil, fmt.Errorf("gzip write: %w", err)
	}
	if err := gz.Close(); err != nil {
		return nil, fmt.Errorf("gzip close: %w", err)
	}
	return buf.Bytes(), nil
}

// UnmarshalGzip decompresses gzip bytes then JSON-decodes into target.
func UnmarshalGzip(data []byte, target any) error { //goverifier:ignore:any-type
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer func() { _ = gr.Close() }()

	raw, err := io.ReadAll(gr)
	if err != nil {
		return fmt.Errorf("gzip read: %w", err)
	}
	return json.Unmarshal(raw, target)
}
