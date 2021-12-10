package main

import (
	"os"

	"github.com/xylonx/config2go/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
