package hls

import (
	"io"
	"strconv"
	"strings"

	"github.com/ingest/manifest"
)

//Parse reads a Master Playlist file and converts it to a MasterPlaylist object
func (p *MasterPlaylist) Parse(reader io.Reader) error {
	buf := manifest.NewBufWrapper()
	buf.ReadFrom(reader)
	if buf.Err != nil {
		return buf.Err
	}
	var eof bool
	var streamInfLastTag bool
	var line string
	var renditions []*Rendition
	key := &Key{}
	variant := &Variant{}
	r := &Rendition{}

	//until EOF, read every line and decode into an object
	for !eof {
		line = buf.ReadString('\n')
		if buf.Err != nil {
			if buf.Err == io.EOF {
				eof = true
			} else {
				break
			}
		}

		line = strings.TrimSpace(line)
		size := len(line)
		//if empty line, skip
		if size <= 1 {
			continue
		}

		if line[0] == '#' {
			streamInfLastTag = false
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
				key = decodeKey(line[index+1:size], true)
				p.SessionKeys = append(p.SessionKeys, key)

			case line[0:index] == "#EXT-X-SESSION-DATA":
				data := decodeSessionData(line[index+1 : size])
				p.SessionData = append(p.SessionData, data)

			case line[0:index] == "#EXT-X-MEDIA":
				r = decodeRendition(line[index+1 : size])
				renditions = append(renditions, r)

			case line[0:index] == "#EXT-X-STREAM-INF":
				variant, buf.Err = decodeVariant(line[index+1:size], false)
				streamInfLastTag = true

			//Case line is EXT-X-I-FRAME-STREAM-INF, it means it's the end of a variant
			//append variant to MasterPlaylist and restart variables
			case line[0:index] == "#EXT-X-I-FRAME-STREAM-INF":
				if variant, buf.Err = decodeVariant(line[index+1:size], true); buf.Err == nil {
					variant.Renditions = renditions
					p.Variants = append(p.Variants, variant)
					variant = &Variant{}
					renditions = []*Rendition{}
				}
			}
			//Case line doesn't start with '#', check if last tag was EXT-X-STREAM-INF.
			//Which means this line is variant URI
			//Append variant to MasterPlaylist and restart variables
		} else {
			if streamInfLastTag {
				variant.URI = line
				variant.Renditions = renditions
				p.Variants = append(p.Variants, variant)
				variant = &Variant{}
				renditions = []*Rendition{}
				streamInfLastTag = false
			}
		}

	}

	return buf.Err
}

//Parse reads a Media Playlist file and convert it to MediaPlaylist object
func (p *MediaPlaylist) Parse(reader io.Reader) error {
	buf := manifest.NewBufWrapper()
	buf.ReadFrom(reader)
	if buf.Err != nil {
		return buf.Err
	}
	var eof bool
	var line string
	//count indicates the segment sequence number
	count := 0
	key := &Key{}
	segment := &Segment{}

	//Until EOF, read every line and decode into an object
	for !eof {
		line = buf.ReadString('\n')
		if buf.Err != nil {
			if buf.Err == io.EOF {
				eof = true
			} else {
				break
			}
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
			count = p.MediaSequence
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
		//case below this point refers to segment tags, if line is uri, it reached the end of a segment.
		//append segment to p.Segments and restart variable
		case line[0:index] == "#EXT-X-KEY":
			key = decodeKey(line[index+1:size], false)
			segment.Keys = append(segment.Keys, key)
		case line[0:index] == "#EXT-X-MAP":
			segment.Map, buf.Err = decodeMap(line[index+1 : size])
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
			segment.ID = count
			p.Segments = append(p.Segments, segment)
			segment = &Segment{}
			count++
		}
	}

	return buf.Err
}
