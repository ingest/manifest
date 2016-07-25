package dash

import (
	"strings"
	"testing"
	"time"

	"stash.redspace.com/ing/manifest"
)

func TestNoPeriodError(t *testing.T) {
	mpd := NewMPD("profile", time.Second*5)
	_, err := mpd.Encode()
	if err.Error() != "MPD must have at least one Period element.\n" {
		t.Fatalf("Expecting Period requirement error, but got %s", err.Error())
	}
}

func TestError(t *testing.T) {
	mpd := NewMPD("", time.Second*0)
	period := &Period{
		ID: "1",
	}
	mpd.Periods = append(mpd.Periods, period)
	_, err := mpd.Encode()
	if !strings.Contains(err.Error(), "MPD field Profiles is required") {
		t.Errorf("Expecting Profiles requirement error, but got %s", err.Error())
	}
	if !strings.Contains(err.Error(), "MPD field MinBufferTime is required") {
		t.Errorf("Expecting MinBufferTime requirement error, but got %s", err.Error())
	}
}

func TestRepresentation(t *testing.T) {
	rep := Representation{
		SegmentBase: &SegmentBase{Timescale: 1, IndexRangeExact: true},
		SegmentList: &SegmentList{Timescale: 2},
	}
	buf := manifest.NewBufWrapper()
	rep.validate(buf)
	if !strings.Contains(buf.Buf.String(), "IndexRangeExact must not be present") {
		t.Error("Expecting 'IndexRangeExact must not be present' error")
	}
	if !strings.Contains(buf.Buf.String(), "At most one of the three") {
		t.Error("Expecting 'At most one of the three, SegmentBase, SegmentTemplate and SegmentList' error")
	}
}
