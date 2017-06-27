package manifest

import (
	"bytes"
	"io"
)

//BufWrapper is a wrapper type for bytes.Buffer.
type BufWrapper struct {
	Buf *bytes.Buffer
	Err error
}

//NewBufWrapper returns an instance of BufWrapper.
func NewBufWrapper() *BufWrapper {
	return &BufWrapper{
		Buf: new(bytes.Buffer),
		Err: nil,
	}
}

//WriteValidString receives an interface and performs a buffer.Write if data is set.
//It returns true if value is set, and false if it isn't.
func (b *BufWrapper) WriteValidString(data interface{}, write string) bool {
	if b.Err != nil {
		return false
	}
	switch data.(type) {
	case string:
		if data.(string) != "" {
			_, b.Err = b.Buf.WriteString(write)
			return true
		}
	case float64:
		if data.(float64) > float64(0) {
			_, b.Err = b.Buf.WriteString(write)
			return true
		}
	case int64:
		if data.(int64) > 0 {
			_, b.Err = b.Buf.WriteString(write)
			return true
		}
	case int:
		if data.(int) > 0 {
			_, b.Err = b.Buf.WriteString(write)
			return true
		}
	case bool:
		if data.(bool) {
			_, b.Err = b.Buf.WriteString(write)
			return true
		}
	}

	return false
}

//WriteString wraps buffer.WriteString
func (b *BufWrapper) WriteString(s string) {
	if b.Err != nil {
		return
	}
	_, b.Err = b.Buf.WriteString(s)
}

//WriteRune wraps buffer.WriteRune
func (b *BufWrapper) WriteRune(r rune) {
	if b.Err != nil {
		return
	}
	_, b.Err = b.Buf.WriteRune(r)
}

//Write wraps buffer.Write
func (b *BufWrapper) Write(p []byte) {
	if b.Err != nil {
		return
	}
	_, b.Err = b.Buf.Write(p)
}

// ReadFrom wraps buffer.ReadFrom
func (b *BufWrapper) ReadFrom(r io.Reader) (int64, error) {
	if b.Err != nil {
		return 0, b.Err
	}

	var count int64
	count, b.Err = b.Buf.ReadFrom(r)
	return count, b.Err
}

// ReadString wraps buffer.ReadString
func (b *BufWrapper) ReadString(delim byte) (line string) {
	if b.Err != nil {
		return
	}

	line, b.Err = b.Buf.ReadString(delim)
	return
}
