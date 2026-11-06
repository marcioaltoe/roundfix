package cli

import (
	"context"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
	"syscall"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
)

const diskFullLogTailBytes int64 = 64 * 1024

var diskFullPath = regexp.MustCompile(`(?i)(?:^|[\s"'(])((?:/|[a-z]:[\\/])[^\s]*): no space left on device`)

// diskFullOperationPath reads the path after the file operation Go and the
// shell name in an ENOSPC error, so a path that holds spaces, such as
// "/Users/me/Google Drive/repo", is read whole.
var diskFullOperationPath = regexp.MustCompile(`(?i)\b(?:write|open|mkdir|create|rename|link|symlink|truncate|sync|close|copy)\s+((?:/|[a-z]:[\\/])[^\n]*?): no space left on device`)

// diskFullFailure reads only retained local diagnostics; it never changes the
// original error text or measures capacity without the command dependency.
func diskFullFailure(ctx context.Context, loaded roundconfig.Loaded, err error, logPath string) *delivery.DiskFullError {
	tail := diskFullLogTail(logPath)
	if !errors.Is(err, syscall.ENOSPC) && !strings.Contains(strings.ToLower(tail), "no space left on device") {
		return nil
	}
	path := lastDiskFullPath(tail)
	if path == "" {
		path = lastDiskFullPath(err.Error())
	}
	full := &delivery.DiskFullError{Path: path, Err: err}
	if path == "" {
		found := false
		for _, location := range inspectFreeSpace(ctx, loaded) {
			if location.err == nil && (!found || location.available < full.AvailableBytes) {
				full.Path, full.AvailableBytes = location.path, location.available
				found = true
			}
		}
		return full
	}
	ancestor, statErr := doctorDiskSpaceAncestor(path)
	if statErr == nil {
		available, capacityErr := commandDependenciesForContext(ctx).availableBytes(ancestor)
		if capacityErr == nil {
			full.AvailableBytes = available
		}
	}
	return full
}

func lastDiskFullPath(text string) string {
	path, end := "", -1
	// The last match in the text wins; at the same end, the operation form,
	// which keeps spaces, wins over the token form.
	for _, pattern := range []*regexp.Regexp{diskFullPath, diskFullOperationPath} {
		for _, match := range pattern.FindAllStringSubmatchIndex(text, -1) {
			if match[1] >= end {
				path, end = text[match[2]:match[3]], match[1]
			}
		}
	}
	return path
}

func diskFullLogTail(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil || !info.Mode().IsRegular() {
		return ""
	}
	// Reading a bounded section also avoids loading the rest of a large build log.
	data, err := io.ReadAll(io.NewSectionReader(file, max(0, info.Size()-diskFullLogTailBytes), diskFullLogTailBytes))
	if err != nil {
		return ""
	}
	return string(data)
}
