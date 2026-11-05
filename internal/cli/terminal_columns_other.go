//go:build !unix

package cli

import "os"

func terminalColumns(_ *os.File) (int, bool) {
	return 0, false
}
