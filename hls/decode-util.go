package hls

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func decodeVariant(line string, isIframe bool) (*Variant, error) {
	fmt.Println(line)
	vMap, err := splitParams(line)
	if err != nil {
		return nil, err
	}
	variant := &Variant{}

	for k, v := range vMap {
		switch k {
		case "BANDWIDTH":
			variant.Bandwidth, err = strconv.ParseInt(v, 10, 64)
		case "AVERAGE-BANDWIDTH":
			variant.AvgBandwidth, err = strconv.ParseInt(v, 10, 64)
		case "CODECS":
			variant.Codecs = v
		case "RESOLUTION":
			variant.Resolution = v
		case "FRAME-RATE":
			variant.FrameRate, err = strconv.ParseFloat(v, 64)
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

func decodeRendition(line string) (*Rendition, error) {
	rMap, err := splitParams(line)
	if err != nil {
		return nil, err
	}

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
			if v == boolYes {
				rendition.Default = true
			}
		case "AUTOSELECT":
			if v == boolYes {
				rendition.AutoSelect = true
			}
		case "FORCED":
			if v == boolYes {
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
	return rendition, err
}

func decodeSessionData(line string) (*SessionData, error) {
	sdMap, err := splitParams(line)
	if err != nil {
		return nil, err
	}
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
	return sd, nil
}

func decodeInf(line string) (*Inf, error) {
	d, err := strconv.ParseFloat(stringBefore(line, ","), 64)
	if err != nil {
		return nil, err
	}
	i := &Inf{Duration: d, Title: stringAfter(line, ",")}
	return i, err
}

func decodeDateRange(line string) (*DateRange, error) {
	drMap, err := splitParams(line)
	if err != nil {
		return nil, err
	}
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
			if d, err := strconv.ParseFloat(v, 64); err == nil {
				dr.Duration = &d
			} else {
				return nil, err
			}
		case k == "PLANNED-DURATION":
			if pd, err := strconv.ParseFloat(v, 64); err == nil {
				dr.PlannedDuration = &pd
			} else {
				return nil, err
			}
		case strings.HasPrefix(k, "X-"):
			dr.XClientAttribute = append(dr.XClientAttribute, fmt.Sprintf("%s=%s", k, v))
		case strings.HasPrefix(k, "SCTE35"):
			dr.SCTE35 = decodeSCTE(k, v)
		case k == "END-ON-NEXT" && v == boolYes:
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
	mMap, err := splitParams(line)
	if err != nil {
		return nil, err
	}

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

//TODO: Improve this
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

func decodeKey(line string) (*Key, error) {
	keyMap, err := splitParams(line)
	if err != nil {
		return nil, err
	}

	key := &Key{}
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
	return key, nil
}

func decodeStartPoint(line string) (*StartPoint, error) {
	spMap, err := splitParams(line)
	if err != nil {
		return nil, err
	}

	sp := &StartPoint{}
	for k, v := range spMap {
		switch k {
		case "TIME-OFFSET":
			if to, err := strconv.ParseFloat(v, 64); err == nil {
				sp.TimeOffset = to
			} else {
				return nil, err
			}
		case "PRECISE":
			if v == boolYes {
				sp.Precise = true
			}
		}
	}
	return sp, err
}

func stringAfter(line string, char string) (ret string) {
	index := strings.Index(line, char)
	if index != -1 {
		ret = line[index+1 : len(line)]
	}
	return
}

func stringBefore(line string, char string) (ret string) {
	index := strings.Index(line, char)
	if index != -1 {
		ret = line[0:index]
	}
	return
}

//TODO: Try to improve this. Regex?
//add regex to ignore comma inside double quotes. IE. CODECS att.
func splitParams(line string) (map[string]string, error) {
	params := strings.Split(line, ",")
	m := make(map[string]string)
	for _, p := range params {
		att := strings.Split(p, "=")
		if len(att) == 2 {
			k, v := att[0], att[1]
			m[strings.ToUpper(k)] = strings.Trim(v, "\"")
		} else {
			return nil, errors.New("Badly formatted attribute list")
		}
	}
	return m, nil
}
