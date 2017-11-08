package hls

import (
	"fmt"
	"net/http"
	"time"
)

//Segment represents the Media Segment and its tags
type Segment struct {
	ID              int //Sequence number
	URI             string
	Inf             *Inf //Required.
	Byterange       *Byterange
	Discontinuity   bool //Represents tag #EXT-X-DISCONTINUITY. MUST be present if there's change in file format; number, type and identifiers of tracks or timestamp sequence
	Keys            []*Key
	Map             *Map
	ProgramDateTime time.Time //Represents tag #EXT-X-PROGRAM-DATE-TIME
	DateRange       *DateRange

	mediaPlaylist *MediaPlaylist // MediaPlaylist is included to be used internally for resolving relative resource locations
}

// Request creates a new http request ready to send to retrieve the segment
func (s *Segment) Request() (*http.Request, error) {
	uri, err := s.AbsoluteURL()
	if err != nil {
		return nil, fmt.Errorf("failed building resource url: %v", err)
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return req, fmt.Errorf("failed to construct request: %v", err)
	}
	return req, nil
}

// AbsoluteURL will resolve the segment URI to a absolute path, given it is a relative URL.
func (s *Segment) AbsoluteURL() (string, error) {
	return resolveURLReference(s.mediaPlaylist.URI, s.URI)
}

// Segments implements golang/sort interface to sort a Segment slice by Segment ID
type Segments []*Segment

func (s Segments) Len() int {
	return len(s)
}
func (s Segments) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Segments) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}

// Inf represents tag
// 		#EXTINF: <duration>,[<title>]
type Inf struct {
	Duration float64
	Title    string
}

// Byterange represents tag #EXT-X-BYTERANGE.
// Introduced in HLSv4.
// Format: length[@offset].
type Byterange struct {
	Length int64
	Offset *int64
}

// Equal determines if the two byterange objects are equal and contain the same values
func (b *Byterange) Equal(other *Byterange) bool {
	if b != nil && other != nil {
		if b.Length != other.Length {
			return false
		}

		if b.Offset != nil && other.Offset != nil {
			return *b.Offset == *other.Offset
		}
	}

	return b == other
}

// Key represents tags #EXT-X-KEY:<attribute=value> and #EXT-X-SESSION-KEY. Specifies how to decrypt an encrypted media segment.
// #EXT-X-SESSION-KEY is exclusively a Master Playlist tag (HLS V7) and it SHOULD be used if multiple Variant Streams use the same encryption keys.
// TODO(jstackhouse): Split SESSION-KEY into it's own type as it's got different validation properties, and is part of the master playlist, not media playlist.
type Key struct {
	IsSession         bool   //Identifies if #EXT-X-KEY or #EXT-X-SESSION-KEY. If #EXT-X-SESSION-KEY, Method MUST NOT be NONE.
	Method            string //Required. Possible Values: NONE, AES-128, SAMPLE-AES. If NONE, other attributes MUST NOT be present.
	URI               string //Required unless the method is NONE. Specifies how to get the key for the encryption method.
	IV                string //Optional. Hexadecimal that specifies a 128-bit int Initialization Vector to be used with the key.
	Keyformat         string //Optional. Specifies how the key is represented in the resource. V5 or higher
	Keyformatversions string //Optional. Indicates which Keyformat versions this instance complies with. Default value is 1. V5 or higher

	masterPlaylist *MasterPlaylist // MasterPlaylist is included to be used internally for resolving relative resource locations for Session keys
	mediaPlaylist  *MediaPlaylist  // MediaPlaylist is included to be used internally for resolving relative resource locations
}

// Request creates a new http request ready to retrieve the segment
func (k *Key) Request() (*http.Request, error) {
	uri, err := k.AbsoluteURL()
	if err != nil {
		return nil, fmt.Errorf("failed building resource url: %v", err)
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return req, fmt.Errorf("failed to construct request: %v", err)
	}
	return req, nil
}

// AbsoluteURL will resolve the Key URI to a absolute path, given it is a URL.
func (k *Key) AbsoluteURL() (string, error) {
	var uri string
	var err error

	if k.IsSession {
		uri, err = resolveURLReference(k.masterPlaylist.URI, k.URI)
	} else {
		uri, err = resolveURLReference(k.mediaPlaylist.URI, k.URI)
	}

	return uri, err
}

// Equal checks whether all public fields are equal in a Key with the exception of the IV field.
func (k *Key) Equal(other *Key) bool {
	if k != nil && other != nil {
		if k.IsSession != other.IsSession {
			return false
		}

		if k.Keyformat != other.Keyformat {
			return false
		}

		if k.Keyformatversions != other.Keyformatversions {
			return false
		}

		if k.Method != other.Method {
			return false
		}

		if k.URI != other.URI {
			return false
		}

		return true
	}

	// are they both nil
	return k == other
}

//Map represents tag #EXT-X-MAP:<attribute=value>. Specifies how to get the Media Initialization Section
type Map struct {
	URI       string     //Required.
	Byterange *Byterange //Optional. Indicates the byte range into the URI resource containing the Media Initialization Section.

	mediaPlaylist *MediaPlaylist // MediaPlaylist is included to be used internally for resolving relative resource locations
}

// Equal determines if the two maps are equal, does not check private fields for equality so this does not guarantee that two maps will act identically.
// Works on nil structures, if both m and other are nil, they are considered equal.
func (m *Map) Equal(other *Map) bool {
	if m != nil && other != nil {
		return m.URI == other.URI && m.Byterange.Equal(other.Byterange)
	}

	return m == other
}

// Request creates a new http request ready to retrieve the segment
func (m *Map) Request() (*http.Request, error) {
	uri, err := m.AbsoluteURL()
	if err != nil {
		return nil, fmt.Errorf("failed building resource url: %v", err)
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return req, fmt.Errorf("failed to construct request: %v", err)
	}
	return req, nil
}

// AbsoluteURL will resolve the EXT-X-MAP URI to a absolute path, given it is a URL.
func (m *Map) AbsoluteURL() (string, error) {
	return resolveURLReference(m.mediaPlaylist.URI, m.URI)
}

//DateRange represents tag #EXT-X-DATERANGE:<attribute=value>.
//
//If present, playlist MUST also contain at least one EXT-X-PROGRAM-DATE-TIME tag.
//Tags with the same Class MUST NOT indicate ranges that overlap.
type DateRange struct {
	ID               string    //Required. If more than one tag with same ID exists, att values MUST be the same.
	Class            string    //Optional. Specifies some set of attributes and their associated value semantics.
	StartDate        time.Time //Required.
	EndDate          time.Time //Optional.
	Duration         *float64  //Optional. If both EndDate and Duration present, check EndDate equal to Duration + StartDate
	PlannedDuration  *float64  //Optional. Expected duration.
	XClientAttribute []string  //Optional. Namespace reserved for client-defined att. eg. X-COM-EXAMPLE="example".
	EndOnNext        bool      //Optional. Possible Value: YES. Indicates the end of the current date range is equal to the start date of the following range of the samePROGRAM-DATE-TIME class.
	SCTE35           *SCTE35
}

//SCTE35 represents a DateRange attribute SCTE35-OUT, SCTE35-IN or SCTE35-CMD
type SCTE35 struct {
	Type  string //Possible Values: IN, OUT, CMD
	Value string //big-endian binary representation of the splice_info_section(), expressed as a hexadecimal-sequence.
}
