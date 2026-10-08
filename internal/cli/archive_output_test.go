package cli

// Suite: Archive output transcripts.
// Invariant: promotions and keyless advice appear in the public confirmation.
// Boundary IN: CLI dispatch and committed temporary repositories.
// Boundary OUT: provider HTTP (a transport that rejects every request).

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/judge"
	"roundfix/internal/spec"
)

func TestArchiveConfirmationNamesPromotedFiles(t *testing.T) {
	t.Parallel()
	for _, override := range []bool{false, true} {
		t.Run(fmt.Sprint(override), func(t *testing.T) {
			home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
			verdict := spec.VerdictPass
			if override {
				verdict = spec.VerdictFail
			}
			writeArchiveQAReport(t, repo, verdict)
			reference := filepath.Join(repo, "docs/specs", implementTestSlug, "references/knowledge.md")
			mustMkdir(t, filepath.Dir(reference))
			mustWrite(t, reference, "# Knowledge\nA reusable lesson.\n")
			commitArchiveFixture(t)
			args := []string{"archive", implementTestSlug, "--promote", "references/knowledge.md"}
			if override {
				args = append(args, "--qa-override", "--approval", "maintainer", "--reason", strings.Repeat("r", 300))
			}
			fake := &specJudgeTransport{t: t}
			env := commandEnvironment{homeDir: home, workDir: repo, environ: []string{}, dependencies: defaultCommandDependencies()}
			env.dependencies.judgeTransport = fake
			var out, errout bytes.Buffer
			if code := runWithContext(t.Context(), args, &out, &errout, env); code != exitOK {
				t.Fatalf("exit=%d stderr=%q", code, &errout)
			}
			path := archiveTestRepositoryPath(repo, spec.ArchiveKindSpec, implementTestSlug) + ".md"
			want := strings.TrimSuffix(archiveConfirmation(t, repo, path, override), "\n") + "; promoted 1 file(s) to docs/references/\n"
			if out.String() != want || errout.Len() != 0 {
				t.Fatalf("stdout=%q want=%q stderr=%q", &out, want, &errout)
			}
			if fake.calls != 0 {
				t.Fatalf("provider requests=%d, want zero", fake.calls)
			}
		})
	}
}

func TestArchivePlanWithoutAKeyPrintsNoAdvice(t *testing.T) {
	// Sequential: clears process-wide provider API key environment variables.
	questions, err := judge.Load()
	if err != nil {
		t.Fatal(err)
	}
	keys := questions.KeyVariables()
	if len(keys) == 0 {
		t.Fatal("judge has no key variables")
	}
	for _, key := range keys {
		t.Setenv(key, "")
	}
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", status: string(spec.StatusCompleted)}})
	const content = "# Knowledge\nA reusable lesson.\n"
	path := filepath.Join(repo, "docs/specs", implementTestSlug, "references/knowledge.md")
	mustMkdir(t, filepath.Dir(path))
	mustWrite(t, path, content)
	commitArchiveFixture(t)
	fake := &specJudgeTransport{t: t}
	env := commandEnvironment{homeDir: home, workDir: repo, environ: []string{}, dependencies: defaultCommandDependencies()}
	for _, key := range keys {
		env.environ = append(env.environ, key+"=")
	}
	env.dependencies.judgeTransport = fake
	var out, errout bytes.Buffer
	if code := runWithContext(t.Context(), []string{"archive", implementTestSlug, "--plan"}, &out, &errout, env); code != exitOK {
		t.Fatalf("exit=%d stderr=%q", code, &errout)
	}
	want := fmt.Sprintf("candidate references/knowledge.md %d bytes: no advice (%s is not set)", len(content), keys[0])
	found := false
	for _, line := range strings.Split(out.String(), "\n") {
		if line == want {
			found = true
		}
	}
	if !found || errout.Len() != 0 {
		t.Fatalf("stdout=%q want line=%q stderr=%q", &out, want, &errout)
	}
	if fake.calls != 0 {
		t.Fatalf("provider requests=%d, want zero", fake.calls)
	}
}
