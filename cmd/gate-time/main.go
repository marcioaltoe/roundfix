package main

import (
	"os"
	"roundfix/internal/gatetime"
)

func main() { os.Exit(gatetime.Run(os.Args[1:], os.Stdout, os.Stderr)) }
