package hls

import (
	"bytes"
	"strings"
)

//GenerateManifest writes a Master Playlist file
func (p *MasterPlaylist) GenerateManifest() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	//Write header tags
	writeHeader(p.Version, buf)

	//Write Session Data tags if enabled
	if p.SessionData != nil {
		for _, sd := range p.SessionData {
			if err := sd.writeSessionData(buf); err != nil {
				return buf, err
			}
		}
	}

	//Write Independent Segments tag if enabled
	writeIndependentSegment(p.IndependentSegments, buf)

	//write Start tag if enabled
	if err := writeStartPoint(p.StartPoint, buf); err != nil {
		return buf, err
	}

	//For every Variant, write rendition and variant tags if enabled
	if p.Variants != nil {
		for _, variant := range p.Variants {
			if variant.Renditions != nil {
				for _, rendition := range variant.Renditions {
					//Check backwards compatibility issue before continuing
					if strings.HasPrefix(rendition.InstreamID, "SERVICE") && p.Version < 7 {
						return buf, backwardsCompatibilityError(p.Version, "#EXT-X-MEDIA")
					}
					rendition.writeXMedia(buf)
				}
			}
			variant.writeStreamInf(p.Version, buf)
		}
	}

	return buf, nil
}

//GenerateManifest writes a Media Playlist file
func (p *MediaPlaylist) GenerateManifest() (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	//write header tags
	writeHeader(p.Version, buf)
	//write Target Duration tag
	p.writeTargetDuration(buf)
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
	//write I-Frames Only tag if enabled and version > = 4
	if p.IFramesOnly && p.Version < 4 {
		return buf, backwardsCompatibilityError(p.Version, "#EXT-X-I-FRAMES-ONLY")
	}
	p.writeIFramesOnly(buf)

	if p.Segments != nil {
		for _, segment := range p.Segments {
			if err := p.checkCompatibility(segment); err != nil {
				return buf, err
			}
			segment.writeSegmentTags(buf)
		}
	}
	//write End List tag if enabled
	p.writeEndList(buf)

	return buf, nil
}
