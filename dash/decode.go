package dash

import (
	"encoding/xml"
	"io"
)

//Parse decodes a MPD file into an MPD element.
//It validates MPD according to specs and returns an error if validation fails.
func (m *MPD) Parse(reader io.Reader) error {
	err := xml.NewDecoder(reader).Decode(&m)
	if err != nil {
		return err
	}

	return m.validate()
}
