package main

import (
	"context"
	"os"

	"roundfix/internal/verifyselect"
)

func main() {
	os.Exit(verifyselect.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}
