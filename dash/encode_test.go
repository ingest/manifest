package dash

import (
	"bufio"
	"os"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		Case string
		File string
	}{
		{"Static MPD", "./playlists/static.mpd"},
		{"Static2 MPD", "./playlists/static2.mpd"},
		{"Dynamic MPD", "./playlists/dynamic.mpd"},
		{"Encrypted MPD", "./playlists/encrypted.mpd"},
		{"Event Message MPD", "./playlists/eventmessage.mpd"},
		{"Multiple Periods MPD", "./playlists/multipleperiods.mpd"},
		{"Content Protection MPD", "./playlists/contentprotection.mpd"},
		{"Trick Play MPD", "./playlists/trickplay.mpd"},
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

		// buf := new(bytes.Buffer)
		// buf.ReadFrom(o)
		// fmt.Println(buf.String())
		// t.Error("Err")

		//Parse from previous encoded result into new struct
		newM := &MPD{}
		newM.Parse(o)

		//Both structs must be the same
		if !reflect.DeepEqual(m, newM) {
			t.Errorf("Case: %s - Expected newM:\n %v \n but got: \n%v", tt.Case, m, newM)
		}
	}
}
