package dash

import "encoding/xml"

const (
	dashNS = "urn:mpeg:dash:schema:mpd:2011"
	cencNS = "urn:mpeg:cenc:2013"
	msprNS = "urn:microsoft:playready"
)

//MPD represents a Media Presentation Description.
type MPD struct {
	XMLNS                 string                `xml:"xmlns,attr,omitempty"`
	SchemaLocation        string                `xml:"http://www.w3.org/2001/XMLSchema-instance schemaLocation,attr,omitempty"`
	ID                    string                `xml:"id,attr,omitempty"`                         //Optional.
	Profiles              string                `xml:"profiles,attr,omitempty"`                   //Required
	Type                  string                `xml:"type,attr,omitempty"`                       //Optional. Default:"static". Possible Values: static, dynamic
	PublishTime           *CustomTime           `xml:"publishTime,attr,omitempty"`                //Must be present for type "dynamic".
	AvStartTime           *CustomTime           `xml:"availabilityStartTime,attr,omitempty"`      //Must be present for type "dynamic". In UTC
	AvEndTime             *CustomTime           `xml:"availabilityEndTime,attr,omitempty"`        //Optional
	MediaPresDuration     *CustomDuration       `xml:"mediaPresentationDuration,attr,omitempty"`  //Optional. Shall be present if MinUpdatePeriod and Period.Duration aren't set.
	MinUpdatePeriod       *CustomDuration       `xml:"minimumUpdatePeriod,attr,omitempty"`        //Optional. Must not be present for type "static". Specifies the frequency in which clients must check for updates.
	MinBufferTime         *CustomDuration       `xml:"minBufferTime,attr,omitempty"`              //Required.
	TimeShiftBuffer       *CustomDuration       `xml:"timeShiftBufferDepth,attr,omitempty"`       //Optional for type "dynamic". If type "static", value is undefined.
	SuggestedPresDelay    *CustomDuration       `xml:"suggestedPresentationDelay,attr,omitempty"` //Optional for type "dynamic". If type "static", value is undefined.
	MaxSegmentDuration    *CustomDuration       `xml:"maxSegmentDuration,attr,omitempty"`         //Optional.
	MaxSubsegmentDuration *CustomDuration       `xml:"maxSubsegmentDuration,attr,omitempty"`      //Optional.
	ProgramInformation    []*ProgramInformation `xml:"ProgramInformation,omitempty"`
	BaseURL               []*BaseURL            `xml:"BaseURL,omitempty"`
	Location              []string              `xml:"Location,omitempty"`
	Metrics               []*Metrics            `xml:"Metrics,omitempty"`
	Periods               Periods               `xml:"Period,omitempty"`
}

//ProgramInformation specifies descriptive information about the program
type ProgramInformation struct {
	Lang               string `xml:"lang,attr,omitempty"`               //Optional
	MoreInformationURL string `xml:"moreInformationURL,attr,omitempty"` //Optional
	Title              string `xml:"Title,omitempty"`                   //Optional
	Source             string `xml:"Source,omitempty"`                  //Optional
	Copyright          string `xml:"Copyright,omitempty"`               //Optional
}

//BaseURL can be used for reference resolution and alternative URL selection.
type BaseURL struct {
	URL             string  `xml:",innerxml"`
	ServiceLocation string  `xml:"serviceLocation,attr,omitempty"` //Optional
	ByteRange       string  `xml:"byteRange,attr,omitempty"`       //
	AvTimeOffset    float64 `xml:"availabilityTimeOffset,attr,omitempty"`
	AvTimeComplete  bool    `xml:"availabilityTimeComplete,attr,omitempty"`
}

//Metrics ...
type Metrics struct {
	Metrics   string        `xml:"metrics,attr,omitempty"` //Required
	Range     []*Range      `xml:"Range,omitempty"`        //Optional
	Reporting []*Descriptor `xml:"Reporting,omitempty"`    //Required
}

//Range ...
type Range struct {
	StartTime float64 `xml:"starttime,attr,omitempty"`
	Duration  float64 `xml:"duration,attr,omitempty"`
}

//Descriptor ...
type Descriptor struct {
	SchemeIDURI string `xml:"schemeIdUri,attr,omitempty"`
	Value       string `xml:"value,attr,omitempty"`
	ID          string `xml:"id,attr,omitempty"`
}

