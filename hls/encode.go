package hls

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

//GenerateManifest writes a Master Playlist file
func (p *MasterPlaylist) GenerateManifest() (io.Reader, error) {
	buf := NewBufWriter()

	//Write header tags
	if err := writeHeader(p.Version, buf); err != nil {
		return nil, err
	}

	//Write Session Data tags if enabled
	if p.SessionData != nil {
		for _, sd := range p.SessionData {
			if err := sd.writeSessionData(buf); err != nil {
				return nil, err
			}
		}
	}
	//write session keys tags if enabled
	if p.SessionKeys != nil {
		for _, sk := range p.SessionKeys {
			if err := sk.writeKey(buf); err != nil {
				return nil, err
			}
		}
	}

	//Write Independent Segments tag if enabled
	writeIndependentSegment(p.IndependentSegments, buf)

	//write Start tag if enabled
	if err := writeStartPoint(p.StartPoint, buf); err != nil {
		return nil, err
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
					if err := rendition.writeXMedia(buf); err != nil {
						return nil, err
					}
				}
			}
			if err := variant.writeStreamInf(p.Version, buf); err != nil {
				return nil, err
			}
		}
	}

	return bytes.NewReader(buf.buf.Bytes()), buf.err
}

//GenerateManifest writes a Media Playlist file
func (p *MediaPlaylist) GenerateManifest() (io.Reader, error) {
	buf := NewBufWriter()

	//write header tags
	if err := writeHeader(p.Version, buf); err != nil {
		return nil, err
	}
	//write Target Duration tag
	if err := p.writeTargetDuration(buf); err != nil {
		return nil, err
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
	if err := p.writeIFramesOnly(buf); err != nil {
		return nil, err
	}

	//write segment tags
	if p.Segments != nil {
		for _, segment := range p.Segments {
			if err := p.checkCompatibility(segment); err != nil {
				return nil, err
			}
			if err := segment.writeSegmentTags(buf); err != nil {
				return nil, err
			}
		}
	} else {
		return nil, errors.New("MediaPlaylist must have at least one Segment")
	}
	//write End List tag if enabled
	p.writeEndList(buf)

	return bytes.NewReader(buf.buf.Bytes()), buf.err
}
