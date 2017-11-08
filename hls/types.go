package hls

import (
	"context"
	"fmt"
	"io"
	"net/url"
)

const (
	sub     = "SUBTITLES"
	aud     = "AUDIO"
	vid     = "VIDEO"
	cc      = "CLOSED-CAPTIONS"
	aes     = "AES-128"
	none    = "NONE"
	sample  = "SAMPLE-AES"
	boolYes = "YES"
	boolNo  = "NO"
)

// Source represents how you can fetch the components of a HLS manifest from different locations
type Source interface {
	Master(ctx context.Context, uri string) (*MasterPlaylist, error)
	Media(ctx context.Context, variant *Variant) (*MediaPlaylist, error)
	Resource(ctx context.Context, uri string) (io.ReadCloser, error)
}

func resolveURLReference(base, sub string) (string, error) {
	ref, err := url.Parse(sub)
	if err != nil {
		return "", fmt.Errorf("failed to parse subresource uri: %v", err)
	}
	if ref.IsAbs() {
		return ref.String(), nil
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	return baseURL.ResolveReference(ref).String(), nil
}
