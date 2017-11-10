package source_test

import (
	"context"
	"net/http"

	"github.com/ingest/manifest/hls/source"
)

// HTTP implements the source interface for HTTP for retrieving HLS data.
// By default uses the http.DefaultClient if a nil pointer is passed.
func ExampleHTTP() {
	ctx := context.Background()
	src := source.HTTP(http.DefaultClient)

	// Fetching master and media playlists
	master, _ := src.Master(ctx, "https://example.com/hls-master.m3u8")
	media, _ := src.Media(ctx, master.Variants[0])

	// Fetching Session Key, could be used to decrypt segment below
	sURL, _ := master.SessionKeys[0].AbsoluteURL()
	sessionKey, _ := src.Resource(ctx, sURL)
	sessionKey.Close()

	// Fetch segment content
	sURL, _ = media.Segments[0].AbsoluteURL()
	segment, _ := src.Resource(ctx, sURL)
	segment.Close()
}
