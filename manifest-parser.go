package m3u8

import (
	"bytes"
	"io"
)

//ManifestParser is the interface for the Generate and Read manifest methods
type ManifestParser interface {
	GenerateManifest() (*bytes.Buffer, error)
	ReadManifest(reader io.Reader) error
}
