package dash

import "encoding/xml"

//TODO:Go through every duration field, make sure it's seconds or h:m:s
//Go through every int field, make sure int is enough or int64 better
//Make sure int field of value 0 will appear on playlist when marshalling to xml
//Look up if order of certain elements matter
//Figure out how to separate common attr into generic structs. There's a lot of duplicated fields at the moment.

//DashNS is the XML schema for the MPD
const DashNS = "urn:mpeg:dash:schema:mpd:2011"

//MPD represents a Media Presentation Description.
//TODO:xsi namespace prefix unmarshals successfully, but errors on marshalling:
//It renders: xmlns:XMLSchema-instance="http://www.w3.org/2001/XMLSchema-instance"
//XMLSchema-instance:schemaLocation="urn:mpeg:dash:schema:mpd:2011 DASH-MPD.xsd"
type MPD struct {
	XMLNS                 string                `xml:"xmlns,attr,omitempty"`
	Xsi                   string                `xml:"http://www.w3.org/2001/XMLSchema-instance xsi,attr,omitempty"`
	SchemaLocation        string                `xml:"xsi schemaLocation,attr,omitempty"`
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

//Metrics represents DASH Metrics.
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
	Start              *CustomDuration  `xml:"start,attr,omitempty"`                                //Optional. Used as anchor to determine the start of each Media Segment. TODO:Check when not present
	Duration           *CustomDuration  `xml:"duration,attr,omitempty"`                             //Optional. Determine the Start time of next Period. TODO:check when not present
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

//EventStream ...
type EventStream struct { //TODO:check specs for validation
	XlinkHref    string   `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
	XlinkActuate string   `xml:"http://www.w3.org/1999/xlink actuate,attr,omitempty"`
	SchemeIdURI  string   `xml:"schemeIdUri,attr,omitempty"`
	Value        string   `xml:"value,attr,omitempty"`
	Timescale    int      `xml:"timescale,attr,omitempty"`
	Event        []*Event `xml:"Event,omitempty"`
}

//Event ...
type Event struct {
	PresTime int64 `xml:"presentationTime,attr,omitempty"`
	Duration int64 `xml:"duration,attr,omitempty"`
	ID       int   `xml:"id,attr,omitempty"`
}

//URLType ...
type URLType struct {
	SourceURL string `xml:"sourceURL,attr,omitempty"`
	Range     string `xml:"range,attr,omitempty"`
}

//SegmentBase represents a media file played by a DASH client
type SegmentBase struct {
	XMLName             xml.Name
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

//MultipleSegmentBase ...
type MultipleSegmentBase struct {
	SegmentBase        *SegmentBase
	Duration           int              `xml:"duration,attr,omitempty"`
	StartNumber        int              `xml:"startNumber,attr,omitempty"`
	SegmentTimeline    *SegmentTimeline `xml:"SegmentTimeline,omitempty"`
	BitstreamSwitching *URLType         `xml:"BitstreamSwitching,omitempty"`
}

//SegmentList ... SegmentBase + MultipleSegmentBase
type SegmentList struct {
	Timescale           int              `xml:"timescale,attr,omitempty"`                //Optional. . If not present, it must be set to 1.
	PresTimeOffset      int64            `xml:"presentationTimeOffset,attr,omitempty"`   //Optional.
	TimeShiftBuffer     *CustomDuration  `xml:"timeShiftBufferDepth,attr,omitempty"`     //Optional.
	IndexRange          string           `xml:"indexRange,attr,omitempty"`               //Optional. ByteRange that contains the Segment Index in all Segments of the Representation.
	IndexRangeExact     bool             `xml:"indexRangeExact,attr,omitempty"`          //Default: false. Must not be present if IndexRange isn't present.
	AvTimeOffset        float64          `xml:"availabilityTimeOffset,attr,omitempty"`   //Optional.
	AvTmeComplete       bool             `xml:"availabilityTimeComplete,attr,omitempty"` //Optional.
	Initialization      *URLType         `xml:"Initialization,omitempty"`
	RepresentationIndex *URLType         `xml:"RepresentationIndex,omitempty"`
	Duration            int              `xml:"duration,attr,omitempty"`
	StartNumber         int              `xml:"startNumber,attr,omitempty"`
	SegmentTimeline     *SegmentTimeline `xml:"SegmentTimeline,omitempty"`
	BitstreamSwitching  *URLType         `xml:"BitstreamSwitching,omitempty"`
	XlinkHref           string           `xml:"xlink:href,attr,omitempty"`
	XlinkActuate        string           `xml:"xlink:actuate,attr,omitempty"`
	SegmentURLs         []*SegmentURL    `xml:"SegmentURL,omitempty"`
}

//SegmentURL ...
type SegmentURL struct {
	Media      string `xml:"media,attr,omitempty"`      //Optional. Combined with MediaRange, specifies HTTP-URL for Media Segment. If not present, MediaRange must be present and it's combined with BaseURL
	MediaRange string `xml:"mediaRange,attr,omitempty"` //Optional. If not present, Media Segment is the entire resource in Media
	Index      string `xml:"index,attr,omitempty"`      //Optional.
	IndexRange string `xml:"indexRange,attr,omitempty"`
}

//SegmentTemplate ... SegmentBase + MultipleSegmentBase
type SegmentTemplate struct {
	Timescale              int              `xml:"timescale,attr,omitempty"`                //Optional. . If not present, it must be set to 1.
	PresTimeOffset         int64            `xml:"presentationTimeOffset,attr,omitempty"`   //Optional.
	TimeShiftBuffer        *CustomDuration  `xml:"timeShiftBufferDepth,attr,omitempty"`     //Optional.
	IndexRange             string           `xml:"indexRange,attr,omitempty"`               //Optional. ByteRange that contains the Segment Index in all Segments of the Representation.
	IndexRangeExact        bool             `xml:"indexRangeExact,attr,omitempty"`          //Default: false. Must not be present if IndexRange isn't present.
	AvTimeOffset           float64          `xml:"availabilityTimeOffset,attr,omitempty"`   //Optional.
	AvTmeComplete          bool             `xml:"availabilityTimeComplete,attr,omitempty"` //Optional.
	Initialization         *URLType         `xml:"Initialization,omitempty"`
	RepresentationIndex    *URLType         `xml:"RepresentationIndex,omitempty"`
	Duration               int              `xml:"duration,attr,omitempty"`
	StartNumber            int              `xml:"startNumber,attr,omitempty"`
	SegmentTimeline        *SegmentTimeline `xml:"SegmentTimeline,omitempty"`
	BitstreamSwitching     *URLType         `xml:"BitstreamSwitching,omitempty"`
	Media                  string           `xml:"media,attr,omitempty"`              //Optional. Template to create Media Segment List
	Index                  string           `xml:"index,attr,omitempty"`              //Optional. Template to create the Index Segment List. If neither $Number% nor %Time% is included, it provides the URL to a Representation Index
	InitializationAttr     string           `xml:"initialization,attr,omitempty"`     //Optional. Template to create Initialization Segment. $Number% and %Time% must not be included.
	BitstreamSwitchingAttr string           `xml:"bitstreamSwitching,attr,omitempty"` //Optional. Template to create Bitstream Switching Segment. $Number% and %Time% must not be included.
}

//SegmentTimeline represents the earliest presentation time and duration for each Segment in the Representation.
//It contains a list of S elements, each describing a sequence of continguous segments of identical MPD duration.
//The order of the S elements must match the numbering order (time) of the corresponding Media Segments.
type SegmentTimeline struct {
	Segments Segments `xml:"S"` //Must have at least 1 S element.
}

//S is contained in a SegmentTimeline tag. TODO:Check specs for validation.
type S struct {
	T int `xml:"t,attr"` //Optional. Specifies MPD start time, in timescale units. Relative to the befinning of the Period.
	D int `xml:"d,attr"` //Required. Segment duration int timescale units. Must not exceed the value of MPD.MaxSegmentDuration.
	R int `xml:"r,attr"` //Default: 0. Specifies repeat count of number of following continguous segments with same duration as D.
}

//Subset ..
type Subset struct {
	Contains CustomInt `xml:"contains,attr"` //Required. Whitespace separated list.
	ID       string    `xml:"id,attr,omitempty"`
}

//AdaptationSet represents a set of versions of one or more media streams.
type AdaptationSet struct {
	XlinkHref               string              `xml:"http://www.w3.org/1999/xlink href,attr,omitempty"`
	XlinkActuate            string              `xml:"http://www.w3.org/1999/xlink actuate,attr,omitempty"` //Possible Values: 'onLoad', 'onRequest'. Default: onRequest.
	ID                      int                 `xml:"id,attr,omitempty"`
	Group                   int                 `xml:"group,attr,omitempty"`
	Lang                    string              `xml:"lang,attr,omitempty"`
	ContentType             string              `xml:"contentType,attr,omitempty"`
	Par                     string              `xml:"par,attr,omitempty"` //Optional. Picture Aspect Ratio. TODO:check specs for validation (regex)
	MinBandwidth            int                 `xml:"minBandwith,attr,omitempty"`
	MaxBandwidth            int                 `xml:"maxBandwidth,attr,omitempty"`
	MinWidth                int                 `xml:"minWidth,attr,omitempty"`
	MaxWidth                int                 `xml:"maxWidth,attr,omitempty"`
	MinHeight               int                 `xml:"minHeight,attr,omitempty"`
	MaxHeight               int                 `xml:"maxHeight,attr,omitempty"`
	MinFrameRate            string              `xml:"minFrameRate,attr,omitempty"` //TODO:Check specs for validation (regex)
	MaxFrameRate            string              `xml:"maxFrameRate,attr,omitempty"`
	SegmentAlignment        bool                `xml:"segmentAlignment,attr,omitempty"`        //Default: false. TODO: check specs for validation. Accepts 0,1 or true,false
	BitstreamSwitching      bool                `xml:"bitstreamSwitching,attr,omitempty"`      //TODO: check specs for validation. Accepts 0,1 or true,false
	SubsegmentAlignment     bool                `xml:"subsegmentAlignment,attr,omitempty"`     //Default: false. TODO: check specs for validation
	SubsegmentStartsWithSAP int                 `xml:"subsegmentStartsWithSap,attr,omitempty"` //Default: 0. TODO: check specs for validation
	Accessibility           []*Descriptor       `xml:"Accessibility,omitempty"`
	Role                    []*Descriptor       `xml:"Role,omitempty"`
	Rating                  []*Descriptor       `xml:"Rating,omitempty"`
	ViewPoint               []*Descriptor       `xml:"Viewpoint,omitempty"`
	ContentComponent        []*ContentComponent `xml:"ContentComponent,omitempty"`
	BaseURL                 []*BaseURL          `xml:"BaseURL,omitempty"`
	SegmentBase             *SegmentBase        `xml:"SegmentBase,omitempty"`
	SegmentList             *SegmentList        `xml:"SegmentList,omitempty"`
	SegmentTemplate         *SegmentTemplate    `xml:"SegmentTemplate,omitempty"`
	Representations         Representations     `xml:"Representation,omitempty"`
	Profiles                string              `xml:"profiles,attr,omitempty"`
	Width                   int                 `xml:"width,attr,omitempty"`
	Height                  int                 `xml:"height,attr,omitempty"`
	Sar                     string              `xml:"sar,attr,omitempty"`       //RatioType
	FrameRate               string              `xml:"frameRate,attr,omitempty"` //FrameRateType
	AudioSamplingRate       string              `xml:"audioSamplingRate,attr,omitempty"`
	MimeType                string              `xml:"mimeType,attr,omitempty"`
	SegmentProfiles         string              `xml:"segmentProfiles,attr,omitempty"`
	Codecs                  string              `xml:"codecs,attr,omitempty"`
	MaxSAPPeriod            float64             `xml:"maximumSAPPeriod,attr,omitempty"` //seconds
	StartWithSAP            int                 `xml:"startWithSAP,attr,omitempty"`     //SAPType
	MaxPlayoutRate          float64             `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency        bool                `xml:"codingDependency,attr,omitempty"`
	ScanType                string              `xml:"scanType,attr,omitempty"` //VideoScanType
	FramePacking            []*Descriptor       `xml:"FramePacking,omitempty"`
	AudioChannelConfig      []*Descriptor       `xml:"AudioChannelConfiguration,omitempty"`
	ContentProtection       []*Descriptor       `xml:"ContentProtection,omitempty"`
	EssentialProperty       []*Descriptor       `xml:"EssentialProperty,omitempty"`
	SupplementalProperty    []*Descriptor       `xml:"SupplementalProperty,omitempty"`
	InbandEventStream       []*Descriptor       `xml:"InbandEventStream,omitempty"`
	//CommonComponents        *CommonComponents
}

//ContentComponent ...
type ContentComponent struct {
	ID            int           `xml:"id,attr,omitempty"`
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
	ID                      string               `xml:"id,attr"` //Required. TODO:Check validation (regex)
	Bandwidth               int                  `xml:"bandwidth,attr,omitempty"`
	QualityRanking          int                  `xml:"qualityRanking,attr,omitempty"`
	DependencyID            string               `xml:"dependencyId,attr,omitempty"`            //Whitespace separated list of int
	MediaStreamsStructureID string               `xml:"mediaStreamsStructureId,attr,omitempty"` //Whitespace separated list of int
	BaseURL                 []*BaseURL           `xml:"BaseURL,omitempty"`
	SubRepresentation       []*SubRepresentation `xml:"SubRepresentation,omitempty"`
	SegmentBase             *SegmentBase         `xml:"SegmentBase,omitempty"`
	SegmentList             *SegmentList         `xml:"SegmentList,omitempty"`
	SegmentTemplate         *SegmentTemplate     `xml:"SegmentTemplate,omitempty"`
	Profiles                string               `xml:"profiles,attr,omitempty"`
	Width                   int                  `xml:"width,attr,omitempty"`
	Height                  int                  `xml:"height,attr,omitempty"`
	Sar                     string               `xml:"sar,attr,omitempty"`       //RatioType
	FrameRate               string               `xml:"frameRate,attr,omitempty"` //FrameRateType
	AudioSamplingRate       string               `xml:"audioSamplingRate,attr,omitempty"`
	MimeType                string               `xml:"mimeType,attr,omitempty"`
	SegmentProfiles         string               `xml:"segmentProfiles,attr,omitempty"`
	Codecs                  string               `xml:"codecs,attr,omitempty"`
	MaxSAPPeriod            float64              `xml:"maximumSAPPeriod,attr,omitempty"`
	StartWithSAP            int                  `xml:"startWithSAP,attr,omitempty"` //SAPType
	MaxPlayoutRate          float64              `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency        bool                 `xml:"codingDependency,attr,omitempty"`
	ScanType                string               `xml:"scanType,attr,omitempty"` //VideoScanType
	FramePacking            []*Descriptor        `xml:"FramePacking,omitempty"`
	AudioChannelConfig      []*Descriptor        `xml:"AudioChannelConfiguration,omitempty"`
	ContentProtection       []*Descriptor        `xml:"ContentProtection,omitempty"`
	EssentialProperty       []*Descriptor        `xml:"EssentialProperty,omitempty"`
	SupplementalProperty    []*Descriptor        `xml:"SupplementalProperty,omitempty"`
	InbandEventStream       []*Descriptor        `xml:"InbandEventStream,omitempty"`
	//CommonComponents *CommonComponents `xml:"Representation"`
}

//SubRepresentation represents SubRepresentation elements. Describes properties of one or several media
//content components that are embedded in the Representation. TODO: check specs for validation. check validation
//for common attributes with Representation
type SubRepresentation struct {
	Level                int           `xml:"level,attr,omitempty"`
	DependencyLevel      CustomInt     `xml:"dependencyLevel,attr,omitempty"` //Whitespace separated list of int
	Bandwidth            int           `xml:"bandwidth,attr,omitempty"`
	ContentComponent     []string      `xml:"contentComponent,attr,omitempty"` //Whitespace separated list of string
	Profiles             string        `xml:"profiles,attr,omitempty"`
	Width                int           `xml:"width,attr,omitempty"`
	Height               int           `xml:"height,attr,omitempty"`
	Sar                  string        `xml:"sar,attr,omitempty"`       //RatioType
	FrameRate            string        `xml:"frameRate,attr,omitempty"` //FrameRateType
	AudioSamplingRate    string        `xml:"audioSamplingRate,attr,omitempty"`
	MimeType             string        `xml:"mimeType,attr,omitempty"`
	SegmentProfiles      string        `xml:"segmentProfiles,attr,omitempty"`
	Codecs               string        `xml:"codecs,attr,omitempty"`
	MaxSAPPeriod         float64       `xml:"maximumSAPPeriod,attr,omitempty"`
	StartWithSAP         int           `xml:"startWithSAP,attr,omitempty"` //SAPType
	MaxPlayoutRate       float64       `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency     bool          `xml:"codingDependency,attr,omitempty"`
	ScanType             string        `xml:"scanType,attr,omitempty"` //VideoScanType
	FramePacking         []*Descriptor `xml:"FramePacking,omitempty"`
	AudioChannelConfig   []*Descriptor `xml:"AudioChannelConfiguration,omitempty"`
	ContentProtection    []*Descriptor `xml:"ContentProtection,omitempty"`
	EssentialProperty    []*Descriptor `xml:"EssentialProperty,omitempty"`
	SupplementalProperty []*Descriptor `xml:"SupplementalProperty,omitempty"`
	InbandEventStream    []*Descriptor `xml:"InbandEventStream,omitempty"`
	//	CommonComponents *CommonComponents
}

//CommonComponents are attributes and elements present in AdaptationSet, Representation and SubRepresentation elements
type CommonComponents struct {
	Profiles             string        `xml:"profiles,attr,omitempty"`
	Width                int           `xml:"width,attr,omitempty"`
	Height               int           `xml:"height,attr,omitempty"`
	Sar                  string        `xml:"sar,attr,omitempty"`       //RatioType
	FrameRate            string        `xml:"frameRate,attr,omitempty"` //FrameRateType
	AudioSamplingRate    string        `xml:"audioSamplingRate,attr,omitempty"`
	MimeType             string        `xml:"mimeType,attr,omitempty"`
	SegmentProfiles      string        `xml:"segmentProfiles,attr,omitempty"`
	Codecs               string        `xml:"codecs,attr,omitempty"`
	MaxSAPPeriod         float64       `xml:"maximumSAPPeriod,attr,omitempty"`
	StartWithSAP         int           `xml:"startWithSAP,attr,omitempty"` //SAPType
	MaxPlayoutRate       float64       `xml:"maxPlayoutRate,attr,omitempty"`
	CodingDependency     bool          `xml:"codingDependency,attr,omitempty"`
	ScanType             string        `xml:"scanType,attr,omitempty"` //VideoScanType
	FramePacking         []*Descriptor `xml:"FramePacking,omitempty"`
	AudioChannelConfig   []*Descriptor `xml:"AudioChannelConfiguration,omitempty"`
	ContentProtection    []*Descriptor `xml:"ContentProtection,omitempty"`
	EssentialProperty    []*Descriptor `xml:"EssentialProperty,omitempty"`
	SupplementalProperty []*Descriptor `xml:"SupplementalProperty,omitempty"`
	InbandEventStream    []*Descriptor `xml:"InbandEventStream,omitempty"`
}
