package dash

import (
	"bufio"
	"os"
	"reflect"
	"testing"
	"time"
)

//TODO:Check for every struct field possible.
func TestStaticParse(t *testing.T) {
	f, err := os.Open("./playlists/static.mpd")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

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
	defer f.Close()

	mpd := &MPD{}

	if err := mpd.Parse(bufio.NewReader(f)); err != nil {
		t.Fatal(err)
	}

	pt, _ := time.Parse(time.RFC3339Nano, "2013-08-10T22:03:00Z")
	if !reflect.DeepEqual(mpd.PublishTime.Time, pt) {
		t.Errorf("Expected PublishTime %v, but got %v", pt, mpd.PublishTime)
	}
	d, _ := time.ParseDuration("0h10m54.00s")
	if !reflect.DeepEqual(mpd.MediaPresDuration.Duration, d) {
		t.Errorf("Expecting MediaPresDuration to be %v, but got %v", d, mpd.MediaPresDuration)
	}

	if len(mpd.Periods) != 1 {
		t.Errorf("Expecting 1 Period element, but got %d", len(mpd.Periods))
	}

	if len(mpd.Periods[0].AdaptationSets) != 2 {
		t.Errorf("Expecting 2 AdaptationSets, but got %d", len(mpd.Periods[0].AdaptationSets))
	}

	if len(mpd.Periods[0].AdaptationSets[0].Representations) != 3 {
		t.Errorf("Expecting 3 Representations of AdaptationSets[0], but got %d", len(mpd.Periods[0].AdaptationSets[0].Representations))
	}

	if len(mpd.Periods[0].AdaptationSets[1].Representations) != 1 {
		t.Errorf("Expecting 3 Representations of AdaptationSets[1], but got %d", len(mpd.Periods[0].AdaptationSets[1].Representations))
	}

	if len(mpd.Periods[0].AdaptationSets[1].Representations[0].AudioChannelConfig) != 1 {
		t.Errorf("Expecting 1 AudioChannelConfig, but got %d", len(mpd.Periods[0].AdaptationSets[1].Representations[0].AudioChannelConfig))
	}
}
