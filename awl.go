package main

import (
	"fmt"
	"os"
)

func main() {
	app := prepareCLI()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
