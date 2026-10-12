package cli

// Suite: cancelled Archive Record readers.
// Boundary IN: public commands, disposable Git repositories and temporary homes.
// Boundary OUT: live providers, detached owners and daemon Task settlement.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/gittest"
	"roundfix/internal/runcause"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

func cancelledReaderRecord(slug, revision string) spec.ArchiveRecord {
	return spec.ArchiveRecord{
		Spec: slug, Title: "Cancelled reader", Created: "2026-10-10", Archived: "2026-10-12",
		Source: "docs/specs/" + slug, SourceRevision: revision, Disposition: spec.ArchiveCancelled,
		CancelReason: "no consumer", CompletedTasks: []string{"task_01"}, PendingTasks: []string{"task_02"},
		Outcome: "Original intent.",
	}
}

func assertCancelledReaderRefusal(t *testing.T, command string, args []string, reason, repo, home string) {
	t.Helper()
	beforeRepo, beforeHome := historySnapshot(t, repo), historySnapshot(t, home)
	status := gittest.Run(t, repo, "status", "--porcelain", "--untracked-files=all")
	var out, diag, expected bytes.Buffer
	code := runCLI(t, args, &out, &diag)
	printPreflightFailure(command, fmt.Errorf("%s", reason), &expected)
	if code != exitPreflight || out.Len() != 0 || diag.String() != expected.String() {
		t.Fatalf("exit=%d stdout=%q stderr=%q; want exit 2 and stderr=%q", code, &out, &diag, &expected)
	}
	if !reflect.DeepEqual(beforeRepo, historySnapshot(t, repo)) || !reflect.DeepEqual(beforeHome, historySnapshot(t, home)) {
		t.Fatal("refusal changed repository or home bytes")
	}
	if after := gittest.Run(t, repo, "status", "--porcelain", "--untracked-files=all"); after != status {
		t.Fatalf("refusal changed Git status: %q -> %q", status, after)
	}
}

func TestSupersedeRefusesACancelledDeliverer(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	writeReaderRecord(t, repo, cancelledReaderRecord(readerRecordSlug, strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))))
	assertCancelledReaderRefusal(t, "supersede", []string{"supersede", "--spec", implementTestSlug, "--by", readerRecordSlug, "--reason", "replacement"},
		fmt.Sprintf("superseding Spec %q is cancelled; a cancelled Spec delivers no content", readerRecordSlug), repo, home)
}

func TestDeliverStartRefusesARequiredCancelledSpec(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	writeReaderRecord(t, repo, cancelledReaderRecord(readerRecordSlug, strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))))
	manifest := filepath.Join(repo, "docs/specs", implementTestSlug, "_tasks.md")
	mustWrite(t, manifest, strings.Replace(mustRead(t, manifest), "schema: spec-tasks/v1\n", "schema: spec-tasks/v1\nrequires: ["+readerRecordSlug+"]\n", 1))
	assertCancelledReaderRefusal(t, "deliver start", []string{"deliver", "start", implementTestSlug},
		fmt.Sprintf("Delivery Queue Spec %s requires cancelled Spec %s", implementTestSlug, readerRecordSlug), repo, home)
	assertPathMissing(t, store.DatabasePath(home))
}

func TestDeliverRetryRefusesACancelledSpec(t *testing.T) {
	t.Parallel()
	for _, database := range []string{"parked queue", "absent", "malformed"} {
		t.Run(database, func(t *testing.T) {
			home, repo := newParkedDeliveryQueueForRetry(t, 0, "")
			record := cancelledReaderRecord(implementTestSlug, strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")))
			if err := os.RemoveAll(filepath.Join(repo, record.Source)); err != nil {
				t.Fatal(err)
			}
			recordPath := writeReaderRecord(t, repo, record)
			var before store.DeliveryQueue
			if database == "parked queue" {
				db, err := store.Open(t.Context(), home)
				if err != nil {
					t.Fatal(err)
				}
				before, _, err = db.DeliveryQueue(t.Context(), repo)
				if err != nil {
					t.Fatal(err)
				}
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			} else if database == "absent" {
				if err := os.Remove(store.DatabasePath(home)); err != nil {
					t.Fatal(err)
				}
			} else {
				mustWrite(t, store.DatabasePath(home), "not a database\n")
			}
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.newDeliveryEngine = func(*store.Store, roundconfig.Loaded) deliveryEngine {
					t.Fatal("cancelled retry reached the Delivery Engine")
					return nil
				}
			})
			assertCancelledReaderRefusal(t, "deliver retry", []string{"deliver", "retry", implementTestSlug},
				fmt.Sprintf("Spec %q is cancelled in %s; a cancelled Spec is never delivered", implementTestSlug, filepath.Join(repo, recordPath)), repo, home)
			if database == "parked queue" {
				db, err := store.Open(t.Context(), home)
				if err != nil {
					t.Fatal(err)
				}
				defer db.Close()
				after, found, err := db.DeliveryQueue(t.Context(), repo)
				if err != nil || !found || !reflect.DeepEqual(before, after) {
					t.Fatalf("queue rows changed: before=%+v after=%+v found=%v err=%v", before, after, found, err)
				}
			}
		})
	}
}

