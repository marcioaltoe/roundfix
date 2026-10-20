package cli

// Suite: the complete history conversion Surface Transcript 1.
// Boundary IN: CLI dispatch and an isolated, tagged local Git repository.
// Boundary OUT: every stdout line, empty stderr and exit status.

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func historyTranscriptTestHome(t *testing.T) bool {
	t.Helper()
	const key = "ROUNDFIX_HISTORY_TRANSCRIPT_TEST"
	if os.Getenv(key) == t.Name() {
		return true
	}
	ctx, cancel := context.WithCancel(t.Context())
	cmd := exec.CommandContext(ctx, os.Args[0], "-test.run=^"+regexp.QuoteMeta(t.Name())+"$", "-test.count=1")
	for _, entry := range gittest.IsolatedEnv() {
		if !strings.HasPrefix(entry, "HOME=") && !strings.HasPrefix(entry, key+"=") {
			cmd.Env = append(cmd.Env, entry)
		}
	}
	cmd.Env = append(cmd.Env, "HOME="+t.TempDir(), key+"="+t.Name())
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	waited := false
	t.Cleanup(func() {
		cancel()
		if cmd.Process != nil && !waited {
			_ = cmd.Wait()
		}
	})
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	err := cmd.Wait()
	waited = true
	if err != nil {
		t.Fatalf("isolated history transcript: %v\n%s", err, output.String())
	}
	return false
}

