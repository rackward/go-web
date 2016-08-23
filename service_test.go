package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/mock"
)

func TestService(t *testing.T) {
	str := `<html><body><h1>Hello World</h1></body></html>`

	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, str)
	}

	reg := mock.NewRegistry()

	service := NewService(
		Name("go.micro.web.test"),
		Registry(reg),
	)

	service.HandleFunc("/", fn)

	go func() {
		if err := service.Run(); err != nil {
			t.Fatal(err)
		}
	}()

	var s []*registry.Service

	eventually(func() bool {
		var err error
		s, err = reg.GetService("go.micro.web.test")
		return err == nil
	}, t.Fatal)

	if have, want := len(s), 1; have != want {
		t.Fatalf("Expected %d but got %d services", want, have)
	}

	rsp, err := http.Get(fmt.Sprintf("http://%s:%d", s[0].Nodes[0].Address, s[0].Nodes[0].Port))
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()

	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(b) != str {
		t.Errorf("Expected %s got %s", str, string(b))
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM)

	syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	<-ch

	eventually(func() bool {
		_, err := reg.GetService("go.micro.web.test")
		return err == registry.ErrNotFound
	}, t.Error)
}

func TestOptions(t *testing.T) {
	var (
		name             = "service-name"
		id               = "service-id"
		version          = "service-version"
		address          = "service-addr"
		advertise        = "service-adv"
		registry         = mock.NewRegistry()
		registerTTL      = 123 * time.Second
		registerInterval = 456 * time.Second
		handler          = http.NewServeMux()
		metadata         = map[string]string{"key": "val"}
	)

	service := NewService(
		Name(name),
		Id(id),
		Version(version),
		Address(address),
		Advertise(advertise),
		Registry(registry),
		RegisterTTL(registerTTL),
		RegisterInterval(registerInterval),
		Handler(handler),
		Metadata(metadata),
	)

	opts := service.Options()

	tests := []struct {
		subject string
		want    interface{}
		have    interface{}
	}{
		{"name", name, opts.Name},
		{"version", version, opts.Version},
		{"id", id, opts.Id},
		{"address", address, opts.Address},
		{"advertise", advertise, opts.Advertise},
		{"registry", registry, opts.Registry},
		{"registerTTL", registerTTL, opts.RegisterTTL},
		{"registerInterval", registerInterval, opts.RegisterInterval},
		{"handler", handler, opts.Handler},
		{"metadata", metadata["key"], opts.Metadata["key"]},
	}

	for _, tc := range tests {
		if tc.want != tc.have {
			t.Errorf("unexpected %s: want %v, have %v", tc.subject, tc.want, tc.have)
		}
	}
}

func eventually(pass func() bool, fail func(...interface{})) {
	tick := time.NewTicker(10 * time.Millisecond)
	defer tick.Stop()

	timeout := time.After(time.Second)

	for {
		select {
		case <-timeout:
			fail("timed out")
			return
		case <-tick.C:
			if pass() {
				return
			}
		}
	}
}
