package m3u8

import "io"

//Manifest is the interface for the Generate and Read manifest methods
type Manifest interface {
	GenerateManifest() (io.Reader, error)
	ReadManifest(reader io.Reader) error
}
