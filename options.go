package web

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/divisionone/cli"
	"github.com/divisionone/go-micro/cmd"
	"github.com/divisionone/go-micro/registry"
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

	TLSOptions *TLSOptions

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

// TLSOptions are options related to using TLS.
type TLSOptions struct {
	Enabled bool
	// TimeToLive is how long until the self-signed certificate used for TLS expires.
	TimeToLive   time.Duration
	// WiggleRoom defines a level of leeway for valid / expired certificates to help handle time differences between
	// servers.
	WiggleRoom time.Duration
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
		TLSOptions: 	  &TLSOptions{
			TimeToLive: 1 * time.Hour,
			WiggleRoom: 10 * time.Minute,
		},
	}

	for _, o := range opts {
		o(&opt)
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
func WithTLSEnabled(enabled bool) Option {
	return func(o *Options) {
		o.TLSOptions.Enabled = enabled
	}
}

// WithTLSTimeToLive sets the time to live on a self-signed certificate used for TLS.
func WithTLSTimeToLive(ttl time.Duration) Option {
	return func(o *Options) {
		o.TLSOptions.TimeToLive = ttl
	}
}

// WithTLSWiggleRoom sets the level of leeway for valid / expired certificates to help handle time differences between
// servers.
func WithTLSWiggleRoom(wiggle time.Duration) Option {
	return func(o *Options) {
		o.TLSOptions.WiggleRoom = wiggle
	}
}
