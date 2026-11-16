package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/agent"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/phaserecord"
	"roundfix/internal/runevent"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

var phasePromptNow = time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)

func TestEveryEngineTheCLIBuildsRecordsItsPrompts(t *testing.T) {
	t.Parallel()
	fset := token.NewFileSet()
	packages, err := parser.ParseDir(fset, ".", func(info os.FileInfo) bool { return !strings.HasSuffix(info.Name(), "_test.go") }, 0)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {
				literal, ok := node.(*ast.CompositeLit)
				if !ok {
					return true
				}
				selector, ok := literal.Type.(*ast.SelectorExpr)
				if !ok || selector.Sel.Name != "Dependencies" {
					return true
				}
				name, ok := selector.X.(*ast.Ident)
				if !ok || name.Name != "daemon" {
					return true
				}
				count++
				for _, field := range literal.Elts {
					pair, ok := field.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					key, ok := pair.Key.(*ast.Ident)
					if !ok || key.Name != "PhaseRecord" {
						continue
					}
					if value, ok := pair.Value.(*ast.Ident); ok && value.Name == "nil" {
						t.Errorf("nil PhaseRecord at %s", fset.Position(pair.Pos()))
					}
					return true
				}
				t.Errorf("daemon.Dependencies at %s sets no PhaseRecord", fset.Position(literal.Pos()))
				return true
			})
		}
	}
	if count < 3 {
		t.Fatalf("found only %d engines", count)
	}
}

func promptRecordRows(t *testing.T, home string) []map[string]string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(home, ".roundfix", "phase-record", phasePromptNow.Format("2006-01-02")+".jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasSuffix(data, []byte("\n")) {
		t.Fatal("missing final newline")
	}
	var rows []map[string]string
	for _, line := range bytes.Split(bytes.TrimSuffix(data, []byte("\n")), []byte("\n")) {
		var row map[string]string
		if err := json.Unmarshal(line, &row); err != nil {
			t.Fatal(err)
		}
		rows = append(rows, row)
	}
	return rows
}

func assertPromptRow(t *testing.T, row map[string]string, req agent.ExecuteRequest, runID, specSlug, kind, scopeID, repository, session string) {
	t.Helper()
	started, err := time.Parse(time.RFC3339, row["started"])
	if err != nil {
		t.Fatal(err)
	}
	ended, err := time.Parse(time.RFC3339, row["ended"])
	if err != nil || started.After(ended) {
		t.Fatalf("times = %v..%v, %v", started, ended, err)
	}
	want := map[string]string{"schema": phaserecord.Schema, "kind": "prompt", "time": phasePromptNow.Format(time.RFC3339), "started": row["started"], "ended": row["ended"], "run_id": runID, "spec": specSlug, "scope_kind": kind, "scope_id": scopeID, "repository": repository, "work_dir": req.Session.WorkDir, "harness": req.Runtime.ID, "session": session}
	if !reflect.DeepEqual(row, want) {
		t.Fatalf("row = %+v, want %+v", row, want)
	}
}

type phaseImplementRunner struct {
	*implementFakeRunner
	requests []agent.ExecuteRequest
}

func (r *phaseImplementRunner) Run(ctx context.Context, req agent.ExecuteRequest, sink runevent.Sink) (agent.ExecuteResult, error) {
	r.requests = append(r.requests, req)
	result, err := r.implementFakeRunner.Run(ctx, req, sink)
	result.ACPSessionID = "implement-adapter-session"
	return result, err
}