//Period represents a media content period.
type Period struct {
	XlinkHref          string           `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`    //Optional
	XlinkActuate       string           `xml:"http://www.w3.org/1999/xlink actuate,attr,omitempty"` //Optional. Possible Values: onDemand, onRequest
	ID                 string           `xml:"id,attr,omitempty"`                                   //Optional. Must be unique. If type "dynamic", id must be present and not updated.
	Start              *CustomDuration  `xml:"start,attr,omitempty"`                                //Optional. Used as anchor to determine the start of each Media Segment.
	Duration           *CustomDuration  `xml:"duration,attr,omitempty"`                             //Optional. Determine the Start time of next Period.
	BitstreamSwitching bool             `xml:"bitstreamSwitching,attr,omitempty"`                   //Optional. Default: false. If 'true', means that every AdaptationSet.BitstreamSwitching is set to 'true'. TODO: check if there's 'false' on AdaptationSet
	BaseURL            []*BaseURL       `xml:"BaseURL,omitempty"`                                   //Optional
	SegmentBase        *SegmentBase     `xml:"SegmentBase,omitempty"`                               //Optional. Default Segment Base information. Overidden by AdaptationSet.SegmentBase and Representation.SegmentBase
	SegmentList        *SegmentList     `xml:"SegmentList,omitempty"`
	SegmentTemplate    *SegmentTemplate `xml:"SegmentTemplate,omitempty"`
	AssetIdentifier    *Descriptor      `xml:"AssetIdentifier,omitempty"`
	EventStream        []*EventStream   `xml:"EventStream,omitempty"`
	AdaptationSets     `xml:"AdaptationSet,omitempty"`
	Subsets            `xml:"Subset,omitempty"`
}

//EventStream represents a sequence of related events.
type EventStream struct {
	XlinkHref    string   `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
	XlinkActuate string   `xml:"http://www.w3.org/1999/xlink actuate,attr,omitempty"`
	SchemeIDURI  string   `xml:"schemeIdUri,attr,omitempty"`
	Value        string   `xml:"value,attr,omitempty"`
	Timescale    int      `xml:"timescale,attr,omitempty"`
	Event        []*Event `xml:"Event,omitempty"`
}

//Event represents aperiodic sparse media-time related auxiliary information to DASH
//client or an application.
type Event struct {
	Message  string `xml:",innerxml"`
	PresTime int64  `xml:"presentationTime,attr,omitempty"`
	Duration int64  `xml:"duration,attr,omitempty"`
	ID       int    `xml:"id,attr,omitempty"`
}

//URLType ...
type URLType struct {
	SourceURL string `xml:"sourceURL,attr,omitempty"`
	Range     string `xml:"range,attr,omitempty"`
}

//SegmentBase represents a media file played by a DASH client
type SegmentBase struct {
	Timescale           int             `xml:"timescale,attr,omitempty"`                //Optional. If not present, it must be set to 1.
	PresTimeOffset      int64           `xml:"presentationTimeOffset,attr,omitempty"`   //Optional.
	TimeShiftBuffer     *CustomDuration `xml:"timeShiftBufferDepth,attr,omitempty"`     //Optional.
	IndexRange          string          `xml:"indexRange,attr,omitempty"`               //Optional. ByteRange that contains the Segment Index in all Segments of the Representation.
	IndexRangeExact     bool            `xml:"indexRangeExact,attr,omitempty"`          //Default: false. Must not be present if IndexRange isn't present.
	AvTimeOffset        float64         `xml:"availabilityTimeOffset,attr,omitempty"`   //Optional.
	AvTmeComplete       bool            `xml:"availabilityTimeComplete,attr,omitempty"` //Optional.
	Initialization      *URLType        `xml:"Initialization,omitempty"`
	RepresentationIndex *URLType        `xml:"RepresentationIndex,omitempty"`
}

//SegmentList contains a list of SegmentURL elements.
type SegmentList struct {
	XlinkHref           string           `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
	XlinkActuate        string           `xml:"http://www.w3.org/1999/xlink actuate,attr,omitempty"`
	Timescale           int              `xml:"timescale,attr,omitempty"`                //Optional. . If not present, it must be set to 1.
	PresTimeOffset      int64            `xml:"presentationTimeOffset,attr,omitempty"`   //Optional.
	TimeShiftBuffer     *CustomDuration  `xml:"timeShiftBufferDepth,attr,omitempty"`     //Optional.
	IndexRange          string           `xml:"indexRange,attr,omitempty"`               //Optional. ByteRange that contains the Segment Index in all Segments of the Representation.
	IndexRangeExact     bool             `xml:"indexRangeExact,attr,omitempty"`          //Default: false. Must not be present if IndexRange isn't present.
	AvTimeOffset        float64          `xml:"availabilityTimeOffset,attr,omitempty"`   //Optional.
	AvTmeComplete       bool             `xml:"availabilityTimeComplete,attr,omitempty"` //Optional.
	Duration            int              `xml:"duration,attr,omitempty"`
	StartNumber         int              `xml:"startNumber,attr,omitempty"`
	Initialization      *URLType         `xml:"Initialization,omitempty"`
	RepresentationIndex *URLType         `xml:"RepresentationIndex,omitempty"`
	SegmentTimeline     *SegmentTimeline `xml:"SegmentTimeline,omitempty"`
	BitstreamSwitching  *URLType         `xml:"BitstreamSwitching,omitempty"`
	SegmentURLs         []*SegmentURL    `xml:"SegmentURL,omitempty"`
}

