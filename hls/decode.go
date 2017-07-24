package hls

import (
	"io"
	"strconv"
	"strings"

	"github.com/ingest/manifest"
)

type masterPlaylistParseState struct {
	eof              bool
	streamInfLastTag bool
	variant          *Variant
}

//Parse reads a Master Playlist file and converts it to a MasterPlaylist object
func (p *MasterPlaylist) Parse(reader io.Reader) error {
	buf := manifest.NewBufWrapper()

	// Populate buffer into memory, we could change this to a tokenizer with a line-by-line scanner?
	if _, err := buf.ReadFrom(reader); err != nil {
		return err
	}

	s := masterPlaylistParseState{
		variant: &Variant{masterPlaylist: p},
	}

	// Raw line
	var line string

	// Runs until io.EOF, reads line-by-line from buffer and decode into an object
	for !s.eof {
		line = buf.ReadString('\n')
		if buf.Err == io.EOF {
			buf.Err = nil
			s.eof = true
		}

		if buf.Err != nil {
			return buf.Err
		}

		line = strings.TrimSpace(line)
		size := len(line)
		//if empty line, skip
		if size <= 1 {
			continue
		}

		if line[0] == '#' {
			s.streamInfLastTag = false
			index := stringsIndex(line, ":")
			switch {
			case line == "#EXTM3U":
				p.M3U = true

			case line[0:index] == "#EXT-X-VERSION":
				p.Version, buf.Err = strconv.Atoi(line[index+1 : size])

			case line[0:index] == "#EXT-X-START":
				p.StartPoint, buf.Err = decodeStartPoint(line[index+1 : size])

			case line == "#EXT-X-INDEPENDENT-SEGMENTS":
				p.IndependentSegments = true

			case line[0:index] == "#EXT-X-SESSION-KEY":
				key := decodeKey(line[index+1:size], true)
				key.masterPlaylist = p
				p.SessionKeys = append(p.SessionKeys, key)

			case line[0:index] == "#EXT-X-SESSION-DATA":
				data := decodeSessionData(line[index+1 : size])
				data.masterPlaylist = p
				p.SessionData = append(p.SessionData, data)

			case line[0:index] == "#EXT-X-MEDIA":
				r := decodeRendition(line[index+1 : size])
				r.masterPlaylist = p
				p.Renditions = append(p.Renditions, r)

			case line[0:index] == "#EXT-X-STREAM-INF":
				s.variant, buf.Err = decodeVariant(line[index+1:size], false)
				s.variant.masterPlaylist = p
				s.streamInfLastTag = true

			//Case line is EXT-X-I-FRAME-STREAM-INF, it means it's the end of a variant
			//append variant to MasterPlaylist and restart variables
			case line[0:index] == "#EXT-X-I-FRAME-STREAM-INF":
				variant, err := decodeVariant(line[index+1:size], true)
				if err != nil {
					buf.Err = err
					continue // shouldn't include a partially decoded iframe playlist
				}
				variant.masterPlaylist = p

				p.Variants = append(p.Variants, variant)
			}
			//Case line doesn't start with '#', check if last tag was EXT-X-STREAM-INF.
			//Which means this line is variant URI
			//Append variant to MasterPlaylist and restart variables
		} else if s.streamInfLastTag {
			s.variant.URI = line
			p.Variants = append(p.Variants, s.variant)
			// Reset state
			s.variant = &Variant{masterPlaylist: p}
			s.streamInfLastTag = false
		}

	}

	if buf.Err != nil {
		return buf.Err
	}

	// Check master playlist compatibility
	if err := p.checkCompatibility(); err != nil {
		return err
	}

	return nil
}

// Holds the state while parsing a media playlist
type mediaPlaylistParseState struct {
	eof             bool
	previousMap     *Map
	previousKey     *Key
	segmentSequence int
}

