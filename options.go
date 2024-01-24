package web

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/divisionone/cli"
	"github.com/divisionone/go-micro/cmd"
	"github.com/divisionone/go-micro/registry"
	"github.com/rackward/go-web/tls"
)

type Options struct {
	Name      string
	Version   string
	Id        string
	Metadata  map[string]string
	Address   string
	Advertise string

	RegisterTTL      time.Duration
	RegisterInterval time.Duration

	// NetServer is a hack that allows to optionally support TLS with our existing http.Server and net.Listener.
	NetServer tls.NetServer

	Listen  func(network, address string) (net.Listener, error)
	Server  *http.Server
	Handler http.Handler

	// Alternative Options
	Context context.Context

	Cmd         cmd.Cmd
	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error
}

func newOptions(opts ...Option) Options {
	opt := Options{
		Name:             DefaultName,
		Version:          DefaultVersion,
		Id:               DefaultId,
		Address:          DefaultAddress,
		RegisterTTL:      DefaultRegisterTTL,
		RegisterInterval: DefaultRegisterInterval,
		Cmd:              cmd.DefaultCmd,
		Context:          context.TODO(),
		Listen:           net.Listen,
	}

	for _, o := range opts {
		o(&opt)
	}

	// Setup default NetServer if one hasn't been set.
	// This will wrap around the http.Server and optionally allows TLS if TLSNetServer is used.
	if opt.NetServer == nil {
		opt.NetServer = &tls.DefaultNetServer{}
	}

	return opt
}

func Listen(fn func(network, address string) (net.Listener, error)) Option {
	return func(o *Options) {
		o.Listen = fn
	}
}

// Server name
func Name(n string) Option {
	return func(o *Options) {
		o.Name = n
	}
}

// Unique server id
func Id(id string) Option {
	return func(o *Options) {
		o.Id = id
	}
}

// Version of the service
func Version(v string) Option {
	return func(o *Options) {
		o.Version = v
	}
}

// Metadata associated with the service
func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

// Address to bind to - host:port
func Address(a string) Option {
	return func(o *Options) {
		o.Address = a
	}
}

// The address to advertise for discovery - host:port
func Advertise(a string) Option {
	return func(o *Options) {
		o.Advertise = a
	}
}

// Context specifies a context for the service.
// Can be used to signal shutdown of the service.
// Can be used for extra option values.
func Context(ctx context.Context) Option {
	return func(o *Options) {
		o.Context = ctx
	}
}

func RegisterTTL(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = t
	}
}

func RegisterInterval(t time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = t
	}
}

func Handler(h http.Handler) Option {
	return func(o *Options) {
		o.Handler = h
	}
}

func Server(srv *http.Server) Option {
	return func(o *Options) {
		o.Server = srv
	}
}

// Set registry to be used by the service
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		registry.DefaultRegistry = r
	}
}

// Cmd sets the command instance.
func Cmd(c cmd.Cmd) Option {
	return func(o *Options) {
		o.Cmd = c
	}
}

// Flags sets the command flags.
func Flags(flags ...cli.Flag) Option {
	return func(o *Options) {
		o.Cmd.App().Flags = append(o.Cmd.App().Flags, flags...)
	}
}

// Action sets the command action.
func Action(a func(*cli.Context)) Option {
	return func(o *Options) {
		o.Cmd.App().Action = a
	}
}

// BeforeStart is executed before the server starts.
func BeforeStart(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStart = append(o.BeforeStart, fn)
	}
}

// BeforeStop is executed before the server stops.
func BeforeStop(fn func() error) Option {
	return func(o *Options) {
		o.BeforeStop = append(o.BeforeStop, fn)
	}
}

// AfterStart is executed after server start.
func AfterStart(fn func() error) Option {
	return func(o *Options) {
		o.AfterStart = append(o.AfterStart, fn)
	}
}

// AfterStop is executed after server stop.
func AfterStop(fn func() error) Option {
	return func(o *Options) {
		o.AfterStop = append(o.AfterStop, fn)
	}
}

// WithTLSEnabled can be used to enable/disable TLS.
func WithNetServer(server tls.NetServer) Option {
	return func(o *Options) {
		o.NetServer = server
	}
}
