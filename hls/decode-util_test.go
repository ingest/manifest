package hls

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestReadMasterPlaylistFile(t *testing.T) {
	f, err := os.Open("./playlists/masterp.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := &MasterPlaylist{}
	err = p.ReadManifest(bufio.NewReader(f))
	if err != io.EOF {
		t.Errorf("Expected err to be EOF, but got %s", err)
	}
	if len(p.SessionData) != 2 {
		t.Errorf("Expected SessionData len 2, but got %d", len(p.SessionData))
	}
	//TODO:add more checks
}

func TestReadMediaPlaylistFile(t *testing.T) {
	f, err := os.Open("./playlists/mediap.m3u8")
	if err != nil {
		t.Fatal(err)
	}
	p := &MediaPlaylist{}
	p.ReadManifest(bufio.NewReader(f))
	for _, s := range p.Segments {
		if s.DateRange != nil {
			fmt.Println(s.DateRange.StartDate)
		}
	}
}

func TestReadMediaPlaylist(t *testing.T) {
	offset := int64(700)
	pt, _ := time.Parse(time.RFC3339Nano, "2016-06-22T15:33:52.199039986Z")
	seg := &Segment{
		URI:       "segment.com",
		Inf:       &Inf{Duration: 9.052},
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
		URI:       "segment2.com",
		Inf:       &Inf{Duration: 8.052, Title: "seg title"},
		Byterange: &Byterange{Length: 4000},
		Keys:      []*Key{&Key{Method: "sample-aes", URI: "keyuri"}},
		Map:       &Map{URI: "map2"},
		DateRange: &DateRange{ID: "test", StartDate: pt},
	}

	seg3 := &Segment{
		URI:             "segment3.com",
		Inf:             &Inf{Duration: 9.500},
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
	p.M3U = true

	buf, err := p.GenerateManifest()
	fmt.Println(err)
	newP := NewMediaPlaylist(7)
	err = newP.ReadManifest(buf)
	if err != io.EOF {
		t.Fatalf("Expected err to be EOF, but got %s", err.Error())
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
		if s.Map != nil && !reflect.DeepEqual(s.Map, newP.Segments[i].Map) {
			t.Errorf("Expected %d Segment Map tp be %v, but got %v", i, s.Map, newP.Segments[i].Map)
		}
		if s.DateRange != nil && !reflect.DeepEqual(s.DateRange, newP.Segments[i].DateRange) {
			t.Errorf("Expected %d Segment DateRange to be %v, but got %v", i, s.DateRange, newP.Segments[i].DateRange)
		}
	}
}
