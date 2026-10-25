// Suite: repository boundary leak characterization
// Invariant: Git fingerprinting preserves every existing leak detection class.
// Boundary IN: guarded subprocess exit, diagnostics, and Go's test input log.
// Boundary OUT: installing the guard in other packages.
package suiteguard_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func fixtureGit(t *testing.T, root string, args ...string) {
	t.Helper()
	command := exec.CommandContext(t.Context(), "git", append([]string{"-C", root,
		"-c", "core.hooksPath=" + os.DevNull, "-c", "commit.gpgsign=false",
		"-c", "user.name=Suite Guard Fixture", "-c", "user.email=suiteguard@example.invalid"}, args...)...)
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(strings.ToUpper(entry), "GIT_") {
			command.Env = append(command.Env, entry)
		}
	}
	command.Env = append(command.Env, "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+os.DevNull)
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("fixture git %s: %v: %s", args[0], err, output)
	}
}

func newGitRepository(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	fixtureGit(t, root, "init")
	return root
}

func leakRepository(t *testing.T) string {
	t.Helper()
	root := newGitRepository(t)
	for name, contents := range map[string]string{
		".gitignore": "/ignored/\n*.log\n!keep.log\n", "sub/.gitignore": "local.tmp\n",
		"tracked.txt": "tracked before\n", "dirty.txt": "committed\n", "sub/nested.txt": "nested before\n",
	} {
		mustWrite(t, filepath.Join(root, name), contents)
	}
	if err := os.Symlink("tracked.txt", filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}
	fixtureGit(t, root, "add", ".")
	fixtureGit(t, root, "commit", "-m", "fixture")
	mustWrite(t, filepath.Join(root, "dirty.txt"), "dirty before\n")
	mustWrite(t, filepath.Join(root, "scratch.txt"), "scratch before\n")
	if err := os.Mkdir(filepath.Join(root, "emptydir"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

type leakCase struct {
	name, path, operation string
	want                  []string
}

var leakCases = []leakCase{
	{"tracked edit", "tracked.txt", "write", []string{"modified: tracked.txt"}},
	{"tracked deletion", "tracked.txt", "remove", []string{"removed: tracked.txt"}},
	{"tracked empty file", "sub/nested.txt", "empty", []string{"modified: sub/nested.txt"}},
	{"dirty tracked edit", "dirty.txt", "write", []string{"modified: dirty.txt"}},
	{"untracked creation", "created.txt", "write", []string{"created: created.txt"}},
	{"new directory file", "newdir/deep/created.txt", "write", []string{"created: newdir", "created: newdir/deep/created.txt"}},
	{"beside tracked file", "sub/created.txt", "write", []string{"created: sub/created.txt"}},
	{"untracked edit", "scratch.txt", "write", []string{"modified: scratch.txt"}},
	{"untracked deletion", "scratch.txt", "remove", []string{"removed: scratch.txt"}},
	{"empty directory creation", "createddir", "mkdir", []string{"created: createddir"}},
	{"empty directory removal", "emptydir", "remove", []string{"removed: emptydir"}},
	{"tracked executable", "tracked.txt", "executable", []string{"modified: tracked.txt"}},
	{"untracked executable", "scratch.txt", "executable", []string{"modified: scratch.txt"}},
	{"symlink retarget", "link", "symlink", []string{"modified: link"}},
	{"file becomes symlink", "tracked.txt", "symlink", []string{"modified: tracked.txt"}},
	{"file permission narrowed", "tracked.txt", "private", []string{"modified: tracked.txt"}},
	{"directory permission narrowed", "sub", "private directory", []string{"modified: sub"}},
	{"directory with ignored files", "onlylogs/x.log", "write", []string{"created: onlylogs"}},
	{"negated ignore", "keep.log", "write", []string{"created: keep.log"}},
	{"directory ignore", "ignored/build.log", "write", nil},
	{"root ignore", "sub/trace.log", "write", nil},
	{"nested ignore", "sub/local.tmp", "write", nil},
}

func mutateLeak(t *testing.T, root, name string) {
	t.Helper()
	for _, tc := range leakCases {
		if tc.name != name {
			continue
		}
		target := filepath.Join(root, tc.path)
		var err error
		switch tc.operation {
		case "write":
			mustWrite(t, target, "changed\n")
		case "empty":
			mustWrite(t, target, "")
			var info os.FileInfo
			info, err = os.Stat(target)
			if err == nil && info.Size() != 0 {
				t.Fatalf("emptied file size = %d, want zero", info.Size())
			}
		case "remove":
			err = os.Remove(target)
		case "mkdir":
			err = os.Mkdir(target, 0o755)
		case "executable":
			err = os.Chmod(target, 0o755)
		case "private":
			err = os.Chmod(target, 0o600)
		case "private directory":
			err = os.Chmod(target, 0o700)
		case "symlink":
			if err = os.Remove(target); err == nil {
				err = os.Symlink("dirty.txt", target)
			}
		default:
			t.Fatalf("unknown leak operation %q", tc.operation)
		}
		if err != nil {
			t.Fatal(err)
		}
		return
	}
	t.Fatalf("unknown leak case %q", name)
}

func TestGuardReportsEveryLeakClass(t *testing.T) {
	t.Parallel()
	for _, tc := range leakCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			output, code := runGuardedFixture(t, leakRepository(t), "leak", "", tc.name)
			if len(tc.want) == 0 {
				if code != 0 {
					t.Fatalf("ignored write exit = %d, want 0:\n%s", code, output)
				}
				if strings.Contains(output, "suiteguard: repository boundary violated:") {
					t.Fatalf("ignored write reported:\n%s", output)
				}
				return
			}
			if code == 0 {
				t.Fatalf("leak exit = 0, want failure:\n%s", output)
			}
			for _, want := range tc.want {
				if !strings.Contains(output, "suiteguard: "+want+"\n") {
					t.Errorf("missing diagnostic %q:\n%s", want, output)
				}
			}
		})
	}
}

func TestGuardLeavesNoRepositoryPathInTheTestInputLog(t *testing.T) {
	t.Parallel()
	repository := leakRepository(t)
	logPath := filepath.Join(t.TempDir(), "test-inputs.log")
	output, code := runGuardedFixture(t, repository, "isolated", "", "", "-test.testlogfile="+logPath)
	if code != 0 {
		t.Fatalf("guarded log fixture exit = %d:\n%s", code, output)
	}
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(content) == 0 {
		t.Fatal("empty test input log")
	}
	resolved, err := filepath.EvalSymlinks(repository)
	if err != nil {
		t.Fatal(err)
	}
	for _, line := range strings.Split(string(content), "\n") {
		operation, name, _ := strings.Cut(line, " ")
		if operation != "open" && operation != "stat" {
			continue
		}
		for _, root := range []string{repository, resolved} {
			if name == root || strings.HasPrefix(name, root+string(filepath.Separator)) {
				t.Errorf("repository input logged: %s", line)
			}
		}
	}
}

func TestGuardRefusesARootOutsideAGitWorkTree(t *testing.T) {
	t.Parallel()
	output, code := runGuardedFixture(t, t.TempDir(), "isolated", "")
	if code == 0 || !strings.Contains(output, "suiteguard: fingerprint repository before tests:") {
		t.Fatalf("non-Git root exit = %d, want refusal before tests:\n%s", code, output)
	}
}
