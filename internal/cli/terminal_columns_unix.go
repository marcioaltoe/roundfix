//go:build unix

package cli

import (
	"os"

	"golang.org/x/sys/unix"
)

func terminalColumns(file *os.File) (int, bool) {
	size, err := unix.IoctlGetWinsize(int(file.Fd()), unix.TIOCGWINSZ)
	if err != nil || size.Col == 0 {
		return 0, false
	}
	return int(size.Col), true
}
