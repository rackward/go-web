package web

import (
	"net"
	"net/http"

	log "github.com/divisionone/micro-go-log"
)

// NetServer handles the initialisation of the http.Server by making sure it calls the correct startup method
// based on the provided http.Server configuration.
type NetServer struct{}

// NewNetServer returns a new initialised NetServer.
func NewNetServer() *NetServer {
	return &NetServer{}
}

// Serve helps to ensure that the correct server startup method is called based on the provided http.
// Server configuration.
func (s *NetServer) Serve(server *http.Server, listener net.Listener) {
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
