//go:build !windows

package baseline

import (
	"errors"
	"os"

	"golang.org/x/sys/unix"
)

func openCitationFileNoFollow(path string) (*os.File, error) {
	fd, err := unix.Open(
		path,
		unix.O_RDONLY|unix.O_NOFOLLOW|unix.O_NONBLOCK|unix.O_CLOEXEC,
		0,
	)
	if err != nil {
		return nil, err
	}
	file := os.NewFile(uintptr(fd), path)
	if file == nil {
		return nil, errors.Join(errors.New("create citation file from descriptor"), unix.Close(fd))
	}
	return file, nil
}
