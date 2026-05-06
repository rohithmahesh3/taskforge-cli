package main

import (
	"os"

	"github.com/rohithmahesh3/taskforge-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
