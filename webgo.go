/*
WebGo is a lightweight framework for building web apps. It has a multiplexer,
middleware plugging mechanism & context management of its own. The primary goal
of WebGo is to get out of the developer's way as much as possible. i.e. it does
not enfore you to build your app in any particular pattern instead just helps you
get all the trivial things done faster and easier.

e.g.
1. Sharing your database handler among all the HTTP handlers.
2. Getting named URI parameters.
3. Multiplexer for regex matching of URI and such.
*/
package webgo

import (
	"net/http"
	"time"
)

// Start the server with the appropriate configurations
func Start(cfg *Config, router *Router, readTimeout, writeTimeout time.Duration) {
	host := cfg.Host
	httpshost := cfg.Host

	if len(cfg.Port) > 0 {
		host += ":" + cfg.Port
	}

	if len(cfg.HTTPSPort) > 0 {
		httpshost += ":" + cfg.HTTPSPort
	}

	if cfg.HTTPSOnly {
		if cfg.CertFile == "" {
			println("No certificate provided for HTTPS")
			return
		}

		if cfg.KeyFile == "" {
			println("No key file provided for HTTPS")
			return
		}

		httpsServer := &http.Server{
			Addr:         httpshost,
			Handler:      router,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
		}

		println("\nStarting HTTPS server, listening on '" + httpshost + "'")
		err := httpsServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			println("HTTPS Server exited with error:", err.Error())
		}
		return
	}

	if len(cfg.HTTPSPort) > 0 {
		if cfg.CertFile == "" {
			println("No certificate provided for HTTPS")
			return
		}

		if cfg.KeyFile == "" {
			println("No key file provided for HTTPS")
			return
		}

		if cfg.Port == cfg.HTTPSPort {
			println("HTTP & HTTPS cannot listen on the same port. [" + cfg.Port + "]")
			return
		}

		//Starting HTTPS server
		go func() {
			httpsServer := &http.Server{
				Addr:         httpshost,
				Handler:      router,
				ReadTimeout:  readTimeout,
				WriteTimeout: writeTimeout,
			}

			println("\nStarting HTTPS server, listening on '" + httpshost + "'")
			err := httpsServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
			if err != nil {
				println("HTTPS Server exited with error:", err.Error())
			}
			return
		}()
	}

	httpServer := &http.Server{
		Addr:         host,
		Handler:      router,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	println("\nStarting HTTP server, listening on '" + host + "'")
	err := httpServer.ListenAndServe()
	if err != nil {
		println("HTTP Server exited with error:", err.Error())
	}

}
