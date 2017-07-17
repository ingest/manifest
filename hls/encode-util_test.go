package hls

import (
	"bytes"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ingest/manifest"
)

func TestWriteXMedia(t *testing.T) {
	rendition := &Rendition{
		Type:     "VIDEO",
		GroupID:  "TestID",
		Name:     "Test",
		Language: "English",
		Default:  true,
		URI:      "http://test.com",
	}

	buf := manifest.NewBufWrapper()

	rendition.writeXMedia(buf)
	if buf.Err != nil {
		t.Errorf("Expected err to be nil, but got %s", buf.Err.Error())
	}
}

func TestWriteXMediaTypeError(t *testing.T) {
	rendition := &Rendition{
		GroupID: "TestID",
	}

	buf := manifest.NewBufWrapper()

	rendition.writeXMedia(buf)
	if buf.Err.Error() != attributeNotSetError("EXT-X-MEDIA", "TYPE").Error() {
		t.Errorf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-MEDIA", "TYPE"), buf.Err.Error())
	}
}

func TestWriteXMediaGroupError(t *testing.T) {
	rendition := &Rendition{
		Type: "AUDIO",
	}

	buf := manifest.NewBufWrapper()

	rendition.writeXMedia(buf)
	if buf.Err.Error() != attributeNotSetError("EXT-X-MEDIA", "GROUP-ID").Error() {
		t.Errorf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-MEDIA", "GROUP-ID"), buf.Err.Error())
	}
}

func TestWriteXMediaInvalid(t *testing.T) {
	rendition := &Rendition{
		Type:       "CLOSED-CAPTIONS",
		GroupID:    "TestID",
		InstreamID: "CC3",
		Name:       "Test",
	}

	buf := manifest.NewBufWrapper()

	rendition.writeXMedia(buf)
	if buf.Err != nil {
		t.Errorf("Expected err to be nil")
	}

	rendition.URI = "test"
	buf = manifest.NewBufWrapper()
	rendition.writeXMedia(buf)
	if strings.Contains(buf.Buf.String(), "URI") {
		t.Error("Expected buf to not contain URI")
	}

	rendition.Type = "SUBTITLES"
	rendition.URI = ""
	buf = manifest.NewBufWrapper()
	rendition.writeXMedia(buf)
	if buf.Err.Error() != attributeNotSetError("EXT-X-MEDIA", "URI for SUBTITLES").Error() {
		t.Errorf("Exptected err to be %s, but got %s", attributeNotSetError("EXT-X-MEDIA", "URI for SUBTITLES").Error(), buf.Err.Error())
	}
}

func TestWriteStreamInf(t *testing.T) {
	variant := &Variant{
		IsIframe:   true,
		URI:        "http://test.com",
		Bandwidth:  234000,
		Resolution: "230x400",
	}

	buf := manifest.NewBufWrapper()
	variant.writeStreamInf(7, buf)

	if buf.Err != nil {
		t.Fatalf("Expected err to be nil, but got %s", buf.Err.Error())
	}

	if !strings.Contains(buf.Buf.String(), "#EXT-X-I-FRAME-STREAM-INF") {
		t.Error("Expected buf to contain #EXT-X-I-FRAME-STREAM-INF")
	}

	if strings.Contains(buf.Buf.String(), "#EXT-X-STREAM-INF") {
		t.Error("Expected buf to not contain #EXT-X-STREAM-INF")
	}
}

