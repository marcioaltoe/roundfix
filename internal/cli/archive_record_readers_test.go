package cli

// Suite: queue, review and Run-cause Archive Record readers.
// Boundary IN: disposable Git repositories, local remotes and persisted artifacts.
// Boundary OUT: network, real review providers and daemon Task settlement.

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
	"roundfix/internal/runcause"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

const readerRecordSlug = "0242-record-reader"

func writeReaderRecord(t *testing.T, repo string, r spec.ArchiveRecord) string {
	t.Helper()
	file := spec.ArchiveRecordPath("docs/history/specs", r.Spec)
	content, err := spec.RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	mustMkdir(t, filepath.Dir(filepath.Join(repo, file)))
	mustWrite(t, filepath.Join(repo, file), string(content))
	return file
}

// squashReaderRecord keeps source history on the candidate branch while main
// receives only the record. Missing revisions are explicitly absent objects.
func squashReaderRecord(t *testing.T, repo, slug string, missing bool, override bool) spec.ArchiveRecord {
	t.Helper()
	branch := strings.TrimSpace(gittest.Run(t, repo, "branch", "--show-current"))
	if branch == "main" {
		gittest.Run(t, repo, "checkout", "-b", "feat/archive-record-source")
		branch = "feat/archive-record-source"
	}
	source := "docs/specs/" + slug
	mustMkdir(t, filepath.Join(repo, source))
	mustWrite(t, filepath.Join(repo, source, "_prd.md"), "---\nstatus: approved\n---\n\n# Record reader\n\n## Decisions\n\n- Keep archive context from Git.\n")
	mustWrite(t, filepath.Join(repo, source, "_techspec.md"), "# Source technical design\n\nUse the recorded source revision.\n")
	if _, err := os.Stat(filepath.Join(repo, source, "_tasks.md")); os.IsNotExist(err) {
		mustWrite(t, filepath.Join(repo, source, "_tasks.md"), "---\nschema: spec-tasks/v1\nqa: task_02\ngraph:\n  nodes:\n    - {id: task_01, file: task_01.md}\n    - {id: task_02, file: task_02.md}\n---\n")
	}
	gittest.Run(t, repo, "add", source)
	gittest.Run(t, repo, "commit", "-m", "docs: pre-archive Spec")
	r := spec.ArchiveRecord{Spec: slug, Title: "Record reader", Source: source, SourceRevision: strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")), Disposition: spec.ArchivePass, QATask: "task_02", QAVerdict: "pass"}
	if missing {
		r.SourceRevision = strings.Repeat("a", 40)
	}
	if override {
		r.Disposition = spec.ArchiveQAOverride
		r.QAOverride = &spec.QAArchiveOverrideRecord{Approval: "operator", Reason: "environment unavailable", QAOutcome: "partial", QATaskStatus: "pending", Revision: r.SourceRevision}
	}
	gittest.Run(t, repo, "rm", "-r", "--", source)
	recordPath := writeReaderRecord(t, repo, r)
	gittest.Run(t, repo, "add", recordPath)
	gittest.Run(t, repo, "commit", "-m", "docs: archive record")
	gittest.Run(t, repo, "checkout", "main")
	gittest.Run(t, repo, "merge", "--squash", branch)
	gittest.Run(t, repo, "commit", "-m", "feat: squash archived Spec")
	if history := strings.TrimSpace(gittest.Run(t, repo, "log", "main", "--format=%H", "--", source)); history != "" {
		t.Fatalf("main held Spec folder: %s", history)
	}
	return r
}

func newReaderRecordRepo(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	gittest.InitRepo(t, repo, "-b", "main")
	gittest.PersistIdentity(t, repo)
	mustMkdir(t, filepath.Join(repo, "docs/specs"))
	mustWrite(t, filepath.Join(repo, "docs/specs/.keep"), "Spec root\n")
	gittest.Run(t, repo, "add", ".")
	gittest.Run(t, repo, "commit", "-m", "docs: seed repository")
	return repo
}

func TestInspectItemReadsTheArchiveRecord(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "source present", true: "source absent"}[missing], func(t *testing.T) {
			repo := newReaderRecordRepo(t)
			squashReaderRecord(t, repo, readerRecordSlug, missing, true)
			// Simulate a checkout with no remaining active Specs Root.
			if err := os.RemoveAll(filepath.Join(repo, "docs/specs")); err != nil {
				t.Fatal(err)
			}
			w := retirementWorkflow(repo, t.TempDir())
			state, err := w.InspectItem(t.Context(), repo, readerRecordSlug)
			if err != nil || !state.Archived || !state.QAOverride || len(state.UnfinishedTasks) != 0 || state.Head != strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")) {
				t.Fatalf("state=%+v err=%v", state, err)
			}
		})
	}
}

