package service

import (
	"context"
	"time"
)

type Options struct {
	Name      string
	Version   string
	Id        string
	Metadata  map[string]string
	Address   string
	Advertise string
	Namespace string

	RegistryAddr     string
	RegisterTTL      time.Duration
	RegisterInterval time.Duration

	// Alternative Options
	Context context.Context

	BeforeStart []func() error
	BeforeStop  []func() error
	AfterStart  []func() error
	AfterStop   []func() error
}

type Option func(*Options)

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

func newOptions(name, version string) Options {
	opt := Options{
		Name:    name,
		Version: version,
		Id:      DefaultId,
		Address: DefaultAddress,
		Context: context.TODO(),
	}
	return opt
}
