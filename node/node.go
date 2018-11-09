package node

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"github.com/srchain/srcd/accounts"
	"github.com/srchain/srcd/database"
	"github.com/srchain/srcd/log"
	"github.com/srchain/srcd/p2p"
	"github.com/prometheus/prometheus/util/flock"
)

// Node is a container on which services can be registered.
type Node struct {
	// eventmux *event.TypeMux		// Event multiplexer used between the services of a stack
	config            *Config
	accman            *accounts.Manager

	ephemeralKeystore string                   // if non-empty, the key directory that will be removed by Stop
	instanceDirLock   flock.Releaser           // prevents concurrent use of instance directory

	serverConfig      p2p.Config
	server            *p2p.Server		   // Currently running P2P networking layer

	serviceFuncs      []ServiceConstructor     // Service constructors (in dependency order)
	services          map[reflect.Type]Service // Currently running services

	// rpcAPIs       []rpc.API   // List of APIs currently provided by the node
	// inprocHandler *rpc.Server // In-process RPC request handler to process the API requests

	// ipcEndpoint string       // IPC endpoint to listen at (empty = IPC disabled)
	// ipcListener net.Listener // IPC RPC listener socket to serve API requests
	// ipcHandler  *rpc.Server  // IPC RPC request handler to process the API requests

	// httpEndpoint  string       // HTTP endpoint (interface + port) to listen at (empty = HTTP disabled)
	// httpWhitelist []string     // HTTP RPC modules to allow through this endpoint
	// httpListener  net.Listener // HTTP RPC listener socket to server API requests
	// httpHandler   *rpc.Server  // HTTP RPC request handler to process the API requests

	// wsEndpoint string       // Websocket endpoint (interface + port) to listen at (empty = websocket disabled)
	// wsListener net.Listener // Websocket RPC listener socket to server API requests
	// wsHandler  *rpc.Server  // Websocket RPC request handler to process the API requests

	stop chan struct{} // Channel to wait for termination notifications
	lock sync.RWMutex

	log log.Logger
}

// New creates a new P2P node, ready for protocol registration.
func New(conf *Config) (*Node, error) {
	// Copy config and resolve the datadir so future changes to the current
	// working directory don't affect the node.
	confCopy := *conf
	conf = &confCopy
	if conf.DataDir != "" {
		absdatadir, err := filepath.Abs(conf.DataDir)
		if err != nil {
			return nil, err
		}
		conf.DataDir = absdatadir
	}

	// Ensure that the instance name doesn't cause weird conflicts with
	// other files in the data directory.
	if strings.ContainsAny(conf.Name, `/\`) {
		return nil, errors.New(`Config.Name must not contain '/' or '\'`)
	}
	if conf.Name == datadirDefaultKeyStore {
		return nil, errors.New(`Config.Name cannot be "` + datadirDefaultKeyStore + `"`)
	}
	if strings.HasSuffix(conf.Name, ".ipc") {
		return nil, errors.New(`Config.Name cannot end in ".ipc"`)
	}

	// Ensure that the AccountManager method works before the node has started.
	am, ephemeralKeystore, err := makeAccountManager(conf)
	if err != nil {
		return nil, err
	}

	return &Node{
		accman:            am,
		ephemeralKeystore: ephemeralKeystore,
		config:            conf,
		serviceFuncs:      []ServiceConstructor{},
		// ipcEndpoint:       conf.IPCEndpoint(),
		// httpEndpoint:      conf.HTTPEndpoint(),
		// wsEndpoint:        conf.WSEndpoint(),
		// eventmux:          new(event.TypeMux),
	}, nil
}

// Register injects a new service into the node's stack. The service created by
// the passed constructor must be unique in its type with regard to sibling ones.
func (n *Node) Register(constructor ServiceConstructor) error {
	n.lock.Lock()
	defer n.lock.Unlock()

	// if n.server != nil {
	// return ErrNodeRunning
	// }

	n.serviceFuncs = append(n.serviceFuncs, constructor)
	return nil
}

