package hls

import (
	"bufio"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestProtocolCompat(t *testing.T) {
	tests := []struct {
		expectErr bool
		mediaType string
		file      string
	}{
		{
			expectErr: true,
			mediaType: "media",
			file:      "fixture-v3-fail.m3u8",
		},
		{
			expectErr: false,
			mediaType: "media",
			file:      "fixture-v3.m3u8",
		},
		{
			expectErr: true,
			mediaType: "media",
			file:      "fixture-v4-byterange-fail.m3u8",
		},
		{
			expectErr: true,
			mediaType: "media",
			file:      "fixture-v4-iframes-fail.m3u8",
		},
		{
			expectErr: false,
			mediaType: "media",
			file:      "fixture-v4.m3u8",
		},
		{
			expectErr: true,
			mediaType: "media",
			file:      "fixture-v5-keyformat-fail.m3u8",
		},
		{
			expectErr: true,
			mediaType: "media",
			file:      "fixture-v5-map-fail.m3u8",
		},
		{
			expectErr: false,
			mediaType: "media",
			file:      "fixture-v5.m3u8",
		},
		{
			expectErr: true,
			mediaType: "media",
			file:      "fixture-v6-map-no-iframes-fail.m3u8",
		},
		{
			expectErr: false,
			mediaType: "media",
			file:      "fixture-v6.m3u8",
		},
		{
			expectErr: true,
			mediaType: "master",
			file:      "fixture-v7-media-service-fail.m3u8",
		},
		{
			expectErr: false,
			mediaType: "master",
			file:      "fixture-v7.m3u8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			f, err := os.Open(filepath.Join("./testdata", tt.mediaType, tt.file))
			if err != nil {
				t.Fatal(err)
			}

			switch tt.mediaType {
			case "media":
				p := NewMediaPlaylist(0)
				if err := p.Parse(f); (err != nil) != tt.expectErr {
					t.Errorf("expected (%t) err: %v", tt.expectErr, err)
				}
			case "master":
				p := NewMasterPlaylist(0)
				if err := p.Parse(f); (err != nil) != tt.expectErr {
					t.Errorf("expected (%t) err: %v", tt.expectErr, err)
				}
			}

		})
	}
}