func TestPrerequisitesCountAnArchiveRecord(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "source present", true: "source absent"}[missing], func(t *testing.T) {
			repo := newReaderRecordRepo(t)
			manifest := filepath.Join(repo, "docs/specs/dependent/_tasks.md")
			mustMkdir(t, filepath.Dir(manifest))
			mustWrite(t, manifest, "---\nschema: spec-tasks/v1\nrequires: ["+readerRecordSlug+", missing]\n---\n")
			gittest.Run(t, repo, "add", "docs/specs/dependent")
			gittest.Run(t, repo, "commit", "-m", "docs: dependent Spec")
			squashReaderRecord(t, repo, readerRecordSlug, missing, false)
			graphs := []*spec.Graph{{Spec: spec.Spec{Slug: "dependent"}, Requires: []string{readerRecordSlug}}}
			root := roundconfig.SpecsRoot{Path: filepath.Join(repo, "docs/specs"), BuiltInRoot: true}
			if err := validateDeliveryPrerequisites(root, graphs); err != nil {
				t.Fatal(err)
			}
			checkout := filepath.Join(t.TempDir(), "checkout")
			gittest.Run(t, "", "clone", "--single-branch", "--no-local", repo, checkout)
			gittest.Harden(t, checkout)
			unmet, err := retirementWorkflow(checkout, t.TempDir()).UnmetPrerequisites(t.Context(), checkout, "dependent")
			if err != nil || !slices.Equal(unmet, []string{"missing"}) {
				t.Fatalf("unmet=%v err=%v", unmet, err)
			}
		})
	}
}

func TestReviewCorrectionAllowsTheArchiveRecordPath(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "source present", true: "source absent"}[missing], func(t *testing.T) {
			f := newReviewCorrectionFixture(t)
			gittest.Run(t, f.review.repository, "rm", "--", correctionArchivePath)
			gittest.Run(t, f.review.repository, "commit", "-m", "docs: replace legacy fixture")
			r := squashReaderRecord(t, f.review.repository, "0225-example", missing, false)
			f.candidate = strings.TrimSpace(gittest.Run(t, f.review.repository, "rev-parse", "HEAD"))
			r.Outcome = "Corrected archive outcome."
			file := writeReaderRecord(t, f.review.repository, r)
			gittest.Run(t, f.review.repository, "add", file)
			gittest.Run(t, f.review.repository, "commit", "-m", "docs: correct record")
			f.head = strings.TrimSpace(gittest.Run(t, f.review.repository, "rev-parse", "HEAD"))
			f.record.HeadCommit = f.candidate
			for i := range f.dispositions {
				f.dispositions[i].HeadCommit = f.candidate
				if f.dispositions[i].Disposition == "fixed" {
					f.dispositions[i].FixedBy = f.head
				}
			}
			f.persist(t)
			if proof := f.prove(t); !proof.Accepted {
				t.Fatalf("proof=%+v", proof)
			}
			f.head = commitReviewCorrectionFile(t, f.review, strings.TrimSuffix(file, ".md")+"-sibling.md", "unrelated\n")
			assertReviewCorrectionRefused(t, f, "outside the named archived Specs")
		})
	}
}

type archiveRecordReviewFlow struct {
	*parkTestDeliveryFlow
	result delivery.ReviewResult
}

func (f *archiveRecordReviewFlow) Review(context.Context, string, string, string) (delivery.ReviewResult, error) {
	return f.result, nil
}

func TestReviewParksForACorrectiveSpecOnAnArchiveRecord(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "source present", true: "source absent"}[missing], func(t *testing.T) {
			runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "Findings:\n- review.txt:2: correction needed", StopReason: "end_turn"}}}}
			f := newReviewCommandFixture(t, "codex", runner)
			squashReaderRecord(t, f.repository, readerRecordSlug, missing, false)
			f.headCommit = strings.TrimSpace(gittest.Run(t, f.repository, "rev-parse", "HEAD"))
			code, record, diag := f.run(t)
			if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || !slices.Equal(record.ArchivedSpecs, []string{readerRecordSlug}) || !strings.Contains(diag, "corrective Spec") {
				t.Fatalf("exit=%d record=%+v stderr=%s", code, record, diag)
			}
			if !missing && (!strings.Contains(runner.request.Prompt, "Keep archive context from Git") || !strings.Contains(runner.request.Prompt, "Source technical design")) {
				t.Fatalf("source context absent: %s", runner.request.Prompt)
			}
			if missing && !slices.Contains(record.SkippedSpecs, readerRecordSlug) {
				t.Fatalf("missing source not reported skipped: %+v", record)
			}
			result, err := deliveryReviewResult(record, f.headCommit)
			if err != nil {
				t.Fatal(err)
			}
			db, err := store.Open(t.Context(), t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			queue, err := db.CreateDeliveryQueue(t.Context(), f.repository, []string{readerRecordSlug})
			if err != nil {
				t.Fatal(err)
			}
			item := queue.Items[0]
			item.Stage, item.Worktree, item.Branch = store.DeliveryStageReviewing, f.repository, "main"
			item.WorktreeProvisioned = true
			item.CandidateCommits = []string{f.headCommit}
			if _, _, _, err := db.RecordDeliveryQueueItemWorktree(t.Context(), f.repository, item.SpecSlug, item.Branch, item.Worktree); err != nil {
				t.Fatal(err)
			}
			if err := db.UpdateDeliveryQueueItem(t.Context(), f.repository, item); err != nil {
				t.Fatal(err)
			}
			flow := &archiveRecordReviewFlow{parkTestDeliveryFlow: &parkTestDeliveryFlow{}, result: result}
			engine := delivery.NewEngine(db, delivery.EngineDependencies{Workspace: flow, Runner: flow, Reviewer: flow, Archiver: flow, Gate: flow, Authorizer: flow, Publication: flow, PullRequests: flow, Revalidator: flow})
			if _, err := engine.Run(t.Context(), f.repository); err != nil {
				t.Fatal(err)
			}
			parked := readDeliveryItemForCLI(t, db, f.repository)
			if parked.Stage != store.DeliveryStageParked || parked.Blocker != "corrective-spec-required: "+readerRecordSlug {
				t.Fatalf("parked=%+v", parked)
			}
		})
	}
}

