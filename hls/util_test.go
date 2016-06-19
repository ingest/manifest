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

	err := rendition.writeXMedia(buf)
	//do something with buf
	fmt.Println(err)
	fmt.Println(buf.String())
	// t.Fatalf("ERR")
}

func TestWriteStreamInf(t *testing.T) {
	variant := &Variant{
		IsIframe:   true,
		URI:        "http://test.com",
		Bandwidth:  234000,
		Resolution: "230x400",
	}

	buf := new(bytes.Buffer)
	err := variant.writeStreamInf(7, buf)
	//do something with buf
	fmt.Println(err)
	fmt.Println(buf.String())
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
	rend := &Rendition{
		Type:     "AUDIO",
		GroupID:  "Testing",
		Name:     "Another test",
		Language: "English",
		Default:  false,
	}
	variant := &Variant{
		Renditions: []*Rendition{rendition, rend},
		IsIframe:   false,
		URI:        "http://test.com",
		Bandwidth:  234000,
		Resolution: "230x400",
		Codecs:     "These codescs",
	}

	rend3 := &Rendition{
		Type:     "VIDEO",
		GroupID:  "Test",
		Name:     "Bla",
		Language: "Portuguese",
	}

	variant2 := &Variant{
		Renditions: []*Rendition{rend3},
		IsIframe:   true,
		URI:        "thistest.com",
		Bandwidth:  145000,
	}

	p := NewMasterPlaylist(7)
	p.Variants = append(p.Variants, variant)
	p.Variants = append(p.Variants, variant2)
	p.SessionData = []*SessionData{&SessionData{DataID: "test", Value: "this is the session data"}}
	p.SessionKeys = []*Key{&Key{IsSession: true, Method: "aes-128", URI: "key url"}}
	p.IndependentSegments = true
	buf, err := p.GenerateManifest()
	//do something with buf
	fmt.Println(buf.String())
	fmt.Println(err)
	//t.Fatal("ERR")
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
