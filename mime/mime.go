// Package mimetype enumerates MIME types.
package mime

import "strings"

const (
	GIF          = "image/gif"
	HTML         = "text/html; charset=utf-8"
	Icon         = "image/x-icon"
	JavaScript   = "text/javascript"
	JPEG         = "image/jpeg"
	JSON         = "application/json"
	PNG          = "image/png"
	SVG          = "image/svg+xml"
	Text         = "text/plain; charset=utf-8"
	XML          = "text/xml"
	Raw          = "application/octet-stream"
	StyleSheet   = "text/css"
	OpenType     = "font/opentype"
	WebOpenFont  = "application/woff"
	WebOpenFont2 = "application/woff2"
	Compressed   = "application/x-compressed"
)

// Infer MIME type from file extension. Ignore added GZip extension if present.
func Infer(fileName string) string {
	parts := strings.Split(strings.ToLower(fileName), ".")
	ext := parts[len(parts)-1]

	if ext == "gz" && len(parts) > 2 {
		ext = parts[len(parts)-2]
	}

	switch strings.ToLower(ext) {
	case "css":
		return StyleSheet
	case "gif":
		return GIF
	case "htm":
		fallthrough
	case "html":
		return HTML
	case "ico":
		return Icon
	case "jpg":
		fallthrough
	case "jpeg":
		return JPEG
	case "js":
		return JavaScript
	case "json":
		return JSON
	case "otf":
		return OpenType
	case "svg":
		return SVG
	case "txt":
		return Text
	case "png":
		return PNG
	case "wof":
		fallthrough
	case "woff":
		return WebOpenFont
	case "woff2":
		return WebOpenFont2
	case "xml":
		return XML
	case "gz":
		fallthrough
	case "zip":
		return Compressed
	default:
		return Raw
	}
}
