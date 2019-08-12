package hls

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func decodeVariant(line string, isIframe bool) (*Variant, error) {
	var err error
	vMap := splitParams(line)
	variant := &Variant{}

	for k, v := range vMap {
		switch k {
		case "BANDWIDTH":
			if variant.Bandwidth, err = strconv.ParseInt(v, 10, 64); err != nil {
				return nil, err
			}
		case "AVERAGE-BANDWIDTH":
			if variant.AvgBandwidth, err = strconv.ParseInt(v, 10, 64); err != nil {
				return nil, err
			}
		case "CODECS":
			variant.Codecs = v
		case "RESOLUTION":
			variant.Resolution = v
		case "FRAME-RATE":
			if variant.FrameRate, err = strconv.ParseFloat(v, 64); err != nil {
				return nil, err
			}
		case "AUDIO":
			variant.Audio = v
		case "VIDEO":
			variant.Video = v
		case "SUBTITLES":
			variant.Subtitles = v
		case "CLOSED-CAPTIONS":
			variant.ClosedCaptions = v
		case "URI":
			variant.URI = v
		}

		if err != nil {
			return nil, err
		}
	}

	variant.IsIframe = isIframe
	return variant, err
}

func decodeRendition(line string) *Rendition {
	rMap := splitParams(line)

	rendition := &Rendition{}
	for k, v := range rMap {
		switch k {
		case "TYPE":
			if isValidType(strings.ToUpper(v)) {
				rendition.Type = v
			}
		case "URI":
			rendition.URI = v
		case "GROUP-ID":
			rendition.GroupID = v
		case "LANGUAGE":
			rendition.Language = v
		case "ASSOC-LANGUAGE":
			rendition.AssocLanguage = v
		case "NAME":
			rendition.Name = v
		case "DEFAULT":
			if strings.EqualFold(v, boolYes) {
				rendition.Default = true
			}
		case "AUTOSELECT":
			if strings.EqualFold(v, boolYes) {
				rendition.AutoSelect = true
			}
		case "FORCED":
			if strings.EqualFold(v, boolYes) {
				rendition.Forced = true
			}
		case "INSTREAM-ID":
			if isValidInstreamID(strings.ToUpper(v)) {
				rendition.InstreamID = v
			}
		case "CHARACTERISTICS":
			rendition.Characteristics = v
		}
	}
	return rendition
}

func decodeSessionData(line string) *SessionData {
	sdMap := splitParams(line)
	sd := &SessionData{}
	for k, v := range sdMap {
		switch k {
		case "DATA-ID":
			sd.DataID = v
		case "VALUE":
			sd.Value = v
		case "URI":
			sd.URI = v
		case "LANGUAGE":
			sd.Language = v
		}
	}
	return sd
}

func decodeInf(line string) (*Inf, error) {
	var err error
	i := &Inf{}
	index := strings.Index(line, ",")
	if index == -1 {
		return nil, fmt.Errorf("no comma was found when decoding #EXTINF: %s", line)
	}
	if i.Duration, err = strconv.ParseFloat(line[0:index], 64); err != nil {
		return nil, err
	}
	i.Title = line[index+1 : len(line)]
	return i, err
}

func decodeDateRange(line string) (*DateRange, error) {
	var err error
	var d float64
	var pd float64
	drMap := splitParams(line)

	dr := &DateRange{}
	for k, v := range drMap {
		switch {
		case k == "ID":
			dr.ID = v
		case k == "CLASS":
			dr.Class = v
		case k == "START-DATE":
			dr.StartDate, err = decodeDateTime(v)
		case k == "END-DATE":
			dr.EndDate, err = decodeDateTime(v)
		case k == "DURATION":
			if d, err = strconv.ParseFloat(v, 64); err == nil {
				dr.Duration = &d
			} else {
				return nil, err
			}
		case k == "PLANNED-DURATION":
			if pd, err = strconv.ParseFloat(v, 64); err == nil {
				dr.PlannedDuration = &pd
			} else {
				return nil, err
			}
		case strings.HasPrefix(k, "X-"):
			dr.XClientAttribute = append(dr.XClientAttribute, fmt.Sprintf("%s=%s", k, v))
		case strings.HasPrefix(k, "SCTE35"):
			dr.SCTE35 = decodeSCTE(k, v)
		case k == "END-ON-NEXT" && strings.EqualFold(v, boolYes):
			dr.EndOnNext = true
		}
	}

	return dr, err
}

func decodeSCTE(att string, value string) *SCTE35 {
	if strings.HasSuffix(att, "IN") {
		return &SCTE35{Type: "IN", Value: value}
	} else if strings.HasSuffix(att, "OUT") {
		return &SCTE35{Type: "OUT", Value: value}
	} else if strings.HasSuffix(att, "CMD") {
		return &SCTE35{Type: "CMD", Value: value}
	}
	return nil
}

func decodeDateTime(line string) (time.Time, error) {
	line = strings.Trim(line, "\"")
	t, err := time.Parse(time.RFC3339Nano, line)
	return t, err
}

func decodeMap(line string) (*Map, error) {
	mMap := splitParams(line)
	var err error
	m := &Map{}
	for k, v := range mMap {
		switch k {
		case "URI":
			m.URI = v
		case "BYTERANGE":
			m.Byterange, err = decodeByterange(v)
		}
	}
	return m, err
}

func decodeByterange(value string) (*Byterange, error) {
	params := strings.Split(value, "@")
	l, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return nil, err
	}
	b := &Byterange{Length: l}
	if len(params) == 2 {
		o, err := strconv.ParseInt(params[1], 10, 64)
		if err != nil {
			return b, err
		}
		b.Offset = &o
	}
	return b, nil
}

func decodeKey(line string, isSession bool) *Key {
	keyMap := splitParams(line)

	key := Key{IsSession: isSession}
	for k, v := range keyMap {
		switch k {
		case "METHOD":
			key.Method = v
		case "URI":
			key.URI = v
		case "IV":
			key.IV = v
		case "KEYFORMAT":
			key.Keyformat = v
		case "KEYFORMATVERSIONS":
			key.Keyformatversions = v
		}
	}
	return &key
}

func decodeStartPoint(line string) (*StartPoint, error) {
	spMap := splitParams(line)
	var err error
	sp := &StartPoint{}
	for k, v := range spMap {
		switch k {
		case "TIME-OFFSET":
			if sp.TimeOffset, err = strconv.ParseFloat(v, 64); err != nil {
				return nil, err
			}
		case "PRECISE":
			if strings.EqualFold(v, boolYes) {
				sp.Precise = true
			}
		}
	}
	return sp, err
}

//splitParams receives the comma-separated list of attributes and maps attribute-value pairs
func splitParams(line string) map[string]string {
	//regex to recognize att=val format and split on comma, unless comma is inside quotes
	re := regexp.MustCompile(`([a-zA-Z\d_-]+)=("[^"]+"|[^",]+)`)
	m := make(map[string]string)
	for _, kv := range re.FindAllStringSubmatch(line, -1) {
		k, v := kv[1], kv[2]
		m[strings.ToUpper(k)] = strings.Trim(v, "\"")
	}
	return m
}

//stringsIndex wraps string.Index and sets index = 0 if not found
func stringsIndex(line string, char string) int {
	index := strings.Index(line, char)
	if index == -1 {
		index = 0
	}
	return index
}
