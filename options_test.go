package web

import (
	"testing"

	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/mock"
)

func TestDefaultRegistry(t *testing.T) {
	opts := newOptions()
	if want, have := registry.DefaultRegistry, opts.Registry; want != have {
		t.Errorf("unexpected registry: want %v, have %v", want, have)
	}
}

func TestRegistryOption(t *testing.T) {
	registry := mock.NewRegistry()
	opts := newOptions(Registry(registry))
	if want, have := registry, opts.Registry; want != have {
		t.Errorf("unexpected registry: want %v, have %v", want, have)
	}
}
