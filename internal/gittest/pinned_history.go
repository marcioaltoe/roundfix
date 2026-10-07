package gittest

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// PinnedHistoryRevision holds the pre-cut archived Spec corpus.
const PinnedHistoryRevision = "40a7893d872c8a6705f6d7745e6efe430ae9deeb"

// PinnedHistory materializes only the requested pinned paths in a disposable
// directory. It never writes the source repository.
func PinnedHistory(t testing.TB, repoRoot string, pathspecs ...string) string {
	t.Helper()
	if len(pathspecs) == 0 {
		t.Fatal("pinned history requires at least one path")
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
	args := append([]string{"archive", "--format=tar", PinnedHistoryRevision, "--"}, pathspecs...)
	content, err := run(args...)
	if err != nil {
		t.Fatalf("archive pinned history %s: %v", PinnedHistoryRevision, err)
	}
	root := t.TempDir()
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
	return root
}
