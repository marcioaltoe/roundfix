// Command gocache-bound applies the Build Cache Bound to the Makefile's Go build cache.
package main

import (
	"os"
	"time"

	"roundfix/internal/sharedcache"
)

func main() {
	os.Exit(sharedcache.BoundMain(os.Args[1:], os.Stderr, time.Now))
}
