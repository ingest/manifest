package dash

import (
	"bufio"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestStaticParse(t *testing.T) {
	f, err := os.Open("./playlists/static.mpd")
	if err != nil {
		t.Fatal(err)
	}
	mpd := &MPD{}
	if err := mpd.Parse(bufio.NewReader(f)); err != nil {
		t.Fatal(err)
	}

	if len(mpd.Periods) != 3 {
		t.Errorf("Expecting 3 Period elements, but got %d", len(mpd.Periods))
	}

	if len(mpd.Periods[0].AdaptationSets) != 2 {
		t.Errorf("Expecting 2 AdaptationSet element on first Period, but got %d", len(mpd.Periods[0].AdaptationSets))
	}
}

func TestDynamicParse(t *testing.T) {
	f, err := os.Open("./playlists/dynamic.mpd")
	if err != nil {
		t.Fatal(err)
	}
	mpd := &MPD{}

	if err := mpd.Parse(bufio.NewReader(f)); err != nil {
		t.Fatal(err)
	}

	pt, _ := time.Parse(time.RFC3339Nano, "2013-08-10T22:03:00Z")
	if !reflect.DeepEqual(mpd.PublishTime.Time, pt) {
		t.Errorf("Expected PublishTime %v, but got %v", pt, mpd.PublishTime)
	}
}

func TestMultiplePeriods(t *testing.T) {
	f, err := os.Open("./playlists/multipleperiods.mpd")
	if err != nil {
		t.Fatal(err)
	}
	mpd := &MPD{}

	if err := mpd.Parse(bufio.NewReader(f)); err != nil {
		t.Fatal(err)
	}
}
