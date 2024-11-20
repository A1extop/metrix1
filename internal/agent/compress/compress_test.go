package compress

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompressData(t *testing.T) {
	inputData := []byte("test data for compression")

	compressedData, err := CompressData(inputData)

	assert.NoError(t, err, "expected no error from CompressData")

	assert.NotEmpty(t, compressedData, "compressed data should not be empty")

	gzipReader, err := gzip.NewReader(bytes.NewReader(compressedData))
	assert.NoError(t, err, "expected no error creating gzip reader")
	defer gzipReader.Close()

	decompressedData, err := io.ReadAll(gzipReader)
	assert.NoError(t, err, "expected no error reading decompressed data")

	assert.Equal(t, inputData, decompressedData, "decompressed data should match original input")
}
