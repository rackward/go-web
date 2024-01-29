package web

import (
	"net"
	"net/http"

	"github.com/sirupsen/logrus"
)

// NetServer handles the initialisation of the http.Server by making sure it calls the correct startup method
// based on the provided http.Server configuration.
type NetServer struct {
	log logrus.FieldLogger
}

// NewNetServer returns a new initialised NetServer.
func NewNetServer(log logrus.FieldLogger) *NetServer {
	return &NetServer{
		log: log,
	}
}

// Serve helps to ensure that the correct server startup method is called based on the provided http.
// Server configuration.
func (s *NetServer) Serve(server *http.Server, listener net.Listener) {
	if server.TLSConfig != nil {
		go func() {
			s.log.Info("serving with TLS")

			if err := server.ServeTLS(listener, "", ""); err != nil {
				s.log.WithError(err).Error("Error serving on TLS server")
			}
		}()
	} else {
		go server.Serve(listener)
	}
}
