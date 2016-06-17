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
		if !bufWriteString(buf, sp.TimeOffset, fmt.Sprintf("#EXT-X-START:TIME-OFFSET=%s", strconv.FormatFloat(sp.TimeOffset, 'f', 3, 32))) {
			return attributeNotSetError("EXT-X-START", "TIME-OFFSET")
		}
		bufWriteString(buf, sp.Precise, ",PRECISE=YES")
		buf.WriteRune('\n')
	}
	return nil
}

func bufWriteString(buf *bytes.Buffer, data interface{}, write string) bool {
	switch data.(type) {
	case string:
		if data.(string) != "" {
			buf.WriteString(write)
			return true
		}
	case float64:
		if data.(float64) > float64(0) {
			buf.WriteString(write)
			return true
		}
	case int, int64:
		if data.(int64) > 0 {
			buf.WriteString(write)
			return true
		}
	case bool:
		if data.(bool) {
			buf.WriteString(write)
			return true
		}
	}

	return false
}

//writeSessionData sets the EXT-X-SESSION-DATA tag on Master Playlist file
func (s *SessionData) writeSessionData(buf *bytes.Buffer) error {
	if s != nil {
		if !bufWriteString(buf, s.DataID, fmt.Sprintf("#EXT-X-SESSION-DATA:DATA-ID=\"%s\"", s.DataID)) {
			return attributeNotSetError("EXT-X-SESSION-DATA", "DATA-ID")
		}

		if s.Value != "" && s.URI != "" {
			return errors.New("EXT-X-SESSION-DATA must have attributes URI or VALUE, not both")
		} else if s.Value != "" || s.URI != "" {
			bufWriteString(buf, s.Value, fmt.Sprintf(",VALUE=\"%s\"", s.Value))
			bufWriteString(buf, s.URI, fmt.Sprintf(",URI=\"%s\"", s.URI))
		} else {
			return errors.New("EXT-X-SESSION-DATA must have either URI or VALUE attributes set")
		}

		bufWriteString(buf, s.Language, fmt.Sprintf(",LANGUAGE=\"%s\"", s.Language))
		buf.WriteRune('\n')
	}
	return nil
}

//writeXMedia sets the EXT-X-MEDIA tag on Master Playlist file
func (r *Rendition) writeXMedia(buf *bytes.Buffer) error {
	if r != nil {

		if !bufWriteString(buf, r.Type, fmt.Sprintf("#EXT-X-MEDIA:TYPE=%s", r.Type)) {
			return attributeNotSetError("EXT-X-MEDIA", "TYPE")
		}
		if !bufWriteString(buf, r.GroupID, fmt.Sprintf(",GROUP-ID=\"%s\"", r.GroupID)) {
			return attributeNotSetError("EXT-X-MEDIA", "GROUP-ID")
		}
		if !bufWriteString(buf, r.Name, fmt.Sprintf(",NAME=\"%s\"", r.Name)) {
			return attributeNotSetError("EXT-X-MEDIA", "NAME")
		}
		bufWriteString(buf, r.Language, fmt.Sprintf(",LANGUAGE=\"%s\"", r.Language))
		bufWriteString(buf, r.AssocLanguage, fmt.Sprintf(",ASSOC-LANGUAGE=\"%s\"", r.AssocLanguage))
		bufWriteString(buf, r.Default, ",DEFAULT=YES")
		if r.Forced && strings.ToUpper(r.Type) == "SUBTITLES" {
			bufWriteString(buf, r.Forced, ",FORCED=YES")
		}
		if r.InstreamID != "" && strings.ToUpper(r.Type) == "CLOSED-CAPTIONS" && isValidInstreamID(strings.ToUpper(r.InstreamID)) {
			bufWriteString(buf, r.InstreamID, fmt.Sprintf(",INSTREAM-ID=\"%s\"", r.InstreamID))
		}
		bufWriteString(buf, r.Characteristics, fmt.Sprintf("CHARACTERISTICS=\"%s\"", r.Characteristics))
		bufWriteString(buf, r.URI, fmt.Sprintf(",URI=\"%s\"", r.URI))
		buf.WriteRune('\n')
	}
	return nil
}

func isValidInstreamID(instream string) bool {
	return instream == "CC1" || instream == "CC2" || instream == "CC3" || instream == "CC4" || strings.HasPrefix(instream, "SERVICE")
}

