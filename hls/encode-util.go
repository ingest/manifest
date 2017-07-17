package hls

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ingest/manifest"
)

// NewMediaPlaylist returns an instance of a MediaPlaylist with a set version
func NewMediaPlaylist(version int) *MediaPlaylist {
	return &MediaPlaylist{
		Version: version,
	}
}

// WithVariant supplies the data which was processed from the master playlist.
func (p *MediaPlaylist) WithVariant(v *Variant) *MediaPlaylist {
	p.Variant = v
	return p
}

// NewMasterPlaylist returns an instance of a MasterPlaylist with a set version
func NewMasterPlaylist(version int) *MasterPlaylist {
	return &MasterPlaylist{Version: version}
}

func backwardsCompatibilityError(version int, tag string) error {
	return fmt.Errorf("Backwards compatibility error on tag %s with version %d", tag, version)
}

func attributeNotSetError(tag string, attribute string) error {
	return fmt.Errorf("%s attribute %s must be set", tag, attribute)
}

//writeHeader sets the initial tags for both Media and Master Playlists files
func writeHeader(version int, buf *manifest.BufWrapper) {
	if version > 0 {
		buf.WriteString("#EXTM3U\n#EXT-X-VERSION:")
		buf.WriteString(strconv.Itoa(version))
		buf.WriteRune('\n')
		return
	}
	buf.Err = attributeNotSetError("Playlist", "Version")
}

//writeIndependentSegment sets the #EXT-X-INDEPENDENT-SEGMENTS tag on Media and Master Playlist file
func writeIndependentSegment(isIndSeg bool, buf *manifest.BufWrapper) {
	if isIndSeg {
		buf.WriteString("#EXT-X-INDEPENDENT-SEGMENTS\n")
	}
}

//writeStartPoint sets the #EXT-X-START tag on Media and Master Playlist file
func writeStartPoint(sp *StartPoint, buf *manifest.BufWrapper) {
	if sp != nil {
		buf.WriteString(fmt.Sprintf("#EXT-X-START:TIME-OFFSET=%s", strconv.FormatFloat(sp.TimeOffset, 'f', 3, 32)))
		buf.WriteValidString(sp.Precise, ",PRECISE=YES")
		buf.WriteRune('\n')
	}
}

//writeSessionData sets the EXT-X-SESSION-DATA tag on Master Playlist file
func (s *SessionData) writeSessionData(buf *manifest.BufWrapper) {
	if s != nil {
		if !buf.WriteValidString(s.DataID, fmt.Sprintf("#EXT-X-SESSION-DATA:DATA-ID=\"%s\"", s.DataID)) {
			buf.Err = attributeNotSetError("EXT-X-SESSION-DATA", "DATA-ID")
			return
		}

		if s.Value != "" && s.URI != "" {
			buf.Err = errors.New("EXT-X-SESSION-DATA must have attributes URI or VALUE, not both")
			return
		} else if s.Value != "" || s.URI != "" {
			buf.WriteValidString(s.Value, fmt.Sprintf(",VALUE=\"%s\"", s.Value))
			buf.WriteValidString(s.URI, fmt.Sprintf(",URI=\"%s\"", s.URI))
		} else {
			buf.Err = errors.New("EXT-X-SESSION-DATA must have either URI or VALUE attributes set")
			return
		}

		buf.WriteValidString(s.Language, fmt.Sprintf(",LANGUAGE=\"%s\"", s.Language))
		buf.WriteRune('\n')
	}
}

