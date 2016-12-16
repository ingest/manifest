package dash

import (
	"bytes"
	"encoding/xml"
	"io"

	"github.com/ingest/manifest"
)

//Encode marshals an MPD structure into an MPD XML structure.
func (m *MPD) Encode() (io.Reader, error) {
	//Validates MPD structure according to specs before encoding
	if err := m.validate(); err != nil {
		return nil, err
	}

	output, err := xml.MarshalIndent(m, "", " ")
	if err != nil {
		return nil, err
	}

	buf := manifest.NewBufWrapper()
	buf.WriteString(xml.Header)
	buf.Write(output)

	return bytes.NewReader(buf.Buf.Bytes()), buf.Err
}
