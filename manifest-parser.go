package m3u8

import "io"

//ManifestParser is the interface for the Generate and Read manifest methods
type ManifestParser interface {
	GenerateManifest() (io.Reader, error)
	ReadManifest(reader io.Reader) error
}