//writeXMedia sets the EXT-X-MEDIA tag on Master Playlist file
func (r *Rendition) writeXMedia(buf *manifest.BufWrapper) {
	if r != nil {

		if !isValidType(strings.ToUpper(r.Type)) || !buf.WriteValidString(r.Type, fmt.Sprintf("#EXT-X-MEDIA:TYPE=%s", r.Type)) {
			buf.Err = attributeNotSetError("EXT-X-MEDIA", "TYPE")
			return
		}
		if !buf.WriteValidString(r.GroupID, fmt.Sprintf(",GROUP-ID=\"%s\"", r.GroupID)) {
			buf.Err = attributeNotSetError("EXT-X-MEDIA", "GROUP-ID")
			return
		}
		if !buf.WriteValidString(r.Name, fmt.Sprintf(",NAME=\"%s\"", r.Name)) {
			buf.Err = attributeNotSetError("EXT-X-MEDIA", "NAME")
			return
		}
		buf.WriteValidString(r.Language, fmt.Sprintf(",LANGUAGE=\"%s\"", r.Language))
		buf.WriteValidString(r.AssocLanguage, fmt.Sprintf(",ASSOC-LANGUAGE=\"%s\"", r.AssocLanguage))
		buf.WriteValidString(r.Default, ",DEFAULT=YES")
		if r.Forced && strings.EqualFold(r.Type, sub) {
			buf.WriteValidString(r.Forced, ",FORCED=YES")
		}
		if strings.EqualFold(r.Type, cc) && isValidInstreamID(strings.ToUpper(r.InstreamID)) {
			buf.WriteValidString(r.InstreamID, fmt.Sprintf(",INSTREAM-ID=\"%s\"", r.InstreamID))
		}
		buf.WriteValidString(r.Characteristics, fmt.Sprintf(",CHARACTERISTICS=\"%s\"", r.Characteristics))

		//URI is required for SUBTITLES and MUST NOT be present for CLOSED-CAPTIONS, other types URI is optinal
		if strings.EqualFold(r.Type, sub) {
			if !buf.WriteValidString(r.URI, fmt.Sprintf(",URI=\"%s\"", r.URI)) {
				buf.Err = attributeNotSetError("EXT-X-MEDIA", "URI for SUBTITLES")
				return
			}
		} else if !strings.EqualFold(r.Type, cc) {
			buf.WriteValidString(r.URI, fmt.Sprintf(",URI=\"%s\"", r.URI))
		}

		buf.WriteRune('\n')
	}
}

//isValidType checks rendition Type is supported value (AUDIO, VIDEO, CLOSED-CAPTIONS or SUBTITLES)
func isValidType(t string) bool {
	return t == aud || t == vid || t == cc || t == sub
}

//isValidInstreamID checks rendition InstreamID is supported value
func isValidInstreamID(instream string) bool {
	return instream == "CC1" || instream == "CC2" || instream == "CC3" || instream == "CC4" || strings.HasPrefix(instream, "SERVICE")
}