// Start create a live P2P node and starts running it.
func (n *Node) Start() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	// Short circuit if the node's already running
	if n.server != nil {
		return ErrNodeRunning
	}
	if err := n.openDataDir(); err != nil {
		return err
	}

	// Initialize the p2p server. This creates the node key and
	// discovery databases.
	n.serverConfig = n.config.P2P
	n.serverConfig.PrivateKey = n.config.NodeKey()
	n.serverConfig.Name = n.config.NodeName()
	// n.serverConfig.Logger = n.log
	if n.serverConfig.StaticNodes == nil {
		n.serverConfig.StaticNodes = n.config.StaticNodes()
	}
	if n.serverConfig.TrustedNodes == nil {
		n.serverConfig.TrustedNodes = n.config.TrustedNodes()
	}
	if n.serverConfig.NodeDatabase == "" {
		n.serverConfig.NodeDatabase = n.config.NodeDB()
	}
	running := &p2p.Server{Config: n.serverConfig}
	// n.log.Info("Starting peer-to-peer node", "instance", n.serverConfig.Name)

	// Otherwise copy and specialize the P2P configuration
	services := make(map[reflect.Type]Service)
	for _, constructor := range n.serviceFuncs {
		// Create a new context for the particular service
		ctx := &ServiceContext{
			config:         n.config,
			services:       make(map[reflect.Type]Service),
			AccountManager: n.accman,
		}
		// copy needed for threaded access
		for kind, s := range services {
			ctx.services[kind] = s
		}
		// Construct and save the service
		service, err := constructor(ctx)
		if err != nil {
			return err
		}
		kind := reflect.TypeOf(service)
		if _, exists := services[kind]; exists {
			return &DuplicateServiceError{Kind: kind}
		}
		services[kind] = service
	}

	// Gather the protocols and start the freshly assembled P2P server
	for _, service := range services {
		running.Protocols = append(running.Protocols, service.Protocols()...)
	}
	if err := running.Start(); err != nil {
		return convertFileLockError(err)
	}

	// Start each of the services
	started := []reflect.Type{}
	for kind, service := range services {
		// Start the next service, stopping all previous upon failure
		// if err := service.Start(running); err != nil {
		if err := service.Start(running); err != nil {
			for _, kind := range started {
				services[kind].Stop()
			}
			// running.Stop()

			return err
		}
		// Mark the service started for potential cleanup
		started = append(started, kind)
	}

	//TODO rpc to be start

	// Finish initializing the startup
	n.services = services
	n.server = running
	n.stop = make(chan struct{})

	return nil
}

func (n *Node) openDataDir() error {
	if n.config.DataDir == "" {
		return nil
	}

	instdir := filepath.Join(n.config.DataDir, n.config.name())
	if err := os.MkdirAll(instdir, 0700); err != nil {
		return err
	}
	// Lock the instance directory to prevent concurrent use by another instance as well as
	// accidental use of the instance directory as a database.
	release, _, err := flock.New(filepath.Join(instdir, "LOCK"))
	if err != nil {
		return convertFileLockError(err)
	}
	n.instanceDirLock = release
	return nil
}

// Stop terminates a running peer along with all it's services. In the peer was
// not started, an error is returned.
func (n *Node) Stop() error {
	n.lock.Lock()
	defer n.lock.Unlock()

	// Short circuit if the node's not running
	if n.server == nil {
		return ErrNodeStopped
	}

	// Terminate the API, services and the p2p server.
	// n.stopWS()
	// n.stopHTTP()
	// n.stopIPC()
	// n.rpcAPIs = nil
	failure := &StopError{
		Services: make(map[reflect.Type]error),
	}
	for kind, service := range n.services {
		if err := service.Stop(); err != nil {
			failure.Services[kind] = err
		}
	}
	n.server.Stop()
	n.services = nil
	n.server = nil

	// Release instance directory lock.
	if n.instanceDirLock != nil {
		if err := n.instanceDirLock.Release(); err != nil {
			n.log.Error("Can't release datadir lock", "err", err)
		}
		n.instanceDirLock = nil
	}

	// unblock n.Wait
	close(n.stop)

	// Remove the keystore if it was created ephemerally.
	var keystoreErr error
	if n.ephemeralKeystore != "" {
		keystoreErr = os.RemoveAll(n.ephemeralKeystore)
	}

	if len(failure.Services) > 0 {
		return failure
	}
	if keystoreErr != nil {
		return keystoreErr
	}
	return nil
}

// Wait blocks the thread until the node is stopped. If the node is not running
// at the time of invocation, the method immediately returns.
func (n *Node) Wait() {
	n.lock.RLock()
	if n.server == nil {
		n.lock.RUnlock()
		return
	}
	stop := n.stop
	n.lock.RUnlock()

	<-stop
}

// Service retrieves a currently running service registered of a specific type.
func (n *Node) Service(service interface{}) error {
	n.lock.RLock()
	defer n.lock.RUnlock()

	// Short circuit if the node's not running
	if n.server == nil {
		return ErrNodeStopped
	}
	// Otherwise try to find the service to return
	element := reflect.ValueOf(service).Elem()
	if running, ok := n.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return ErrServiceUnknown
}

// AccountManager retrieves the account manager used by the protocol stack.
func (n *Node) AccountManager() *accounts.Manager {
	return n.accman
}

// OpenDatabase opens an existing database with the given name (or creates one if no
// previous can be found) from within the node's instance directory. If the node is
// ephemeral, a memory database is returned.
func (n *Node) OpenDatabase(name string, cache, handles int) (database.Database, error) {
	if n.config.DataDir == "" {
		return database.NewMemDatabase(), nil
	}
	return database.NewLDBDatabase(n.config.resolvePath(name), cache, handles)
}
