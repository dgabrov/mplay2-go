package main

import (
	"fmt"
	"github.com/amanagement24/mplay2-go/internal/app"
	"os"
)

func main() {
	err := app.Start()

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