func TestAnImplementRunLeavesOnePromptEntryPerPrompt(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01", title: "First"}, {id: "task_02", title: "Second", needs: []string{"task_01"}}})
	writeUserConfig(t, home, "verification:\n  repository_at_settlement: false\n")
	runner := &phaseImplementRunner{implementFakeRunner: &implementFakeRunner{gitRoot: repo, statusByTask: map[string]spec.Status{"task_01": spec.StatusCompleted, "task_02": spec.StatusCompleted}}}
	withImplementCollaborators(t, runner)
	withImplementOwnerIdentity(t, "phase-record-owner")
	// The Daemon clock stamps started and ended and renews the Run Budget
	// deadline, whose context runs on the wall clock, so it is fixed at the
	// present: a fixed past date expired the Run once the wall clock passed it
	// by the two-hour budget.
	daemonNow := time.Now().UTC().Truncate(time.Second)
	updateCommandDependenciesForTest(t, func(d *commandDependencies) {
		d.phaseNow = func() time.Time { return phasePromptNow }
		d.implementBudgetNow = func() time.Time { return daemonNow }
	})
	var out, diagnostics bytes.Buffer
	code := runCLIContext(t, t.Context(), []string{"implement", "--spec", implementTestSlug, "--no-input"}, &out, &diagnostics)
	if code != exitOK {
		t.Fatalf("exit = %d, stdout=%s stderr=%s", code, &out, &diagnostics)
	}
	rows := promptRecordRows(t, home)
	if len(rows) != 2 || len(runner.requests) != 2 {
		t.Fatalf("rows=%+v requests=%d", rows, len(runner.requests))
	}
	runID := implementRunIDFromStderr(t, diagnostics.String())
	for i, row := range rows {
		assertPromptRow(t, row, runner.requests[i], runID, implementTestSlug, "task", runner.taskIDs[i], repo, "implement-adapter-session")
		if row["started"] != daemonNow.Format(time.RFC3339) || row["ended"] != daemonNow.Format(time.RFC3339) {
			t.Fatalf("unexpected times: %+v", row)
		}
	}
	run := implementRunFromStore(t, home, runID)
	if run.State != store.StateClean {
		t.Fatalf("Run = %+v", run)
	}
}