func TestReadMasterPlaylistFile(t *testing.T) {
	f, err := os.Open("./testdata/masterp.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := &MasterPlaylist{}
	err = p.Parse(bufio.NewReader(f))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(p.SessionData) != 2 {
		t.Errorf("Expected SessionData len 2, but got %d", len(p.SessionData))
	}

	if len(p.Variants) != 14 {
		t.Errorf("Expected Variants len 14, but got %d", len(p.Variants))
	}

	if len(p.Renditions) != 5 {
		t.Errorf("Expected Renditions len 5, but got %d", len(p.Renditions))
	}

	k := &Key{IsSession: true,
		Method:            "SAMPLE-AES",
		IV:                "0x29fd9eba3735966ddfca572e51e68ff2",
		URI:               "com.keyuri.example",
		Keyformat:         "com.apple.streamingkeydelivery",
		Keyformatversions: "1"}
	if p.SessionKeys != nil {
		if !reflect.DeepEqual(k, p.SessionKeys[0]) {
			t.Errorf("Expected SessionKeys to be %v, but got %v", k, p.SessionKeys[0])
		}
	}
}

func TestReadMediaPlaylistFile(t *testing.T) {
	f, err := os.Open("./testdata/mediap.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := &MediaPlaylist{}
	p.Parse(bufio.NewReader(f))
	if p.TargetDuration != 10 {
		t.Errorf("Expected TargetDuration 10, but got %d", p.TargetDuration)
	}

	if p.StartPoint.TimeOffset != 8.345 {
		t.Errorf("Expected StartPoint to be 8.345, but got %v", p.StartPoint.TimeOffset)
	}
	if p.Segments != nil {
		if len(p.Segments) != 6 {
			t.Errorf("Expected len Segments 6, but got %d", len(p.Segments))
		}
		sd, _ := time.Parse(time.RFC3339Nano, "2010-02-19T14:54:23.031+08:00")
		dr := &DateRange{ID: "6FFF00", StartDate: sd, SCTE35: &SCTE35{Type: "OUT", Value: "0xFC002F0000000000FF0"}}
		if !reflect.DeepEqual(p.Segments[0].DateRange, dr) {
			t.Errorf("Expected DateRange to be %v, but got %v", dr, p.Segments[0].DateRange)
		}
		c := p.MediaSequence
		for i := range p.Segments {
			if p.Segments[i].ID != c {
				t.Errorf("Expected Segments %d ID to be %d, but got %d", i, c, p.Segments[i].ID)
			}
			c++
		}
	}

}

func TestReadMediaPlaylist(t *testing.T) {
	offset := int64(700)
	duration := float64(200)
	pt, _ := time.Parse(time.RFC3339Nano, "2016-06-22T15:33:52.199039986Z")
	seg := &Segment{
		URI: "segment.com",
		Inf: &Inf{
			Duration: 9.052,
		},
		Byterange: &Byterange{Length: 6000, Offset: &offset},
		Keys:      []*Key{&Key{Method: "sample-aes", URI: "keyuri"}, &Key{Method: "sample-aes", URI: "secondkeyuri"}},
		Map:       &Map{URI: "mapuri"},
		DateRange: &DateRange{ID: "TEST",
			StartDate:        pt,
			EndDate:          pt.Add(1 * time.Hour),
			SCTE35:           &SCTE35{Type: "IN", Value: "bla"},
			XClientAttribute: []string{"X-THIS-TAG=TEST", "X-THIS-OTHER-TAG=TESTING"}},
	}

	seg2 := &Segment{
		URI: "segment2.com",
		Inf: &Inf{
			Duration: 8.052,
			Title:    "seg title",
		},
		Byterange: &Byterange{Length: 4000},
		Keys:      []*Key{&Key{Method: "sample-aes", URI: "keyuri"}},
		Map:       &Map{URI: "map2"},
		DateRange: &DateRange{ID: "test", StartDate: pt, Duration: &duration},
	}

	seg3 := &Segment{
		URI: "segment3.com",
		Inf: &Inf{
			Duration: 9.500,
		},
		ProgramDateTime: time.Now(),
		Discontinuity:   true,
	}

	p := NewMediaPlaylist(7)
	p.Segments = append(p.Segments, seg, seg2, seg3)
	p.DiscontinuitySequence = 2
	p.TargetDuration = 10
	p.EndList = true
	p.MediaSequence = 1
	p.StartPoint = &StartPoint{TimeOffset: 10.543}
	buf, err := p.Encode()

	newP := NewMediaPlaylist(0)
	err = newP.Parse(buf)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if newP.Version != 7 {
		t.Errorf("expected version to be 7, got %d", newP.Version)
	}

	if newP.TargetDuration != p.TargetDuration {
		t.Errorf("Expected TargetDuration to be %d, but got %d", p.TargetDuration, newP.TargetDuration)
	}

	if newP.DiscontinuitySequence != p.DiscontinuitySequence {
		t.Errorf("Expected DiscontinuitySequence to be %d, but got %d", p.DiscontinuitySequence, newP.DiscontinuitySequence)
	}

	if !reflect.DeepEqual(newP.StartPoint, p.StartPoint) {
		t.Errorf("Expected StartPoint to be %v, but got %v", p.StartPoint, newP.StartPoint)
	}

	for i, s := range p.Segments {
		if !reflect.DeepEqual(s.Inf, newP.Segments[i].Inf) {
			t.Errorf("Expected %d Segment Inf to be %v, but got %v", i, s.Inf, newP.Segments[i].Inf)
		}
		if s.URI != newP.Segments[i].URI {
			t.Errorf("Expected URI to be %s, but got %s", s.URI, newP.Segments[i].URI)
		}
		if s.Map != nil && s.Map.URI != newP.Segments[i].Map.URI {
			t.Errorf("Expected %d Segment Map to be %v, but got %v", i, s.Map, newP.Segments[i].Map)
		}
		// if s.DateRange != nil && !reflect.DeepEqual(s.DateRange, newP.Segments[i].DateRange) {
		// 	t.Errorf("Expected %d Segment DateRange to be %v, but got %v", i, s.DateRange, newP.Segments[i].DateRange)
		// }
	}
}
