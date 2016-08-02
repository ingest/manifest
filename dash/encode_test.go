package dash

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"
)

func getDefaultMPD() *MPD {
	mpd := NewMPD("urn:mpeg:dash:profile:isoff-live:2011", time.Second*2)
	period := &Period{
		AdaptationSets: AdaptationSets{
			&AdaptationSet{
				MaxWidth:  1280,
				MaxHeight: 720,
				Representations: Representations{&Representation{
					ID: "1", MimeType: "video/mp4", Codecs: "avc1.4d01f", Width: 1280, Height: 720,
					FrameRate: "24", Bandwidth: 980104,
					SegmentTemplate: &SegmentTemplate{
						Timescale: 12288,
						Duration:  24576,
						Media:     "video_$Number$.mp4"}},
				},
			},
		}}
	mpd.Periods = append(mpd.Periods, period)
	return mpd
}

func TestEncode(t *testing.T) {
	tests := []struct {
		Case string
		File string
	}{
		{"Static MPD", "./testdata/static.mpd"},
		{"Static2 MPD", "./testdata/static2.mpd"},
		{"Dynamic MPD", "./testdata/dynamic.mpd"},
		{"Encrypted MPD", "./testdata/encrypted.mpd"},
		{"Event Message MPD", "./testdata/eventmessage.mpd"},
		{"Multiple Periods MPD", "./testdata/multipleperiods.mpd"},
		{"Trick Play MPD", "./testdata/trickplay.mpd"},
	}

	for _, tt := range tests {
		f, err := os.Open(tt.File)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		//Parse from file
		m := &MPD{}
		if err = m.Parse(bufio.NewReader(f)); err != nil {
			t.Fatalf("%s - %s", tt.Case, err)
		}

		//Encode from m struct
		o, err := m.Encode()
		if err != nil {
			t.Fatalf("%s - %s", tt.Case, err)
		}

		//Parse from previous encoded result into new struct
		newM := &MPD{}
		newM.Parse(o)

		//Both structs must be the same
		if !reflect.DeepEqual(m, newM) {
			t.Errorf("Case: %s - Expected newM:\n %v \n but got: \n%v", tt.Case, m, newM)
		}
	}
}

func TestContentProtection(t *testing.T) {
	mpd := getDefaultMPD()

	cp := NewContentProtection("test.uri", "cenc", "somehash", "", "")
	mpd.Periods[0].AdaptationSets[0].CENCContentProtections = append(mpd.Periods[0].AdaptationSets[0].CENCContentProtections, cp)

	o, err := mpd.Encode()
	if err != nil {
		t.Fatal(err)
	}

	expect := `<?xml version="1.0" encoding="UTF-8"?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" profiles="urn:mpeg:dash:profile:isoff-live:2011" type="static" minBufferTime="PT2S">
 <Period>
  <AdaptationSet maxWidth="1280" maxHeight="720">
   <ContentProtection xmlns:cenc="urn:mpeg:cenc:2013" schemeIdUri="test.uri" value="cenc" cenc:default_KID="somehash"></ContentProtection>
   <Representation id="1" bandwidth="980104" width="1280" height="720" frameRate="24" mimeType="video/mp4" codecs="avc1.4d01f">
    <SegmentTemplate timescale="12288" duration="24576" media="video_$Number$.mp4"></SegmentTemplate>
   </Representation>
  </AdaptationSet>
 </Period>
</MPD>`

	buf := new(bytes.Buffer)
	buf.ReadFrom(o)

	if expect != buf.String() {
		t.Errorf("Expecting:\n%s \n but got \n%s", expect, buf.String())
	}
}

func TestContentProtectionPlayready(t *testing.T) {
	mpd := getDefaultMPD()

	cp := NewContentProtection("test.uri", "cenc", "somehash", "psshhash", "msprhash")
	mpd.Periods[0].AdaptationSets[0].CENCContentProtections = append(mpd.Periods[0].AdaptationSets[0].CENCContentProtections, cp)

	o, err := mpd.Encode()
	if err != nil {
		t.Fatal(err)
	}

	expect := `<?xml version="1.0" encoding="UTF-8"?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" profiles="urn:mpeg:dash:profile:isoff-live:2011" type="static" minBufferTime="PT2S">
 <Period>
  <AdaptationSet maxWidth="1280" maxHeight="720">
   <ContentProtection xmlns:cenc="urn:mpeg:cenc:2013" xmlns:mspr="urn:microsoft:playready" schemeIdUri="test.uri" value="cenc" cenc:default_KID="somehash">
    <cenc:pssh>psshhash</cenc:pssh>
    <mspr:pro>msprhash</mspr:pro>
   </ContentProtection>
   <Representation id="1" bandwidth="980104" width="1280" height="720" frameRate="24" mimeType="video/mp4" codecs="avc1.4d01f">
    <SegmentTemplate timescale="12288" duration="24576" media="video_$Number$.mp4"></SegmentTemplate>
   </Representation>
  </AdaptationSet>
 </Period>
</MPD>`

	buf := new(bytes.Buffer)
	buf.ReadFrom(o)

	if expect != buf.String() {
		t.Errorf("Expecting:\n%s \n but got \n%s", expect, buf.String())
	}
}

func TestContentProtectionTrackEncryption(t *testing.T) {
	mpd := getDefaultMPD()

	cp := NewContentProtection("test.uri", "cenc", "", "", "")
	cp.SetTrackEncryptionBox(8, "kidhexacode")
	mpd.Periods[0].AdaptationSets[0].CENCContentProtections = append(mpd.Periods[0].AdaptationSets[0].CENCContentProtections, cp)

	o, err := mpd.Encode()
	if err != nil {
		t.Fatal(err)
	}

	expect := `<?xml version="1.0" encoding="UTF-8"?>
<MPD xmlns="urn:mpeg:dash:schema:mpd:2011" profiles="urn:mpeg:dash:profile:isoff-live:2011" type="static" minBufferTime="PT2S">
 <Period>
  <AdaptationSet maxWidth="1280" maxHeight="720">
   <ContentProtection xmlns:mspr="urn:microsoft:playready" schemeIdUri="test.uri" value="cenc">
    <mspr:IsEncrypted>1</mspr:IsEncrypted>
    <mspr:IV_size>8</mspr:IV_size>
    <mspr:kid>kidhexacode</mspr:kid>
   </ContentProtection>
   <Representation id="1" bandwidth="980104" width="1280" height="720" frameRate="24" mimeType="video/mp4" codecs="avc1.4d01f">
    <SegmentTemplate timescale="12288" duration="24576" media="video_$Number$.mp4"></SegmentTemplate>
   </Representation>
  </AdaptationSet>
 </Period>
</MPD>`

	buf := new(bytes.Buffer)
	buf.ReadFrom(o)

	if expect != buf.String() {
		t.Errorf("Expecting:\n%s \n but got \n%s", expect, buf.String())
	}
}
