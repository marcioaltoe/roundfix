// Boundary: CLI configuration loading, temporary repository and home.
package cli

import (
	"bytes"
	"path/filepath"
	"roundfix/internal/gittest"
	"testing"
)

const retiredLaneWarnings = "config: delivery.linux_lane is deprecated and ignored; the Linux Lane was retired, so remove it\nconfig: delivery.linux_lane_image is deprecated and ignored; the Linux Lane was retired, so remove it\n"

func TestRunsListReproducesSurfaceTranscriptRetiredLaneKeys(t *testing.T) {
	t.Parallel()
	_, repo := withCLIWorkspace(t)
	gittest.InitRepo(t, repo, "-b", "main")
	mustWrite(t, filepath.Join(repo, ".roundfixrc.yml"), "delivery:\n  linux_lane: true\n  linux_lane_image: golang:1.26\n")
	var out, diagnostic bytes.Buffer
	code := runCLI(t, []string{"runs", "list"}, &out, &diagnostic)
	if code != 0 || out.String() != "No Runs found.\n" || diagnostic.String() != retiredLaneWarnings {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, out.String(), diagnostic.String())
	}
}
