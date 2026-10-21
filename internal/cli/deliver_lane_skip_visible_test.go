// Boundary: real gate and status, disposable repositories and compiled lane tools.
package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/store"
)

func TestGateReturnsTheLaneEnvironmentSkip(t *testing.T) {
	t.Parallel()
	docker, goBinary := linuxLaneTools(t)
	for _, tc := range []struct{ name, reason string }{
		{"no go", "go is not on PATH"},
		{"no docker", "docker is not on PATH"},
		{"daemon", "docker info failed: daemon unavailable"},
		{"image", "image golang:1.26.1 is not present locally"},
		{"start refused", "docker could not start the container: fixture output"},
		{"no version", "go.mod declares no go version"},
		{"unreadable default", "the fetched default refs/remotes/origin/main cannot be resolved"},
		{"no module", "no go.mod at the repository root"},
		{"no package", "no changed Go package"},
		{"disabled", ""}, {"passed", ""}, {"failed repository gate", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := newLinuxLaneFixture(t, docker, goBinary, map[string]string{"go.mod": "module fixture\n\ngo 1.26.1\n", "pkg/code.go": "package pkg\n"})
			f.commit(t, "pkg/changed.go", "package pkg\n")
			switch tc.name {
			case "no go":
				f.environment = withEnvValue(f.environment, "PATH", t.TempDir())
			case "no docker":
				f.environment = withEnvValue(f.environment, "PATH", filepath.Dir(goBinary))
			case "daemon":
				f.environment = withEnvValue(f.environment, "LANE_INFO_FAIL", "1")
			case "image":
				f.environment = withEnvValue(f.environment, "LANE_IMAGE_FAIL", "1")
			case "start refused":
				f.environment = withEnvValue(f.environment, "LANE_EXIT", "125")
				f.environment = withEnvValue(f.environment, "LANE_PACKAGES", "")
			case "no version":
				f.commit(t, "go.mod", "module fixture\n")
			case "unreadable default":
				gittest.Run(t, f.repo, "update-ref", "-d", "refs/remotes/origin/main")
			case "no module":
				if err := os.Remove(filepath.Join(f.repo, "go.mod")); err != nil {
					t.Fatal(err)
				}
				f.commit(t, "docs/readme", "docs")
			case "no package":
				gittest.Run(t, f.repo, "update-ref", "refs/remotes/origin/main", f.head)
				f.commit(t, "docs/readme", "docs")
			case "disabled":
				f.commit(t, ".roundfixrc.yml", "delivery:\n  linux_lane: false\n")
				gittest.Run(t, f.repo, "update-ref", "refs/remotes/origin/main", f.head)
			case "failed repository gate":
				f.workflow.loaded.Config.Defaults.Verification = "exit 1"
			}
			result, err := f.workflow.Gate(f.context(t.Context()), f.repo, "widget", f.head)
			wantSkip := tc.reason
			if tc.name == "no module" || tc.name == "no package" {
				wantSkip = ""
			}
			if err != nil || result.Passed != (tc.name != "failed repository gate") || result.LaneSkip != wantSkip {
				t.Fatalf("gate=%+v err=%v, want lane skip %q", result, err, wantSkip)
			}
			wantConsole := ""
			if tc.reason != "" {
				wantConsole = "roundfix: notice: Delivery Queue item widget: Linux Lane skipped: " + tc.reason + "\n"
			}
			if tc.name == "passed" {
				wantConsole = "roundfix: Delivery Queue item widget: Linux Lane passed 1 package(s) on golang:1.26.1\n"
			}
			if f.log.String() != wantConsole {
				t.Fatalf("console=%q want=%q", f.log.String(), wantConsole)
			}
			if tc.reason != "" && !strings.HasSuffix(mustRead(t, f.logPath()), wantConsole) {
				t.Fatal("lane log lost skip notice")
			}
			if (tc.name == "disabled" || tc.name == "failed repository gate") && len(f.starts(t)) != 0 {
				t.Fatal("gate started lane tools")
			}
		})
	}
}

func TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	runStore, err := store.Open(t.Context(), home)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = runStore.Close() })
	queue, err := runStore.CreateDeliveryQueue(t.Context(), repo, []string{"0285-example"})
	if err != nil {
		t.Fatal(err)
	}
	item := queue.Items[0]
	item.Stage = store.DeliveryStageMerged
	item.Warning = delivery.WarningLinuxLaneSkipped + ": image golang:1.26 is not present locally"
	if err := runStore.UpdateDeliveryQueueItem(t.Context(), repo, item); err != nil {
		t.Fatal(err)
	}
	if err := runStore.Close(); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	code := runCLI(t, []string{"deliver", "status"}, &stdout, &stderr)
	want := "0285-example\tmerged\t-\t-\nWarning: 0285-example linux-lane-skipped: image golang:1.26 is not present locally\nLimits: deadline none, retries per item none, concurrency 1, tokens none\nUsage: no Runs recorded\n"
	if code != exitOK || stderr.Len() != 0 || stdout.String() != want {
		t.Fatalf("exit=%d stderr=%q stdout=%q want=%q", code, stderr.String(), stdout.String(), want)
	}
}
