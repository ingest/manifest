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
			variant.Bandwidth, err = parseInt(v, 10, 64)
		case "AVERAGE-BANDWIDTH":
			variant.AvgBandwidth, err = parseInt(v, 10, 64)
		case "CODECS":
			variant.Codecs = v
		case "RESOLUTION":
			variant.Resolution = v
		case "FRAME-RATE":
			variant.FrameRate, err = parseFloat(v, 64)
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

//parseInt wraps strconv.ParseInt
func parseInt(s string, base int, bitSize int) (int64, error) {
	d, err := strconv.ParseInt(s, base, bitSize)
	if err != nil {
		return 0, err
	}
	return d, err
}

//parseFloat wraps strconv.ParseFloat
func parseFloat(s string, bitSize int) (float64, error) {
	f, err := strconv.ParseFloat(s, bitSize)
	if err != nil {
		return float64(0), err
	}
	return f, err
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
	i.Duration, err = parseFloat(stringBefore(line, ","), 64)
	i.Title = stringAfter(line, ",")
	return i, err
}

func decodeDateRange(line string) (*DateRange, error) {
	drMap := splitParams(line)
	var err error
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

//TODO:Improve this
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

func decodeKey(line string) *Key {
	keyMap := splitParams(line)

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
	return key
}

func decodeStartPoint(line string) (*StartPoint, error) {
	spMap := splitParams(line)
	var err error
	sp := &StartPoint{}
	for k, v := range spMap {
		switch k {
		case "TIME-OFFSET":
			sp.TimeOffset, err = parseFloat(v, 64)
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

//TODO: Improve this to check for badly formatted attributes?
func splitParams(line string) map[string]string {
	re := regexp.MustCompile(`([a-zA-Z\d_-]+)=("[^"]+"|[^",]+)`)
	m := make(map[string]string)
	for _, kv := range re.FindAllStringSubmatch(line, -1) {
		k, v := kv[1], kv[2]
		m[strings.ToUpper(k)] = strings.Trim(v, "\"")
	}
	return m
}