func TestTheReviewRecordsEachPromptInThePhaseRecord(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"one_prompt", "one_spec_named", "several_specs_named", "unwritable_record", "segmented", "failed_before_prompt", "retry_before_prompt", "plain_runner", "plain_runner_pre_prompt"} {
		t.Run(name, func(t *testing.T) {
			answer := agent.ExecuteResult{Message: "No findings.", StopReason: "end_turn", ACPSessionID: "review-adapter-session"}
			runner := &segmentReviewRunner{reviewCommandRunner: reviewCommandRunner{results: []reviewCommandRunResult{{result: answer}, {result: answer}, {result: answer}}}}
			fixture := newReviewCommandFixture(t, "codex", runner)
			updateCommandDependenciesForTest(t, func(d *commandDependencies) { d.phaseNow = func() time.Time { return phasePromptNow } })
			wantSpec := ""
			switch name {
			case "one_spec_named", "several_specs_named":
				slugs := []string{"0397-first"}
				if name == "several_specs_named" {
					slugs = append(slugs, "0398-second")
				} else {
					wantSpec = slugs[0]
				}
				for _, slug := range slugs {
					commitScopeFile(t, fixture, "docs/specs/"+slug+"/_prd.md", "# Prompt record\n\n## Decisions\n\nKeep the review read-only.\n")
					commitScopeFile(t, fixture, "docs/specs/"+slug+"/_techspec.md", "# Prompt record design\n")
				}
			case "segmented":
				commitReviewSegments(t, fixture, 3)
			case "failed_before_prompt":
				runner.prepareErrors = []error{errors.New("prepare failed")}
			case "retry_before_prompt":
				runner.results[0] = reviewCommandRunResult{err: retryProtocolFailure(false)}
			case "plain_runner":
				withAgentRunner(t, phasePlainReviewRunner{result: answer})
			case "plain_runner_pre_prompt":
				withAgentRunner(t, phasePlainReviewRunner{err: &agent.ModelNotAdvertisedError{Runtime: "codex", Model: "not-offered"}})
			}
			var baseline reviewRecord
			if name == "unwritable_record" {
				code, record, diagnostics := fixture.run(t)
				if code != exitOK || diagnostics != "" {
					t.Fatalf("baseline exit=%d stderr=%q", code, diagnostics)
				}
				baseline = record
				// Delete only disposable evidence so the same candidate opens a fresh session.
				recordDir := reviewCheckoutDir(fixture.artifactDir, fixture.repository)
				if err := os.RemoveAll(recordDir); err != nil {
					t.Fatal(err)
				}
				if err := os.RemoveAll(filepath.Join(fixture.homeDir, ".roundfix", "phase-record")); err != nil {
					t.Fatal(err)
				}
				mustWrite(t, filepath.Join(fixture.homeDir, ".roundfix", "phase-record"), "blocked")
			}
			code, record, diagnostics := fixture.run(t)
			if name == "failed_before_prompt" || name == "plain_runner_pre_prompt" {
				if code != exitPreflight || runner.preparedCalls != 0 {
					t.Fatalf("exit=%d runner=%+v stderr=%q", code, runner, diagnostics)
				}
				if _, err := os.Stat(filepath.Join(fixture.homeDir, ".roundfix", "phase-record")); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("record written before prompt: %v", err)
				}
				return
			}
			if code != exitOK || record.Outcome != reviewOutcomeReviewed {
				t.Fatalf("exit=%d record=%+v stderr=%q", code, record, diagnostics)
			}
			if name == "unwritable_record" {
				// The fresh invocation generates a new session name.
				baseline.Lineage.Session = record.Lineage.Session
				if !reflect.DeepEqual(record, baseline) {
					t.Fatalf("record changed: %+v; baseline %+v", record, baseline)
				}
				err := phaserecord.New(fixture.homeDir, func() time.Time { return phasePromptNow }).RecordPrompt(t.Context(), phaserecord.PromptEntry{})
				want := fmt.Sprintf("roundfix: warning: phase record not written for pre-pr-review %s: %v\n", record.Lineage.Session, err)
				if err == nil || diagnostics != want {
					t.Fatalf("warning=%q want=%q", diagnostics, want)
				}
				return
			}
			if name == "retry_before_prompt" {
				if runner.preparedCalls != 2 {
					t.Fatalf("prompt calls=%d", runner.preparedCalls)
				}
				assertReviewRetry(t, record, diagnostics, "session/set_model", "adapter refused the request", 0)
			} else if diagnostics != "" {
				t.Fatalf("stderr=%q", diagnostics)
			}
			rows := promptRecordRows(t, fixture.homeDir)
			wantCount := 1
			if name == "segmented" {
				wantCount = 3
			}
			if len(rows) != wantCount {
				t.Fatalf("rows=%+v", rows)
			}
			for i, row := range rows {
				req := runner.request
				if name == "plain_runner" {
					req = agent.ExecuteRequest{Runtime: agent.RuntimeSpec{ID: "codex"}, Session: agent.SessionRef{Name: record.Lineage.Session, WorkDir: fixture.repository}}
				} else {
					requestIndex := i
					if name == "retry_before_prompt" {
						requestIndex++
					}
					req = runner.requests[requestIndex]
				}
				assertPromptRow(t, row, req, "", wantSpec, "pre-pr-review", record.Lineage.Session, fixture.repository, "review-adapter-session")
				if row["started"] != phasePromptNow.Format(time.RFC3339) || row["ended"] != phasePromptNow.Format(time.RFC3339) {
					t.Fatalf("times=%+v", row)
				}
			}
			if name == "one_spec_named" && !reflect.DeepEqual(record.Specs, []string{wantSpec}) {
				t.Fatalf("Specs=%v", record.Specs)
			}
			if name == "several_specs_named" && len(record.Specs) != 2 {
				t.Fatalf("Specs=%v", record.Specs)
			}
		})
	}
}

type phasePlainReviewRunner struct {
	result agent.ExecuteResult
	err    error
}

