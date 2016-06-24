package hls

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

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
	pt, bla := time.Parse(time.RFC3339Nano, "2016-06-22T15:33:52.199039986Z")
	fmt.Println(pt)
	fmt.Println(bla)
	seg := &Segment{
		URI:       "segment.com",
		Inf:       &Inf{Duration: 9.052},
		Byterange: &Byterange{Length: 6000, Offset: &offset},
		Keys:      []*Key{&Key{Method: "sample-aes", URI: "keyuri"}},
		Map:       &Map{URI: "mapuri"},
		DateRange: &DateRange{ID: "test",
			StartDate:        pt,
			EndDate:          pt.Add(1 * time.Hour),
			SCTE35:           &SCTE35{Type: "in", Value: "blablabla"},
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

	p := NewMediaPlaylist(7)
	p.Segments = append(p.Segments, seg, seg2)
	p.TargetDuration = 10
	p.EndList = true
	p.MediaSequence = 1
	p.StartPoint = &StartPoint{TimeOffset: 10.543}

	buf, _ := p.GenerateManifest()

	newP := NewMediaPlaylist(7)
	err := newP.ReadManifest(buf)
	fmt.Println(err)
	fmt.Println(newP)
}
