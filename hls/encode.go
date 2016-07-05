package hls

import (
	"bytes"
	"errors"
	"io"
	"sort"
	"strings"
)

//Encode writes a Master Playlist file
func (p *MasterPlaylist) Encode() (io.Reader, error) {
	buf := NewBufWrapper()

	//Write header tags
	writeHeader(p.Version, buf)
	if buf.err != nil {
		return nil, buf.err
	}

	//Write Session Data tags if enabled
	if p.SessionData != nil {
		for _, sd := range p.SessionData {
			sd.writeSessionData(buf)
			if buf.err != nil {
				return nil, buf.err
			}
		}
	}
	//write session keys tags if enabled
	if p.SessionKeys != nil {
		for _, sk := range p.SessionKeys {
			sk.writeKey(buf)
			if buf.err != nil {
				return nil, buf.err
			}
		}
	}

	//Write Independent Segments tag if enabled
	writeIndependentSegment(p.IndependentSegments, buf)

	//write Start tag if enabled
	writeStartPoint(p.StartPoint, buf)
	if buf.err != nil {
		return nil, buf.err
	}

	//For every Variant, write rendition and variant tags if enabled
	if p.Variants != nil {
		for _, variant := range p.Variants {
			if variant.Renditions != nil {
				for _, rendition := range variant.Renditions {
					//Check backwards compatibility issue before continuing
					if strings.HasPrefix(strings.ToUpper(rendition.InstreamID), "SERVICE") && p.Version < 7 {
						return nil, backwardsCompatibilityError(p.Version, "#EXT-X-MEDIA")
					}
					rendition.writeXMedia(buf)
					if buf.err != nil {
						return nil, buf.err
					}
				}
			}
			variant.writeStreamInf(p.Version, buf)
			if buf.err != nil {
				return nil, buf.err
			}
		}
	}

	return bytes.NewReader(buf.buf.Bytes()), buf.err
}

//Encode writes a Media Playlist file
func (p *MediaPlaylist) Encode() (io.Reader, error) {
	buf := NewBufWrapper()

	//write header tags
	writeHeader(p.Version, buf)
	if buf.err != nil {
		return nil, buf.err
	}
	//write Target Duration tag
	p.writeTargetDuration(buf)
	if buf.err != nil {
		return nil, buf.err
	}
	//write Media Sequence tag if enabled
	p.writeMediaSequence(buf)
	//write Independent Segment tag if enabled
	writeIndependentSegment(p.IndependentSegments, buf)
	//write Start tag if enabled
	writeStartPoint(p.StartPoint, buf)
	//write Discontinuity Sequence tag if enabled
	p.writeDiscontinuitySequence(buf)
	//write Playlist Type tag if enabled
	p.writePlaylistType(buf)
	//write Allow Cache tag if enabled
	p.writeAllowCache(buf)
	//write I-Frames Only if enabled
	p.writeIFramesOnly(buf)
	if buf.err != nil {
		return nil, buf.err
	}

	//write segment tags
	if p.Segments != nil {
		sort.Sort(p.Segments)
		for _, segment := range p.Segments {
			if err := p.checkCompatibility(segment); err != nil {
				return nil, err
			}
			segment.writeSegmentTags(buf)
			if buf.err != nil {
				return nil, buf.err
			}
		}
	} else {
		return nil, errors.New("MediaPlaylist must have at least one Segment")
	}
	//write End List tag if enabled
	p.writeEndList(buf)

	return bytes.NewReader(buf.buf.Bytes()), buf.err
}
