package util

import "github.com/micro/go-micro/registry"
import "github.com/micro/go-plugins/registry/consul"

func ConsulRegistry() registry.Registry {
	return consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
}
