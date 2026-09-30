//go:build windows

package baseline

import (
	"os"
	"path/filepath"
)

func openCitationFileNoFollow(root *os.Root, relative string) (*os.File, error) {
	return root.Open(filepath.FromSlash(relative))
}
