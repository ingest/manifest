package hls_test

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/ingest/manifest/hls"
)

const chunkSize = 5120

func equal(r1, r2 io.Reader) bool {
	// Compare
	for {
		b1 := make([]byte, chunkSize)
		_, err1 := r1.Read(b1)

		b2 := make([]byte, chunkSize)
		_, err2 := r2.Read(b2)

		if err1 != nil || err2 != nil {
			if err1 == io.EOF && err2 == io.EOF {
				return bytes.Equal(b1, b2)
			} else if err1 == io.EOF || err2 == io.EOF {
				return false
			} else {
				log.Fatal(err1, err2)
			}
		}

		if !bytes.Equal(b1, b2) {
			return false
		}
	}
}

func TestIdempotentDecodeEncodeCycle(t *testing.T) {
	tests := []struct {
		file string
	}{
		{
			file: "apple-ios5-macOS10_7.m3u8",
		},
		{
			file: "apple-ios6-tvOS9.m3u8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			f, e := os.Open("./testdata/" + tt.file)
			if e != nil {
				t.Fatal(e)
			}
			defer f.Close()

			p := hls.NewMasterPlaylist(0)
			if err := p.Parse(f); err != nil && err != io.EOF {
				t.Fatal(err)
			}

			if _, err := f.Seek(0, 0); err != nil {
				t.Fatal(err)
			}

			output, err := p.Encode()
			if err != nil {
				t.Fatal(err)
			}

			if !equal(f, output) {
				t.Fatal("parse/decode not idempotent")
			}
		})
	}
}