func TestReviewReadsACancelledArchiveRecord(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings", StopReason: "end_turn"}}}}
	f := newReviewCommandFixture(t, "codex", runner)
	source := "docs/specs/" + readerRecordSlug
	mustMkdir(t, filepath.Join(f.repository, source))
	mustWrite(t, filepath.Join(f.repository, source, "_prd.md"), "# Cancelled intent\n\n## Decisions\n\n- Keep cancellation context from Git.\n")
	mustWrite(t, filepath.Join(f.repository, source, "_techspec.md"), "# Cancelled technical design\n")
	f.baseCommit = historyCommit(t, f.repository, "docs: source before cancellation")
	gittest.Run(t, f.repository, "rm", "-r", "--", source)
	path := writeReaderRecord(t, f.repository, cancelledReaderRecord(readerRecordSlug, f.baseCommit))
	f.headCommit = historyCommit(t, f.repository, "docs: cancelled record")
	before := mustRead(t, filepath.Join(f.repository, path))
	code, record, diag := f.run(t)
	if code != exitOK || diag != "" || record.Outcome != reviewOutcomeReviewed || runner.preparedCalls != 1 || !slices.Equal(record.ArchivedSpecs, []string{readerRecordSlug}) || !slices.Equal(record.Specs, []string{readerRecordSlug}) || len(record.SkippedSpecs) != 0 {
		t.Fatalf("exit=%d record=%+v calls=%d stderr=%s", code, record, runner.preparedCalls, diag)
	}
	for _, want := range []string{"diff --git a/" + path, "disposition: cancelled", "Keep cancellation context from Git", "Cancelled technical design"} {
		if !strings.Contains(runner.request.Prompt, want) {
			t.Fatalf("provider prompt missing %q: %s", want, runner.request.Prompt)
		}
	}
	if strings.Contains(runner.request.Prompt, "diff --git a/"+source+"/") || len(record.OmittedPaths) != 2 {
		t.Fatalf("removed Spec folder reached provider diff: prompt=%s omitted=%v", runner.request.Prompt, record.OmittedPaths)
	}
	for _, omitted := range record.OmittedPaths {
		if !strings.HasPrefix(omitted.Path, source+"/") || omitted.Reason != "archived-spec-source" {
			t.Fatalf("unexpected omission: %+v", omitted)
		}
	}
	if after := mustRead(t, filepath.Join(f.repository, path)); after != before {
		t.Fatal("review rewrote cancelled record")
	}
}

func TestRunCausesReportACancelledRecord(t *testing.T) {
	t.Parallel()
	for _, missing := range []bool{false, true} {
		t.Run(map[bool]string{false: "source present", true: "source absent"}[missing], func(t *testing.T) {
			env, source := archivedCausesCLI(t, missing)
			record := cancelledReaderRecord(source.Spec, source.SourceRevision)
			path := writeReaderRecord(t, env.workDir, record)
			before := mustRead(t, filepath.Join(env.workDir, path))
			code, out, diag := runCausesCLI(t, env, "--format", "json")
			var report runcause.Report
			if err := json.Unmarshal([]byte(out), &report); err != nil {
				t.Fatal(err)
			}
			if code != exitOK || diag != "" || len(report.SpecsNotFound) != 0 || len(report.ArchivedSpecs) != 1 || report.ArchivedSpecs[0].Disposition != spec.ArchiveCancelled || report.ArchivedSpecs[0].SourceAvailable == missing || len(report.Tasks) != 4 {
				t.Fatalf("exit=%d report=%+v stderr=%s", code, report, diag)
			}
			code, out, diag = runCausesCLI(t, env)
			if code != exitOK || diag != "" || (missing && !strings.Contains(out, "Spec "+record.Spec+" archived (cancelled); source_revision unavailable")) {
				t.Fatalf("exit=%d stdout=%s stderr=%s", code, out, diag)
			}
			if after := mustRead(t, filepath.Join(env.workDir, path)); after != before {
				t.Fatal("Run causes rewrote cancelled record")
			}
		})
	}
}