func TestReviewOverrideConventionReadsTheArchiveRecord(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "source present", true: "source absent"}[missing], func(t *testing.T) {
			repo := newReaderRecordRepo(t)
			squashReaderRecord(t, repo, readerRecordSlug, missing, true)
			r := reviewRepository{Root: repo, Head: strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")), SpecRoots: []string{"docs/specs", "docs/history/specs"}, Git: preflight.ExecGitRunner{}}
			rules, err := conventionRegions(t.Context(), r, reviewFindingAnchor{Path: "docs/history/specs/" + readerRecordSlug + ".md", StartLine: 1, EndLine: 1})
			if err != nil || !slices.Contains(rules, "C5") {
				t.Fatalf("rules=%v err=%v", rules, err)
			}
		})
	}
}

func TestReviewDiffOmitsTheRemovedSpecFolder(t *testing.T) {
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "source present", true: "source absent"}[missing], func(t *testing.T) {
			repo := newReaderRecordRepo(t)
			r := squashReaderRecord(t, repo, readerRecordSlug, missing, false)
			head := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "feat/archive-record-source"))
			base := strings.TrimSpace(gittest.Run(t, repo, "rev-parse", head+"^"))
			diff, omitted, err := reviewScopedDiff(t.Context(), repo, base, head, []string{"docs/specs", "docs/history/specs"}, preflight.ExecGitRunner{})
			if err != nil || strings.Contains(diff, r.Source+"/") || !strings.Contains(diff, "diff --git a/docs/history/specs/"+readerRecordSlug+".md") || len(omitted) != 3 {
				t.Fatalf("diff=%s omitted=%v err=%v", diff, omitted, err)
			}
		})
	}
}

func archivedCausesCLI(t *testing.T, missing bool) (commandEnvironment, spec.ArchiveRecord) {
	t.Helper()
	env := seedCausesCLI(t)
	gittest.PersistIdentity(t, env.workDir)
	// The seeded graph belongs to the delivery branch, never main.
	gittest.Run(t, env.workDir, "commit", "--allow-empty", "-m", "docs: base")
	r := squashReaderRecord(t, env.workDir, "0300-example", missing, false)
	return env, r
}

func TestRunCausesReadAnArchivedTaskGraphFromGit(t *testing.T) {
	env, _ := archivedCausesCLI(t, false)
	code, out, diag := runCausesCLI(t, env, "--format", "json")
	var report runcause.Report
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("decode JSON: %v; exit=%d stderr=%s", err, code, diag)
	}
	if code != 0 || diag != "" || len(report.Items) != 4 || len(report.SpecsNotFound) != 0 || len(report.Tasks) != 4 || report.Tasks[3].Corrective == nil || !*report.Tasks[3].Corrective || len(report.ArchivedSpecs) != 1 || !report.ArchivedSpecs[0].SourceAvailable {
		t.Fatalf("exit=%d report=%+v stderr=%s", code, report, diag)
	}
}

func TestRunCausesReportAnArchiveRecordWithoutItsRevision(t *testing.T) {
	env, _ := archivedCausesCLI(t, true)
	code, out, diag := runCausesCLI(t, env, "--format", "json")
	var report runcause.Report
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatalf("decode JSON: %v; exit=%d stderr=%s", err, code, diag)
	}
	if code != 0 || diag != "" || len(report.SpecsNotFound) != 0 || len(report.ArchivedSpecs) != 1 || report.ArchivedSpecs[0].Disposition != spec.ArchivePass || report.ArchivedSpecs[0].SourceAvailable || len(report.Items) != 3 || report.Tasks[3].Corrective != nil {
		t.Fatalf("exit=%d report=%+v stderr=%s", code, report, diag)
	}
	code, text, diag := runCausesCLI(t, env)
	if code != 0 || diag != "" || !strings.Contains(text, "Spec 0300-example archived (pass); source_revision unavailable") {
		t.Fatalf("exit=%d text=%s stderr=%s", code, text, diag)
	}
}
