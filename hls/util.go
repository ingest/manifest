package hls

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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
	return fmt.Errorf(tag + " attribute " + attribute + " must be set.")
}

//writeHeader sets the initial tags for both Media and Master Playlists files
func writeHeader(version int, buf *bytes.Buffer) {
	buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
	buf.WriteString(strconv.Itoa(version))
	buf.WriteRune('\n')
}

//writeIndependentSegment sets the #EXT-X-INDEPENDENT-SEGMENTS tag on Media and Master Playlist file
func writeIndependentSegment(isIndSeg bool, buf *bytes.Buffer) {
	if isIndSeg {
		buf.WriteString("EXT-X-INDEPENDENT-SEGMENTS\n")
	}
}

//writeStartPoint sets the #EXT-X-START tag on Media and Master Playlist file
func writeStartPoint(sp *StartPoint, buf *bytes.Buffer) error {
	if sp != nil {
		var attributes []string
		if sp.TimeOffset == float64(0) {
			return attributeNotSetError("EXT-X-START", "TIME-OFFSET")
		}
		attributes = append(attributes, fmt.Sprintf("TIME-OFFSET=%s", strconv.FormatFloat(sp.TimeOffset, 'f', 3, 32)))
		if sp.Precise {
			attributes = append(attributes, "PRECISE=YES")
		}
		buf.WriteString(fmt.Sprintf("#EXT-X-START:%s", strings.Join(attributes, ",")))
	}
	return nil
}

//writeSessionData sets the EXT-X-SESSION-DATA tag on Master Playlist file
func (s *SessionData) writeSessionData(buf *bytes.Buffer) error {
	if s != nil {
		//Initiate a slice of string of size 4 (number of fields in SessionData)
		var attributes []string

		if s.DataID != "" {
			attributes = append(attributes, fmt.Sprintf("DATA-ID=\"%s\"", s.DataID))
		} else {
			return attributeNotSetError("EXT-X-SESSION-DATA", "DATA-ID")
		}

		if s.Value != "" && s.URI != "" {
			return errors.New("EXT-X-SESSION-DATA must have attributes URI or VALUE, not both.")
		} else if s.Value != "" {
			attributes = append(attributes, fmt.Sprintf("VALUE=\"%s\"", s.Value))
		} else if s.URI != "" {
			attributes = append(attributes, fmt.Sprintf("URI=\"%s\"", s.URI))
		} else {
			return errors.New("EXT-X-SESSION-DATA must have either URI or VALUE attributes set.")
		}

		if s.Language != "" {
			attributes = append(attributes, fmt.Sprintf("LANGUAGE=\"%s\"", s.Language))
		}

		buf.WriteString(fmt.Sprintf("#EXT-X-SESSION-DATA:%s\n", strings.Join(attributes, ",")))
	}
	return nil
}

//writeXMedia sets the EXT-X-MEDIA tag on Master Playlist file
func (r *Rendition) writeXMedia(buf *bytes.Buffer) error {
	if r != nil {
		//Initiate a slice of string of size 11 (number of fields in Rendition)
		var attributes []string

		if r.Type != "" {
			attributes = append(attributes, fmt.Sprintf("TYPE=%s", r.Type))
		} else {
			return attributeNotSetError("EXT-X-MEDIA", "TYPE")
		}
		if r.GroupID != "" {
			attributes = append(attributes, fmt.Sprintf("GROUP-ID=\"%s\"", r.GroupID))
		} else {
			return attributeNotSetError("EXT-X-MEDIA", "GROUP-ID")
		}
		if r.Name != "" {
			attributes = append(attributes, fmt.Sprintf("NAME=\"%s\"", r.Name))
		} else {
			return attributeNotSetError("EXT-X-MEDIA", "NAME")
		}
		if r.Language != "" {
			attributes = append(attributes, fmt.Sprintf("LANGUAGE=\"%s\"", r.Language))
		}
		if r.AssocLanguage != "" {
			attributes = append(attributes, fmt.Sprintf("ASSOC-LANGUAGE=\"%s\"", r.AssocLanguage))
		}
		if r.Default {
			attributes = append(attributes, "DEFAULT=YES")
		}
		if r.AutoSelect {
			attributes = append(attributes, "AUTOSELECT=YES")
		}
		if r.Forced && strings.ToUpper(r.Type) == "SUBTITLES" {
			attributes = append(attributes, "FORCED=YES")
		}
		if r.InstreamID != "" && strings.ToUpper(r.Type) == "CLOSED-CAPTIONS" && isValidInstreamID(strings.ToUpper(r.InstreamID)) {
			attributes = append(attributes, fmt.Sprintf("INSTREAM-ID=\"%s\"", r.InstreamID))
		}
		if r.Characteristics != "" {
			attributes = append(attributes, fmt.Sprintf("CHARACTERISTICS=\"%s\"", r.Characteristics))
		}
		if r.URI != "" {
			attributes = append(attributes, fmt.Sprintf("URI=\"%s\"", r.URI))
		}
		buf.WriteString(fmt.Sprintf("#EXT-X-MEDIA:%s\n", strings.Join(attributes, ",")))
	}
	return nil
}

