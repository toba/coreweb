package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"toba.tech/app/lib/config"
	"toba.tech/app/lib/db"
	"toba.tech/app/lib/host"
	"toba.tech/app/lib/ldap"
	"toba.tech/app/lib/license"
	"toba.tech/app/lib/module"
	"toba.tech/app/lib/web"
	"toba.tech/app/lib/web/file"
	"toba.tech/app/lib/web/socket"
	"toba.tech/app/modules/cogs"
	"toba.tech/app/modules/person"
	"toba.tech/app/modules/setup"
	"toba.tech/app/modules/system"
)

var (
	freeModules = []*module.ModuleInfo{
		system.Module,
		person.Module,
		setup.Module,
		cogs.Module,
	}

	flagSilent   = flag.Bool("silent", false, "Do not launch browser window when running setup mode.")
	flagDebug    = flag.Bool("debug", false, "Debug mode.")
	flagFiles    = flag.String("files", "", "File path")
	flagLocalURL = flag.String("local", "localhost", "Local URL to use for setup")
)

// main runs database migrations and initializes dependencies for the HTTP and
// web socket handlers. Web sockets handle all service calls defined by modules
// while module paths are made HTTP endpoints that load client React components.
func main() {
	flag.Parse()

	debug := *flagDebug

	if debug {
		log.Println("Starting in debug mode")
		file.Resolve(func() (string, error) {
			return os.Getenv("TOBA_PATH"), nil
		})
	}

	c, err := config.Load()
	web.ExitIfError(err)

	c.HTTP.SyncFileAccess = debug

	if debug {
		c.HTTP.FromFolder = "static"
	}

	license, err := license.Load()
	web.ExitIfError(err)

	if flagFiles != nil {
		folder := *flagFiles
		if folder != "" {
			c.HTTP.FromFolder = folder
		}
	}

	web.ExitIfError(err)
	web.ExitIfError(db.Initialize(c.Database))

	if license.Valid {
		serveLicensed(c)
	} else {
		serveSetup(c)
	}
	// TODO: add timeout properties to server
}

// serveLicensed runs the web server in licensed mode with LDAP integration and
// module endpoints.
func serveLicensed(c config.Server) {
	web.ExitIfError(system.Module.Migrate())
	web.ExitIfError(ldap.Initialize(c.LDAP))

	services, modulePaths := module.Amalgamate(freeModules)
	serviceHandler := module.Handle(services)

	log.Printf("Initializing %d modules and %d services", len(modulePaths), len(services))

	http.HandleFunc("/ws", socket.Handle(c.HTTP, serviceHandler))
	http.HandleFunc("/", web.Handle(c.HTTP, modulePaths, nil))
	log.Printf("Server starting on port %d", c.HTTP.Port)
	addr := fmt.Sprintf(":%d", c.HTTP.Port)

	// Consider http://goroutines.com/ssl
	web.ExitIfError(http.ListenAndServeTLS(addr, c.HTTP.SslCert, c.HTTP.SslKey, nil))
}

// serveSetup runs the web server in setup or trial mode.
func serveSetup(c config.Server) {
	services, modulePaths := module.Amalgamate(freeModules)
	serviceHandler := module.Handle(services)

	log.Printf("Initializing %d modules and %d services", len(modulePaths), len(services))

	http.HandleFunc("/ws", socket.Handle(c.HTTP, serviceHandler))
	http.HandleFunc("/", web.Handle(c.HTTP, modulePaths, nil))

	port := host.FirstAvailablePort(80, 8000, 3000)
	if port == 0 {
		log.Fatal("Unable to find bindable port")
	}
	addr := fmt.Sprintf(":%d", port)

	if !*flagSilent {
		url := *flagLocalURL
		go func() {
			time.Sleep(time.Second * 2)
			log.Printf("Launching browser for setup or demo")
			host.Start("http://" + url + addr + "/setup/")
		}()
	}

	log.Printf("Server starting on port %d", port)
	web.ExitIfError(http.ListenAndServe(addr, nil))
}
