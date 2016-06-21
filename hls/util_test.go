package hls

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"
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

	buf := NewBufWriter()

	err := rendition.writeXMedia(buf)
	if err != nil {
		t.Errorf("Expected err to be nil, but got %s", err.Error())
	}
}

func TestWriteXMediaTypeError(t *testing.T) {
	rendition := &Rendition{
		GroupID: "TestID",
	}

	buf := NewBufWriter()

	err := rendition.writeXMedia(buf)
	if err.Error() != attributeNotSetError("EXT-X-MEDIA", "TYPE").Error() {
		t.Errorf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-MEDIA", "TYPE"), err.Error())
	}
}

func TestWriteXMediaGroupError(t *testing.T) {
	rendition := &Rendition{
		Type: "AUDIO",
	}

	buf := NewBufWriter()

	err := rendition.writeXMedia(buf)
	if err.Error() != attributeNotSetError("EXT-X-MEDIA", "GROUP-ID").Error() {
		t.Errorf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-MEDIA", "GROUP-ID"), err.Error())
	}
}

func TestWriteXMediaInvalid(t *testing.T) {
	rendition := &Rendition{
		Type:       "CLOSED-CAPTIONS",
		GroupID:    "TestID",
		InstreamID: "CC3",
		Name:       "Test",
	}

	buf := NewBufWriter()

	err := rendition.writeXMedia(buf)
	if err != nil {
		t.Errorf("Expected err to be nil")
	}

	rendition.URI = "test"
	buf = NewBufWriter()
	_ = rendition.writeXMedia(buf)
	if strings.Contains(buf.buf.String(), "URI") {
		t.Error("Expected buf to not contain URI")
	}

	rendition.Type = "SUBTITLES"
	rendition.URI = ""
	buf = NewBufWriter()
	err = rendition.writeXMedia(buf)
	if err.Error() != attributeNotSetError("EXT-X-MEDIA", "URI for SUBTITLES").Error() {
		t.Errorf("Exptected err to be %s, but got %s", attributeNotSetError("EXT-X-MEDIA", "URI for SUBTITLES").Error(), err.Error())
	}
}

func TestWriteStreamInf(t *testing.T) {
	variant := &Variant{
		IsIframe:   true,
		URI:        "http://test.com",
		Bandwidth:  234000,
		Resolution: "230x400",
	}

	buf := NewBufWriter()
	err := variant.writeStreamInf(7, buf)

	if err != nil {
		t.Fatalf("Expected err to be nil, but got %s", err.Error())
	}

	if !strings.Contains(buf.buf.String(), "#EXT-X-I-FRAME-STREAM-INF") {
		t.Error("Expected buf to contain #EXT-X-I-FRAME-STREAM-INF")
	}

	if strings.Contains(buf.buf.String(), "#EXT-X-STREAM-INF") {
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
	variant := &Variant{
		Renditions: []*Rendition{rend, rend2},
		IsIframe:   false,
		URI:        "http://test.com",
		Bandwidth:  234000,
		Resolution: "230x400",
		Codecs:     "These codescs",
	}

	rend3 := &Rendition{
		Type:     "VIDEO",
		GroupID:  "Test",
		Name:     "Bla",
		Language: "Portuguese",
	}

	variant2 := &Variant{
		Renditions: []*Rendition{rend3},
		IsIframe:   false,
		URI:        "thistest.com",
		Bandwidth:  145000,
	}

	p := NewMasterPlaylist(5)
	p.Variants = append(p.Variants, variant, variant2)
	p.SessionData = []*SessionData{&SessionData{DataID: "test", Value: "this is the session data"}}
	p.SessionKeys = []*Key{&Key{IsSession: true, Method: "sample-aes", URI: "keyuri"}}
	p.IndependentSegments = true
	buf, err := p.GenerateManifest()

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

	buf, err := p.GenerateManifest()

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
}

func TestDateRange(t *testing.T) {
	buf := NewBufWriter()
	d := &DateRange{}
	err := d.writeDateRange(buf)
	if err.Error() != attributeNotSetError("EXT-X-DATERANGE", "ID").Error() {
		t.Errorf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-DATERANGE", "ID"), err)
	}

	buf = NewBufWriter()
	d.ID = "test"
	err = d.writeDateRange(buf)
	if err.Error() != attributeNotSetError("EXT-X-DATERANGE", "START-DATE").Error() {
		t.Errorf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-DATERANGE", "START-DATE"), err)
	}

	buf = NewBufWriter()
	d.StartDate = time.Now()
	d.EndOnNext = true
	err = d.writeDateRange(buf)
	if err == nil {
		t.Error("EndOnNext without Class should return error")
	}

	buf = NewBufWriter()
	d.EndDate = time.Now().Add(-1 * time.Hour)
	d.EndOnNext = false
	err = d.writeDateRange(buf)
	if err == nil {
		t.Error("EndDate before StartDate should return error")
	}
}

func TestMap(t *testing.T) {
	buf := NewBufWriter()
	m := &Map{
		Byterange: &Byterange{
			Length: 100,
		},
	}
	err := m.writeMap(buf)
	if err.Error() != attributeNotSetError("EXT-X-MAP", "URI").Error() {
		t.Fatalf("Expected err to be %s, but got %s", attributeNotSetError("EXT-X-MAP", "URI").Error(), err.Error())
	}

	buf = NewBufWriter()
	m.URI = "test"
	err = m.writeMap(buf)
	if err != nil {
		t.Error("Expected err to be nil")
	}
	fmt.Println(buf.buf.String())
	if !strings.Contains(buf.buf.String(), "#EXT-X-MAP:URI=\"test\",BYTERANGE=\"100@0\"") {
		t.Error("Expected buf to contain #EXT-X-MAP")
	}
}

func TestCompatibilityCheck(t *testing.T) {
	p := NewMediaPlaylist(4)
	s := &Segment{
		Keys: []*Key{&Key{Method: "sample-aes", URI: "keyuri"}},
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
