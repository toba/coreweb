// Package header enumerates HTTP headers keys.
package header

const (
	Accept       = "Accept"
	CacheControl = "Cache-Control"
	Connection   = "Connection"
	DoNotTrack   = "dnt"
	eTag         = "ETag"
	Host         = "Host"
	// LastModified is the RFC1123 time the file was modified.
	// Example: Tue, 15 Nov 1994 12:45:26 GMT
	LastModified   = "Last-Modified"
	Origin         = "Origin"
	Referer        = "Referer"
	ResponseTime   = "Response-Time"
	RequestedWidth = "X-Requested-With"
	UserAgent      = "User-Agent"
	// Vary indicates header keys whose values can vary while still considering
	// the page to be cached.
	Vary = "Vary"
)
