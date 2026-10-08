// Suite: compiled script fixtures.
// Invariant: fixtures execute the compiled test binary, never a written executable.
// Boundary IN: symlinks, private shell sidecars, arguments, environment and real subprocesses.
// Boundary OUT: installed adapters and providers; Darwin's signed macro Codex fixture is compiled separately.
package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"roundfix/internal/agent"
)

const scriptFixtureSuffix = ".fixture.sh"

var scriptFixtureBinaryPath string

func writeScriptFixture(t testing.TB, path string, body string) {
	t.Helper()
	if !filepath.IsAbs(scriptFixtureBinaryPath) {
		t.Fatal("compiled test binary path was not initialized by TestMain")
	}
	if err := os.WriteFile(path+scriptFixtureSuffix, []byte(body), 0o600); err != nil {
		t.Fatalf("write script fixture sidecar: %v", err)
	}
	if err := os.Symlink(scriptFixtureBinaryPath, path); err != nil {
		t.Fatalf("link script fixture to compiled test binary: %v", err)
	}
}

func runScriptFixture() {
	path := os.Args[0]
	if !strings.ContainsRune(path, os.PathSeparator) {
		resolved, err := exec.LookPath(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "resolve script fixture: %v\n", err)
			os.Exit(2)
		}
		path = resolved
	}
	sidecar := path + scriptFixtureSuffix
	if _, err := os.Stat(sidecar); errors.Is(err, os.ErrNotExist) {
		return
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "inspect script fixture sidecar: %v\n", err)
		os.Exit(2)
	}
	args := append([]string{"/bin/sh", sidecar}, os.Args[1:]...)
	if err := syscall.Exec("/bin/sh", args, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "execute script fixture: %v\n", err)
		os.Exit(2)
	}
}

func TestScriptFixtureIsTheCompiledTestBinary(t *testing.T) {
	// Sequential: sets the process-wide carry-forward hook marker environment variable.
	t.Run("script arguments and environment", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "script")
		writeScriptFixture(t, path, "printf '%s|%s|%s\\n' \"$1\" \"$2\" \"$FIXTURE_VALUE\"\nexit 7\n")
		assertScriptFixture(t, path, []string{"one space", "two"}, "one space|two|inherited\n", 7,
			[]string{"FIXTURE_VALUE=inherited", cliTestHelperEnv + "=1", detachTestChildModeEnv + "=invalid"})
	})
	t.Run("adapter version", func(t *testing.T) {
		path := writeAdapterVersionFixture(t, agent.CodexAdapterPackage, agent.PinnedCodexAdapterVersion)
		assertScriptFixture(t, path, []string{"--version"}, agent.CodexAdapterPackage+" "+agent.PinnedCodexAdapterVersion+"\n", 0, nil)
	})
	t.Run("fake ACPX and adapters", func(t *testing.T) {
		setCommandHomeDirForTest(t, t.TempDir())
		path := fakeACPXCommand(t)
		for _, fixture := range []struct{ name, output string }{
			{"acpx", agent.MinimumACPXVersion},
			{"codex-acp", agent.CodexAdapterPackage + " " + agent.PinnedCodexAdapterVersion},
			{"claude-agent-acp", agent.ClaudeAdapterPackage + " " + agent.PinnedClaudeAdapterVersion},
		} {
			t.Run(fixture.name, func(t *testing.T) {
				assertScriptFixture(t, filepath.Join(filepath.Dir(path), fixture.name), []string{"--version"}, fixture.output+"\n", 0, nil)
			})
		}
	})
	t.Run("macro ACPX and adapters", func(t *testing.T) {
		fake := newMacroFakeACPX(t)
		for _, fixture := range []struct{ name, output string }{
			{"acpx", agent.MinimumACPXVersion + "\n"}, {"codex-acp", ""}, {"claude-agent-acp", ""},
			{"opencode", ""}, {"npx", agent.CodexAdapterPackage + " " + agent.PinnedCodexAdapterVersion + "\n"},
		} {
			t.Run(fixture.name, func(t *testing.T) {
				assertScriptFixture(t, filepath.Join(fake.binDir, fixture.name), []string{"--version"}, fixture.output, 0, nil)
			})
		}
		if runtime.GOOS != "darwin" {
			assertScriptFixture(t, fake.codexPath, nil, "", 0, nil)
		}
	})
	t.Run("carry-forward hooks", func(t *testing.T) {
		hooks := writeCarryForwardHookFixtures(t, true)
		for _, name := range []string{"pre-commit", "prepare-commit-msg", "commit-msg", "post-commit"} {
			t.Run(name, func(t *testing.T) {
				wantExit := 0
				if name == "pre-commit" || name == "commit-msg" {
					wantExit = 1
				}
				assertScriptFixture(t, filepath.Join(hooks.directory, name), nil, "", wantExit, nil)
				assertCarryForwardHookMarkers(t, hooks.marker, []string{name, name, name})
				if err := os.Remove(hooks.marker); err != nil {
					t.Fatalf("clear hook markers: %v", err)
				}
			})
		}
	})
	t.Run("settle hook", func(t *testing.T) {
		root := t.TempDir()
		gitImplement(t, root, "init")
		mustWrite(t, filepath.Join(root, "oversized.txt"), strings.Repeat("line\n", 501))
		gitImplement(t, root, "add", "oversized.txt")
		path := filepath.Join(root, "pre-commit")
		writeSettleHookFixture(t, path)
		assertScriptFixture(t, path, nil, "oversized.txt: 501 lines exceeds the 500-line limit\n", 1, nil)
	})
}

func assertScriptFixture(t *testing.T, path string, args []string, want string, wantExit int, extraEnv []string) {
	t.Helper()
	target, err := os.Readlink(path)
	if err != nil {
		t.Fatalf("fixture must be a symlink: %v", err)
	}
	if target != scriptFixtureBinaryPath || !filepath.IsAbs(target) {
		t.Fatalf("fixture target = %q, want absolute test binary %q", target, scriptFixtureBinaryPath)
	}
	info, err := os.Stat(path + scriptFixtureSuffix)
	if err != nil {
		t.Fatalf("stat sidecar: %v", err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("sidecar mode = %o, want 600", info.Mode().Perm())
	}
	for _, invocation := range []string{path, "./" + filepath.Base(path), filepath.Base(path)} {
		t.Run("invoke "+invocation, func(t *testing.T) {
			// Use the shell's PATH lookup so the fixture receives a bare argv[0].
			cmd := &exec.Cmd{Path: path, Args: append([]string{invocation}, args...), Dir: filepath.Dir(path)}
			if invocation == filepath.Base(path) {
				cmd = exec.Command("/bin/sh", append([]string{"-c", `exec "$@"`, "sh", invocation}, args...)...)
				cmd.Dir = filepath.Dir(path)
			}
			cmd.Env = withEnvValue(isolatedGitEnvForTest(), "PATH", filepath.Dir(path)+string(os.PathListSeparator)+os.Getenv("PATH"))
			for _, entry := range extraEnv {
				key, value, _ := strings.Cut(entry, "=")
				cmd.Env = withEnvValue(cmd.Env, key, value)
			}
			output, err := cmd.CombinedOutput()
			if got := exitCodeFromWait(err); got != wantExit {
				t.Fatalf("fixture exit = %d, want %d: %v\n%s", got, wantExit, err, output)
			}
			if string(output) != want {
				t.Fatalf("fixture output = %q, want %q", output, want)
			}
		})
	}
}
