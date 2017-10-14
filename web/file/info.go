package file

import (
	"bytes"
	"strconv"
	"time"

	"compress/gzip"

	"github.com/toba/goweb/web/encoding"
	"github.com/toba/goweb/web/header"
	"github.com/toba/goweb/web/header/accept"
	"github.com/toba/goweb/web/header/content"
	"github.com/toba/goweb/web/mime"
)

// Info contains information about a file and optionally its byte content and
// a GZipped version.
type Info struct {
	Content    []byte
	Header     map[string]string
	Path       string
	Compressed *Info
	Modified   time.Time // only used if file watching is active (debug mode)
}

// compressibleTypes lists the MIME types that can be GZipped.
var compressibleTypes = [...]string{
	mime.HTML,
	mime.JavaScript,
	mime.JSON,
	mime.StyleSheet,
	mime.SVG,
	mime.Text,
	mime.XML,
}

// Replace creates new File Info where some content has been replaced.
func (info *Info) Replace(token, name string) *Info {
	i := &Info{
		Content:  bytes.Replace(info.Content, []byte(token), []byte(name), -1),
		Header:   info.Header,
		Path:     info.Path,
		Modified: info.Modified,
	}
	_ = i.Compress()

	return i
}

// copyHeader retrieves all header values as a string map.
func (info *Info) copyHeader() map[string]string {
	h := make(map[string]string)
	for k, v := range info.Header {
		h[k] = v
	}
	return h
}

// Compress GZips file content if it's compatible. Copy header values from
// uncompressed file. Compare
//
// https://github.com/gin-contrib/gzip/blob/master/gzip.go
func (info *Info) Compress() error {
	if info.Compressible() {
		var buffer bytes.Buffer
		gz := gzip.NewWriter(&buffer)

		if _, err := gz.Write(info.Content); err != nil {
			return err
		}
		if err := gz.Close(); err != nil {
			return err
		}

		zipped := buffer.Bytes()
		head := info.copyHeader()
		head[content.Encoding] = encoding.GZip
		head[header.Vary] = accept.Encoding
		head[content.Length] = strconv.FormatInt(int64(len(zipped)), 10)

		info.Compressed = &Info{
			Content:  zipped,
			Path:     info.Path,
			Header:   head,
			Modified: info.Modified,
		}
	}
	return nil
}

// Compressible indicates whether the file content can be compressed. Do not
// re-compress and do not compress types that are already compact.
func (info *Info) Compressible() bool {
	if info.Content == nil || info.Compressed != nil {
		return false
	}

	if enc, ok := info.Header[content.Encoding]; ok {
		if enc == encoding.GZip {
			return false
		}
	}

	mimeType := info.Header[content.Type]

	for _, t := range compressibleTypes {
		if mimeType == t {
			return true
		}
	}
	return false
}
