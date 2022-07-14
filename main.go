package main

import (
	"os"
)

func main() {
	if err := Run(os.Args); err != nil {
		os.Exit(1)
	}
}
