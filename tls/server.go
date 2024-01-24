package tls

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

// NetServer will serve a response on a socket with the provided server and listener,
// and might do some other functions like supporting TLS if configured.
type NetServer interface {
	// Serve will serve the http.Server at the net.Listener address,
	// but will perform optional other actions to support requirements,
	// such as setting up TLSConfig to support TLS.
	Serve(server *http.Server, listener net.Listener)
}

// TLSNetServer is a NetServer that adds TLS support to the provided http.Server and net.Listener.
type TLSNetServer struct {
	options *Options
}

var _ NetServer = (*TLSNetServer)(nil)

// NewTLSNetServer returns a new initialised NetServer or an error.
func NewTLSNetServer(options ...Option) (*TLSNetServer, error) {
	opt := &Options{
		CertificateGetter: nil,
		// InsecureSkipVerify should default to true
		InsecureSkipVerify: true,
	}

	for _, o := range options {
		o(opt)
	}

	// Setup default certificate getter if one is not set.
	if opt.CertificateGetter == nil {
		getter, err := selfSignedCertificateGetter(1*time.Hour, 10*time.Minute)
		if err != nil {
			return nil, err
		}

		opt.CertificateGetter = getter
	}

	return &TLSNetServer{
		options: opt,
	}, nil
}

// Serve will serve the http.Server at the net.Listener address,
// and will set up the TLSConfig to support TLS, then serves the server in a goroutine.
func (s *TLSNetServer) Serve(server *http.Server, listener net.Listener) {
	// Setup the config required for TLS.
	if server.TLSConfig == nil {
		server.TLSConfig = &tls.Config{}
	}

	server.TLSConfig.GetCertificate = s.options.CertificateGetter
	server.TLSConfig.InsecureSkipVerify = s.options.InsecureSkipVerify

	go server.ServeTLS(listener, "", "")
}

type DefaultNetServer struct {
}

var _ NetServer = (*DefaultNetServer)(nil)

// Serve serves the http server on the listener address in a goroutine.
func (s *DefaultNetServer) Serve(server *http.Server, listener net.Listener) {
	go server.Serve(listener)
}

type Options struct {
	CertificateGetter  func(*tls.ClientHelloInfo) (*tls.Certificate, error)
	InsecureSkipVerify bool
}

type Option func(*Options)

func WithSelfSignedCertificate(ttl time.Duration, wiggleRoom time.Duration) (Option, error) {
	var outerErr error
	return func(option *Options) {
		getter, err := selfSignedCertificateGetter(ttl, wiggleRoom)
		if err != nil {
			outerErr = err
		}

		option.CertificateGetter = getter
	}, outerErr
}
