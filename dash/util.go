package dash

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

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
				metric.validateMetrics(buf)
			}
		}

		if m.Periods != nil {
			for _, p := range m.Periods {
				p.validatePeriod(buf, m.Type)
			}
		} else {
			buf.WriteString("MPD must have at least one Period element.\n")
		}
	}
}

func (p *Period) validatePeriod(buf *bytes.Buffer, mpdType string) {
	if p != nil {
		if strings.EqualFold(mpdType, "dynamic") && p.ID == "" {
			buf.WriteString("Period must have ID when Type = 'dynamic'.\n")
		}
		if p.AdaptationSets != nil {
			for _, a := range p.AdaptationSets {
				if p.BitstreamSwitching && !a.BitstreamSwitching {
					buf.WriteString("Period element with field BitstreamSwitching true cannot contain AdaptaionSet with BitstreamSwitching false.\n")
				}
			}
		}
		if p.SegmentBase != nil {

		}
	}
}

func (m *Metrics) validateMetrics(buf *bytes.Buffer) {
	if m != nil {
		if m.Metrics == "" {
			buf.WriteString("Metrics field Metrics is required.\n")
		}

		if m.Reporting != nil {
			for _, r := range m.Reporting {
				r.validateDescriptor(buf, "Reporting")
			}
		} else {
			buf.WriteString("Metrics must have at least one Reporting element.\n")
		}
	}
}

func (d *Descriptor) validateDescriptor(buf *bytes.Buffer, element string) {
	if d.SchemeIDURI == "" {
		buf.WriteString(fmt.Sprintf("%s field SchemeIdURI is required.\n", element))
	}
}
