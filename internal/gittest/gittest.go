// Package gittest bootstraps the disposable Git repositories that tests
// create under t.TempDir().
//
// Git runs automatic maintenance after commands that write to a repository
// (commit, merge, fetch, receive-pack), and that maintenance detaches into a
// background process. The detached child can still be writing inside .git
// after the test returns, so t.TempDir() cleanup fails with
// "directory not empty". Repositories created through InitRepo carry the
// maintenance switches in their own config, so every git invocation against
// them — including the ones Roundfix itself runs — inherits the setting and
// no background process outlives the test that created the repository.
package gittest

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// autoDetachEnv is Git's override for detached automatic maintenance. Setting
// it to 0 keeps maintenance in the foreground for repositories the caller did
// not create through InitRepo, such as clones.
const autoDetachEnv = "GIT_TEST_MAINT_AUTO_DETACH"

// ConfigArgs returns the -c arguments every test git invocation carries: a
// deterministic committer identity and no commit signing.
func ConfigArgs() []string {
	return []string{
		"-c", "user.name=Roundfix Test",
		"-c", "user.email=test@example.com",
		"-c", "commit.gpgsign=false",
		// Reads must not write: a filesystem monitor can refresh the index on
		// otherwise read-only commands, and an index rename bumps the .git
		// directory's mtime — which discards any concurrently running test
		// result that recorded .git as an input.
		"-c", "core.fsmonitor=false",
	}
}

// IsolatedEnv returns the environment every test git invocation carries: no
// inherited GIT_CONFIG_* overrides, no user or system config, and no detached
// automatic maintenance.
func IsolatedEnv() []string {
	env := make([]string, 0, len(os.Environ())+3)
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if strings.HasPrefix(key, "GIT_CONFIG_") || key == autoDetachEnv {
			continue
		}
		env = append(env, entry)
	}
	return append(env,
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		autoDetachEnv+"=0",
		// Read-only commands must stay read-only: without this, git may take
		// the index lock to refresh it during reads such as ls-files, and the
		// resulting .git mtime bump invalidates every concurrently running
		// test that recorded .git as an input.
		"GIT_OPTIONAL_LOCKS=0",
	)
}

// InitRepo creates a disposable repository at the absolute repoPath and
// disables its automatic maintenance. Extra arguments are passed to git init,
// for example "--bare" or "-b", "main".
func InitRepo(t *testing.T, repoPath string, args ...string) {
	t.Helper()
	if !filepath.IsAbs(repoPath) {
		t.Fatalf("gittest: repository path must be absolute, got %q", repoPath)
	}
	initArgs := append([]string{"init"}, args...)
	Run(t, "", append(initArgs, repoPath)...)
	Harden(t, repoPath)
}

// hardenSnippet is the config fragment Harden appends. Git's config file is
// plain INI, so appending a section is equivalent to three `git config`
// invocations — and the whole suite creates hundreds of disposable
// repositories, so at four subprocess spawns per repository the bootstrap cost
// was mostly process creation, not git work. Measured across the suite, kernel
// time dwarfed user time three to one.
//
// maintenance.auto stops commit, merge, fetch and receive-pack from launching
// "git maintenance run --auto" at all; gc.auto disables the older automatic
// "git gc --auto" threshold that the maintenance gc task also consults; and
// maintenance.autoDetach keeps any maintenance that still starts in the
// foreground, so it cannot outlive the git command that started it.
const hardenSnippet = "[maintenance]\n\tauto = false\n\tautoDetach = false\n[gc]\n\tauto = 0\n"

// Harden disables Git's automatic maintenance for the repository at repoPath
// by appending directly to its config file. Use it for repositories that were
// not created through InitRepo, such as clones.
func Harden(t *testing.T, repoPath string) {
	t.Helper()
	configPath := filepath.Join(repoPath, ".git", "config")
	if info, err := os.Stat(filepath.Join(repoPath, ".git")); err != nil || !info.IsDir() {
		if _, headErr := os.Stat(filepath.Join(repoPath, "HEAD")); headErr == nil {
			// A bare repository keeps its config at the top level.
			configPath = filepath.Join(repoPath, "config")
		} else {
			// A worktree or separate git dir: let git resolve the config
			// location rather than guessing.
			for _, setting := range [][2]string{
				{"maintenance.auto", "false"},
				{"gc.auto", "0"},
				{"maintenance.autoDetach", "false"},
			} {
				Run(t, repoPath, "config", setting[0], setting[1])
			}
			return
		}
	}
	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("gittest: open %s to harden: %v", configPath, err)
	}
	if _, err := file.WriteString(hardenSnippet); err != nil {
		_ = file.Close()
		t.Fatalf("gittest: harden %s: %v", configPath, err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("gittest: close %s: %v", configPath, err)
	}
}

// AppendConfig appends a config fragment to the repository's own config file.
// Git's config file is plain INI, so one file append replaces one subprocess
// spawn per key — and the suite was paying roughly nine hundred `git config`
// spawns per run across its fixtures.
func AppendConfig(t *testing.T, repoPath string, snippet string) {
	t.Helper()
	if info, err := os.Stat(filepath.Join(repoPath, ".git")); err != nil || !info.IsDir() {
		t.Fatalf("gittest: %s is not a repository with a .git directory", repoPath)
	}
	configPath := filepath.Join(repoPath, ".git", "config")
	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("gittest: open %s to append config: %v", configPath, err)
	}
	if _, err := file.WriteString(snippet); err != nil {
		_ = file.Close()
		t.Fatalf("gittest: append config to %s: %v", configPath, err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("gittest: close %s: %v", configPath, err)
	}
}

// identitySnippet persists the deterministic test identity in a repository's
// own config, for repositories whose git commands are spawned by production
// code rather than through Run — production inherits the process environment,
// so per-command -c arguments cannot reach it.
const identitySnippet = "[user]\n\tname = Roundfix Test\n\temail = test@example.com\n[commit]\n\tgpgsign = false\n"

// PersistIdentity appends the deterministic test identity to the repository's
// config file, replacing three `git config` subprocess spawns with one file
// append.
func PersistIdentity(t *testing.T, repoPath string) {
	t.Helper()
	configPath := filepath.Join(repoPath, ".git", "config")
	if info, err := os.Stat(filepath.Join(repoPath, ".git")); err != nil || !info.IsDir() {
		t.Fatalf("gittest: %s is not a repository with a .git directory", repoPath)
	}
	file, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("gittest: open %s to persist identity: %v", configPath, err)
	}
	if _, err := file.WriteString(identitySnippet); err != nil {
		_ = file.Close()
		t.Fatalf("gittest: persist identity in %s: %v", configPath, err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("gittest: close %s: %v", configPath, err)
	}
}

// Run runs git inside dir with the shared test config and isolated
// environment and returns its combined output. An empty dir leaves the
// working directory of the test process untouched.
func Run(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmdArgs := make([]string, 0, len(args)+len(ConfigArgs())+2)
	if strings.TrimSpace(dir) != "" {
		cmdArgs = append(cmdArgs, "-C", dir)
	}
	cmdArgs = append(cmdArgs, ConfigArgs()...)
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("git", cmdArgs...)
	cmd.Env = IsolatedEnv()
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}
