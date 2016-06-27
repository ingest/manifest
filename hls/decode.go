package hls

import (
	"io"
	"strconv"
	"strings"
)

//ReadManifest reads a Master Playlist file and converts it to a MasterPlaylist object
func (p *MasterPlaylist) ReadManifest(reader io.Reader) error {
	buf := NewBufWrapper()
	buf.ReadFrom(reader)
	if buf.err != nil {
		return buf.err
	}
	var eof bool
	var line string
	key := &Key{}
	variant := &Variant{}
	r := &Rendition{}
	var renditions []*Rendition

	for !eof {
		line = buf.ReadString('\n')
		if buf.err != nil {
			if buf.err == io.EOF {
				eof = true
			} else {
				break
			}
		}
		if len(line) <= 1 {
			continue
		}
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "#EXTM3U"):
			p.M3U = true

		case strings.HasPrefix(line, "#EXT-X-VERSION"):
			if version := stringAfter(line, ":"); version != "" {
				p.Version, buf.err = strconv.Atoi(version)
			}

		case strings.HasPrefix(line, "#EXT-X-START"):
			p.StartPoint, buf.err = decodeStartPoint(stringAfter(line, ":"))

		case strings.HasPrefix(line, "#EXT-X-INDEPENDENT-SEGMENTS"):
			p.IndependentSegments = true

		case strings.HasPrefix(line, "#EXT-X-SESSION-KEY"):
			key = decodeKey(stringAfter(line, ":"))
			p.SessionKeys = append(p.SessionKeys, key)

		case strings.HasPrefix(line, "#EXT-X-SESSION-DATA"):
			data := decodeSessionData(stringAfter(line, ":"))
			p.SessionData = append(p.SessionData, data)

		case strings.HasPrefix(line, "#EXT-X-MEDIA"):
			r = decodeRendition(stringAfter(line, ":"))
			renditions = append(renditions, r)

		//Case line is playlist uri, tag before is EXT-X-STREAM-INF.
		//Append variant to MasterPlaylist and restart variables
		case !strings.HasPrefix(line, "#"):
			variant.URI = line
			variant.Renditions = renditions
			p.Variants = append(p.Variants, variant)
			variant = &Variant{}
			renditions = []*Rendition{}

		case strings.HasPrefix(line, "#EXT-X-STREAM-INF"):
			variant, buf.err = decodeVariant(stringAfter(line, ":"), false)

		//case line is EXT-X-I-FRAME-STREAM-INF, append variant to MasterPlaylist and restart variables
		case strings.HasPrefix(line, "#EXT-X-I-FRAME-STREAM-INF"):
			if variant, buf.err = decodeVariant(stringAfter(line, ":"), true); buf.err == nil {
				variant.Renditions = renditions
				p.Variants = append(p.Variants, variant)
				variant = &Variant{}
				renditions = []*Rendition{}
			}
		}

	}

	return buf.err
}

//ReadManifest reads a Media Playlist file and convert it to MediaPlaylist object
func (p *MediaPlaylist) ReadManifest(reader io.Reader) error {
	buf := NewBufWrapper()
	buf.ReadFrom(reader)
	if buf.err != nil {
		return buf.err
	}
	var eof bool
	var line string
	key := &Key{}
	segment := &Segment{}

	for !eof {
		line = buf.ReadString('\n')
		if buf.err != nil {
			if buf.err == io.EOF {
				eof = true
			} else {
				break
			}
		}

		if len(line) <= 1 {
			continue
		}
		line = strings.TrimSpace(line)

		switch {
		case strings.HasPrefix(line, "#EXTM3U"):
			p.M3U = true
		case strings.HasPrefix(line, "#EXT-X-VERSION"):
			if version := stringAfter(line, ":"); version != "" {
				p.Version, buf.err = strconv.Atoi(version)
			}
		case strings.HasPrefix(line, "#EXT-X-TARGETDURATION"):
			if duration := stringAfter(line, ":"); duration != "" {
				p.TargetDuration, buf.err = strconv.Atoi(duration)
			}
		case strings.HasPrefix(line, "#EXT-X-MEDIA-SEQUENCE"):
			if sequence := stringAfter(line, ":"); sequence != "" {
				p.MediaSequence, buf.err = strconv.Atoi(sequence)
			}
		case strings.HasPrefix(line, "#EXT-X-DISCONTINUITY-SEQUENCE"):
			if disc := stringAfter(line, ":"); disc != "" {
				p.DiscontinuitySequence, buf.err = strconv.Atoi(disc)
			}
		case strings.HasPrefix(line, "#EXT-X-I-FRAMES-ONLY"):
			p.IFramesOnly = true
		case strings.HasPrefix(line, "#EXT-X-ALLOW-CACHE"):
			if allow := stringAfter(line, ":"); allow == boolYes {
				p.AllowCache = true
			}
		case strings.HasPrefix(line, "#EXT-X-INDEPENDENT-SEGMENTS"):
			p.IndependentSegments = true
		case strings.HasPrefix(line, "#EXT-X-PLAYLIST-TYPE"):
			if t := stringAfter(line, ":"); t == "VOD" || t == "EVENT" {
				p.Type = t
			}
		case strings.HasPrefix(line, "#EXT-X-ENDLIST"):
			p.EndList = true
		case strings.HasPrefix(line, "#EXT-X-START"):
			p.StartPoint, buf.err = decodeStartPoint(stringAfter(line, ":"))

			//check segment tags, if line is uri, append segment to p.Segments and restart segment
		case strings.HasPrefix(line, "#EXT-X-KEY"):
			key = decodeKey(stringAfter(line, ":"))
			segment.Keys = append(segment.Keys, key)
		case strings.HasPrefix(line, "#EXT-X-MAP"):
			segment.Map, buf.err = decodeMap(stringAfter(line, ":"))
		case strings.HasPrefix(line, "#EXT-X-PROGRAM-DATE-TIME"):
			segment.ProgramDateTime, buf.err = decodeDateTime(stringAfter(line, ":"))
		case strings.HasPrefix(line, "#EXT-X-DATERANGE"):
			segment.DateRange, buf.err = decodeDateRange(stringAfter(line, ":"))
		case strings.HasPrefix(line, "#EXT-X-BYTERANGE"):
			segment.Byterange, buf.err = decodeByterange(stringAfter(line, ":"))
		case strings.HasPrefix(line, "#EXTINF"):
			segment.Inf, buf.err = decodeInf(stringAfter(line, ":"))
		case !strings.HasPrefix(line, "#"):
			segment.URI = line
			p.Segments = append(p.Segments, segment)
			segment = &Segment{}
		}
	}

	return buf.err
}
