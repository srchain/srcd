package utils

import (
	"fmt"
	"io"
	"os"
	"runtime"
)

// Fatalf formats a message to standard error and exits the program.
// The message is also printed to standard output if standard error
// is redirected to a different file.
func Fatalf(format string, args ...interface{}) {
	w := io.MultiWriter(os.Stdout, os.Stderr)
	if runtime.GOOS == "windows" {
		// The SameFile check below doesn't work on Windows.
		// stdout is unlikely to get redirected though, so just print there.
		w = os.Stdout
	} else {
		outf, _ := os.Stdout.Stat()
		errf, _ := os.Stderr.Stat()
		if outf != nil && errf != nil && os.SameFile(outf, errf) {
			w = os.Stderr
		}
	}
	fmt.Fprintf(w, "Fatal: "+format+"\n", args...)
	os.Exit(1)
}

// func StartPeer(node *peer.Peer) {
	// if err := peer.Start(); err != nil {
		// Fatalf("Error starting protocol peer: %v", err)
	// }
	// go func() {
		// sigc := make(chan os.Signal, 1)
		// signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		// defer signal.Stop(sigc)
		// <-sigc
		// log.Info("Got interrupt, shutting down...")
		// go peer.Stop()
		// for i := 10; i > 0; i-- {
			// <-sigc
			// if i > 1 {
				// log.Warn("Already shutting down, interrupt more to panic.", "times", i-1)
			// }
		// }
	// }()
// }
