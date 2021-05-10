package command

import (
	"sync"

	"github.com/mitchellh/cli"

	"deploy/restful/service"
)

var (
	mutex       sync.RWMutex
	serviceList = make(map[string]cli.CommandFactory)
)

// registry service type
func RegistryService(service string, factory cli.CommandFactory) {
	mutex.Lock()
	defer mutex.Unlock()

	if factory == nil {
		panic("invalid service command factory")
	}
	if _, dup := serviceList[service]; dup {
		panic("register called twice for service " + service)
	}
	serviceList[service] = factory
}

// Commands returns the mapping of CLI commands for App. The meta
// parameter lets you set meta options for all commands.
func MakeCommands(metaPtr *service.Meta) map[string]cli.CommandFactory {
	return serviceList
}