//SegmentURL may contain the Media Segment URL.
type SegmentURL struct {
	Media      string `xml:"media,attr,omitempty"`      //Optional. Combined with MediaRange, specifies HTTP-URL for Media Segment. If not present, MediaRange must be present and it's combined with BaseURL
	MediaRange string `xml:"mediaRange,attr,omitempty"` //Optional. If not present, Media Segment is the entire resource in Media
	Index      string `xml:"index,attr,omitempty"`      //Optional.
	IndexRange string `xml:"indexRange,attr,omitempty"`
}

//SegmentTemplate specifies identifiers that are substituted by dynamic values assigned
//to Segments, to create a list of Segments.
type SegmentTemplate struct {
	Timescale              int              `xml:"timescale,attr,omitempty"`                //Optional. If not present, it must be set to 1.
	PresTimeOffset         int64            `xml:"presentationTimeOffset,attr,omitempty"`   //Optional.
	TimeShiftBuffer        *CustomDuration  `xml:"timeShiftBufferDepth,attr,omitempty"`     //Optional.
	IndexRange             string           `xml:"indexRange,attr,omitempty"`               //Optional. ByteRange that contains the Segment Index in all Segments of the Representation.
	IndexRangeExact        bool             `xml:"indexRangeExact,attr,omitempty"`          //Default: false. Must not be present if IndexRange isn't present.
	AvTimeOffset           float64          `xml:"availabilityTimeOffset,attr,omitempty"`   //Optional.
	AvTmeComplete          bool             `xml:"availabilityTimeComplete,attr,omitempty"` //Optional.
	Duration               int              `xml:"duration,attr,omitempty"`
	StartNumber            int              `xml:"startNumber,attr,omitempty"`
	Media                  string           `xml:"media,attr,omitempty"`              //Optional. Template to create Media Segment List
	Index                  string           `xml:"index,attr,omitempty"`              //Optional. Template to create the Index Segment List. If neither $Number% nor %Time% is included, it provides the URL to a Representation Index
	InitializationAttr     string           `xml:"initialization,attr,omitempty"`     //Optional. Template to create Initialization Segment. $Number% and %Time% must not be included.
	BitstreamSwitchingAttr string           `xml:"bitstreamSwitching,attr,omitempty"` //Optional. Template to create Bitstream Switching Segment. $Number% and %Time% must not be included.
	Initialization         *URLType         `xml:"Initialization,omitempty"`
	RepresentationIndex    *URLType         `xml:"RepresentationIndex,omitempty"`
	SegmentTimeline        *SegmentTimeline `xml:"SegmentTimeline,omitempty"`
	BitstreamSwitching     *URLType         `xml:"BitstreamSwitching,omitempty"`
}

//SegmentTimeline represents the earliest presentation time and duration for each Segment in the Representation.
//It contains a list of S elements, each describing a sequence of continguous segments of identical MPD duration.
//The order of the S elements must match the numbering order (time) of the corresponding Media Segments.
type SegmentTimeline struct {
	Segments Segments `xml:"S"` //Must have at least 1 S element.
}

