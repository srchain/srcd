package node

import (

)

// ServiceContext is a collection of service independent options inherited from
// the protocol stack, that is passed to all constructors to be optionally used;
// as well as utility methods to operate on the service environment.
type ServiceContext struct {
	config         *Config
	services       map[reflect.Type]Service // Index of the already constructed services
	EventMux       *event.TypeMux           // Event multiplexer used for decoupled notifications
	AccountManager *accounts.Manager        // Account manager created by the node.
}

// ServiceConstructor is the function signature of the constructors needed to be
// registered for service instantiation.
type ServiceConstructor func(ctx *ServiceContext) (Service, error)

// Service is an individual protocol that can be registered into a node.
type Service interface {
	// // Protocols retrieves the P2P protocols the service wishes to start.
	// Protocols() []p2p.Protocol

	// // APIs retrieves the list of RPC descriptors the service provides
	// APIs() []rpc.API

	// // Start is called after all services have been constructed and the networking
	// // layer was also initialized to spawn any goroutines required by the service.
	// Start(server *p2p.Server) error

	// // Stop terminates all goroutines belonging to the service, blocking until they
	// // are all terminated.
	// Stop() error
}
