package main

import (
	"os"

	"github.com/uttam-li/rutter/internal/cli"
)

func main() {
	os.Exit(cli.Execute())
}