//S is contained in a SegmentTimeline tag.
type S struct {
	T int `xml:"t,attr"` //Optional. Specifies MPD start time, in timescale units. Relative to the befinning of the Period.
	D int `xml:"d,attr"` //Required. Segment duration int timescale units. Must not exceed the value of MPD.MaxSegmentDuration.
	R int `xml:"r,attr"` //Default: 0. Specifies repeat count of number of following continguous segments with same duration as D.
}

//Subset restricts the combination of active AdaptationSets where an active
//Adaptation Set is one for which the DASH client is presenting at least one of the
//contained Representation. No subset should contain all the Adaptaion Sets.
type Subset struct {
	Contains CustomInt `xml:"contains,attr"` //Required. Whitespace separated list.
	ID       string    `xml:"id,attr,omitempty"`
}

//AdaptationSet represents a set of versions of one or more media streams.
type AdaptationSet struct {
	XlinkHref               string                 `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
	XlinkActuate            string                 `xml:"http://www.w3.org/1999/xlink actuate,attr,omitempty"` //Possible Values: 'onLoad', 'onRequest'. Default: onRequest.
	ID                      int                    `xml:"id,attr,omitempty"`
	Group                   int                    `xml:"group,attr,omitempty"`
	Lang                    string                 `xml:"lang,attr,omitempty"`
	ContentType             string                 `xml:"contentType,attr,omitempty"`
	Par                     string                 `xml:"par,attr,omitempty"` //Optional. Picture Aspect Ratio. TODO:check specs for validation (regex)
	MinBandwidth            int                    `xml:"minBandwith,attr,omitempty"`
	MaxBandwidth            int                    `xml:"maxBandwidth,attr,omitempty"`
	MinWidth                int                    `xml:"minWidth,attr,omitempty"`
	MaxWidth                int                    `xml:"maxWidth,attr,omitempty"`
	MinHeight               int                    `xml:"minHeight,attr,omitempty"`
	MaxHeight               int                    `xml:"maxHeight,attr,omitempty"`
	MinFrameRate            string                 `xml:"minFrameRate,attr,omitempty"` //TODO:Check specs for validation (regex)
	MaxFrameRate            string                 `xml:"maxFrameRate,attr,omitempty"`
	SegmentAlignment        bool                   `xml:"segmentAlignment,attr,omitempty"`        //Default: false. TODO: check specs for validation. Accepts 0,1 or true,false
	BitstreamSwitching      bool                   `xml:"bitstreamSwitching,attr,omitempty"`      //TODO: check specs for validation. Accepts 0,1 or true,false
	SubsegmentAlignment     bool                   `xml:"subsegmentAlignment,attr,omitempty"`     //Default: false. TODO: check specs for validation
	SubsegmentStartsWithSAP int                    `xml:"subsegmentStartsWithSap,attr,omitempty"` //Default: 0. TODO: check specs for validation
	Profiles                string                 `xml:"profiles,attr,omitempty"`
	Width                   int                    `xml:"width,attr,omitempty"`
	Height                  int                    `xml:"height,attr,omitempty"`
	Sar                     string                 `xml:"sar,attr,omitempty"`       //RatioType
	FrameRate               string                 `xml:"frameRate,attr,omitempty"` //FrameRateType
	AudioSamplingRate       string                 `xml:"audioSamplingRate,attr,omitempty"`
	MimeType                string                 `xml:"mimeType,attr,omitempty"`
	SegmentProfiles         string                 `xml:"segmentProfiles,attr,omitempty"`
	Codecs                  string                 `xml:"codecs,attr,omitempty"`
	MaxSAPPeriod            float64                `xml:"maximumSAPPeriod,attr,omitempty"` //seconds
	StartWithSAP            int                    `xml:"startWithSAP,attr,omitempty"`     //SAPType
	MaxPlayoutRate          float64                `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency        bool                   `xml:"codingDependency,attr,omitempty"`
	ScanType                string                 `xml:"scanType,attr,omitempty"` //VideoScanType
	CENCContentProtections  CENCContentProtections `xml:"ContentProtection,omitempty"`
	FramePacking            []*Descriptor          `xml:"FramePacking,omitempty"`
	AudioChannelConfig      []*Descriptor          `xml:"AudioChannelConfiguration,omitempty"`
	EssentialProperty       []*Descriptor          `xml:"EssentialProperty,omitempty"`
	SupplementalProperty    []*Descriptor          `xml:"SupplementalProperty,omitempty"`
	InbandEventStream       []*Descriptor          `xml:"InbandEventStream,omitempty"`
	Accessibility           []*Descriptor          `xml:"Accessibility,omitempty"`
	Role                    []*Descriptor          `xml:"Role,omitempty"`
	Rating                  []*Descriptor          `xml:"Rating,omitempty"`
	ViewPoint               []*Descriptor          `xml:"Viewpoint,omitempty"`
	ContentComponent        []*ContentComponent    `xml:"ContentComponent,omitempty"`
	BaseURL                 []*BaseURL             `xml:"BaseURL,omitempty"`
	SegmentBase             *SegmentBase           `xml:"SegmentBase,omitempty"`
	SegmentList             *SegmentList           `xml:"SegmentList,omitempty"`
	SegmentTemplate         *SegmentTemplate       `xml:"SegmentTemplate,omitempty"`
	Representations         Representations        `xml:"Representation,omitempty"`
}