func TestHistorySanitizeLeavesACancelledRecord(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, false)
	path := writeReaderRecord(t, h.repo, cancelledReaderRecord(readerRecordSlug, h.revision))
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/aaa/_prd.md"), "---\nspec: aaa\ncreated: 2026-10-10\n---\n\n# Other archive\n\nA delivered outcome.\n")
	historyCommit(t, h.repo, "docs: other archive (#42)")
	historyTag(t, h)
	before := historySnapshot(t, h.repo)
	code, out, diag := historyRun(t, h, "sanitize")
	if code != exitOK || diag != "" || !strings.Contains(out, "1 unit(s) pending") || !strings.Contains(out, "folder docs/history/specs/aaa:") || strings.Contains(out, readerRecordSlug) {
		t.Fatalf("plan exit=%d stdout=%s stderr=%s", code, out, diag)
	}
	historyUnchanged(t, h, before, "")
	code, out, diag = historyRun(t, h, "sanitize", "--apply", "--batch", "1")
	if code != exitOK || diag != "" || !strings.Contains(out, "applied 1 unit(s): wrote 1 Archive Record(s)") {
		t.Fatalf("apply exit=%d stdout=%s stderr=%s", code, out, diag)
	}
	assertPathMissing(t, filepath.Join(h.repo, "docs/history/specs/aaa"))
	archiveRecordAt(t, filepath.Join(h.repo, "docs/history/specs/aaa.md"))
	if after := mustRead(t, filepath.Join(h.repo, path)); !bytes.Equal(before[path], []byte(after)) {
		t.Fatal("sanitize rewrote cancelled record")
	}
}

func TestSupersedeAcceptsOtherArchiveRecordDispositions(t *testing.T) {
	t.Parallel()
	for _, disposition := range []spec.ArchiveDisposition{spec.ArchivePass, spec.ArchivePartial, spec.ArchiveQAOverride, spec.ArchiveFailedQA, spec.ArchiveSuperseded, spec.ArchiveNoQA} {
		t.Run(string(disposition), func(t *testing.T) {
			_, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			record := spec.ArchiveRecord{
				Spec: readerRecordSlug, Title: "Delivered reader", Source: "docs/specs/" + readerRecordSlug,
				SourceRevision: strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD")), Disposition: disposition,
				QAVerdict: spec.VerdictPass, QAReport: "qa-report.md",
			}
			switch disposition {
			case spec.ArchiveFailedQA:
				record.QAVerdict = spec.VerdictFail
			case spec.ArchivePartial:
				record.QAVerdict = spec.VerdictPartial
			case spec.ArchiveQAOverride:
				record.QAOverride = &spec.QAArchiveOverrideRecord{Approval: "maintainer", Reason: "environment unavailable"}
			case spec.ArchiveSuperseded:
				record.SupersededBy = "successor"
			case spec.ArchiveNoQA:
				record.QAVerdict, record.QAReport = "", ""
			}
			writeReaderRecord(t, repo, record)
			var out, diag bytes.Buffer
			code := runCLI(t, []string{"supersede", "--spec", implementTestSlug, "--by", readerRecordSlug, "--reason", "delivered"}, &out, &diag)
			if code != exitOK || diag.Len() != 0 || out.String() != fmt.Sprintf("superseded %s by %s\n", implementTestSlug, readerRecordSlug) {
				t.Fatalf("disposition=%s exit=%d stdout=%s stderr=%s", disposition, code, &out, &diag)
			}
			root := roundconfig.SpecsRoot{Path: filepath.Join(repo, "docs/specs"), BuiltInRoot: true}
			if err := validateDeliveryPrerequisites(root, []*spec.Graph{{Spec: spec.Spec{Slug: implementTestSlug}, Requires: []string{readerRecordSlug}}}); err != nil {
				t.Fatalf("disposition=%s prerequisite refused: %v", disposition, err)
			}
		})
	}
}
