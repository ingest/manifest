package dash

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	f, err := os.Open("./playlists/testmpd.mpd")
	if err != nil {
		t.Fatal(err)
	}
	mpd := &MPD{}
	err = mpd.Parse(bufio.NewReader(f))
	fmt.Println(err)
	fmt.Println(mpd)
	if mpd.BaseURL != nil {
		for _, b := range mpd.BaseURL {
			fmt.Println(b)
		}
	}
	if mpd.ProgramInformation != nil {
		for _, p := range mpd.ProgramInformation {
			fmt.Println(p)
			fmt.Println(p.Title)
			fmt.Println(p.Source)
			fmt.Println(p.Copyright)
		}
	}
	if mpd.Metrics != nil {
		for _, m := range mpd.Metrics {
			fmt.Println(m.Metrics)
			for _, r := range m.Range {
				fmt.Println(r)
			}
			for _, re := range m.Reporting {
				fmt.Println(re)
			}
		}

	}

	if mpd.Periods != nil {
		for _, pe := range mpd.Periods {
			if pe.Subsets != nil {
				for _, ss := range pe.Subsets {
					fmt.Println(ss.Contains)
				}
			}
		}
	}
	//t.Error("Err")
}
