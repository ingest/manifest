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
		attributes := make([]string, 0, 2)
		if sp.TimeOffset == float64(0) {
			return errors.New("StartPoint must have attribute TIME-OFFSET set.")
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
		attributes := make([]string, 0, 4)

		if s.DataID != "" {
			attributes = append(attributes, fmt.Sprintf("DATA-ID=\"%s\"", s.DataID))
		} else {
			return errors.New("SessionData attribute DATA-ID must be set.")
		}

		if s.Value != "" && s.URI != "" {
			return errors.New("SessionData must have attributes URI or VALUE, not both.")
		} else if s.Value != "" {
			attributes = append(attributes, fmt.Sprintf("VALUE=\"%s\"", s.Value))
		} else if s.URI != "" {
			attributes = append(attributes, fmt.Sprintf("URI=\"%s\"", s.URI))
		} else {
			return errors.New("SessionData must have either URI or VALUE attributes set.")
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
		attributes := make([]string, 0, 11)

		if r.Type != "" {
			attributes = append(attributes, fmt.Sprintf("TYPE=%s", r.Type))
		} else {
			return errors.New("Rendition must have attribute TYPE set.")
		}
		if r.GroupID != "" {
			attributes = append(attributes, fmt.Sprintf("GROUP-ID=\"%s\"", r.GroupID))
		} else {
			return errors.New("Rendition must have attribute GROUP-ID set.")
		}
		if r.Name != "" {
			attributes = append(attributes, fmt.Sprintf("NAME=\"%s\"", r.Name))
		} else {
			return errors.New("Rendition must have attribute NAME set.")
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
		attributes := make([]string, 0, 11)

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

func (p *MediaPlaylist) writeTargetDuration(buf *bytes.Buffer) {
	if p.TargetDuration > 0 {
		buf.WriteString(fmt.Sprintf("#EXT-X-TARGETDURATION:%s\n", strconv.Itoa(p.TargetDuration)))
	}
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
		buf.WriteString("#EXT-X-I-FRAMES-ONLY")
	}
}

func (s *Segment) writeSegmentTags(buf *bytes.Buffer) error {
	if s != nil {
		if s.Inf != nil && s.Inf.Duration > float64(0) {
			buf.WriteString(fmt.Sprintf("#EXTINF:%s,%s", strconv.FormatFloat(s.Inf.Duration, 'f', 3, 32), s.Inf.Title))
		} else {
			return errors.New("Segment must have EXTINF duration set")
		}

	}
	return nil
}

func (p *MediaPlaylist) writeEndList(buf *bytes.Buffer) {
	if p.EndList {
		buf.WriteString("#EXT-X-ENDLIST")
	}
}

//checkCompatibility checks backwards compatibility issues according to the Media Playlist version
func (p *MediaPlaylist) checkCompatibility(s *Segment) error {
	if s != nil {
		if s.Key != nil {
			if (s.Key.Method == "SAMPLE-AES" || s.Key.Keyformat != "" || s.Key.Keyformatversions != "") && p.Version < 5 {
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