//writeStreamInf sets the EXT-X-STREAM-INF or EXT-X-I-FRAME-STREAM-INF tag on Master Playlist file
func (v *Variant) writeStreamInf(version int, buf *manifest.BufWrapper) {
	if v != nil {
		if v.IsIframe {
			buf.WriteString("#EXT-X-I-FRAME-STREAM-INF:")
		} else {
			buf.WriteString("#EXT-X-STREAM-INF:")
		}

		if !buf.WriteValidString(v.Bandwidth, fmt.Sprintf("BANDWIDTH=%s", strconv.FormatInt(v.Bandwidth, 10))) {
			buf.Err = attributeNotSetError("Variant", "BANDWIDTH")
			return
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
		if !v.IsIframe {
			buf.WriteValidString(v.Audio, fmt.Sprintf(",AUDIO=\"%s\"", v.Audio))
			buf.WriteValidString(v.Subtitles, fmt.Sprintf(",SUBTITLES=\"%s\"", v.Subtitles))
			buf.WriteValidString(v.ClosedCaptions, fmt.Sprintf(",CLOSED-CAPTIONS=\"%s\"", v.ClosedCaptions))
			//If not IFrame, URI is in its own line
			buf.WriteString(fmt.Sprintf("\n%s\n", v.URI))
		} else {
			//If Iframe, URI is a param
			buf.WriteValidString(v.URI, fmt.Sprintf(",URI=\"%s\"\n", v.URI))
		}
	}
}

func (p *MediaPlaylist) writeTargetDuration(buf *manifest.BufWrapper) {
	if !buf.WriteValidString(p.TargetDuration, fmt.Sprintf("#EXT-X-TARGETDURATION:%s\n", strconv.Itoa(p.TargetDuration))) {
		buf.Err = attributeNotSetError("EXT-X-TARGETDURATION", "")
	}
}

func (p *MediaPlaylist) writeMediaSequence(buf *manifest.BufWrapper) {
	if p.MediaSequence > 0 {
		buf.WriteString(fmt.Sprintf("#EXT-X-MEDIA-SEQUENCE:%s\n", strconv.Itoa(p.MediaSequence)))
	}
}

func (p *MediaPlaylist) writeDiscontinuitySequence(buf *manifest.BufWrapper) {
	if p.DiscontinuitySequence > 0 {
		buf.WriteString(fmt.Sprintf("#EXT-X-DISCONTINUITY-SEQUENCE:%s\n", strconv.Itoa(p.DiscontinuitySequence)))
	}
}

func (p *MediaPlaylist) writeAllowCache(buf *manifest.BufWrapper) {
	if p.Version < 7 && p.AllowCache {
		buf.WriteString("#EXT-X-ALLOW-CACHE:YES\n")
	}
}

func (p *MediaPlaylist) writePlaylistType(buf *manifest.BufWrapper) {
	if p.Type != "" {
		buf.WriteString(fmt.Sprintf("#EXT-X-PLAYLIST-TYPE:%s\n", p.Type))
	}
}

func (p *MediaPlaylist) writeIFramesOnly(buf *manifest.BufWrapper) {
	if p.IFramesOnly {
		if p.Version < 4 {
			buf.Err = backwardsCompatibilityError(p.Version, "#EXT-X-I-FRAMES-ONLY")
			return
		}
		buf.WriteString("#EXT-X-I-FRAMES-ONLY\n")
	}
}

func (s *Segment) writeSegmentTags(buf *manifest.BufWrapper, previousSegment *Segment, version int) {
	if s != nil {
		for _, key := range s.Keys {

			found := false
			// If the previous segment we printed contains the same key, we shouldn't output it again
			if previousSegment != nil {
				for _, oldKey := range previousSegment.Keys {
					if key == oldKey {
						found = true
						break
					}
				}
			}

			if !found {
				key.writeKey(buf)
			}

			if buf.Err != nil {
				return
			}
		}

		if previousSegment == nil || previousSegment.Map == nil || (previousSegment.Map != s.Map) {
			s.Map.writeMap(buf)
		}

		if buf.Err != nil {
			return
		}

		if !s.ProgramDateTime.IsZero() {
			buf.WriteString(fmt.Sprintf("#EXT-X-PROGRAM-DATE-TIME:%s\n", s.ProgramDateTime.Format(time.RFC3339Nano)))
		}
		buf.WriteValidString(s.Discontinuity, "#EXT-X-DISCONTINUITY\n")

		s.DateRange.writeDateRange(buf)
		if buf.Err != nil {
			return
		}

		if s.Inf == nil {
			buf.Err = attributeNotSetError("EXTINF", "DURATION")
			return
		}

		if version < 3 {
			var duration int
			if s.Inf.Duration < 0.5 {
				duration = 0
			}
			// s.Inf.Duration is always > 0, so no need to use math.Abs or math.Copysign on the + 0.5
			duration = int(s.Inf.Duration + 0.5)
			buf.WriteString(fmt.Sprintf("#EXTINF:%d,%s\n", duration, s.Inf.Title))
		} else {
			buf.WriteString(fmt.Sprintf("#EXTINF:%s,%s\n", strconv.FormatFloat(s.Inf.Duration, 'f', 3, 32), s.Inf.Title))
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
			buf.Err = attributeNotSetError("Segment", "URI")
			return
		}
	}
}

func (k *Key) writeKey(buf *manifest.BufWrapper) {
	if k != nil {
		if k.IsSession {
			buf.WriteString("#EXT-X-SESSION-KEY:")
		} else {
			buf.WriteString("#EXT-X-KEY:")
		}

		if !isValidMethod(k.IsSession, strings.ToUpper(k.Method)) ||
			!buf.WriteValidString(k.Method, fmt.Sprintf("METHOD=%s", strings.ToUpper(k.Method))) {
			buf.Err = attributeNotSetError("KEY", "METHOD")
			return
		}
		if k.URI != "" && strings.ToUpper(k.Method) != none {
			buf.WriteValidString(k.URI, fmt.Sprintf(",URI=\"%s\"", k.URI))
		} else {
			buf.Err = attributeNotSetError("EXT-X-KEY", "URI")
			return
		}
		buf.WriteValidString(k.IV, fmt.Sprintf(",IV=%s", k.IV))
		buf.WriteValidString(k.Keyformat, fmt.Sprintf(",KEYFORMAT=\"%s\"", k.Keyformat))
		buf.WriteValidString(k.Keyformatversions, fmt.Sprintf(",KEYFORMATVERSIONS=\"%s\"", k.Keyformatversions))
		buf.WriteRune('\n')
	}
}

//isValidMethod checks Key Method value is supported. Session Key Method can't be NONE
func isValidMethod(isSession bool, method string) bool {
	return (method == aes || method == sample) || (!isSession && method == none)
}

func (m *Map) writeMap(buf *manifest.BufWrapper) {
	if m != nil {
		if !buf.WriteValidString(m.URI, fmt.Sprintf("#EXT-X-MAP:URI=\"%s\"", m.URI)) {
			buf.Err = attributeNotSetError("EXT-X-MAP", "URI")
			return
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
}

func (d *DateRange) writeDateRange(buf *manifest.BufWrapper) {
	if d != nil {
		if !buf.WriteValidString(d.ID, fmt.Sprintf("#EXT-X-DATERANGE:ID=%s", d.ID)) {
			buf.Err = attributeNotSetError("EXT-X-DATERANGE", "ID")
			return
		}
		buf.WriteValidString(d.Class, fmt.Sprintf(",CLASS=\"%s\"", d.Class))
		if !d.StartDate.IsZero() {
			buf.WriteString(fmt.Sprintf(",START-DATE=\"%s\"", d.StartDate.Format(time.RFC3339Nano)))
		} else {
			buf.Err = attributeNotSetError("EXT-X-DATERANGE", "START-DATE")
			return
		}
		if !d.EndDate.IsZero() {
			if d.EndDate.Before(d.StartDate) {
				buf.Err = errors.New("DateRange attribute EndDate must be equal or later than StartDate")
				return
			}
			buf.WriteString(fmt.Sprintf(",END-DATE=\"%s\"", d.EndDate.Format(time.RFC3339Nano)))
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
					buf.Err = errors.New("EXT-X-DATERANGE client-defined attributes must start with X-")
					return
				}
				buf.WriteString(",")
				buf.WriteString(strings.ToUpper(customTag))
			}
		}

		d.SCTE35.writeSCTE(buf)

		if buf.WriteValidString(d.EndOnNext, ",END-ON-NEXT=YES") {
			if d.Class == "" {
				buf.Err = errors.New("EXT-X-DATERANGE tag must have a CLASS attribute when END-ON-NEXT attribue is present")
				return
			}
			if d.Duration != nil || !d.EndDate.IsZero() {
				buf.Err = errors.New("EXT-X-DATERANGE tag must not have DURATION or END-DATE attributes when END-ON-NEXT attribute is present")
				return
			}
		}
		buf.WriteRune('\n')
	}
}

func (s *SCTE35) writeSCTE(buf *manifest.BufWrapper) {
	if s != nil {
		t := strings.ToUpper(s.Type)
		if t == "IN" || t == "OUT" || t == "CMD" {
			if !buf.WriteValidString(s.Value, fmt.Sprintf(",SCTE35-%s=%s", t, s.Value)) {
				buf.Err = attributeNotSetError("SCTE35", "Value")
			}
			return
		}
		buf.Err = errors.New("SCTE35 type must be IN, OUT or CMD")
	}
}

func (p *MediaPlaylist) writeEndList(buf *manifest.BufWrapper) {
	if p.EndList {
		buf.WriteString("#EXT-X-ENDLIST\n")
	}
}

//checkCompatibility checks backwards compatibility issues according to the Media Playlist version
func (p *MediaPlaylist) checkCompatibility(s *Segment) error {
	if s != nil {
		if s.Inf != nil && p.Version < 3 {
			if s.Inf.Duration != float64(int64(s.Inf.Duration)) {
				return backwardsCompatibilityError(p.Version, "#EXTINF")
			}
		}

		if s.Byterange != nil && p.Version < 4 {
			return backwardsCompatibilityError(p.Version, "#EXT-X-BYTERANGE")
		}

		for _, key := range s.Keys {
			if key.IV != "" && p.Version < 2 {
				return backwardsCompatibilityError(p.Version, "#EXT-X-KEY")
			}

			if (key.Keyformat != "" || key.Keyformatversions != "") && p.Version < 5 {
				return backwardsCompatibilityError(p.Version, "#EXT-X-KEY")
			}
		}

		if s.Map != nil {
			if p.Version < 5 || (!p.IFramesOnly && p.Version < 6) {
				return backwardsCompatibilityError(p.Version, "#EXT-X-MAP")
			}
		}
	} else {
		if p.IFramesOnly && p.Version < 4 {
			return backwardsCompatibilityError(p.Version, "#EXT-X-I-FRAMES-ONLY")
		}
	}

	return nil
}

func (p *MasterPlaylist) checkCompatibility() error {
	switch {
	case p.Version < 7:
		for _, rendition := range p.Renditions {
			if rendition.Type == cc {
				if strings.HasPrefix(rendition.InstreamID, "SERVICE") {
					return backwardsCompatibilityError(p.Version, "#EXT-X-MEDIA")
				}
			}
		}
	}

	return nil
}

//TODO:(sliding window) - MediaPlaylist constructor receiving sliding window size. In the case of sliding window playlist, we
//must not include a EXT-X-PLAYLIST-TYPE tag since EVENT and VOD don't support removing segments from the manifest.
//TODO:(sliding window/live streaming) - Public method to add segment to a MediaPlaylist. This method would need helper methods, to check
//sliding window size and remove segment when necessary. If playlist type is EVENT, only adds without removing.
//Also helper methods to update MediaSequence and DiscontinuitySequence values
//TODO:(sliding window/live streaming) - Public method to insert EXT-X-ENDLIST tag when EVENT or sliding window playlist reaches its end
//TODO:(sliding window) - Figure out a way to control tags like KEY, MAP etc, that can be applicable to following segments.
//What to do when that segment is removed (tags would in theory be removed with it)? Add methods to slide these tags along with the window.
