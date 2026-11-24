// Boundary: production gate, disposable Git history and a recording Docker stand-in.
package cli

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func recordingRetiredDocker(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	record := filepath.Join(dir, "started")
	mustWrite(t, filepath.Join(dir, "docker"), "#!/bin/sh\nprintf started > '"+record+"'\nexit 99\n")
	if err := os.Chmod(filepath.Join(dir, "docker"), 0700); err != nil {
		t.Fatal(err)
	}
	env := commandEnvironmentForTest(t).dependencies.environ()
	setCommandEnvForTest(t, "PATH", dir+string(os.PathListSeparator)+environmentValue(env, "PATH"))
	return record
}
func environmentValue(env []string, key string) string {
	for _, v := range env {
		if s, ok := strings.CutPrefix(v, key+"="); ok {
			return s
		}
	}
	return ""
}

type retiredLaneGateGit struct{ t *testing.T }

func (runner retiredLaneGateGit) RunGit(context.Context, string, ...string) (string, error) {
	runner.t.Fatal("the repository gate must not read the default branch configuration")
	return "", nil
}

func TestTheRepositoryGateEndsAtTheRepositoryVerification(t *testing.T) {
	t.Parallel()
	f := newDeliveryItemFixture(t)
	record := recordingRetiredDocker(t)
	mustWrite(t, filepath.Join(f.repo, ".roundfixrc.yml"), "delivery:\n  linux_lane: true\n  linux_lane_image: golang:1.26\n")
	gittest.Run(t, f.repo, "add", ".")
	gittest.Run(t, f.repo, "commit", "-m", "legacy configuration")
	gittest.Run(t, f.repo, "update-ref", "refs/remotes/origin/main", "HEAD")
	f.workflow.loaded.Config.Defaults.Verification = "exit 0"
	f.workflow.git = retiredLaneGateGit{t: t}
	result, err := f.workflow.Gate(commandContextForTest(t, t.Context()), f.repo, "widget", f.head)
	if err != nil || !result.Passed {
		t.Fatalf("gate=%+v error=%v", result, err)
	}
	if _, err := os.Stat(record); !os.IsNotExist(err) {
		t.Fatalf("docker started: %v", err)
	}
	log := filepath.Join(f.workflow.loaded.Config.Defaults.ArtifactDir, "delivery", "widget", "linux-lane.log")
	if _, err := os.Stat(log); !os.IsNotExist(err) {
		t.Fatalf("lane log: %v", err)
	}
	if strings.Contains(strings.ToLower(f.log.String()), "lane") {
		t.Fatalf("console=%q", f.log.String())
	}
}
