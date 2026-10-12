package gittest

import (
	"os/exec"
	"path/filepath"
	"testing"
)

// PullRequestClone makes a main-only clone and a local bare forge that holds
// the source history exclusively under refs/pull/7/head. It proves the clone
// lacks the revision before returning it to a reader test.
func PullRequestClone(t *testing.T, source, head, revision string) (string, string) {
	t.Helper()
	clone := filepath.Join(t.TempDir(), "clone")
	Run(t, "", "clone", "--no-local", "--single-branch", "--branch", "main", source, clone)
	Harden(t, clone)
	PersistIdentity(t, clone)
	forge := filepath.Join(t.TempDir(), "forge.git")
	InitRepo(t, forge, "--bare")
	Run(t, forge, "fetch", "--no-tags", "--no-write-fetch-head", source, head+":refs/pull/7/head")
	Run(t, clone, "remote", "set-url", "origin", forge)
	cmd := exec.Command("git", "-C", clone, "cat-file", "-e", revision+"^{commit}")
	cmd.Env = IsolatedEnv()
	if err := cmd.Run(); err == nil {
		t.Fatal("clone already holds the source revision")
	}
	return clone, forge
}