func TestHistorySanitizeRewritesEveryCitationTranscript(t *testing.T) {
	t.Parallel()
	if !historyTranscriptTestHome(t) {
		return
	}
	h := historyFixture{repo: t.TempDir(), home: os.Getenv("HOME")}
	gittest.InitRepo(t, h.repo, "-b", "main")
	files := map[string]string{
		"docs/history/findings/2026-01-01-x.md":   "---\nstatus: done\n---\n# Finding\n",
		"docs/history/backlog/2026-01-02-b.md":    "---\nstatus: done\n---\n# Backlog\n",
		"docs/history/adr/0001-old.md":            "---\nstatus: deprecated\n---\n# Old\n",
		"docs/adr/0002-live.md":                   "---\nstatus: accepted\n---\n# Live\n\n`docs/history/findings/2026-01-01-x.md` [old](../history/adr/0001-old.md)\n",
		"docs/agents/deploy.md":                   "`../history/findings/2026-01-01-x.md`\n",
		"docs/reference/notes.md":                 "[b](../history/backlog/2026-01-02-b.md)\n",
		"infra/README.md":                         "[b](../docs/history/backlog/2026-01-02-b.md) `../docs/history/findings/2026-01-01-x.md`\n",
		"Makefile":                                "all:\n# see docs/history/findings/2026-01-01-x.md\n",
		".agents/skills/demo/SKILL.md":            "docs/history/findings/2026-01-01-x.md\n",
		"internal/testdata/case.md":               "docs/history/findings/2026-01-01-x.md\n",
		"docs/adr/2026-09-18-old-cluster.md":      "---\nstatus: superseded\nupdated_at: 2026-10-01T00:00:00Z\nsuperseded_by: 2026-10-01-new-cluster\n---\n# Old cluster\n",
		"docs/adr/2026-10-01-new-cluster.md":      "---\nstatus: accepted\n---\n\n# New cluster\n\nCurrent cluster.\nLine eight.\n",
		"docs/specs/0001-demo/_prd.md":            "---\nstatus: active\n---\n\n# Demo\n\nA future measurement.\n\n[m](references/medicao.md)\n",
		"docs/specs/0001-demo/task_01.md":         "---\ntask: task_01\nstatus: pending\ntype: docs\n---\n## Context\n- creates: docs/specs/0001-demo/references/medicao.md\n",
		"docs/specs/0001-demo/references/note.md": "[f](../findings/2026-01-01-x.md)\n",
	}
	for path, content := range files {
		historyWrite(t, filepath.Join(h.repo, path), content)
	}
	gittest.Run(t, h.repo, "add", "--all")
	ctx, cancel := context.WithCancel(t.Context())
	args := append([]string{"-C", h.repo}, gittest.ConfigArgs()...)
	cmd := exec.CommandContext(ctx, "git", append(args, "commit", "-m", "docs: transcript fixture")...)
	cmd.Env = append(gittest.IsolatedEnv(), "GIT_AUTHOR_DATE=2026-10-19T12:00:00Z", "GIT_COMMITTER_DATE=2026-10-19T12:00:00Z")
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	waited := false
	t.Cleanup(func() {
		cancel()
		if cmd.Process != nil && !waited {
			_ = cmd.Wait()
		}
	})
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	err := cmd.Wait()
	waited = true
	if err != nil {
		t.Fatalf("commit: %v\n%s", err, output.String())
	}
	h.revision = strings.TrimSpace(gittest.Run(t, h.repo, "rev-parse", "HEAD"))
	historyTag(t, h)
	h.env = commandEnvironment{workDir: h.repo, homeDir: h.home, environ: []string{"PATH=" + os.Getenv("PATH"), "HOME=" + h.home}}
	code, stdout, stderr := historyRun(t, h, "sanitize")
	const want = "history sanitize plan: 4 unit(s) pending; 4 file(s) (<b> bytes) leave docs/history\n" +
		"findings: removes 1 file(s) (<b> bytes)\n" +
		"backlog: removes 1 file(s) (<b> bytes)\n" +
		"adr: removes 2 file(s) (<b> bytes)\n" +
		"live-links: rewrites 1 link(s) in 1 file(s)\n" +
		"cites Makefile:2 names docs/history/findings/2026-01-01-x.md; rewrites to docs/history/findings/2026-01-01-x.md@<commit>\n" +
		"cites docs/adr/0002-live.md:6 names docs/history/adr/0001-old.md; rewrites to `docs/history/adr/0001-old.md@<commit>`\n" +
		"cites docs/adr/0002-live.md:6 names docs/history/findings/2026-01-01-x.md; rewrites to `docs/history/findings/2026-01-01-x.md@<commit>`\n" +
		"cites docs/agents/deploy.md:1 names docs/history/findings/2026-01-01-x.md; rewrites to `docs/history/findings/2026-01-01-x.md@<commit>`\n" +
		"cites docs/reference/notes.md:1 names docs/history/backlog/2026-01-02-b.md; rewrites to `docs/history/backlog/2026-01-02-b.md@<commit>`\n" +
		"cites infra/README.md:1 names docs/history/backlog/2026-01-02-b.md; rewrites to `docs/history/backlog/2026-01-02-b.md@<commit>`\n" +
		"cites infra/README.md:1 names docs/history/findings/2026-01-01-x.md; rewrites to `docs/history/findings/2026-01-01-x.md@<commit>`\n" +
		"cites docs/adr/2026-10-01-new-cluster.md:10 names docs/adr/2026-09-18-old-cluster.md; rewrites to Supersedes 2026-09-18-old-cluster (2026-10-01); Git keeps its text at `docs/adr/2026-09-18-old-cluster.md@<commit>`.\n" +
		"cites docs/specs/0001-demo/references/note.md:1 names docs/specs/0001-demo/findings/2026-01-01-x.md; rewrites to `docs/history/findings/2026-01-01-x.md@<commit>`\n" +
		"apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)\n"
	pattern := regexp.QuoteMeta(strings.ReplaceAll(want, "<commit>", h.revision))
	pattern = strings.ReplaceAll(pattern, regexp.QuoteMeta("<b>"), "[0-9]+")
	if code != 0 || stderr != "" || !regexp.MustCompile("\\A"+pattern+"\\z").MatchString(stdout) {
		t.Fatalf("exit=%d stderr=%q\nstdout:\n%s\nwant:\n%s", code, stderr, stdout, want)
	}
}
