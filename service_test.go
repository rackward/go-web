package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/micro/go-micro/registry/mock"
)

func TestService(t *testing.T) {
	str := `<html><body><h1>Hello World</h1></body></html>`

	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, str)
	}

	registry := mock.NewRegistry()

	service := NewService(
		Name("go.micro.web.test"),
		Registry(registry),
	)

	service.HandleFunc("/", fn)

	go func() {
		if err := service.Run(); err != nil {
			t.Fatal(err)
		}
	}()

	// another ugly hack
	time.Sleep(time.Millisecond * 100)

	s, err := registry.GetService("go.micro.web.test")
	if err != nil {
		t.Fatal(err)
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
