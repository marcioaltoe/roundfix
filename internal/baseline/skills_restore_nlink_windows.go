//go:build windows

package baseline

import "io/fs"

func fileLinkCount(fs.FileInfo) uint64 {
	// os.FileInfo does not expose a Windows hard-link count. The anchored
	// transaction still rejects links, directories, and special files.
	return 1
}