func isValidInstreamID(instream string) bool {
	return instream == "CC1" || instream == "CC2" || instream == "CC3" || instream == "CC4" || strings.HasPrefix(instream, "SERVICE")
}

//writeStreamInf sets the EXT-X-STREAM-INF or EXT-X-I-FRAME-STREAM-INF tag on Master Playlist file
func (v *Variant) writeStreamInf(version int, buf *bytes.Buffer) {
	if v != nil {
		var attributes []string

		if version < 6 && v.ProgramID > 0 {
			attributes = append(attributes, fmt.Sprintf("PROGRAM-ID=%s", strconv.FormatInt(v.ProgramID, 10)))
		}
		if v.Bandwidth > 0 {
			attributes = append(attributes, fmt.Sprintf("BANDWIDTH=%s", strconv.FormatInt(v.Bandwidth, 10)))
		}
		if v.AvgBandwidth > 0 {
			attributes = append(attributes, fmt.Sprintf("AVERAGE-BANDWIDTH=%s", strconv.FormatInt(v.AvgBandwidth, 10)))
		}
		if v.Codecs != "" {
			attributes = append(attributes, fmt.Sprintf("CODECS=\"%s\"", v.Codecs))
		}
		if v.Resolution != "" {
			attributes = append(attributes, fmt.Sprintf("RESOLUTION=%s", v.Resolution))
		}
		if v.FrameRate > float64(0) {
			attributes = append(attributes, fmt.Sprintf("FRAME-RATE=%s", strconv.FormatFloat(v.FrameRate, 'f', 3, 32)))
		}
		if v.Video != "" {
			attributes = append(attributes, fmt.Sprintf("VIDEO=\"%s\"", v.Video))
		}
		//If is not IFrame tag, adds AUDIO, SUBTITLES and CLOSED-CAPTIONS params
		//If IFrame, add URI as a param
		if !v.IsIframe {
			if v.Audio != "" {
				attributes = append(attributes, fmt.Sprintf("AUDIO=\"%s\"", v.Audio))
			}
			if v.Subtitles != "" {
				attributes = append(attributes, fmt.Sprintf("SUBTITLES=\"%s\"", v.Subtitles))
			}
			if v.ClosedCaptions != "" {
				attributes = append(attributes, fmt.Sprintf("CLOSED-CAPTIONS=\"%s\"", v.ClosedCaptions))
			}
			buf.WriteString(fmt.Sprintf("#EXT-X-STREAM-INF:%s\n%s\n\n", strings.Join(attributes, ","), v.URI))
		} else {
			attributes = append(attributes, fmt.Sprintf("URI=\"%s\"", v.URI))
			buf.WriteString(fmt.Sprintf("#EXT-X-I-FRAME-STREAM-INF:%s\n\n", strings.Join(attributes, ",")))
		}
	}
}

