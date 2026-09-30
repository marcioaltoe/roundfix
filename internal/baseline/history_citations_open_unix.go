//go:build !windows

package baseline

import (
	"os"
	"path/filepath"
	"syscall"
)

func openCitationFileNoFollow(root *os.Root, relative string) (*os.File, error) {
	return root.OpenFile(filepath.FromSlash(relative), os.O_RDONLY|syscall.O_NONBLOCK, 0)
}
