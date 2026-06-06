package main

import (
	"fmt"
	"os"

	"github.com/comfucios/vw/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
