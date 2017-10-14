// Package access enumerates HTTP header Access-Control-* keys.
package access

const (
	prefix  = "Access-Control-"
	allow   = prefix + "Allow-"
	request = prefix + "Request-"
)
const (
	MaxAge           = prefix + "Max-Age"
	AllowCredentials = allow + "Credentials"
	AllowHeaders     = allow + "Headers"
	AllowMethods     = allow + "Methods"
	AllowOrigin      = allow + "Origin"
	RequestHeaders   = request + "Headers"
	RequestMethod    = request + "Method"
)
