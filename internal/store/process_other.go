//go:build !unix && !windows

package store

// ProcessAlive reports alive on platforms where orphan recovery is unsupported.
func ProcessAlive(pid int) bool {
	return true
}
