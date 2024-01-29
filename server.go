package web

import (
	"net"
	"net/http"

	log "github.com/divisionone/micro-go-log"
)

// netServer handles the initialisation of the http.Server by making sure it calls the correct startup method
// based on the provided http.Server configuration.
type netServer struct{}

// Serve helps to ensure that the correct server startup method is called based on the provided http.
// Server configuration.
func (s *netServer) Serve(server *http.Server, listener net.Listener) {
	if server.TLSConfig != nil {
		go func() {
			log.Log("serving with TLS")

			if err := server.ServeTLS(listener, "", ""); err != nil {
				log.Logf("Error serving on TLS server: %s", err.Error())
			}
		}()
	} else {
		go server.Serve(listener)
	}
}
