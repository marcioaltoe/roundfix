// Boundary: full archive-plan transcripts in disposable repositories; judge HTTP is counted locally.
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
	"roundfix/internal/spec"
)

func archiveOfflineFixture(t *testing.T) (commandEnvironment, *specJudgeTransport) {
	t.Helper()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	path := filepath.Join(repo, "docs/specs", implementTestSlug, "references/knowledge.md")
	mustMkdir(t, filepath.Dir(path))
	mustWrite(t, path, "# Knowledge\nA reusable lesson.\n")
	commitArchiveFixture(t)
	transport := &specJudgeTransport{t: t}
	env := commandEnvironment{homeDir: home, workDir: repo, environ: []string{}, dependencies: defaultCommandDependencies()}
	env.dependencies.judgeTransport = transport
	return env, transport
}

func TestArchivePlanOfflineTranscript(t *testing.T) {
	t.Parallel()
	env, _ := archiveOfflineFixture(t)
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
	want := fmt.Sprintf("archive plan for %s: removes %d file(s) (%d bytes) and writes docs/history/specs/%s.md\ncore %d file(s): %s\nevidence 0 file(s) (0 bytes) under qa/evidence/\njudge: skipped (offline)\ncandidate references/knowledge.md %d bytes: no advice (offline)\npromote with: roundfix archive %s --promote <path>\n", implementTestSlug, len(before), total, implementTestSlug, len(core), strings.Join(core, ", "), len(before["references/knowledge.md"]), implementTestSlug)
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", implementTestSlug, "--plan", "--offline"}, &out, &errout, env)
	if code != 0 || out.String() != want || errout.Len() != 0 {
		t.Fatalf("exit=%d stdout=%q stderr=%q; want exit=0 stdout=%q stderr=empty", code, &out, &errout, want)
	}
	if !reflect.DeepEqual(before, snapshotDirectoryFiles(t, root)) {
		t.Fatal("offline plan changed Spec files")
	}
}

func TestArchivePlanOfflineSkipsTheJudge(t *testing.T) {
	t.Parallel()
	env, transport := archiveOfflineFixture(t)
	questions, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	keys := questions.KeyVariables()
	if len(keys) == 0 {
		t.Fatal("judge has no key variables")
	}
	for _, key := range keys {
		env.environ = append(env.environ, key+"=offline-test-value")
	}
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", implementTestSlug, "--plan", "--offline"}, &out, &errout, env)
	if code != 0 || errout.Len() != 0 {
		t.Fatalf("exit=%d stderr=%q", code, &errout)
	}
	if transport.calls != 0 {
		t.Fatalf("judge requests=%d, want zero", transport.calls)
	}
	if !strings.Contains(out.String(), "judge: skipped (offline)\n") {
		t.Fatalf("stdout=%q", &out)
	}
}

func TestArchiveOfflineRequiresPlan(t *testing.T) {
	t.Parallel()
	env, _ := archiveOfflineFixture(t)
	var out, errout bytes.Buffer
	code := runWithContext(t.Context(), []string{"archive", implementTestSlug, "--offline"}, &out, &errout, env)
	want := archivePullRequestPreflight("--offline requires --plan")
	if code != 2 || out.Len() != 0 || errout.String() != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q; want exit=2 stderr=%q", code, &out, &errout, want)
	}
}