//ContentProtection represents the root ContentProtection element.
type ContentProtection struct {
	XMLName     xml.Name `xml:"ContentProtection"`
	XMLNsCenc   string   `xml:"xmlns:cenc,attr,omitempty"`
	XMLNsMspr   string   `xml:"xmlns:mspr,attr,omitempty"`
	SchemeIDURI string   `xml:"schemeIdUri,attr,omitempty"`
	Value       string   `xml:"value,attr,omitempty"`
	DefaultKID  string   `xml:"cenc:default_KID,attr,omitempty"`
}

//CENCContentProtection represents the full ContentProtection element.
//
//Note for Playready encryption: the elements defined in the “mspr” namespace for the
//first edition of Common Encryption (mspr:IsEncrypted, mspr:IV_size, and mspr:kid),
//are deprecated and functionally replaced by cenc:default_KID specified in the second
//edition of Common Encryption [CENC].
//The IV_size and IsEncrypted fields in the Track Encryption Box (‘tenc’) are used during decryption,
//but are not needed in MPD ContentProtection Descriptor elements.
type CENCContentProtection struct {
	ContentProtection
	Pssh        *Pssh
	Pro         *Pro
	IsEncrypted string `xml:"mspr:IsEncrypted,omitempty"`
	IVSize      int    `xml:"mspr:IV_size,omitempty"`
	KID         string `xml:"mspr:kid,omitempty"`
}

//Pssh (Protection System Specific Header) represents the optional cenc:pssh element
//that can be used by all DRM ContentProtection Descriptors for improved interoperability.
type Pssh struct {
	XMLName xml.Name `xml:"cenc:pssh,omitempty"`
	Value   string   `xml:",innerxml"`
}

//Pro represents a Playready Header Object.
type Pro struct {
	XMLName xml.Name `xml:"mspr:pro,omitempty"`
	Value   string   `xml:",innerxml"`
}

//ContentComponent describes the properties of each media content component in an
//Adaptation Set. If only one media content component is present, it can be described directly
//in the Adaptation Set.
type ContentComponent struct {
	ID            *int          `xml:"id,attr,omitempty"`
	Lang          string        `xml:"lang,attr,omitempty"`
	ContentType   string        `xml:"contentType,attr,omitempty"`
	Par           string        `xml:"par,attr,omitempty"`
	Accessibility []*Descriptor `xml:"Accessibility,omitempty"`
	Role          []*Descriptor `xml:"Role,omitempty"`
	Rating        []*Descriptor `xml:"Rating,omitempty"`
	ViewPoint     []*Descriptor `xml:"Viewpoint,omitempty"`
}

