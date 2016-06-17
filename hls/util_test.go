package hls

import (
	"bytes"
	"fmt"
	"testing"
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

	buf := new(bytes.Buffer)

	rendition.writeXMedia(buf)
	//do something with buf
	// fmt.Println(buf.String())
	// t.Fatalf("ERR")
}

func TestWriteStreamInf(t *testing.T) {
	variant := &Variant{
		IsIframe:   false,
		URI:        "http://test.com",
		Bandwidth:  234000,
		Resolution: "230x400",
	}

	buf := new(bytes.Buffer)
	variant.writeStreamInf(7, buf)
	//do something with buf
	// fmt.Println(buf.String())
	// t.Fatal("ERR")
}

func TestGenerateMasterPlaylist(t *testing.T) {
	rendition := &Rendition{
		Type:     "VIDEO",
		GroupID:  "TestID",
		Name:     "Test",
		Language: "English",
		Default:  true,
		URI:      "http://test.com",
	}
	variant := &Variant{
		Renditions: []*Rendition{rendition},
		IsIframe:   false,
		URI:        "http://test.com",
		Bandwidth:  234000,
		Resolution: "230x400",
	}

	p := NewMasterPlaylist(7)
	p.Variants = append(p.Variants, variant)
	p.SessionData = []*SessionData{&SessionData{DataID: "test", Value: "this is the session data"}}
	p.SessionKeys = []*Key{&Key{Method: "aes-128"}}

	buf, _ := p.GenerateManifest()
	//do something with buf
	fmt.Println(buf.String())
	// t.Fatal("ERR")
}

func TestCompatibilityCheck(t *testing.T) {
	p := NewMediaPlaylist(4)
	s := &Segment{
		Key: &Key{
			Method: "sample-aes",
		},
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
}
