package hls

import (
	"bytes"
	"errors"
	"io"
	"sort"

	"github.com/ingest/manifest"
)

//Encode writes a Master Playlist file
func (p *MasterPlaylist) Encode() (io.Reader, error) {
	if err := p.checkCompatibility(); err != nil {
		return nil, err
	}

	buf := manifest.NewBufWrapper()

	//Write header tags
	writeHeader(p.Version, buf)
	if buf.Err != nil {
		return nil, buf.Err
	}

	//Write Session Data tags if enabled
	if p.SessionData != nil {
		for _, sd := range p.SessionData {
			sd.writeSessionData(buf)
			if buf.Err != nil {
				return nil, buf.Err
			}
		}
	}
	//write session keys tags if enabled
	if p.SessionKeys != nil {
		for _, sk := range p.SessionKeys {
			sk.writeKey(buf)
			if buf.Err != nil {
				return nil, buf.Err
			}
		}
	}

	//Write Independent Segments tag if enabled
	writeIndependentSegment(p.IndependentSegments, buf)

	//write Start tag if enabled
	writeStartPoint(p.StartPoint, buf)
	if buf.Err != nil {
		return nil, buf.Err
	}

	// For every Rendition, write #EXT-X-MEDIA tags
	if p.Renditions != nil {
		for _, rendition := range p.Renditions {
			rendition.writeXMedia(buf)
			if buf.Err != nil {
				return nil, buf.Err
			}
		}
	}

	//For every Variant, write #EXT-X-STREAM-INF and #EXT-X-I-FRAME-STREAM-INF tags
	if p.Variants != nil {
		for _, variant := range p.Variants {
			variant.writeStreamInf(p.Version, buf)
			if buf.Err != nil {
				return nil, buf.Err
			}
		}
	}

	return bytes.NewReader(buf.Buf.Bytes()), buf.Err
}

//Encode writes a Media Playlist file
func (p *MediaPlaylist) Encode() (io.Reader, error) {
	buf := manifest.NewBufWrapper()

	//write header tags
	writeHeader(p.Version, buf)
	if buf.Err != nil {
		return nil, buf.Err
	}
	//write Target Duration tag
	p.writeTargetDuration(buf)
	if buf.Err != nil {
		return nil, buf.Err
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
	if buf.Err != nil {
		return nil, buf.Err
	}

	//write segment tags
	if p.Segments != nil {
		sort.Sort(p.Segments)
		var prev *Segment
		for _, segment := range p.Segments {
			if err := p.checkCompatibility(segment); err != nil {
				return nil, err
			}
			segment.writeSegmentTags(buf, prev, p.Version)
			if buf.Err != nil {
				return nil, buf.Err
			}
			prev = segment
		}
	} else {
		return nil, errors.New("MediaPlaylist must have at least one Segment")
	}
	//write End List tag if enabled
	p.writeEndList(buf)

	return bytes.NewReader(buf.Buf.Bytes()), buf.Err
}
