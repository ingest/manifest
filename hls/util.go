package hls

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//BufWriter is a wrapper type for bytes.Buffer
type BufWriter struct {
	buf *bytes.Buffer
	err error
}

//NewBufWriter returns an instance of BufWriter
func NewBufWriter() *BufWriter {
	return &BufWriter{
		buf: new(bytes.Buffer),
		err: nil,
	}
}

//NewMediaPlaylist returns an instance of a MediaPlaylist with a set version
func NewMediaPlaylist(version int) *MediaPlaylist {
	return &MediaPlaylist{Version: version}
}

//NewMasterPlaylist returns an instance of a MasterPlaylist with a set version
func NewMasterPlaylist(version int) *MasterPlaylist {
	return &MasterPlaylist{Version: version}
}

func backwardsCompatibilityError(version int, tag string) error {
	return fmt.Errorf("Backwards compatibility error on tag %s with version %d", tag, version)
}

func attributeNotSetError(tag string, attribute string) error {
	return fmt.Errorf("%s attribute %s must be set", tag, attribute)
}

//WriteValidString checks if data is set and writes on buffer using BufWriter wrapper
func (b *BufWriter) WriteValidString(data interface{}, write string) bool {
	if b.err != nil {
		return false
	}
	switch data.(type) {
	case string:
		if data.(string) != "" {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	case float64:
		if data.(float64) > float64(0) {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	case int64:
		if data.(int64) > 0 {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	case int:
		if data.(int) > 0 {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	case bool:
		if data.(bool) {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	}

	return false
}

//WriteString wraps buffer.WriteString
func (b *BufWriter) WriteString(s string) {
	if b.err != nil {
		return
	}
	_, b.err = b.buf.WriteString(s)
}

//WriteRune wraps buffer.WriteRune
func (b *BufWriter) WriteRune(r rune) {
	if b.err != nil {
		return
	}
	_, b.err = b.buf.WriteRune(r)
}

//writeHeader sets the initial tags for both Media and Master Playlists files
func writeHeader(version int, buf *BufWriter) error {
	if version > 0 {
		buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
		buf.WriteString(strconv.Itoa(version))
		buf.WriteRune('\n')
		return buf.err
	}
	return attributeNotSetError("Playlist", "Version")
}

//writeIndependentSegment sets the #EXT-X-INDEPENDENT-SEGMENTS tag on Media and Master Playlist file
func writeIndependentSegment(isIndSeg bool, buf *BufWriter) {
	if isIndSeg {
		buf.WriteString("#EXT-X-INDEPENDENT-SEGMENTS\n")
	}
}

//writeStartPoint sets the #EXT-X-START tag on Media and Master Playlist file
func writeStartPoint(sp *StartPoint, buf *BufWriter) error {
	if sp != nil {
		if !buf.WriteValidString(sp.TimeOffset, fmt.Sprintf("#EXT-X-START:TIME-OFFSET=%s", strconv.FormatFloat(sp.TimeOffset, 'f', 3, 32))) {
			return attributeNotSetError("EXT-X-START", "TIME-OFFSET")
		}
		buf.WriteValidString(sp.Precise, ",PRECISE=YES")
		buf.WriteRune('\n')
	}
	return buf.err
}

//writeSessionData sets the EXT-X-SESSION-DATA tag on Master Playlist file
func (s *SessionData) writeSessionData(buf *BufWriter) error {
	if s != nil {
		if !buf.WriteValidString(s.DataID, fmt.Sprintf("#EXT-X-SESSION-DATA:DATA-ID=\"%s\"", s.DataID)) {
			return attributeNotSetError("EXT-X-SESSION-DATA", "DATA-ID")
		}

		if s.Value != "" && s.URI != "" {
			return errors.New("EXT-X-SESSION-DATA must have attributes URI or VALUE, not both")
		} else if s.Value != "" || s.URI != "" {
			buf.WriteValidString(s.Value, fmt.Sprintf(",VALUE=\"%s\"", s.Value))
			buf.WriteValidString(s.URI, fmt.Sprintf(",URI=\"%s\"", s.URI))
		} else {
			return errors.New("EXT-X-SESSION-DATA must have either URI or VALUE attributes set")
		}

		buf.WriteValidString(s.Language, fmt.Sprintf(",LANGUAGE=\"%s\"", s.Language))
		buf.WriteRune('\n')
	}
	return buf.err
}

//writeXMedia sets the EXT-X-MEDIA tag on Master Playlist file
func (r *Rendition) writeXMedia(buf *BufWriter) error {
	if r != nil {

		if !isValidType(strings.ToUpper(r.Type)) || !buf.WriteValidString(r.Type, fmt.Sprintf("#EXT-X-MEDIA:TYPE=%s", r.Type)) {
			return attributeNotSetError("EXT-X-MEDIA", "TYPE")
		}
		if !buf.WriteValidString(r.GroupID, fmt.Sprintf(",GROUP-ID=\"%s\"", r.GroupID)) {
			return attributeNotSetError("EXT-X-MEDIA", "GROUP-ID")
		}
		if !buf.WriteValidString(r.Name, fmt.Sprintf(",NAME=\"%s\"", r.Name)) {
			return attributeNotSetError("EXT-X-MEDIA", "NAME")
		}
		buf.WriteValidString(r.Language, fmt.Sprintf(",LANGUAGE=\"%s\"", r.Language))
		buf.WriteValidString(r.AssocLanguage, fmt.Sprintf(",ASSOC-LANGUAGE=\"%s\"", r.AssocLanguage))
		buf.WriteValidString(r.Default, ",DEFAULT=YES")
		if r.Forced && strings.ToUpper(r.Type) == sub {
			buf.WriteValidString(r.Forced, ",FORCED=YES")
		}
		if strings.ToUpper(r.Type) == cc && isValidInstreamID(strings.ToUpper(r.InstreamID)) {
			buf.WriteValidString(r.InstreamID, fmt.Sprintf(",INSTREAM-ID=\"%s\"", r.InstreamID))
		}
		buf.WriteValidString(r.Characteristics, fmt.Sprintf(",CHARACTERISTICS=\"%s\"", r.Characteristics))

		//URI is required for SUBTITLES and MUST NOT be present for CLOSED-CAPTIONS, other types URI is optinal
		if strings.ToUpper(r.Type) == sub {
			if !buf.WriteValidString(r.URI, fmt.Sprintf(",URI=\"%s\"", r.URI)) {
				return attributeNotSetError("EXT-X-MEDIA", "URI for SUBTITLES")
			}
		} else if strings.ToUpper(r.Type) != cc {
			buf.WriteValidString(r.URI, fmt.Sprintf(",URI=\"%s\"", r.URI))
		}

		buf.WriteRune('\n')
	}
	return buf.err
}

func isValidType(t string) bool {
	return t == aud || t == vid || t == cc || t == sub
}

func isValidInstreamID(instream string) bool {
	return instream == "CC1" || instream == "CC2" || instream == "CC3" || instream == "CC4" || strings.HasPrefix(instream, "SERVICE")
}

//writeStreamInf sets the EXT-X-STREAM-INF or EXT-X-I-FRAME-STREAM-INF tag on Master Playlist file
func (v *Variant) writeStreamInf(version int, buf *BufWriter) error {
	if v != nil {
		if v.IsIframe {
			buf.WriteString("#EXT-X-I-FRAME-STREAM-INF:")
		} else {
			buf.WriteString("#EXT-X-STREAM-INF:")
		}

		if !buf.WriteValidString(v.Bandwidth, fmt.Sprintf("BANDWIDTH=%s", strconv.FormatInt(v.Bandwidth, 10))) {
			return attributeNotSetError("Variant", "BANDWIDTH")
		}
		if version < 6 && v.ProgramID > 0 {
			buf.WriteValidString(v.ProgramID, fmt.Sprintf(",PROGRAM-ID=%s", strconv.FormatInt(v.ProgramID, 10)))
		}
		buf.WriteValidString(v.AvgBandwidth, fmt.Sprintf(",AVERAGE-BANDWIDTH=%s", strconv.FormatInt(v.AvgBandwidth, 10)))
		buf.WriteValidString(v.Codecs, fmt.Sprintf(",CODECS=\"%s\"", v.Codecs))
		buf.WriteValidString(v.Resolution, fmt.Sprintf(",RESOLUTION=%s", v.Resolution))
		buf.WriteValidString(v.FrameRate, fmt.Sprintf(",FRAME-RATE=%s", strconv.FormatFloat(v.FrameRate, 'f', 3, 32)))
		buf.WriteValidString(v.Video, fmt.Sprintf(",VIDEO=\"%s\"", v.Video))
		//If is not IFrame tag, adds AUDIO, SUBTITLES and CLOSED-CAPTIONS params
		//If IFrame, add URI as a param
		if !v.IsIframe {
			buf.WriteValidString(v.Audio, fmt.Sprintf(",AUDIO=\"%s\"", v.Audio))
			buf.WriteValidString(v.Subtitles, fmt.Sprintf(",SUBTITLES=\"%s\"", v.Subtitles))
			buf.WriteValidString(v.ClosedCaptions, fmt.Sprintf(",CLOSED-CAPTIONS=\"%s\"", v.ClosedCaptions))
			buf.WriteString(fmt.Sprintf("\n%s\n\n", v.URI))
		} else {
			buf.WriteValidString(v.URI, fmt.Sprintf(",URI=\"%s\"\n\n", v.URI))
		}
	}
	return buf.err
}

func (p *MediaPlaylist) writeTargetDuration(buf *BufWriter) error {
	if !buf.WriteValidString(p.TargetDuration, fmt.Sprintf("#EXT-X-TARGETDURATION:%s\n", strconv.Itoa(p.TargetDuration))) {
		return attributeNotSetError("EXT-X-TARGETDURATION", "")
	}
	return buf.err
}

func (p *MediaPlaylist) writeMediaSequence(buf *BufWriter) {
	if p.MediaSequence > 0 {
		buf.WriteString(fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%s\n", strconv.Itoa(p.MediaSequence)))
	}
}

func (p *MediaPlaylist) writeDiscontinuitySequence(buf *BufWriter) {
	if p.DiscontinuitySequence > 0 {
		buf.WriteString(fmt.Sprintf("#EXT-X-DISCONTINUITY-SEQUENCE:%s\n", strconv.Itoa(p.DiscontinuitySequence)))
	}
}

func (p *MediaPlaylist) writeAllowCache(buf *BufWriter) {
	if p.Version < 7 && p.AllowCache {
		buf.WriteString("#EXT-X-ALLOW-CACHE\n")
	}
}

func (p *MediaPlaylist) writePlaylistType(buf *BufWriter) {
	if p.Type != "" {
		buf.WriteString(fmt.Sprintf("#EXT-X-PLAYLIST-TYPE:%s\n", p.Type))
	}
}

func (p *MediaPlaylist) writeIFramesOnly(buf *BufWriter) error {
	if p.IFramesOnly {
		if p.Version < 4 {
			return backwardsCompatibilityError(p.Version, "#EXT-X-I-FRAMES-ONLY")
		}
		buf.WriteString("#EXT-X-I-FRAMES-ONLY\n")
	}
	return buf.err
}

func (s *Segment) writeSegmentTags(buf *BufWriter) error {
	if s != nil {
		if s.Keys != nil {
			for _, key := range s.Keys {
				if err := key.writeKey(buf); err != nil {
					return err
				}
			}
		}

		if err := s.Map.writeMap(buf); err != nil {
			return err
		}

		if !s.ProgramDateTime.IsZero() {
			buf.WriteString(fmt.Sprintf("#EXT-X-PROGRAM-DATE-TIME:%s\n", s.ProgramDateTime.String()))
		}
		buf.WriteValidString(s.Discontinuity, "#EXT-X-DISCONTINUITY\n")

		//TODO: Check if only one DateRange per segment
		if err := s.DateRange.writeDateRange(buf); err != nil {
			return err
		}

		if s.Inf == nil || !buf.WriteValidString(s.Inf.Duration, fmt.Sprintf("#EXTINF:%s,%s\n", strconv.FormatFloat(s.Inf.Duration, 'f', 3, 32), s.Inf.Title)) {
			return attributeNotSetError("EXTINF", "DURATION")
		}

		if s.Byterange != nil {
			buf.WriteString(fmt.Sprintf("#EXT-X-BYTERANGE:%s", strconv.FormatInt(s.Byterange.Length, 10)))
			if s.Byterange.Offset != nil {
				buf.WriteString("@" + strconv.FormatInt(*s.Byterange.Offset, 10))
			}
			buf.WriteRune('\n')
		}

		if s.URI != "" {
			buf.WriteString(s.URI)
			buf.WriteRune('\n')
		} else {
			return attributeNotSetError("Segment", "URI")
		}
	}
	return buf.err
}

func (k *Key) writeKey(buf *BufWriter) error {
	if k != nil {
		if k.IsSession {
			buf.WriteString("#EXT-X-SESSION-KEY:")
		} else {
			buf.WriteString("#EXT-X-KEY:")
		}

		if !isValidMethod(k.IsSession, strings.ToUpper(k.Method)) ||
			!buf.WriteValidString(k.Method, fmt.Sprintf("METHOD=%s", strings.ToUpper(k.Method))) {
			return attributeNotSetError("KEY", "METHOD")
		}
		if k.URI != "" && strings.ToUpper(k.Method) != none {
			buf.WriteValidString(k.URI, fmt.Sprintf(",URI=\"%s\"", k.URI))
		} else {
			return attributeNotSetError("EXT-X-KEY", "URI")
		}
		buf.WriteValidString(k.IV, fmt.Sprintf(",IV=%s", k.IV))
		buf.WriteValidString(k.Keyformat, fmt.Sprintf(",KEYFORMAT=\"%s\"", k.Keyformat))
		buf.WriteValidString(k.Keyformatversions, fmt.Sprintf(",KEYFORMATVERSIONS=\"%s\"", k.Keyformatversions))
		buf.WriteRune('\n')
	}
	return buf.err
}

//Session Key Method can't be NONE
func isValidMethod(isSession bool, method string) bool {
	return (method == aes || method == sample) || (!isSession && method == none)
}

func (m *Map) writeMap(buf *BufWriter) error {
	if m != nil {
		if !buf.WriteValidString(m.URI, fmt.Sprintf("#EXT-X-MAP:URI=\"%s\"", m.URI)) {
			return attributeNotSetError("EXT-X-MAP", "URI")
		}
		if m.Byterange != nil {
			if m.Byterange.Offset == nil {
				o := int64(0)
				m.Byterange.Offset = &o
			}
			buf.WriteString(fmt.Sprintf(",BYTERANGE=\"%s@%s\"",
				strconv.FormatInt(m.Byterange.Length, 10),
				strconv.FormatInt(*m.Byterange.Offset, 10)))
		}
		buf.WriteRune('\n')
	}
	return buf.err
}

func (d *DateRange) writeDateRange(buf *BufWriter) error {
	if d != nil {
		if !buf.WriteValidString(d.ID, fmt.Sprintf("#EXT-X-DATERANGE:ID=%s", d.ID)) {
			return attributeNotSetError("EXT-X-DATERANGE", "ID")
		}
		buf.WriteValidString(d.Class, fmt.Sprintf(",CLASS=\"%s\"", d.Class))
		if !d.StartDate.IsZero() {
			buf.WriteString(fmt.Sprintf(",START-DATE=\"%s\"", d.StartDate.String()))
		} else {
			return attributeNotSetError("EXT-X-DATERANGE", "START-DATE")
		}
		if !d.EndDate.IsZero() {
			if d.EndDate.Before(d.StartDate) {
				return errors.New("DateRange attribute EndDate must be equal or later than StartDate")
			}
			buf.WriteString(fmt.Sprintf(",END-DATE=\"%s\"", d.EndDate.String()))
		}
		if d.Duration != nil && *d.Duration >= float64(0) {
			buf.WriteString(fmt.Sprintf(",DURATION=%s", strconv.FormatFloat(*d.Duration, 'f', 3, 32)))
		}
		if d.PlannedDuration != nil && *d.PlannedDuration >= float64(0) {
			buf.WriteString(fmt.Sprintf(",PLANNED-DURATION=%s", strconv.FormatFloat(*d.PlannedDuration, 'f', 3, 32)))
		}
		if len(d.XClientAttribute) > 0 {
			for _, customTag := range d.XClientAttribute {
				if !strings.HasPrefix(strings.ToUpper(customTag), "X-") {
					return errors.New("EXT-X-DATERANGE client-defined attributes must start with X-")
				}
				buf.WriteString(",")
				buf.WriteString(strings.ToUpper(customTag))
			}
		}

		d.SCTE35.writeSCTE(buf)

		if buf.WriteValidString(d.EndOnNext, ",END-ON-NEXT=YES") {
			if d.Class == "" {
				return errors.New("EXT-X-DATERANGE tag must have a CLASS attribute when END-ON-NEXT attribue is present")
			}
			if d.Duration != nil || !d.EndDate.IsZero() {
				return errors.New("EXT-X-DATERANGE tag must not have DURATION or END-DATE attributes when END-ON-NEXT attribute is present")
			}
		}
		buf.WriteRune('\n')
	}
	return buf.err
}

func (s *SCTE35) writeSCTE(buf *BufWriter) {
	if s != nil {
		t := strings.ToUpper(s.Type)
		if t == "IN" || t == "OUT" || t == "CMD" {
			if !buf.WriteValidString(s.Value, fmt.Sprintf(",SCTE35-%s=%s", t, s.Value)) {
				buf.err = attributeNotSetError("SCTE35", "Value")
			}
			return
		}
		buf.err = errors.New("SCTE35 type must be IN, OUT or CMD")
	}
}

func (p *MediaPlaylist) writeEndList(buf *BufWriter) {
	if p.EndList {
		buf.WriteString("#EXT-X-ENDLIST\n")
	}
}

//checkCompatibility checks backwards compatibility issues according to the Media Playlist version
func (p *MediaPlaylist) checkCompatibility(s *Segment) error {
	if s != nil {
		for _, key := range s.Keys {
			if (strings.ToUpper(key.Method) == sample || key.Keyformat != "" || key.Keyformatversions != "") && p.Version < 5 {
				return backwardsCompatibilityError(p.Version, "#EXT-X-KEY")
			}
		}

		if s.Map != nil {
			if p.Version < 5 || (!p.IFramesOnly && p.Version < 6) {
				return backwardsCompatibilityError(p.Version, "#EXT-X-MAP")
			}
		}

		if s.Byterange != nil && p.Version < 4 {
			return backwardsCompatibilityError(p.Version, "#EXT-X-BYTERANGE")
		}
	}
	return nil
}
