package peer

import (

)

// Peer is a container on which services can be registered.
type Peer struct {
	eventmux *event.TypeMux 		// Event multiplexer used between the services of a stack
	config   *Config
	accman   *accounts.Manager

	ephemeralKeystore string         	// if non-empty, the key directory that will be removed by Stop
	instanceDirLock   flock.Releaser 	// prevents concurrent use of instance directory

	serverConfig p2p.Config
	server       *p2p.Server 		// Currently running P2P networking layer

	// serviceFuncs []ServiceConstructor     	// Service constructors (in dependency order)
	// services     map[reflect.Type]Service 	// Currently running services

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

	// log log.Logger
}

// Start create a live P2P peer and starts running it.
func (peer *Peer) Start() error {
	peer.lock.Lock()
	defer peer.lock.Unlock()

	// Short circuit if the node's already running
	if peer.server != nil {
		return ErrNodeRunning
	}
	if err := peer.openDataDir(); err != nil {
		return err
	}

	// init p2p server config

	// create p2p server
	running := &p2p.Server{Config: peer.serverConfig}

	// run p2p server
	if err := running.Start(); err != nil {
		return convertFileLockError(err)
	}

	// rpc

	// Finish initializing the startup
	peer.server = running
	peer.stop = make(chan struct{})

	return nil
}

func (peer *peer) openDataDir() error {
	if peer.config.DataDir == "" {
		return nil
	}

	instdir := filepath.Join(peer.config.DataDir, peer.config.name())
	if err := os.MkdirAll(instdir, 0700); err != nil {
		return err
	}
	// Lock the instance directory to prevent concurrent use by another instance as well as
	// accidental use of the instance directory as a database.
	release, _, err := flock.New(filepath.Join(instdir, "LOCK"))
	if err != nil {
		return convertFileLockError(err)
	}
	peer.instanceDirLock = release
	return nil
}

// Stop terminates a running peer along with all it's services. In the peer was
// not started, an error is returned.
func (peer *Peer) Stop() error {
	peer.lock.Lock()
	defer peer.lock.Unlock()

	// Short circuit if the node's not running
	if peer.server == nil {
		return ErrNodeStopped
	}

	// Terminate the API, services and the p2p server.
	// peer.stopWS()
	// peer.stopHTTP()
	// peer.stopIPC()
	// peer.rpcAPIs = nil

	// failure := &StopError{
		// Services: make(map[reflect.Type]error),
	// }
	// for kind, service := range n.services {
		// if err := service.Stop(); err != nil {
			// failure.Services[kind] = err
		// }
	// }
	peer.server.Stop()
	peer.services = nil
	peer.server = nil

	// Release instance directory lock.
	if peer.instanceDirLock != nil {
		if err := peer.instanceDirLock.Release(); err != nil {
			peer.log.Error("Can't release datadir lock", "err", err)
		}
		peer.instanceDirLock = nil
	}

	// unblock peer.Wait
	close(peer.stop)

	return nil
}

// Wait blocks the thread until the node is stopped. If the node is not running
// at the time of invocation, the method immediately returns.
func (peer *Peer) Wait() {
	peer.lock.RLock()
	if peer.server == nil {
		peer.lock.RUnlock()
		return
	}
	stop := peer.stop
	peer.lock.RUnlock()

	<-stop
}

// OpenDatabase opens an existing database with the given name (or creates one if no
// previous can be found) from within the node's instance directory. If the node is
// ephemeral, a memory database is returned.
func (peer *Peer) OpenDatabase(name string, cache, handles int) (ethdb.Database, error) {
	if peer.config.DataDir == "" {
		return ethdb.NewMemDatabase(), nil
	}
	return ethdb.NewLDBDatabase(peer.config.resolvePath(name), cache, handles)
}
