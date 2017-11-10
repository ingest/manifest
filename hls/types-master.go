package hls

import (
	"fmt"
	"net/http"
)

//MasterPlaylist represents a Master Playlist and its tags
type MasterPlaylist struct {
	URI                 string       // Location of the master playlist
	M3U                 bool         // Represents tag #EXTM3U. Indicates if present. MUST be present.
	Version             int          // Represents tag #EXT-X-VERSION. MUST be present.
	Variants            []*Variant   // Represents the #EXT-X-I-FRAME-STREAM-INF and #EXT-X-STREAM-INF playlists
	Renditions          []*Rendition // Represents the #EXT-X-MEDIA tags
	SessionData         []*SessionData
	SessionKeys         []*Key
	IndependentSegments bool // Represents tag #EXT-X-INDEPENDENT-SEGMENTS. Applies to every Media Segment of every Media Playlist referenced. V6 or higher.
	StartPoint          *StartPoint
}

// Request creates a new http request ready to retrieve the segment
func (m *MasterPlaylist) Request() (*http.Request, error) {
	req, err := http.NewRequest("GET", m.URI, nil)
	if err != nil {
		return req, fmt.Errorf("failed to construct request: %v", err)
	}

	return req, nil
}

//Rendition represents the tag #EXT-X-MEDIA
//
//Relates Media Playlists with alternative renditions of the same content.
//Eg. Audio only playlists containing English, French and Spanish renditions of the same content.
//One or more X-MEDIA tags with same GroupID and Type sets a group of renditions and MUST meet the following constraints:
//  -Tags in the same group MUST have different Name att.
//  -MUST NOT have more than one member with a Default att of YES
//  -All members whose AutoSelect att is YES MUST have Language att with unique values
//
type Rendition struct {
	Type            string //Possible Values: AUDIO, VIDEO, SUBTITLES, CLOSED-CAPTIONS. Required.
	URI             string //URI containing the media playlist. If type is CLOSED-CAPTIONS, URI MUST NOT be present.
	GroupID         string //Required.
	Language        string //Optional. Identifies the primary language used in the rendition. Must be one of the standard tags RFC5646
	AssocLanguage   string //Optional. Language tag RFC5646
	Name            string //Required. Description of the rendition. SHOULD be written in the same language as Language
	Default         bool   //Possible Values: YES, NO. Optional. Defines if rendition should be played by client if user doesn't choose a rendition. Default: NO
	AutoSelect      bool   //Possible Values: YES, NO. Optional. Client MAY choose this rendition if user doesn't choose one. if present, MUST be YES if Default=YES. Default: NO.
	Forced          bool   //Possible Values: YES, NO. Optional. MUST NOT be present unless Type is SUBTITLES. Default: NO.
	InstreamID      string //Specifies a rendition within the Media Playlist. MUST NOT be present unless Type is CLOSED-CAPTIONS. Possible Values: CC1, CC2, CC3, CC4, or SERVICEn where n is int between 1 - 63
	Characteristics string //Optional. One or more Uniform Type Indentifiers separated by comma. Each UTI indicates an individual characteristic of the Rendition.

	masterPlaylist *MasterPlaylist // MasterPlaylist is included to be used internally for resolving relative resource locations
}

// Request creates a new http request ready to retrieve the segment
func (r *Rendition) Request() (*http.Request, error) {
	uri, err := r.AbsoluteURL()
	if err != nil {
		return nil, fmt.Errorf("failed building resource url: %v", err)
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return req, fmt.Errorf("failed to construct request: %v", err)
	}

	return req, nil
}

// AbsoluteURL will resolve the rendition URI to a absolute path, given it is a URL.
func (r *Rendition) AbsoluteURL() (string, error) {
	return resolveURLReference(r.masterPlaylist.URI, r.URI)
}

