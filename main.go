package main

import (
	"os"

	"github.com/wyx2685/v2node/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