func (r phasePlainReviewRunner) Run(context.Context, agent.ExecuteRequest, runevent.Sink) (agent.ExecuteResult, error) {
	return r.result, r.err
}
func (phasePlainReviewRunner) Probe(context.Context, agent.ProbeRequest) error { return nil }
func (phasePlainReviewRunner) EndSession(context.Context, agent.RuntimeSpec, agent.SessionRef) error {
	return nil
}

func seedPhaseDayFiles(t *testing.T, home string, names []string) string {
	t.Helper()
	dir := filepath.Join(home, ".roundfix", "phase-record")
	mustMkdir(t, dir)
	for _, name := range names {
		mustWrite(t, filepath.Join(dir, name), "entry\n")
	}
	return dir
}

func TestTheRunRetentionSweepRemovesPhaseRecordDayFiles(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"live", "dry_run", "not_a_date"} {
		t.Run(name, func(t *testing.T) {
			home, repo := withCLIWorkspace(t)
			withGCNow(t, phasePromptNow)
			names := []string{"2026-09-01.jsonl", "2026-10-16.jsonl", "2026-10-17.jsonl", "2026-11-16.jsonl"}
			if name == "not_a_date" {
				names = append(names, "notes.jsonl", "2026-09-31.jsonl")
			}
			dir := seedPhaseDayFiles(t, home, names)
			loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: home, WorkDir: repo})
			if err != nil {
				t.Fatal(err)
			}
			s, err := store.Open(t.Context(), home)
			if err != nil {
				t.Fatal(err)
			}
			defer func() {
				if err := s.Close(); err != nil {
					t.Error(err)
				}
			}()
			dry := name == "dry_run"
			report, err := sweepRunDatabase(commandContextForTest(t, t.Context()), s, loaded, runRetentionOptions{dryRun: dry})
			if err != nil || report.phaseRecordDayFiles != 2 {
				t.Fatalf("report=%+v err=%v", report, err)
			}
			var out bytes.Buffer
			printRunRetentionReport(&out, report, dry)
			action := "removed"
			if dry {
				action = "removable"
			}
			if !strings.Contains(out.String(), fmt.Sprintf("  Runs %s: 0\n  Phase Record day files %s: 2\n", action, action)) {
				t.Fatalf("report=%s", &out)
			}
			for i, file := range names {
				_, err := os.Stat(filepath.Join(dir, file))
				if i < 2 && !dry {
					if !errors.Is(err, os.ErrNotExist) {
						t.Fatalf("%s remains: %v", file, err)
					}
				} else if err != nil {
					t.Fatalf("%s removed: %v", file, err)
				}
			}
		})
	}
}

