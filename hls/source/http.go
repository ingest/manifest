package source

import (
	"context"
	"io"
	"net/http"

	"github.com/ingest/manifest/hls"
)

type httpSource struct {
	Client *http.Client
}

// HTTP returns a source interface that fetches content using HTTP
func HTTP(c *http.Client) hls.Source {
	if c == nil {
		c = http.DefaultClient
	}

	return &httpSource{
		Client: c,
	}
}

// Master will download, and attempt to parse the document at the URI into a HLS master playlist.
// It must be a HTTP accessible address using the provided http.Client.
func (s *httpSource) Master(ctx context.Context, uri string) (*hls.MasterPlaylist, error) {
	master := hls.NewMasterPlaylist(0)
	master.URI = uri

	req, err := master.Request()
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if err := master.Parse(res.Body); err != nil {
		return nil, err
	}
	return master, nil
}

// Media will download, and attempt to parse the HLS media playlist from the given variant that was parsed from a master playlist.
// It must be a HTTP accessible address using the provided http.Client.
func (s *httpSource) Media(ctx context.Context, variant *hls.Variant) (*hls.MediaPlaylist, error) {
	req, err := variant.Request()
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	media := hls.NewMediaPlaylist(0).WithVariant(variant)
	if err := media.Parse(res.Body); err != nil {
		return nil, err
	}

	return media, nil
}

// Resource will download, and return the http.Response.Body for further parsing for whatever the structure might be.
// Some examples of a resource might be the actual media segment, or session decryption key.
func (s *httpSource) Resource(ctx context.Context, uri string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
