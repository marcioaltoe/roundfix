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

// Harden disables Git's automatic maintenance for the repository at repoPath.
// Use it for repositories that were not created through InitRepo, such as
// clones.
//
// maintenance.auto stops commit, merge, fetch and receive-pack from launching
// "git maintenance run --auto" at all; gc.auto disables the older automatic
// "git gc --auto" threshold that the maintenance gc task also consults; and
// maintenance.autoDetach keeps any maintenance that still starts in the
// foreground, so it cannot outlive the git command that started it.
func Harden(t *testing.T, repoPath string) {
	t.Helper()
	settings := [][2]string{
		{"maintenance.auto", "false"},
		{"gc.auto", "0"},
		{"maintenance.autoDetach", "false"},
	}
	for _, setting := range settings {
		Run(t, repoPath, "config", setting[0], setting[1])
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
