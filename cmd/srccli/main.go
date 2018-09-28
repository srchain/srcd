package main

import (
	"runtime"

	cmd "srcd/cmd/srccli/commands"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
