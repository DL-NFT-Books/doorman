package main

import (
	"os"

	"github.com/dl-nft-books/doorman/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