func TestGenerateMasterPlaylist(t *testing.T) {
	rend := &Rendition{
		Type:     "VIDEO",
		GroupID:  "TestID",
		Name:     "Test",
		Language: "English",
		Default:  true,
		URI:      "http://test.com",
	}
	rend2 := &Rendition{
		Type:     "AUDIO",
		GroupID:  "Testing",
		Name:     "Another test",
		Language: "English",
		Default:  false,
	}
	rend3 := &Rendition{
		Type:     "VIDEO",
		GroupID:  "Test",
		Name:     "Bla",
		Language: "Portuguese",
	}

	variant := &Variant{
		IsIframe:   false,
		URI:        "http://test.com",
		Bandwidth:  234000,
		Resolution: "230x400",
		Codecs:     "These codescs",
	}

	variant2 := &Variant{
		IsIframe:  false,
		URI:       "thistest.com",
		Bandwidth: 145000,
	}

	p := NewMasterPlaylist(5)
	p.Variants = append(p.Variants, variant, variant2)
	p.Renditions = append(p.Renditions, rend, rend2, rend3)
	p.SessionData = []*SessionData{&SessionData{DataID: "test", Value: "this is the session data"}}
	p.SessionKeys = []*Key{&Key{IsSession: true, Method: "sample-aes", URI: "keyuri"}}
	p.IndependentSegments = true
	buf, err := p.Encode()

	if err != nil {
		t.Fatalf("Expected err to be nil, but got %s", err.Error())
	}

	b := new(bytes.Buffer)
	b.ReadFrom(buf)

	if !strings.Contains(b.String(), "#EXT-X-SESSION-DATA") {
		t.Error("Expected buf to contain #EXT-X-SESSION-DATA")
	}

	if !strings.Contains(b.String(), "#EXT-X-SESSION-KEY") {
		t.Error("Expected buf to contain #EXT-X-SESSION-KEY")
	}

	if !strings.Contains(b.String(), "#EXT-X-INDEPENDENT-SEGMENTS") {
		t.Error("Expected buf to contain #EXT-X-INDEPENDENT-SEGMENTS")
	}
}

func TestGenerateMediaPlaylist(t *testing.T) {
	offset := int64(700)

	seg := &Segment{
		URI:       "segment.com",
		Inf:       &Inf{Duration: 9.052},
		Byterange: &Byterange{Length: 6000, Offset: &offset},
		Keys:      []*Key{&Key{Method: "sample-aes", URI: "keyuri"}},
		Map:       &Map{URI: "mapuri"},
		DateRange: &DateRange{ID: "test",
			StartDate:        time.Now(),
			EndDate:          time.Now().Add(1 * time.Hour),
			SCTE35:           &SCTE35{Type: "in", Value: "blablabla"},
			XClientAttribute: []string{"X-THIS-TAG=TEST", "X-THIS-OTHER-TAG=TESTING"}},
	}

	p := NewMediaPlaylist(7)
	p.Segments = append(p.Segments, seg)
	p.TargetDuration = 10
	p.EndList = true
	p.MediaSequence = 1
	p.StartPoint = &StartPoint{TimeOffset: 10.543}

	buf, err := p.Encode()
	if err != nil {
		t.Fatalf("Expected err to be nil, but got %s", err.Error())
	}

	b := new(bytes.Buffer)
	b.ReadFrom(buf)
	if !strings.Contains(b.String(), "#EXT-X-TARGETDURATION:10") {
		t.Error("Expected buf to contain #EXT-X-TARGETDURATION")
	}

	if !strings.Contains(b.String(), "#EXT-X-BYTERANGE:6000@700") {
		t.Error("Expected buf to contain #EXT-X-BYTERANGE")
	}

	if !strings.Contains(b.String(), "#EXT-X-ENDLIST") {
		t.Error("Expected buf to contain #EXT-X-ENDLIST")
	}

	if !strings.Contains(b.String(), "#EXT-X-MEDIA-SEQUENCE:1") {
		t.Error("Expected buf to contain #EXT-X-MEDIA-SEQUENCE")
	}

	if !strings.Contains(b.String(), "#EXT-X-START:TIME-OFFSET=10.543") {
		t.Error("Expected buf to contain #EXT-X-START")
	}

	if strings.Contains(b.String(), "#EXT-X-I-FRAMES-ONLY") {
		t.Error("Expected buf to not contain #EXT-X-I-FRAMES-ONLY")
	}

	p.Segments[0].Inf.Duration = 0
	p.StartPoint.TimeOffset = 0

	buf, err = p.Encode()
	if err != nil {
		t.Fatalf("Expected err to be nil, but got %s", err.Error())
	}

	b = new(bytes.Buffer)
	b.ReadFrom(buf)

	if !strings.Contains(b.String(), "#EXT-X-START:TIME-OFFSET=0.000") {
		t.Error("Expected buf to contain #EXT-X-START")
	}

	if !strings.Contains(b.String(), "#EXTINF:0.000,") {
		t.Error("Expected buf to contain #EXT-X-START")
	}
}

