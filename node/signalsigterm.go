package node

import (
	"os"
	"syscall"
)

func init() {
	signals = []os.Signal{os.Interrupt, syscall.SIGTERM}
}