//Parse reads a Media Playlist file and convert it to MediaPlaylist object
func (p *MediaPlaylist) Parse(reader io.Reader) error {
	buf := manifest.NewBufWrapper()

	// Populate buffer into memory, we could change this to a tokenizer with a line-by-line scanner?
	if _, err := buf.ReadFrom(reader); err != nil {
		return err
	}

	s := mediaPlaylistParseState{}
	segment := &Segment{
		mediaPlaylist: p,
	}

	//Until EOF, read every line and decode into an object
	var line string
	for !s.eof {
		line = buf.ReadString('\n')
		if buf.Err == io.EOF {
			buf.Err = nil
			s.eof = true
		}

		if buf.Err != nil {
			return buf.Err
		}

		line = strings.TrimSpace(line)
		size := len(line)
		//if empty line, skip
		if size <= 1 {
			continue
		}

		index := stringsIndex(line, ":")

		switch {
		case line[0:index] == "#EXT-X-VERSION":
			p.Version, buf.Err = strconv.Atoi(line[index+1 : size])
		case line[0:index] == "#EXT-X-TARGETDURATION":
			p.TargetDuration, buf.Err = strconv.Atoi(line[index+1 : size])
		case line[0:index] == "#EXT-X-MEDIA-SEQUENCE":
			p.MediaSequence, buf.Err = strconv.Atoi(line[index+1 : size])
			//case MediaSequence is present, first sequence number = MediaSequence
			s.segmentSequence = p.MediaSequence
		case line[0:index] == "#EXT-X-DISCONTINUITY-SEQUENCE":
			p.DiscontinuitySequence, buf.Err = strconv.Atoi(line[index+1 : size])
		case line == "#EXT-X-I-FRAMES-ONLY":
			p.IFramesOnly = true
		case line[0:index] == "#EXT-X-ALLOW-CACHE":
			if line[index+1:size] == boolYes {
				p.AllowCache = true
			}
		case line == "#EXT-X-INDEPENDENT-SEGMENTS":
			p.IndependentSegments = true
		case line[0:index] == "#EXT-X-PLAYLIST-TYPE":
			if strings.EqualFold(line[index+1:size], "VOD") || strings.EqualFold(line[index+1:size], "EVENT") {
				p.Type = line[index+1 : size]
			}
		case line == "#EXT-X-ENDLIST":
			p.EndList = true
		case line[0:index] == "#EXT-X-START":
			p.StartPoint, buf.Err = decodeStartPoint(line[index+1 : size])

		// Cases below this point refers to tags that effect segments, when we reach a line with no leading #, we've reached the end of a segment definition.
		case line[0:index] == "#EXT-X-KEY":
			key := decodeKey(line[index+1:size], false)
			key.mediaPlaylist = p
			s.previousKey = key // we store this key for future reference because every segment between EXT-X-KEYs should use this key for decryption
			segment.Keys = append(segment.Keys, key)
		case line[0:index] == "#EXT-X-MAP":
			s.previousMap, buf.Err = decodeMap(line[index+1 : size])
			s.previousMap.mediaPlaylist = p
			segment.Map = s.previousMap
		case line[0:index] == "#EXT-X-PROGRAM-DATE-TIME":
			segment.ProgramDateTime, buf.Err = decodeDateTime(line[index+1 : size])
		case line[0:index] == "#EXT-X-DATERANGE":
			segment.DateRange, buf.Err = decodeDateRange(line[index+1 : size])
		case line[0:index] == "#EXT-X-BYTERANGE":
			segment.Byterange, buf.Err = decodeByterange(line[index+1 : size])
		case line[0:index] == "#EXTINF":
			segment.Inf, buf.Err = decodeInf(line[index+1 : size])
		case !strings.HasPrefix(line, "#"):
			segment.URI = line
			segment.ID = s.segmentSequence

			// a previous EXT-X-KEY applies to this segment
			if len(segment.Keys) == 0 && s.previousKey != nil && s.previousKey.URI != "" {
				segment.Keys = append(segment.Keys, s.previousKey)
			}

			// a previous EXT-X-MAP applies to this segment
			if segment.Map == nil && s.previousMap != nil {
				segment.Map = s.previousMap
			}

			p.Segments = append(p.Segments, segment)

			// Reset segment
			segment = &Segment{mediaPlaylist: p}
			s.segmentSequence++
		}
	}

	if buf.Err != nil {
		return buf.Err
	}

	// Check media playlist compatibility
	if err := p.checkCompatibility(nil); err != nil {
		return err
	}

	for _, segment := range p.Segments {
		if err := p.checkCompatibility(segment); err != nil {
			return err
		}
	}

	return nil
}
