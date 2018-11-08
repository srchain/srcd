package node

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/srchain/srcd/p2p"
	"github.com/srchain/srcd/p2p/nat"
)

const (
	DefaultHTTPHost = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort = 9090        // Default TCP port for the HTTP RPC server
)

// DefaultConfig contains reasonable default settings.
var DefaultConfig = Config{
	DataDir:          DefaultDataDir(),
	HTTPPort:         DefaultHTTPPort,
	// HTTPModules:      []string{"net", "web3"},
	// HTTPVirtualHosts: []string{"localhost"},
	P2P: p2p.Config{
		ListenAddr: ":10101",
		MaxPeers:   25,
		NAT:        nat.Any(),
	},
}

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "Srcd")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "Srcd")
		} else {
			return filepath.Join(home, ".srcd")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
