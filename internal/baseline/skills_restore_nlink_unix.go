//go:build !windows

package baseline

import (
	"io/fs"
	"syscall"
)

func fileLinkCount(info fs.FileInfo) uint64 {
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0
	}
	return uint64(stat.Nlink)
}
