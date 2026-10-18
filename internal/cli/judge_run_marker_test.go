// Boundary: public command dispatch with injected Run Marker and credentials.
package cli

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/judge"
)

const commandRunMarkerFixture = "run-command-fixture"
const commandRunMarkerReason = "inside Run run-command-fixture, named by ROUNDFIX_RUN_ID"

func TestSpecJudgeSkipsInsideARun(t *testing.T) {
	t.Parallel()
	env, transport := specJudgeFixture(t, "ROUNDFIX_OPENROUTER_API_KEY")
	env.environ = append(env.environ, judge.RunMarkerVariable+"="+commandRunMarkerFixture)
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"spec", "judge", "0300-example"}, &out, &errout, env)
	want := "Judge: skipped: " + commandRunMarkerReason + "; 5 judgment(s) not asked\n"
	if code != 0 || out.String() != want || errout.Len() != 0 || transport.calls != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q calls=%d; want stdout=%q", code, &out, &errout, transport.calls, want)
	}
	assertCommandRunMarkerNoLog(t, env.homeDir)
}

func TestArchivePlanInsideARunMatchesItsTranscript(t *testing.T) {
	t.Parallel()
	// Surface Transcript 1's keyless case and the same transcript with a key.
	for _, withKey := range []bool{false, true} {
		t.Run(fmt.Sprintf("key=%t", withKey), func(t *testing.T) {
			t.Parallel()
			env, transport := archiveOfflineFixture(t)
			env.environ = append(env.environ, judge.RunMarkerVariable+"="+commandRunMarkerFixture)
			if withKey {
				env.environ = append(env.environ, "ROUNDFIX_OPENROUTER_API_KEY=fixture-key")
			}
			root := filepath.Join(env.workDir, "docs/specs", implementTestSlug)
			before := snapshotDirectoryFiles(t, root)
			var core []string
			total := 0
			for path, data := range before {
				total += len(data)
				if path != "references/knowledge.md" {
					core = append(core, path)
				}
			}
			sort.Strings(core)
			want := fmt.Sprintf("archive plan for %s: removes %d file(s) (%d bytes) and writes docs/history/specs/%s.md\ncore %d file(s): %s\nevidence 0 file(s) (0 bytes) under qa/evidence/\ncandidate references/knowledge.md %d bytes: no advice (inside Run %s, named by ROUNDFIX_RUN_ID)\npromote with: roundfix archive %s --promote <path>\n", implementTestSlug, len(before), total, implementTestSlug, len(core), strings.Join(core, ", "), len(before["references/knowledge.md"]), commandRunMarkerFixture, implementTestSlug)
			var out, errout bytes.Buffer
			code := runWithContext(t.Context(), []string{"archive", implementTestSlug, "--plan"}, &out, &errout, env)
			if code != 0 || out.String() != want || errout.Len() != 0 || transport.calls != 0 {
				t.Fatalf("exit=%d stdout=%q stderr=%q calls=%d; want exit=0 stdout=%q stderr=empty", code, &out, &errout, transport.calls, want)
			}
			if !reflect.DeepEqual(before, snapshotDirectoryFiles(t, root)) {
				t.Fatal("plan changed Spec files")
			}
			assertCommandRunMarkerNoLog(t, env.homeDir)
			// Explicit offline output takes precedence even with a marker and key.
			out.Reset()
			errout.Reset()
			code = runWithContext(t.Context(), []string{"archive", implementTestSlug, "--plan", "--offline"}, &out, &errout, env)
			offline := strings.Replace(want, "candidate references/", "judge: skipped (offline)\ncandidate references/", 1)
			offline = strings.Replace(offline, commandRunMarkerReason, "offline", 1)
			if code != 0 || out.String() != offline || errout.Len() != 0 || transport.calls != 0 {
				t.Fatalf("offline exit=%d stdout=%q stderr=%q calls=%d; want %q", code, &out, &errout, transport.calls, offline)
			}
		})
	}
}

func TestHistorySanitizeAdviceSkipsInsideARun(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, true)
	transport := &specJudgeTransport{t: t}
	h.env.dependencies.judgeTransport = transport
	h.env.environ = append(h.env.environ, judge.RunMarkerVariable+"="+commandRunMarkerFixture, "ROUNDFIX_OPENROUTER_API_KEY=fixture-key")
	before := historySnapshot(t, h.repo)
	code, out, errout := historyRun(t, h, "sanitize", "--batch", "1", "--advise")
	want := "candidate docs/history/specs/aaa/references/lesson.md 20 bytes: no advice (" + commandRunMarkerReason + ")\n"
	if code != 0 || errout != "" || !strings.Contains(out, want) || transport.calls != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q calls=%d; want candidate=%q", code, out, errout, transport.calls, want)
	}
	historyUnchanged(t, h, before, "")
	assertCommandRunMarkerNoLog(t, h.home)
}

func assertCommandRunMarkerNoLog(t *testing.T, home string) {
	t.Helper()
	if files, err := filepath.Glob(filepath.Join(home, ".roundfix", "judge", "*")); err != nil || len(files) != 0 {
		t.Fatalf("Judge Log files=%v err=%v", files, err)
	}
}
