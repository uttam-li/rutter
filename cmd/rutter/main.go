package main

import (
	"fmt"
	"os"

	"github.com/uttam-li/rutter/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "rutter:", err)
		os.Exit(1)
	}
}
