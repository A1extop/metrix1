package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func CompressData(data []byte) ([]byte, error) {
	var compressedData bytes.Buffer
	gzipWriter := gzip.NewWriter(&compressedData)

	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error compressing data with gzip: %w", err)
	}

	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf("error closing gzip writer: %w", err)
	}

	return compressedData.Bytes(), nil
}
