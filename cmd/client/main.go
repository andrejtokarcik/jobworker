package main

import (
	"context"
	"os"
)

func main() {
	cxt := context.Background()
	if err := CLI(cxt).Execute(); err != nil {
		// Error printing handled by cobra in these cases
		os.Exit(1)
	}
}