//Representation represents a deliverable encoded version of one or more media components.
type Representation struct {
	ID                      string                 `xml:"id,attr"`        //Required. It must not contain whitespace characters.
	Bandwidth               int64                  `xml:"bandwidth,attr"` //Required.
	QualityRanking          int                    `xml:"qualityRanking,attr,omitempty"`
	DependencyID            string                 `xml:"dependencyId,attr,omitempty"`            //Whitespace separated list of int
	MediaStreamsStructureID string                 `xml:"mediaStreamsStructureId,attr,omitempty"` //Whitespace separated list of int
	Profiles                string                 `xml:"profiles,attr,omitempty"`
	Width                   int                    `xml:"width,attr,omitempty"`
	Height                  int                    `xml:"height,attr,omitempty"`
	Sar                     string                 `xml:"sar,attr,omitempty"`       //RatioType
	FrameRate               string                 `xml:"frameRate,attr,omitempty"` //FrameRateType
	AudioSamplingRate       string                 `xml:"audioSamplingRate,attr,omitempty"`
	MimeType                string                 `xml:"mimeType,attr,omitempty"`
	SegmentProfiles         string                 `xml:"segmentProfiles,attr,omitempty"`
	Codecs                  string                 `xml:"codecs,attr,omitempty"`
	MaxSAPPeriod            float64                `xml:"maximumSAPPeriod,attr,omitempty"`
	StartWithSAP            int                    `xml:"startWithSAP,attr,omitempty"` //SAPType
	MaxPlayoutRate          float64                `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency        bool                   `xml:"codingDependency,attr,omitempty"`
	ScanType                string                 `xml:"scanType,attr,omitempty"` //VideoScanType
	CENCContentProtections  CENCContentProtections `xml:"ContentProtection,omitempty"`
	FramePacking            []*Descriptor          `xml:"FramePacking,omitempty"`
	AudioChannelConfig      []*Descriptor          `xml:"AudioChannelConfiguration,omitempty"`
	EssentialProperty       []*Descriptor          `xml:"EssentialProperty,omitempty"`
	SupplementalProperty    []*Descriptor          `xml:"SupplementalProperty,omitempty"`
	InbandEventStream       []*Descriptor          `xml:"InbandEventStream,omitempty"`
	BaseURL                 []*BaseURL             `xml:"BaseURL,omitempty"`
	SubRepresentation       []*SubRepresentation   `xml:"SubRepresentation,omitempty"`
	SegmentBase             *SegmentBase           `xml:"SegmentBase,omitempty"`
	SegmentList             *SegmentList           `xml:"SegmentList,omitempty"`
	SegmentTemplate         *SegmentTemplate       `xml:"SegmentTemplate,omitempty"`
}

//SubRepresentation describes properties of one or several media content components
//that are embedded in the Representation.
type SubRepresentation struct {
	Level                  *int                   `xml:"level,attr,omitempty"`
	DependencyLevel        CustomInt              `xml:"dependencyLevel,attr,omitempty"` //Whitespace separated list of int
	Bandwidth              int                    `xml:"bandwidth,attr,omitempty"`
	ContentComponent       string                 `xml:"contentComponent,attr,omitempty"` //Whitespace separated list of string
	Profiles               string                 `xml:"profiles,attr,omitempty"`
	Width                  int                    `xml:"width,attr,omitempty"`
	Height                 int                    `xml:"height,attr,omitempty"`
	Sar                    string                 `xml:"sar,attr,omitempty"`       //RatioType
	FrameRate              string                 `xml:"frameRate,attr,omitempty"` //FrameRateType
	AudioSamplingRate      string                 `xml:"audioSamplingRate,attr,omitempty"`
	MimeType               string                 `xml:"mimeType,attr,omitempty"`
	SegmentProfiles        string                 `xml:"segmentProfiles,attr,omitempty"`
	Codecs                 string                 `xml:"codecs,attr,omitempty"`
	MaxSAPPeriod           float64                `xml:"maximumSAPPeriod,attr,omitempty"`
	StartWithSAP           int                    `xml:"startWithSAP,attr,omitempty"` //SAPType
	MaxPlayoutRate         float64                `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency       bool                   `xml:"codingDependency,attr,omitempty"`
	ScanType               string                 `xml:"scanType,attr,omitempty"` //VideoScanType
	CENCContentProtections CENCContentProtections `xml:"ContentProtection,omitempty"`
	FramePacking           []*Descriptor          `xml:"FramePacking,omitempty"`
	AudioChannelConfig     []*Descriptor          `xml:"AudioChannelConfiguration,omitempty"`
	EssentialProperty      []*Descriptor          `xml:"EssentialProperty,omitempty"`
	SupplementalProperty   []*Descriptor          `xml:"SupplementalProperty,omitempty"`
	InbandEventStream      []*Descriptor          `xml:"InbandEventStream,omitempty"`
}
