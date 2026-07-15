//go:build !unix

package store

// ProcessAlive reports alive when the platform cannot prove process liveness.
func ProcessAlive(pid int) bool {
	return true
}
