//Package m3u8 holds the main interface for Manifest encode/parse.
package m3u8

import "io"

//Manifest is the interface for the Generate and Read manifest methods
type Manifest interface {
	Encode() (io.Reader, error)
	Parse(reader io.Reader) error
}
