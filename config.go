package coreweb

type Config struct {
	SslCert    string `json:"sslCert"`    // SslCert is the path and name of the SSL certificate file.
	SslKey     string `json:"sslKey"`     // SslKey is the path and name of the SSL key file.
	Port       int    `json:"port"`       // FromZip is the name of a zip file to serve content from rather than the file system.
	FromFolder string `json:"fromFolder"` // FromFolder is the folder containing web content.

	// SyncFileAccess indicates if RWMutex lock should be used when reading
	// the file cache. It should only be true while debugging when files might
	// be changing while the web server is active.
	SyncFileAccess bool
}
