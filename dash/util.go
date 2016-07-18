package dash

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"time"
)

//NewMPD initiates a MPD struct with the minimum required attributes
func NewMPD(profile string, minBufferTime time.Duration) *MPD {
	return &MPD{XMLNS: DashNS,
		Type:          "static",
		Profiles:      profile,
		MinBufferTime: &CustomDuration{Duration: minBufferTime},
	}
}

func (m *MPD) validate() error {
	buf := new(bytes.Buffer)
	m.validateMPD(buf)
	if buf != nil && buf.String() != "" {
		return errors.New(buf.String())
	}
	return nil
}

func (m *MPD) validateMPD(buf *bytes.Buffer) {
	if m != nil {

		if m.Profiles == "" {
			buf.WriteString("MPD field Profiles is required.\n")
		}

		if !strings.EqualFold(m.Type, "static") && !strings.EqualFold(m.Type, "dynamic") {
			buf.WriteString("Possible values for MPD field Type are 'static' and 'dynamic'.\n")
		}

		if strings.EqualFold(m.Type, "dynamic") {
			if m.AvStartTime == nil {
				buf.WriteString("MPD field AvStartTime must be present when Type = 'dynamic'.\n")
			}
			if m.PublishTime == nil {
				buf.WriteString("MPD field PublishTime must be present when Type = 'dynamic'\n")
			}
		} else {
			if m.MinUpdatePeriod != nil {
				buf.WriteString("MPD field MinUpdatePeriod must not be present when Type = 'static'\n")
			}
		}

		if m.MinBufferTime == nil {
			buf.WriteString("MPD field MinBufferTime is required.\n")
		}

		if m.Metrics != nil {
			for _, metric := range m.Metrics {
				metric.validate(buf)
			}
		}

		if m.Periods != nil {
			for _, period := range m.Periods {
				period.validate(buf, m.Type)
			}
		} else {
			buf.WriteString("MPD must have at least one Period element.\n")
		}
	}
}

func (p *Period) validate(buf *bytes.Buffer, mpdType string) {
	if strings.EqualFold(mpdType, "dynamic") && p.ID == "" {
		buf.WriteString("Period must have ID when Type = 'dynamic'.\n")
	}
	if p.AdaptationSets != nil {
		for _, adaptationSet := range p.AdaptationSets {
			if p.BitstreamSwitching && !adaptationSet.BitstreamSwitching {
				buf.WriteString("Period element with field BitstreamSwitching = 'true' cannot contain AdaptaionSet with BitstreamSwitching = 'false'.\n")
			}
			adaptationSet.validate(buf)
		}
	}
	if p.SegmentBase != nil {
		p.SegmentBase.validate(buf)
	}
	if p.AssetIdentifier != nil {
		p.AssetIdentifier.validate(buf, "AssetIdentifier")
	}
	if p.SegmentList != nil {
		p.SegmentList.validate(buf)
	}
	if p.SegmentTemplate != nil {
		p.SegmentTemplate.validate(buf)
	}
}

func (s *SegmentTemplate) validate(buf *bytes.Buffer) {

}

func (s *SegmentList) validate(buf *bytes.Buffer) {

}

func xlinkActuateError(buf *bytes.Buffer, element string, xlinkActuate string) {
	if xlinkActuate != "onLoad" && xlinkActuate != "onRequest" {
		buf.WriteString(fmt.Sprintf("%s field XlinkActuate accepts values 'onRequest' and 'onLoad'.\n", element))
	}
}

func (a *AdaptationSet) validate(buf *bytes.Buffer) {
	if a.Accessibility != nil {
		for _, acess := range a.Accessibility {
			acess.validate(buf, "Accessibility")
		}
	}
	if a.AudioChannelConfig != nil {
		for _, acc := range a.AudioChannelConfig {
			acc.validate(buf, "AudioChannelConfig")
		}
	}
	if a.ContentProtection != nil {
		for _, cp := range a.ContentProtection {
			cp.validate(buf, "ContentProtection")
		}
	}
	if a.EssentialProperty != nil {
		for _, ep := range a.EssentialProperty {
			ep.validate(buf, "EssentialProperty")
		}
	}
	if a.FramePacking != nil {
		for _, fp := range a.FramePacking {
			fp.validate(buf, "FramePacking")
		}
	}
	if a.InbandEventStream != nil {
		for _, ies := range a.InbandEventStream {
			ies.validate(buf, "InbandEventStream")
		}
	}
	if a.Rating != nil {
		for _, r := range a.Rating {
			r.validate(buf, "Rating")
		}
	}
	if a.Role != nil {
		for _, ro := range a.Role {
			ro.validate(buf, "Role")
		}
	}
	if a.SupplementalProperty != nil {
		for _, sp := range a.SupplementalProperty {
			sp.validate(buf, "SupplementalProperty")
		}
	}
	if a.ViewPoint != nil {
		for _, v := range a.ViewPoint {
			v.validate(buf, "Viewpoint")
		}
	}
	if a.SegmentBase != nil {
		a.SegmentBase.validate(buf)
	}
	if a.Representations != nil {
		for _, re := range a.Representations {
			re.validate(buf)
		}
	}
}

func (r *Representation) validate(buf *bytes.Buffer) {

}

func (s *SegmentBase) validate(buf *bytes.Buffer) {
	if s != nil {
		if s.IndexRange == "" && s.IndexRangeExact {
			buf.WriteString("SegmentBase element field IndexRangeExact must not be present if IndexRange isn't specified.\n")
		}
	}
}

func (m *Metrics) validate(buf *bytes.Buffer) {
	if m != nil {
		if m.Metrics == "" {
			buf.WriteString("Metrics field Metrics is required.\n")
		}

		if m.Reporting != nil {
			for _, r := range m.Reporting {
				r.validate(buf, "Reporting")
			}
		} else {
			buf.WriteString("Metrics must have at least one Reporting element.\n")
		}
	}
}

func (d *Descriptor) validate(buf *bytes.Buffer, element string) {
	if d.SchemeIDURI == "" {
		buf.WriteString(fmt.Sprintf("%s field SchemeIdURI is required.\n", element))
	}
}
