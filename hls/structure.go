package hls

import "time"

const (
	sub     = "SUBTITLES"
	aud     = "AUDIO"
	vid     = "VIDEO"
	cc      = "CLOSED-CAPTIONS"
	aes     = "AES-128"
	none    = "NONE"
	sample  = "SAMPLE-AES"
	boolYes = "YES"
	boolNo  = "NO"
)

//MediaPlaylist represents a Media Playlist and its tags.
//
//TODO:(sliding window) - add field for sliding window to represent either the max amount of segments
//or the max duration of a window (TBD). Also would be useful to add variable to track the current first and last sequence numbers
//as a helper to adding and removing segments and tracking MediaSequence, DiscontinuitySequence etc
//
type MediaPlaylist struct {
	M3U                   bool //Represents tag #EXTM3U. Indicates if present. MUST be present.
	Version               int  //Represents tag #EXT-X-VERSION. MUST be present.
	Segments              Segments
	TargetDuration        int         //Required. Represents tag #EXT-X-TARGETDURATION. MUST BE >= EXTINF
	MediaSequence         int         //Represents tag #EXT-X-MEDIA-SEQUENCE. Number of the first media sequence in the playlist.
	DiscontinuitySequence int         //Represents tag #EXT-X-DISCONTINUITY-SEQUENCE. If present, MUST appear before the first Media Segment. MUST appear before any EXT-X-DISCONTINUITY Media Segment tag.
	EndList               bool        //Represents tag #EXT-X-ENDLIST. Indicates no more media segments will be added to the playlist.
	Type                  string      //Possible Values: EVENT or VOD. Represents tag #EXT-X-PLAYLIST-TYPE. If EVENT - segments can only be added to the end of playlist. If VOD - playlist cannot change. If segments need to be removed from playlist, this tag MUST NOT be present
	IFramesOnly           bool        //Represents tag #EXT-X-I-FRAMES-ONLY. If present, segments MUST begin with either a Media Initialization Section or have a EXT-X-MAP tag.
	AllowCache            bool        //Possible Values: YES or NO. Represents tag #EXT-X-ALLOW-CACHE. Versions 3 - 6 only.
	IndependentSegments   bool        //Represents tag #EXT-X-INDEPENDENT-SEGMENTS. Applies to every Media Segment in the playlist.
	StartPoint            *StartPoint //Represents tag #EXT-X-START
}

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
}

//Inf represents tag #EXTINF<duration>,[<title>]
type Inf struct {
	Duration float64
	Title    string
}

//Byterange represents tag #EXT-X-BYTERANGE (V4 or higher) or a Byterange attribute of tag #EXT-X-MAP.
//Format: length[@offset].
type Byterange struct {
	Length int64
	Offset *int64
}

//Key represents tags #EXT-X-KEY:<attribute=value> and #EXT-X-SESSION-KEY. Specifies how to decrypt an encrypted media segment.
//#EXT-X-SESSION-KEY is exclusively a Master Playlist tag (HLS V7) and it SHOULD be used if multiple Variant Streams use the same encryption keys.
type Key struct {
	IsSession         bool   //Identifies if #EXT-X-KEY or #EXT-X-SESSION-KEY. If #EXT-X-SESSION-KEY, Method MUST NOT be NONE.
	Method            string //Required. Possible Values: NONE, AES-128, SAMPLE-AES. If NONE, other attributes MUST NOT be present.
	URI               string //Required unless the method is NONE. Specifies how to get the key for the encryption method.
	IV                string //Optional. Hexadecimal that specifies a 128-bit int Initialization Vector to be used with the key.
	Keyformat         string //Optional. Specifies how the key is represented in the resource. V5 or higher
	Keyformatversions string //Optional. Indicates which Keyformat versions this instance complies with. Default value is 1. V5 or higher
}

//Map represents tag #EXT-X-MAP:<attribute=value>. Specifies how to get the Media Initialization Section
type Map struct {
	URI       string     //Required.
	Byterange *Byterange //Optional. Indicates the byte range into the URI resource containing the Media Initialization Section.
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
	EndOnNext        bool      //Optional. Possible Value: YES. Indicates the end of the current date range is equal to the start date of the following range of the same class.
	SCTE35           *SCTE35
}

//SCTE35 represents a DateRange attribute SCTE35-OUT, SCTE35-IN or SCTE35-CMD
type SCTE35 struct {
	Type  string //Possible Values: IN, OUT, CMD
	Value string //big-endian binary representation of the splice_info_section(), expressed as a hexadecimal-sequence.
}

//MasterPlaylist represents a Master Playlist and its tags
type MasterPlaylist struct {
	M3U                 bool //Represents tag #EXTM3U. Indicates if present. MUST be present.
	Version             int  //Represents tag #EXT-X-VERSION. MUST be present.
	Variants            []*Variant
	SessionData         []*SessionData
	SessionKeys         []*Key
	IndependentSegments bool //Represents tag #EXT-X-INDEPENDENT-SEGMENTS. Applies to every Media Segment of every Media Playlist referenced. V6 or higher.
	StartPoint          *StartPoint
}

//Variant represents tag #EXT-X-STREAM-INF:<attribute-list> and tag #EXT-X-I-FRAME-STREAM-INF.
//
//EXT-X-STREAM-INF specifies a Variant Stream, which is a set of Renditions which can be combined to play the presentation.
//A URI line following the tag indicates the Media Playlist carrying a rendition of the Variant Stream and it MUST be present.
//
//#EXT-X-I-FRAME-STREAM-INF identifies Media Playlist file containing the I-frames of a multimedia presentation.
//It supports the same parameters as EXT-X-STREAM-INF except Audio, Subtitles and ClosedCaptions.
type Variant struct {
	Renditions     []*Rendition
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
	//If NONE, all EXT-X-STREAM-INF MUST have this attribute as NONE. If quoted-string, MUST match GroupID value of an EXT-X-MEDIA tag whose Type is CLOSED-CAPTIONS.
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
}

//SessionData represents tag #EXT-X-SESSION-DATA. Master Playlist MAY contain more than one tag with the same DataID but the Language MUST be different.
//V7
type SessionData struct {
	DataID   string //Required. SHOULD conform with a reverse DNS naming convention.
	Value    string //Required IF URI is not present. Contains the session data
	URI      string //Required IF Value is not present. Resource with the session data
	Language string //Optional. RFC5646 language tag that identifies the language of the data
}

//StartPoint represents tag #EXT-X-START. Indicates preferred point at which to start playing a Playlist.
type StartPoint struct {
	TimeOffset float64 //Required. If positive, time offset from the beginning of the Playlist. If negative, time offset from the end of the last segment of the playlist
	Precise    bool    //Possible Values: YES or NO.
}