// Variant represents the tag #EXT-X-STREAM-INF: <attribute-list> and tag #EXT-X-I-FRAME-STREAM-INF.
// #EXT-X-STREAM-INF specifies a Variant Stream, which is one of the ren which can be combined to play the presentation.
// A URI line following the tag indicates the Media Playlist carrying a rendition of the Variant Stream and it MUST be present.
//
// #EXT-X-I-FRAME-STREAM-INF identifies Media Playlist file containing the I-frames of a multimedia presentation.
// It supports the same parameters as EXT-X-STREAM-INF except Audio, Subtitles and ClosedCaptions.
type Variant struct {
	IsIframe       bool    //Identifies if #EXT-X-STREAM-INF or #EXT-X-I-FRAME-STREAM-INF
	URI            string  //If #EXT-X-STREAM-INF, URI line MUST follow the tag. If #EXT-X-I-FRAME-STREAM-INF, URI MUST appear as an attribute of the tag.
	ProgramID      int64   //Removed on Version 6
	Bandwidth      int64   //Required. Peak segment bit rate.
	AvgBandwidth   int64   //Optional. Average segment bit rate of the Variant Stream.
	Codecs         string  //Optional. Comma-separated list of formats. Valid formats are the ones specified in RFC6381. SHOULD be present.
	Resolution     string  //Optional. Optimal pixel resolution.
	FrameRate      float64 //Optional. Maximum frame rate. Optional. SHOULD be included if any video exceeds 30 frames per second.
	Audio          string  //Optional. Indicates the set of audio renditions that SHOULD be used. MUST match GroupID value of an EXT-X-MEDIA tag whose Type is AUDIO.
	Video          string  //Optional. Indicates the set of video renditions that SHOULD be used. MUST match GroupID value of an EXT-X-MEDIA tag whose Type is VIDEO.
	Subtitles      string  //Optional. Indicates the set of subtitle renditions that SHOULD be used. MUST match GroupID value of an EXT-X-MEDIA tag whose Type is SUBTITLES.
	ClosedCaptions string  //Optional. Indicates the set of closed-caption renditions that SHOULD be used. Can be quoted-string or NONE.
	// If NONE, all EXT-X-STREAM-INF MUST have this attribute as NONE. If quoted-string, MUST match GroupID value of an EXT-X-MEDIA tag whose Type is CLOSED-CAPTIONS.

	masterPlaylist *MasterPlaylist // MasterPlaylist is included to be used internally for resolving relative resource locations
}

// Request creates a new http request ready to retrieve the segment
func (v *Variant) Request() (*http.Request, error) {
	uri, err := v.AbsoluteURL()
	if err != nil {
		return nil, fmt.Errorf("failed building resource url: %v", err)
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return req, fmt.Errorf("failed to construct request: %v", err)
	}
	return req, nil
}

// AbsoluteURL will resolve the variant URI to a absolute path, given it is a URL.
func (v *Variant) AbsoluteURL() (string, error) {
	return resolveURLReference(v.masterPlaylist.URI, v.URI)
}

// SessionData represents tag #EXT-X-SESSION-DATA.
// Master Playlist MAY contain more than one tag with the same DataID but the Language MUST be different.
// Introduced in HLSv7.
type SessionData struct {
	DataID   string //Required. SHOULD conform with a reverse DNS naming convention.
	Value    string //Required IF URI is not present. Contains the session data
	URI      string //Required IF Value is not present. Resource with the session data
	Language string //Optional. RFC5646 language tag that identifies the language of the data

	masterPlaylist *MasterPlaylist // MasterPlaylist is included to be used internally for resolving relative resource locations
}

// Request creates a new http request ready to retrieve the segment
func (s *SessionData) Request() (*http.Request, error) {
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

// AbsoluteURL will resolve the SessionData URI to a absolute path, given it is a URL.
func (s *SessionData) AbsoluteURL() (string, error) {
	return resolveURLReference(s.masterPlaylist.URI, s.URI)
}

// StartPoint represents tag #EXT-X-START.
// Indicates preferred point at which to start playing a Playlist.
type StartPoint struct {
	TimeOffset float64 //Required. If positive, time offset from the beginning of the Playlist. If negative, time offset from the end of the last segment of the playlist
	Precise    bool    //Possible Values: YES or NO.
}
