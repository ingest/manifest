//Package manifest holds the main interface for Manifest encode/parse.
package manifest

import "io"

// Parser is the interface by which we convert the textual format into a Go based structure format
type Parser interface {
	Parse(reader io.Reader) error
}

// Encoder is the interface by which we convert our Go based structured format into the textual representation
type Encoder interface {
	Encode() (io.Reader, error)
}
