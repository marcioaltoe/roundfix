package gittest

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

// PinnedHistoryRevision holds the pre-cut archived Spec corpus.
const PinnedHistoryRevision = "40a7893d872c8a6705f6d7745e6efe430ae9deeb"

var pinnedHistories = struct {
	sync.Mutex
	roots   map[string]string
	readers map[string]int
}{roots: make(map[string]string), readers: make(map[string]int)}

// PinnedHistory materializes each requested path set once per test binary.
// Callers share a read-only tree; the source repository is never written.
func PinnedHistory(t testing.TB, repoRoot string, pathspecs ...string) string {
	t.Helper()
	if len(pathspecs) == 0 {
		t.Fatal("pinned history requires at least one path")
	}
	paths := slices.Clone(pathspecs)
	slices.Sort(paths)
	paths = slices.Compact(paths)
	repository, err := filepath.Abs(repoRoot)
	if err != nil {
		t.Fatal(err)
	}
	repository, err = filepath.EvalSymlinks(repository)
	if err != nil {
		t.Fatal(err)
	}
	key := repository + "\x00" + strings.Join(paths, "\x00")
	pinnedHistories.Lock()
	defer pinnedHistories.Unlock()
	if root, ok := pinnedHistories.roots[key]; ok {
		retainPinnedHistory(t, root)
		return root
	}
	run := func(args ...string) ([]byte, error) {
		cmd := exec.CommandContext(t.Context(), "git", append(ConfigArgs(), args...)...)
		cmd.Dir = repoRoot
		cmd.Env = IsolatedEnv()
		return cmd.Output()
	}
	if _, err := run("cat-file", "-e", PinnedHistoryRevision+"^{commit}"); err != nil {
		t.Skipf("pinned history commit %s is unavailable: %v", PinnedHistoryRevision, err)
	}
	args := append([]string{"archive", "--format=tar", PinnedHistoryRevision, "--"}, paths...)
	content, err := run(args...)
	if err != nil {
		t.Fatalf("archive pinned history %s: %v", PinnedHistoryRevision, err)
	}
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	// Like the shared fixture binaries, this tree belongs to the test binary,
	// not to the first test that requests it.
	root, err := os.MkdirTemp(filepath.Dir(executable), "roundfix-pinned-history-")
	if err != nil {
		t.Fatal(err)
	}
	reader := tar.NewReader(bytes.NewReader(content))
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read pinned archive: %v", err)
		}
		name := filepath.FromSlash(header.Name)
		if filepath.IsAbs(name) || name == ".." || strings.HasPrefix(filepath.Clean(name), ".."+string(filepath.Separator)) {
			t.Fatalf("pinned archive path escapes destination: %q", header.Name)
		}
		target := filepath.Join(root, name)
		switch header.Typeflag {
		case tar.TypeXGlobalHeader:
			// git archive records the commit in a global PAX header.
			continue
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				t.Fatal(err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				t.Fatal(err)
			}
			data, err := io.ReadAll(reader)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(target, data, os.FileMode(header.Mode).Perm()); err != nil {
				t.Fatal(err)
			}
		case tar.TypeSymlink:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.Symlink(header.Linkname, target); err != nil {
				t.Fatal(err)
			}
		default:
			t.Fatalf("unsupported pinned archive entry %q (type %d)", header.Name, header.Typeflag)
		}
	}
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		return os.Chmod(path, info.Mode().Perm()&^0o222)
	}); err != nil {
		t.Fatal(err)
	}
	pinnedHistories.roots[key] = root
	retainPinnedHistory(t, root)
	return root
}

// The tree stays read-only for the lifetime of every caller. When the last
// reader ends, restore directory permissions so the Go tool can remove its
// build directory. A later test reuses the extraction and protects it again.
// The caller holds pinnedHistories.Mutex.
func retainPinnedHistory(t testing.TB, root string) {
	t.Helper()
	if pinnedHistories.readers[root] == 0 {
		if err := pinnedHistoryDirectoryMode(root, 0o555); err != nil {
			t.Fatal(err)
		}
	}
	pinnedHistories.readers[root]++
	t.Cleanup(func() {
		pinnedHistories.Lock()
		defer pinnedHistories.Unlock()
		pinnedHistories.readers[root]--
		if pinnedHistories.readers[root] == 0 {
			if err := pinnedHistoryDirectoryMode(root, 0o755); err != nil {
				t.Error(err)
			}
		}
	})
}

func pinnedHistoryDirectoryMode(root string, mode os.FileMode) error {
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return os.Chmod(path, mode)
		}
		return nil
	})
}
