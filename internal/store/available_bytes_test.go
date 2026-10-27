package store

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAvailableBytesReadsTheFilesystemOfAPath(t *testing.T) {
	t.Parallel()
	directory := t.TempDir()
	available, err := AvailableBytes(directory)
	if err != nil || available <= 0 {
		t.Fatalf("AvailableBytes(%q) = %d, %v; want positive capacity", directory, available, err)
	}
	missing := filepath.Join(directory, "missing")
	if _, err := AvailableBytes(missing); !errors.Is(err, os.ErrNotExist) || !strings.Contains(err.Error(), missing) {
		t.Fatalf("AvailableBytes(%q) error = %v; want missing-path error naming path", missing, err)
	}
}
