package hls

import (
	"bytes"
	"io"
)

//When support for DASH is implemented we should probably move this to m3u8 package

//BufWrapper is a wrapper type for bytes.Buffer.
type BufWrapper struct {
	buf *bytes.Buffer
	err error
}

//NewBufWrapper returns an instance of BufWrapper.
func NewBufWrapper() *BufWrapper {
	return &BufWrapper{
		buf: new(bytes.Buffer),
		err: nil,
	}
}

//WriteValidString receives an interface and performs a buffer.Write if data is set.
//If returns true if value is set, and false if it isn't.
func (b *BufWrapper) WriteValidString(data interface{}, write string) bool {
	if b.err != nil {
		return false
	}
	switch data.(type) {
	case string:
		if data.(string) != "" {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	case float64:
		if data.(float64) > float64(0) {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	case int64:
		if data.(int64) > 0 {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	case int:
		if data.(int) > 0 {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	case bool:
		if data.(bool) {
			_, b.err = b.buf.WriteString(write)
			return true
		}
	}

	return false
}

//WriteString wraps buffer.WriteString
func (b *BufWrapper) WriteString(s string) {
	if b.err != nil {
		return
	}
	_, b.err = b.buf.WriteString(s)
}

//WriteRune wraps buffer.WriteRune
func (b *BufWrapper) WriteRune(r rune) {
	if b.err != nil {
		return
	}
	_, b.err = b.buf.WriteRune(r)
}

//ReadFrom wraps buffer.ReadFrom
func (b *BufWrapper) ReadFrom(r io.Reader) {
	if b.err != nil {
		return
	}

	_, b.err = b.buf.ReadFrom(r)
}

//ReadString wraps buffer.ReadString
func (b *BufWrapper) ReadString(delim byte) (line string) {
	if b.err != nil {
		return
	}

	line, b.err = b.buf.ReadString(delim)
	return
}

//BySegID implements golang/sort interface to sort a Segment slice by Segment ID
type BySegID []*Segment

func (s BySegID) Len() int {
	return len(s)
}
func (s BySegID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s BySegID) Less(i, j int) bool {
	return s[i].ID < s[j].ID
}