//writeStreamInf sets the EXT-X-STREAM-INF or EXT-X-I-FRAME-STREAM-INF tag on Master Playlist file
func (v *Variant) writeStreamInf(version int, buf *bytes.Buffer) error {
	if v != nil {
		if v.IsIframe {
			buf.WriteString("#EXT-X-I-FRAME-STREAM-INF:")
		} else {
			buf.WriteString("#EXT-X-STREAM-INF:")
		}

		if !bufWriteString(buf, v.Bandwidth, fmt.Sprintf("BANDWIDTH=%s", strconv.FormatInt(v.Bandwidth, 10))) {
			return attributeNotSetError("Variant", "BANDWIDTH")
		}
		if version < 6 && v.ProgramID > 0 {
			bufWriteString(buf, v.ProgramID, fmt.Sprintf(",PROGRAM-ID=%s", strconv.FormatInt(v.ProgramID, 10)))
		}
		bufWriteString(buf, v.AvgBandwidth, fmt.Sprintf(",AVERAGE-BANDWIDTH=%s", strconv.FormatInt(v.AvgBandwidth, 10)))
		bufWriteString(buf, v.Codecs, fmt.Sprintf(",CODECS=\"%s\"", v.Codecs))
		bufWriteString(buf, v.Resolution, fmt.Sprintf(",RESOLUTION=%s", v.Resolution))
		bufWriteString(buf, v.FrameRate, fmt.Sprintf(",FRAME-RATE=%s", strconv.FormatFloat(v.FrameRate, 'f', 3, 32)))
		bufWriteString(buf, v.Video, fmt.Sprintf(",VIDEO=\"%s\"", v.Video))
		//If is not IFrame tag, adds AUDIO, SUBTITLES and CLOSED-CAPTIONS params
		//If IFrame, add URI as a param
		if !v.IsIframe {
			bufWriteString(buf, v.Audio, fmt.Sprintf(",AUDIO=\"%s\"", v.Audio))
			bufWriteString(buf, v.Subtitles, fmt.Sprintf(",SUBTITLES=\"%s\"", v.Subtitles))
			bufWriteString(buf, v.ClosedCaptions, fmt.Sprintf(",CLOSED-CAPTIONS=\"%s\"", v.ClosedCaptions))
			buf.WriteString(fmt.Sprintf("\n%s\n\n", v.URI))
		} else {
			bufWriteString(buf, v.URI, fmt.Sprintf(",URI=\"%s\"\n\n", v.URI))
		}
	}
	return nil
}

func (p *MediaPlaylist) writeTargetDuration(buf *bytes.Buffer) error {
	if !bufWriteString(buf, p.TargetDuration, fmt.Sprintf("#EXT-X-TARGETDURATION:%s\n", strconv.Itoa(p.TargetDuration))) {
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
		bufWriteString(buf, s.Discontinuity, "#EXT-X-DISCONTINUITY\n")

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
		if k.IsSession {
			buf.WriteString("#EXT-X-SESSION-KEY:")
		} else {
			buf.WriteString("#EXT-X-KEY:")
		}

		if !isValidMethod(k.IsSession, strings.ToUpper(k.Method)) ||
			!bufWriteString(buf, k.Method, fmt.Sprintf("METHOD=%s", strings.ToUpper(k.Method))) {
			return attributeNotSetError("KEY", "METHOD")
		}
		if k.URI != "" && strings.ToUpper(k.Method) != "NONE" {
			bufWriteString(buf, k.URI, fmt.Sprintf(",URI=\"%s\"", k.URI))
		} else {
			return attributeNotSetError("EXT-X-KEY", "URI")
		}
		bufWriteString(buf, k.IV, fmt.Sprintf(",IV=%s", k.IV))
		bufWriteString(buf, k.Keyformat, fmt.Sprintf("KEYFORMAT=\"%s\"", k.Keyformat))
		bufWriteString(buf, k.Keyformatversions, fmt.Sprintf("KEYFORMATVERSIONS=\"%s\"", k.Keyformatversions))

	}
	return nil
}

//Session Key Method can't be NONE
func isValidMethod(isSession bool, method string) bool {
	return (method == "AES-128" || method == "SAMPLE-AES") || (!isSession && method == "NONE")
}

func (m *Map) writeMap(buf *bytes.Buffer) error {
	if m != nil {
		if !bufWriteString(buf, m.URI, fmt.Sprintf("#EXT-X-MAP:URI=\"%s\"", m.URI)) {
			return attributeNotSetError("EXT-X-MAP", "URI")
		}
		//TODO:look if offset is included when = 0
		if m.Byterange != nil {
			buf.WriteString(fmt.Sprintf(",BYTERANGE=\"%s@%s\"",
				strconv.FormatInt(m.Byterange.Length, 10),
				strconv.FormatInt(m.Byterange.Offset, 10)))
		}
		buf.WriteRune('\n')
	}
	return nil
}

func (d *DateRange) writeDateRange(buf *bytes.Buffer) error {
	if d != nil {
		if !bufWriteString(buf, d.ID, fmt.Sprintf("#EXT-X-DATERANGE:ID=%s", d.ID)) {
			return attributeNotSetError("EXT-X-DATERANGE", "ID")
		}
		bufWriteString(buf, d.Class, fmt.Sprintf(",CLASS=\"%s\"", d.Class))
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
		if d.XClientAttribute != nil {
			// for _, x := range d.XClientAttribute {
			// 	//TODO:
			// }
		}

		//TODO:SCET35
		bufWriteString(buf, d.EndOnNext, ",END-ON-NEXT=YES")
		buf.WriteRune('\n')
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
