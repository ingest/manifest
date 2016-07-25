package dash

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

//NewMPD initiates a MPD struct with the minimum required attributes
func NewMPD(profile string, minBufferTime time.Duration) *MPD {
	return &MPD{XMLNS: dashNS,
		Type:          "static",
		Profiles:      profile,
		MinBufferTime: &CustomDuration{Duration: minBufferTime},
	}
}

//NewContentProtection sets ContentProtection element with the appropriate
//namespaces.
func NewContentProtection(schemeIDUri string,
	value string,
	defaultKID string,
	pssh string,
	pro string) *CENCContentProtection {
	cp := &CENCContentProtection{
		ContentProtection: ContentProtection{
			SchemeIDURI: schemeIDUri,
			Value:       value,
		},
	}
	if defaultKID != "" {
		cp.ContentProtection.XMLNsCenc = cencNS
		cp.ContentProtection.DefaultKID = defaultKID
	}
	if pssh != "" {
		cp.ContentProtection.XMLNsCenc = cencNS
		cp.Pssh = &Pssh{
			XMLName: xml.Name{Local: "pssh", Space: "cenc"},
			Value:   pssh,
		}
	}
	if pro != "" {
		cp.ContentProtection.XMLNsMspr = msprNS
		cp.Pro = &Pro{
			XMLName: xml.Name{Local: "pro", Space: "mspr"},
			Value:   pro,
		}
	}

	return cp
}

//SetTrackEncryptionBox sets PlayReady's Track Encryption Box fields (tenc).
func (c *CENCContentProtection) SetTrackEncryptionBox(ivSize int, kid string) {
	c.ContentProtection.XMLNsMspr = msprNS
	c.IsEncrypted = "1"
	c.IVSize = ivSize
	c.KID = kid
}

//AddSegment adds a Segment to a SegmentTimeline and sorts it.
func (st *SegmentTimeline) AddSegment(t, d, r int) {
	if st == nil {
		return
	}
	s := &S{T: t,
		D: d,
		R: r,
	}
	st.Segments = append(st.Segments, s)
	sort.Sort(st.Segments)
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
		//validate attributes
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

		if m.MinBufferTime == nil || m.MinBufferTime.Duration == 0 {
			buf.WriteString("MPD field MinBufferTime is required.\n")
		}

		if m.Metrics != nil {
			for _, metric := range m.Metrics {
				metric.validate(buf)
			}
		}
		//validate Period
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
	//validate XlinkActuate
	xlinkActuateError(buf, "Period", p.XlinkActuate)

	if strings.EqualFold(mpdType, "dynamic") && p.ID == "" {
		buf.WriteString("Period must have ID when Type = 'dynamic'.\n")
	}
	if p.AssetIdentifier != nil {
		p.AssetIdentifier.validate(buf, "AssetIdentifier")
	}
	if p.SegmentBase != nil {
		p.SegmentBase.validate(buf)
	}
	if p.SegmentList != nil {
		p.SegmentList.validate(buf)
	}
	//validate AdaptationSets
	if p.AdaptationSets != nil {
		for _, adaptationSet := range p.AdaptationSets {
			if p.BitstreamSwitching && !adaptationSet.BitstreamSwitching {
				buf.WriteString("Period element with field BitstreamSwitching = 'true' cannot contain AdaptaionSet with BitstreamSwitching = 'false'.\n")
			}
			adaptationSet.validate(buf)
		}
	}
	//check if only one out of SegmentBase, SegmentList and SegmentTemplate is present
	if !validateSegmentPresence(p.SegmentBase, p.SegmentTemplate, p.SegmentList) {
		buf.WriteString("At most one of the three, SegmentBase, SegmentTemplate and SegmentList shall be present in a Period element.\n")
	}
}

func (s *SegmentList) validate(buf *bytes.Buffer) {
	xlinkActuateError(buf, "SegmentList", s.XlinkActuate)
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
	if a.Representations != nil {
		for _, re := range a.Representations {
			re.validate(buf)
		}
	}

	if a.SegmentBase != nil {
		a.SegmentBase.validate(buf)
	}
	//check if only one out of SegmentBase, SegmentList and SegmentTemplate is present
	if !validateSegmentPresence(a.SegmentBase, a.SegmentTemplate, a.SegmentList) {
		buf.WriteString("At most one of the three, SegmentBase, SegmentTemplate and SegmentList shall be present in AdaptationSet element.\n")
	}
}

func (r *Representation) validate(buf *bytes.Buffer) {
	r.SegmentBase.validate(buf)
	//check if only one out of SegmentBase, SegmentList and SegmentTemplate is present
	if !validateSegmentPresence(r.SegmentBase, r.SegmentTemplate, r.SegmentList) {
		buf.WriteString("At most one of the three, SegmentBase, SegmentTemplate and SegmentList shall be present in Representation element.\n")
	}
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

func xlinkActuateError(buf *bytes.Buffer, element string, xlinkActuate string) {
	if xlinkActuate != "" && xlinkActuate != "onLoad" && xlinkActuate != "onRequest" {
		buf.WriteString(fmt.Sprintf("%s field XlinkActuate accepts values 'onRequest' and 'onLoad'.\n", element))
	}
}

//checks if only one out of SegmentBase, SegmentList and SegmentTemplate is present
func validateSegmentPresence(sb *SegmentBase, st *SegmentTemplate, sl *SegmentList) bool {
	var segment int
	if sb != nil {
		segment++
	}
	if st != nil {
		segment++
	}
	if sl != nil {
		segment++
	}
	return segment <= 1
}