func (p *MediaPlaylist) writeTargetDuration(buf *bytes.Buffer) error {
	if p.TargetDuration > 0 {
		buf.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%s\n", strconv.Itoa(p.TargetDuration)))
	} else {
		return attributeNotSetError("EXT-X-TARGETDURATION", "")
	}
	return nil
}

func (p *MediaPlaylist) writeMediaSequence(buf *bytes.Buffer) {
	if p.MediaSequence > 0 {
		buf.WriteString(fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%s\n", strconv.Itoa(p.MediaSequence)))
	}
}

func (p *MediaPlaylist) writeDiscontinuitySequence(buf *bytes.Buffer) {
	if p.DiscontinuitySequence > 0 {
		buf.WriteString(fmt.Sprintf("#EXT-X-DISCONTINUITY-SEQUENCE:%s\n", strconv.Itoa(p.DiscontinuitySequence)))
	}
}

func (p *MediaPlaylist) writeAllowCache(buf *bytes.Buffer) {
	if p.Version < 7 && p.AllowCache {
		buf.WriteString("#EXT-X-ALLOW-CACHE\n")
	}
}

func (p *MediaPlaylist) writePlaylistType(buf *bytes.Buffer) {
	if p.Type != "" {
		buf.WriteString(fmt.Sprintf("#EXT-X-PLAYLIST-TYPE:%s\n", p.Type))
	}
}

func (p *MediaPlaylist) writeIFramesOnly(buf *bytes.Buffer) {
	if p.IFramesOnly {
		buf.WriteString("#EXT-X-I-FRAMES-ONLY\n")
	}
}

func (s *Segment) writeSegmentTags(buf *bytes.Buffer) error {
	if s != nil {
		//TODO:confirm if only one key per segment possible
		if err := s.Key.writeKey(buf); err != nil {
			return err
		}

		if err := s.Map.writeMap(buf); err != nil {
			return err
		}

		if !s.ProgramDateTime.IsZero() {
			buf.WriteString(fmt.Sprintf("#EXT-X-PROGRAM-DATE-TIME:%s\n", s.ProgramDateTime.String()))
		}

		if s.Discontinuity {
			buf.WriteString("#EXT-X-DISCONTINUITY\n")
		}

		if s.Inf != nil && s.Inf.Duration > float64(0) {
			buf.WriteString(fmt.Sprintf("#EXTINF:%s,%s\n", strconv.FormatFloat(s.Inf.Duration, 'f', 3, 32), s.Inf.Title))
		} else {
			return attributeNotSetError("EXTINF", "DURATION")
		}

		if s.Byterange != nil {
			buf.WriteString(fmt.Sprintf("#EXT-X-BYTERANGE:%s", strconv.FormatInt(s.Byterange.Length, 10)))
			if s.Byterange.Offset > 0 {
				buf.WriteString("@" + strconv.FormatInt(s.Byterange.Offset, 10))
			}
			buf.WriteRune('\n')
		}

		//write DateRange

		if s.URI != "" {
			buf.WriteString(s.URI)
			buf.WriteRune('\n')
		} else {
			return attributeNotSetError("Segment", "URI")
		}
	}
	return nil
}

func (k *Key) writeKey(buf *bytes.Buffer) error {
	if k != nil {
		var attributes []string

		if k.Method != "" && isValidMethod(k.IsSession, strings.ToUpper(k.Method)) {
			attributes = append(attributes, fmt.Sprintf("METHOD=%s", k.Method))
		} else {
			return attributeNotSetError("KEY", "METHOD")
		}
		if k.URI != "" && strings.ToUpper(k.Method) != "NONE" {
			attributes = append(attributes, fmt.Sprintf("URI=\"%s\"", k.URI))
		} else {
			return attributeNotSetError("EXT-X-KEY", "URI")
		}
		if k.IV != "" {
			attributes = append(attributes, fmt.Sprintf("IV=%s", k.IV))
		}
		if k.Keyformat != "" {
			attributes = append(attributes, fmt.Sprintf("KEYFORMAT=\"%s\"", k.Keyformat))
		}
		if k.Keyformatversions != "" {
			attributes = append(attributes, fmt.Sprintf("KEYFORMATVERSIONS=\"%s\"", k.Keyformatversions))
		}
		if k.IsSession {
			buf.WriteString(fmt.Sprintf("#EXT-X-SESSION-KEY:%s", strings.Join(attributes, ",")))
		} else {
			buf.WriteString(fmt.Sprintf("#EXT-X-KEY:%s", strings.Join(attributes, ",")))
		}
	}
	return nil
}

//Session Key Method can't be NONE
func isValidMethod(isSession bool, method string) bool {
	return (method == "AES-128" || method == "SAMPLE-AES") || (!isSession && method == "NONE")
}

func (m *Map) writeMap(buf *bytes.Buffer) error {
	if m != nil {
		var attributes []string
		if m.URI != "" {
			attributes = append(attributes, fmt.Sprintf("URI=\"%s\"", m.URI))
		} else {
			return attributeNotSetError("EXT-X-MAP", "URI")
		}
		//TODO:look if offset is included when = 0
		if m.Byterange != nil {
			attributes = append(attributes,
				fmt.Sprintf("BYTERANGE=\"%s@%s\"",
					strconv.FormatInt(m.Byterange.Length, 10),
					strconv.FormatInt(m.Byterange.Offset, 10)))
		}

		buf.WriteString(fmt.Sprintf("#EXT-X-MAP:%s", strings.Join(attributes, ",")))
	}
	return nil
}

func (d *DateRange) writeDateRange(buf *bytes.Buffer) error {
	if d != nil {
		var attributes []string

		if d.ID != "" {
			attributes = append(attributes, fmt.Sprintf("ID=\"%s\"", d.ID))
		} else {
			return attributeNotSetError("EXT-X-DATERANGE", "ID")
		}
		if d.Class != "" {
			attributes = append(attributes, fmt.Sprintf("CLASS=\"%s\"", d.Class))
		}
		if !d.StartDate.IsZero() {
			attributes = append(attributes, fmt.Sprintf("START-DATE=\"%s\"", d.StartDate.String()))
		} else {
			return attributeNotSetError("EXT-X-DATERANGE", "START-DATE")
		}
		if !d.EndDate.IsZero() {
			if d.EndDate.Before(d.StartDate) {
				return errors.New("DateRange attribute EndDate must be equal or later than StartDate")
			}
			attributes = append(attributes, fmt.Sprintf("END-DATE=\"%s\"", d.EndDate.String()))
		}
		if d.Duration != nil && *d.Duration >= float64(0) {
			attributes = append(attributes, fmt.Sprintf("DURATION=%s", strconv.FormatFloat(*d.Duration, 'f', 3, 32)))
		}
		if d.PlannedDuration != nil && *d.PlannedDuration >= float64(0) {
			attributes = append(attributes, fmt.Sprintf("PLANNED-DURATION=%s", strconv.FormatFloat(*d.PlannedDuration, 'f', 3, 32)))
		}
		if d.XClientAttribute != nil {
			for _, x := range d.XClientAttribute {
				attributes = append(attributes, x)
			}
		}

		//TODO:SCET35

		if d.EndOnNext {
			attributes = append(attributes, "END-ON-NEXT=YES")
		}

		buf.WriteString(fmt.Sprintf("#EXT-X-DATERANGE:%s\n", strings.Join(attributes, ",")))
	}
	return nil
}

func (p *MediaPlaylist) writeEndList(buf *bytes.Buffer) {
	if p.EndList {
		buf.WriteString("#EXT-X-ENDLIST\n")
	}
}

//checkCompatibility checks backwards compatibility issues according to the Media Playlist version
func (p *MediaPlaylist) checkCompatibility(s *Segment) error {
	if s != nil {
		if s.Key != nil {
			if (strings.ToUpper(s.Key.Method) == "SAMPLE-AES" || s.Key.Keyformat != "" || s.Key.Keyformatversions != "") && p.Version < 5 {
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