func TestDateRange(t *testing.T) {
	buf := manifest.NewBufWrapper()
	d := &DateRange{}
	d.writeDateRange(buf)
	if buf.Err.Error() != attributeNotSetError("EXT-X-DATERANGE", "ID").Error() {
		t.Errorf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-DATERANGE", "ID"), buf.Err)
	}

	buf = manifest.NewBufWrapper()
	d.ID = "test"
	d.writeDateRange(buf)
	if buf.Err.Error() != attributeNotSetError("EXT-X-DATERANGE", "START-DATE").Error() {
		t.Errorf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-DATERANGE", "START-DATE"), buf.Err)
	}

	buf = manifest.NewBufWrapper()
	d.StartDate = time.Now()
	d.EndOnNext = true
	d.writeDateRange(buf)
	if buf.Err == nil {
		t.Error("EndOnNext without Class should return error")
	}

	buf = manifest.NewBufWrapper()
	d.EndDate = time.Now().Add(-1 * time.Hour)
	d.EndOnNext = false
	d.writeDateRange(buf)
	if buf.Err == nil {
		t.Error("EndDate before StartDate should return error")
	}
}

func TestMap(t *testing.T) {
	buf := manifest.NewBufWrapper()
	m := &Map{
		Byterange: &Byterange{
			Length: 100,
		},
	}
	m.writeMap(buf)
	if buf.Err.Error() != attributeNotSetError("EXT-X-MAP", "URI").Error() {
		t.Fatalf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-MAP", "URI").Error(), buf.Err.Error())
	}

	buf = manifest.NewBufWrapper()
	m.URI = "test"
	m.writeMap(buf)
	if buf.Err != nil {
		t.Error("Expected err to be nil")
	}

	if !strings.Contains(buf.Buf.String(), "#EXT-X-MAP:URI=\"test\",BYTERANGE=\"100@0\"") {
		t.Error("Expected buf to contain #EXT-X-MAP")
	}
}

func TestCompatibilityCheck(t *testing.T) {
	p := NewMediaPlaylist(4)
	s := &Segment{
		Keys: []*Key{&Key{Method: "sample-aes", URI: "keyuri", Keyformat: "com.apple.streamingkeydelivery", Keyformatversions: "1"}},
	}

	err := p.checkCompatibility(s)
	if err.Error() != backwardsCompatibilityError(p.Version, "#EXT-X-KEY").Error() {
		t.Errorf("Error should be %s, but got %s", backwardsCompatibilityError(p.Version, "#EXT-X-KEY"), err)
	}

	p = NewMediaPlaylist(5)
	err = p.checkCompatibility(s)
	if err != nil {
		t.Errorf("Expected err to be nil, but got %s", err)
	}

	s = &Segment{
		Map: &Map{
			URI: "test",
		},
	}

	err = p.checkCompatibility(s)
	if err.Error() != backwardsCompatibilityError(p.Version, "#EXT-X-MAP").Error() {
		t.Errorf("Error should be %s, but got %s", backwardsCompatibilityError(p.Version, "#EXT-X-MAP"), err)
	}

	p = NewMediaPlaylist(3)
	s = &Segment{
		Byterange: &Byterange{Length: 100},
	}

	err = p.checkCompatibility(s)
	if err.Error() != backwardsCompatibilityError(p.Version, "#EXT-X-BYTERANGE").Error() {
		t.Errorf("Error should be %s, but got %s", backwardsCompatibilityError(p.Version, "#EXT-X-BYTERANGE").Error(), err.Error())
	}
}

func TestSortSegments(t *testing.T) {
	s := &Segment{
		ID:  1,
		URI: "firstsegment",
	}
	s2 := &Segment{
		ID:  2,
		URI: "secondsegment",
	}
	s3 := &Segment{
		ID:  3,
		URI: "thirdsegment",
	}
	var segs Segments
	segs = append(segs, s3, s, s2)
	sort.Sort(segs)
	for i := range segs {
		if segs[i].ID != i+1 {
			t.Errorf("Expected seg %d ID to be %d, but got %d", i, i+1, segs[i].ID)
		}
	}
}
