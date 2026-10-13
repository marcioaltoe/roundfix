//go:build !darwin && !linux

package suiteguard

import (
	"fmt"
	"runtime"
)

func platformProcessVanished(error) bool { return false }
func helperResidue(string, int) ([]residueProcess, error) {
	return nil, fmt.Errorf("unsupported on %s", runtime.GOOS)
}
func killResidue(residueProcess, string) error { return fmt.Errorf("unsupported on %s", runtime.GOOS) }
