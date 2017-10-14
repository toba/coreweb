// Package web manages the HTTP server and constants.
package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"toba.tech/app/lib/auth"
	"toba.tech/app/lib/config"
	"toba.tech/app/lib/web/encoding"
	"toba.tech/app/lib/web/file"
	"toba.tech/app/lib/web/header/accept"

	// load embedded assets
	_ "toba.tech/app/assets"
)

const (
	osSlash       = string(os.PathSeparator)
	webSlash      = "/"
	templatePath  = "html" + webSlash + "template.html"
	templateToken = "{name}"
)

// webPath converts operating system to URL path.
func webPath(path string) string {
	return strings.Replace(path, osSlash, webSlash, -1)
}

// Handle responds to all HTTP requests. Endpoints are created for all files
// discovered in the configured path or zip file. Endpoints are also created
// for all module paths and authentication provider callbacks.
//
// After initialization, the handler does no routing or file system reads.
// Instead, modules perform client-side routing and retrieve data through web
// socket connections. This simplifies caching and security.
//
// 	https://cryptic.io/go-http/
//
func Handle(c config.HTTP, modulePaths []string, authPaths map[string]*auth.AuthProvider) func(w http.ResponseWriter, r *http.Request) {
	cache := &file.Map{Files: make(map[string]*file.Info)}
	var (
		m   *file.Map
		err error
	)

	if file.HasZipData() && c.FromFolder == "" {
		m, err = file.InZipFile()
	} else {
		// read all files in folder
		m, err = file.InFolder(c.FromFolder, true)
	}
	ExitIfError(err)
	log.Printf("Caching %d static files", len(m.Files))
	err = m.Read(true)
	ExitIfError(err)

	for k, v := range m.Files {
		cache.Files[webPath(k)] = v
	}

	if _, there := cache.Files[templatePath]; len(cache.Files) < 2 || !there {
		log.Fatalf("Invalid Template (%d files) for folder \"%s\"", len(cache.Files), c.FromFolder)
	}

	// make single reference to template and remove it from cache array
	template := cache.Files[templatePath]
	delete(cache.Files, templatePath)

	// add cache entry for template rendered for each module
	for _, name := range modulePaths {
		log.Printf("Adding module endpoint /%s", name)
		cache.Files[name] = template.Replace(templateToken, name)
	}

	if c.SyncFileAccess {
		file.Monitor(cache)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
				debug.PrintStack()
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		path := strings.TrimPrefix(r.RequestURI, webSlash)
		path = strings.TrimSuffix(path, webSlash)

		if c.SyncFileAccess {
			cache.RLock()
			defer cache.RUnlock()
		}

		info, exists := cache.Files[path]

		if !exists {
			// see if request path includes view name like /<app>/<view-name>
			path = strings.Split(path, webSlash)[0]
			info, exists = cache.Files[path]
		}

		if exists {
			allowGZip := strings.Contains(r.Header.Get(accept.Encoding), encoding.GZip)

			if allowGZip && info.Compressed != nil {
				info = info.Compressed
			}

			for k, v := range info.Header {
				w.Header().Set(k, v)
			}
			w.Write(info.Content)
		} else {
			http.Error(w, r.RequestURI+" does not exist", http.StatusNotFound)
		}
	}
}

// ExitIfError logs error if non-nil and exits program.
func ExitIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
