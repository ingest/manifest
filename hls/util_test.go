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

	buf, _ := p.GenerateManifest()
	//do something with buf
	fmt.Println(buf.String())
	// t.Fatal("ERR")
}