func TestGCCountsPhaseRecordDayFilesWithoutARunDatabase(t *testing.T) {
	t.Parallel()
	home, _ := withCLIWorkspace(t)
	withGCNow(t, phasePromptNow)
	names := []string{"2026-09-01.jsonl", "2026-10-16.jsonl", "2026-11-16.jsonl"}
	dir := seedPhaseDayFiles(t, home, names)
	var out, diagnostics bytes.Buffer
	code := runCLIContext(t, t.Context(), []string{"gc", "--dry-run"}, &out, &diagnostics)
	want := "GC dry-run\n" +
		"  Journal Retention: 336h\n  Cutoff: 2026-11-02T12:00:00Z\n" +
		"  Runs eligible: 0\n  Journal rows eligible: 0\n  Artifact bytes reclaimable: 0\n  Orphan artifact dirs eligible: 0\n" +
		"  Run Retention: 30 days\n  Run Retention cutoff: 2026-10-17T12:00:00Z\n" +
		"  Runs removable: 0\n  Phase Record day files removable: 2\n" +
		"  Runs kept past the cutoff: 0 (queue-referenced 0, worktree present 0, artifact root unproven 0)\n" +
		"  Rows removable: runs=0 run_events=0 run_agent_selections=0 run_token_usage=0 active_run_locks=0\n" +
		"  Database bytes reclaimable (estimated): 0\n"
	if code != exitOK || diagnostics.Len() != 0 || out.String() != want {
		t.Fatalf("exit=%d stderr=%q stdout=%q; want %q", code, diagnostics.String(), out.String(), want)
	}
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil || string(data) != "entry\n" {
			t.Fatalf("day file %s: %s, %v", name, data, err)
		}
	}
	if _, err := os.Stat(store.DatabasePath(home)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("GC created database: %v", err)
	}
	// The same no-database branch also handles live removal.
	out.Reset()
	diagnostics.Reset()
	code = runCLIContext(t, t.Context(), []string{"gc"}, &out, &diagnostics)
	if code != exitOK || diagnostics.Len() != 0 || !strings.Contains(out.String(), "  Runs removed: 0\n  Phase Record day files removed: 2\n") {
		t.Fatalf("live exit=%d stderr=%q stdout=%q", code, diagnostics.String(), out.String())
	}
	for i, name := range names {
		_, err := os.Stat(filepath.Join(dir, name))
		if i < 2 {
			if !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("old file remains: %s, %v", name, err)
			}
		} else if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAPhaseRecordRetentionFailureDoesNotKeepEligibleRuns(t *testing.T) {
	t.Parallel()
	f := seedGCRetention(t, false, 2)
	mustWrite(t, filepath.Join(f.home, ".roundfix", "phase-record"), "a file blocks the directory")
	loaded, err := roundconfig.Load(roundconfig.LoadOptions{HomeDir: f.home, WorkDir: commandEnvironmentForTest(t).workDir})
	if err != nil {
		t.Fatal(err)
	}
	s, err := store.Open(t.Context(), f.home)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := s.Close(); err != nil {
			t.Error(err)
		}
	}()
	report, err := sweepRunDatabase(commandContextForTest(t, t.Context()), s, loaded, runRetentionOptions{})
	if err == nil || !strings.Contains(err.Error(), "read phase record directory:") || !reflect.DeepEqual(report.runIDs, []string{f.removed.ID}) {
		t.Fatalf("report=%+v err=%v", report, err)
	}
	if run, found, err := s.Run(t.Context(), f.removed.ID); err != nil || found {
		t.Fatalf("eligible Run remains: %+v, %v", run, err)
	}
	for _, kept := range []store.Run{f.queued, f.worktree, f.recent} {
		if _, found, err := s.Run(t.Context(), kept.ID); err != nil || !found {
			t.Fatalf("protected Run %s removed: %v", kept.ID, err)
		}
	}
	// The rest of the Run sweep still runs: compaction, the sweep record and the final size.
	if sweep, found, err := s.LastRunRetentionSweep(t.Context()); err != nil || !found || sweep.RetentionDays != report.retentionDays {
		t.Fatalf("Run Retention sweep not recorded after a Phase Record failure: %+v found=%v err=%v", sweep, found, err)
	}
	if report.compaction == "not needed" || report.bytesAfter == 0 {
		t.Fatalf("Run Retention stopped at the Phase Record failure: %+v", report)
	}
}

func TestTheSealedConventionCheckLeavesNoPromptEntry(t *testing.T) {
	t.Parallel()
	fixture, runner := validatorFixture(t, "- "+validatorTaskPath+":4 Failure: invalid state", validatorAnswer(validatorVerdict{"F1", "stands", "", "a real state defect"}))
	runner.results[0].result.ACPSessionID = "review-adapter-session"
	updateCommandDependenciesForTest(t, func(d *commandDependencies) { d.phaseNow = func() time.Time { return phasePromptNow } })
	code, record, diagnostics := fixture.run(t)
	if code != exitRunFailed || record.Validation.Validator != "ran" || len(runner.sealedRequests) != 1 {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, diagnostics)
	}
	rows := promptRecordRows(t, fixture.homeDir)
	if len(rows) != 1 {
		t.Fatalf("sealed check recorded a prompt: %+v", rows)
	}
	assertPromptRow(t, rows[0], runner.request, "", "", "pre-pr-review", record.Lineage.Session, fixture.repository, "review-adapter-session")
}
