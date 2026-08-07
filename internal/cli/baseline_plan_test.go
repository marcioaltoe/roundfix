// Suite: Baseline plan preflight CLI
// Invariant: baseline plan exposes deterministic local preflight evidence and actionable exits without prompts or repository writes.
// Boundary IN: plan dispatch, repo/format flags, text and JSON rendering, diagnostics, and exit categories.
// Boundary OUT: preservation decisions, profile alignment, portable Plan Documents, and apply.

package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"roundfix/internal/baseline"
	"roundfix/internal/gittest"
)

func TestBaselinePlanPreflightJSONActionRequired(t *testing.T) {
	t.Parallel()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, "AGENTS.md", "root policy\n")
	commitBaselinePlanTestRepository(t, repo)
	writeBaselinePlanTestFile(t, repo, "scratch.txt", "unrelated dirty bytes\n")

	before := baselinePlanTestTree(t, repo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "plan", "--repo", repo, "--format", "json",
	}, &stdout, &stderr)
	if code != exitUnverified {
		t.Fatalf("baseline plan exit = %d, want %d stdout=%q stderr=%q", code, exitUnverified, stdout.String(), stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("action-required plan stderr = %q, want empty", stderr.String())
	}
	var result baseline.Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode baseline plan JSON: %v\n%s", err, stdout.String())
	}
	if result.SchemaVersion != baselineResultSchema || result.Operation != "plan" ||
		result.State != "action_required" || result.Category != "decision" ||
		result.NextAction == "" {
		t.Fatalf("baseline plan result = %+v", result)
	}
	if strings.Contains(stdout.String(), filepath.ToSlash(repo)) {
		t.Fatalf("portable preflight JSON contains absolute checkout path %q:\n%s", repo, stdout.String())
	}
	after := baselinePlanTestTree(t, repo)
	if before != after {
		t.Fatalf("baseline plan changed repository bytes:\nbefore=%s\nafter=%s", before, after)
	}
}

func TestBaselinePlanPreflightText(t *testing.T) {
	t.Parallel()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, "nested/CLAUDE.md", "nested policy\n")
	commitBaselinePlanTestRepository(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "plan", "--repo=" + repo, "--format=text",
	}, &stdout, &stderr)
	if code != exitUnverified || stderr.Len() != 0 {
		t.Fatalf("baseline plan text exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"Baseline plan: action required",
		"Category: decision",
		"Baseline Profile selection is required",
		"Next action:",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("baseline plan text missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), repo) {
		t.Fatalf("baseline plan text contains checkout path %q:\n%s", repo, stdout.String())
	}
}

func TestBaselinePlanPreflightBlocksUnsafeRepository(t *testing.T) {
	t.Parallel()
	repo := newBaselinePlanTestRepository(t)
	if err := os.Symlink("../outside.md", filepath.Join(repo, "AGENTS.md")); err != nil {
		t.Fatalf("create escaping alias: %v", err)
	}
	commitBaselinePlanTestRepository(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "plan", "--repo", repo, "--format=json",
	}, &stdout, &stderr)
	if code != exitPreflight {
		t.Fatalf("unsafe baseline plan exit = %d, want %d stdout=%q stderr=%q", code, exitPreflight, stdout.String(), stderr.String())
	}
	if !strings.Contains(stderr.String(), "baseline.inventory.unsafe-alias") {
		t.Fatalf("unsafe baseline plan diagnostic = %q", stderr.String())
	}
	var result baseline.Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode unsafe baseline plan JSON: %v\n%s", err, stdout.String())
	}
	if result.State != "action_required" || result.Category != "preflight" ||
		len(result.Warnings) == 0 || result.NextAction == "" {
		t.Fatalf("unsafe baseline plan result = %+v", result)
	}
}

func TestBaselinePlanPreflightRejectsUsageAndUncommittedRepository(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args func(*testing.T) []string
		want string
	}{
		{
			name: "unknown flag",
			args: func(t *testing.T) []string {
				return []string{"baseline", "plan", "--unknown"}
			},
			want: "invalid baseline plan arguments",
		},
		{
			name: "repository without commit",
			args: func(t *testing.T) []string {
				return []string{"baseline", "plan", "--repo", newBaselinePlanTestRepository(t)}
			},
			want: "requires at least one commit",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := RunContext(context.Background(), tt.args(t), &stdout, &stderr)
			if code != exitPreflight || stdout.Len() != 0 || !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("baseline plan exit = %d stdout=%q stderr=%q, want %q", code, stdout.String(), stderr.String(), tt.want)
			}
		})
	}
}

func TestBaselinePlanPreflightHelp(t *testing.T) {
	t.Parallel()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), []string{"baseline", "plan", "--help"}, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("baseline plan help exit = %d stdout=%q stderr=%q", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"roundfix baseline plan",
		"--repo",
		"--profile",
		"--profile-file",
		"--decision",
		"--decision-file",
		"--format",
		"never prompts",
		"Exit codes:",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("baseline plan help missing %q:\n%s", want, stdout.String())
		}
	}
	for _, forbidden := range []string{"--yes", "--interactive", "--no-input"} {
		if strings.Contains(stdout.String(), forbidden) {
			t.Fatalf("baseline plan help advertises interactive flag %q:\n%s", forbidden, stdout.String())
		}
	}
}

func TestBaselinePlanProfileFile(t *testing.T) {
	t.Parallel()
	repository, input, decisions := baselinePlanProfileFileFixture(t)
	draftPath := filepath.Join(t.TempDir(), "guided-backend.json")
	if err := os.WriteFile(draftPath, input.Document, 0o644); err != nil {
		t.Fatalf("write Profile draft input: %v", err)
	}
	args := baselinePlanProfileFileArgs(repository, draftPath, decisions)
	before := baselinePlanTestTree(t, repository)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := RunContext(context.Background(), args, &stdout, &stderr)
	if code != exitOK || stderr.Len() != 0 {
		t.Fatalf("Profile draft plan exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	automationPlan, err := baseline.ParsePlanDocument(stdout.Bytes())
	if err != nil {
		t.Fatalf("parse Profile draft Plan: %v\n%s", err, stdout.String())
	}
	direct, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{
		Repository:   repository,
		ProfileDraft: &input,
		Decisions:    decisions,
		Preservation: baseline.RootPreservationRequest{Mode: baseline.PreservationModeGreenfield},
	})
	if err != nil || direct.Plan == nil {
		t.Fatalf("build direct Profile draft Plan: outcome=%+v error=%v", direct, err)
	}
	if !reflect.DeepEqual(automationPlan.Profile, direct.Plan.Profile) ||
		!reflect.DeepEqual(automationPlan.Postimages, direct.Plan.Postimages) ||
		automationPlan.PlanDigest != direct.Plan.PlanDigest {
		t.Fatalf(
			"Profile draft normalization differs: cli=%s direct=%s",
			automationPlan.PlanDigest,
			direct.Plan.PlanDigest,
		)
	}
	if after := baselinePlanTestTree(t, repository); before != after {
		t.Fatalf("Profile draft planning changed repository bytes")
	}

	t.Run("mutually exclusive Profile inputs", func(t *testing.T) {
		var exclusiveOut bytes.Buffer
		var exclusiveErr bytes.Buffer
		exclusiveArgs := append([]string(nil), args...)
		exclusiveArgs = append(exclusiveArgs, "--profile", "standard-typescript-monorepo")
		exit := RunContext(context.Background(), exclusiveArgs, &exclusiveOut, &exclusiveErr)
		if exit != exitPreflight ||
			!strings.Contains(exclusiveErr.String(), "mutually exclusive") {
			t.Fatalf(
				"mutually exclusive inputs exit=%d stdout=%s stderr=%s",
				exit,
				exclusiveOut.String(),
				exclusiveErr.String(),
			)
		}
	})

	tests := []struct {
		name     string
		document []byte
		want     string
	}{
		{name: "invalid draft", document: []byte(`{"schemaVersion":`), want: "custom.profile.json"},
		{
			name:     "stale draft",
			document: staleProfileDraftDocument(t, input.Document),
			want:     "custom.profile.catalog-schema.invalid",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "draft.json")
			if err := os.WriteFile(path, test.document, 0o644); err != nil {
				t.Fatal(err)
			}
			treeBefore := baselinePlanTestTree(t, repository)
			var failureOut bytes.Buffer
			var failureErr bytes.Buffer
			exit := RunContext(
				context.Background(),
				baselinePlanProfileFileArgs(repository, path, decisions),
				&failureOut,
				&failureErr,
			)
			if exit != exitPreflight || !strings.Contains(failureErr.String(), test.want) {
				t.Fatalf(
					"%s exit=%d stdout=%s stderr=%s",
					test.name,
					exit,
					failureOut.String(),
					failureErr.String(),
				)
			}
			if treeAfter := baselinePlanTestTree(t, repository); treeBefore != treeAfter {
				t.Fatalf("%s changed repository bytes", test.name)
			}
		})
	}

	t.Run("output failure", func(t *testing.T) {
		treeBefore := baselinePlanTestTree(t, repository)
		var outputErr bytes.Buffer
		exit := RunContext(
			context.Background(),
			args,
			failingWriter{err: errors.New("injected Profile Plan output failure")},
			&outputErr,
		)
		if exit != exitRunFailed ||
			!strings.Contains(outputErr.String(), "injected Profile Plan output failure") {
			t.Fatalf("output failure exit=%d stderr=%s", exit, outputErr.String())
		}
		if treeAfter := baselinePlanTestTree(t, repository); treeBefore != treeAfter {
			t.Fatal("Profile Plan output failure changed repository bytes")
		}
	})
}

func TestBaselinePlanPreflightRealCLI(t *testing.T) {
	t.Parallel()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, "AGENTS.md", "real CLI policy\n")
	commitBaselinePlanTestRepository(t, repo)

	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("resolve project root: %v", err)
	}
	binary := filepath.Join(t.TempDir(), "roundfix")
	build := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
	build.Dir = projectRoot
	build.Env = os.Environ()
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build real roundfix CLI: %v\n%s", err, output)
	}

	command := exec.Command(binary, "baseline", "plan", "--repo", repo, "--format=json")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err = command.Run()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != exitUnverified {
		t.Fatalf("real CLI error = %v stdout=%q stderr=%q, want exit %d", err, stdout.String(), stderr.String(), exitUnverified)
	}
	if stderr.Len() != 0 {
		t.Fatalf("real CLI stderr = %q, want empty", stderr.String())
	}
	var result baseline.Result
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("decode real CLI JSON: %v\n%s", err, stdout.String())
	}
	if result.State != "action_required" || result.Category != "decision" {
		t.Fatalf("real CLI result = %+v", result)
	}
}

func TestBaselinePlanCommandEmitsPortableJSONAndNormalizesDecisionFiles(t *testing.T) {
	t.Parallel()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeBaselinePlanTestFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeBaselinePlanTestFile(t, repo, "Makefile", "verify:\n\t@true\n")
	commitBaselinePlanTestRepository(t, repo)

	inline := []string{
		"baseline", "plan", "--repo", repo, "--profile", "go-cli-tui",
		"--decision", "preservation.mode=greenfield",
		"--decision", "language.generated=English",
		"--decision", "verification.gate=make verify",
		"--decision", "spec.scaffold=true",
		"--decision", "domain.layout=single-context",
		"--decision", "triage.external=false",
		"--decision", "autonomous.enabled=true",
		"--decision", "runtime.backend=codex gpt-5.5 xhigh",
		"--decision", "runtime.design=claude opus xhigh",
		"--decision", "secondbrain.enabled=false",
		"--decision", "repository.extension.enabled=false",
		"--format", "json",
	}
	before := baselinePlanTestTree(t, repo)
	var inlineOut bytes.Buffer
	var inlineErr bytes.Buffer
	if code := RunContext(context.Background(), inline, &inlineOut, &inlineErr); code != exitOK {
		t.Fatalf("inline plan exit = %d stdout=%s stderr=%s", code, inlineOut.String(), inlineErr.String())
	}
	inlinePlan, err := baseline.ParsePlanDocument(inlineOut.Bytes())
	if err != nil {
		t.Fatalf("parse inline plan: %v\n%s", err, inlineOut.String())
	}
	if strings.Contains(inlineOut.String(), filepath.ToSlash(repo)) {
		t.Fatalf("plan contains checkout path %q", repo)
	}
	if after := baselinePlanTestTree(t, repo); after != before {
		t.Fatalf("planning mutated repository: before=%s after=%s", before, after)
	}

	document := baseline.DecisionDocument{
		SchemaVersion: baseline.DecisionDocumentSchemaVersion,
		Version:       baseline.DecisionDocumentVersion,
		Decisions: append(
			[]baseline.DecisionValue{{ID: "preservation.mode", Value: "greenfield"}},
			inlinePlan.Decisions...,
		),
	}
	decisionBytes, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	decisionPath := filepath.Join(t.TempDir(), "decisions.json")
	if err := os.WriteFile(decisionPath, append(decisionBytes, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	var fileOut bytes.Buffer
	var fileErr bytes.Buffer
	code := RunContext(context.Background(), []string{
		"baseline", "plan", "--repo", repo, "--profile", "go-cli-tui",
		"--decision-file", decisionPath, "--format=json",
	}, &fileOut, &fileErr)
	if code != exitOK {
		t.Fatalf("file plan exit = %d stdout=%s stderr=%s", code, fileOut.String(), fileErr.String())
	}
	filePlan, err := baseline.ParsePlanDocument(fileOut.Bytes())
	if err != nil {
		t.Fatalf("parse file plan: %v", err)
	}
	if filePlan.PlanDigest != inlinePlan.PlanDigest {
		t.Fatalf("decision source changed digest: inline=%s file=%s",
			inlinePlan.PlanDigest, filePlan.PlanDigest)
	}
}

func TestBaselinePlanAdoptionAndDecisionCharacterizationCorpus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		args                     func(*testing.T, string) []string
		wantExit                 int
		wantPlan                 bool
		wantCategory             string
		wantRetentionDisposition string
		assertResult             func(*testing.T, baseline.Result)
	}{
		{
			name: "first-adoption-greenfield-with-inline-decisions-emits-plan",
			args: func(_ *testing.T, repo string) []string {
				return baselinePlanCharacterizationArgs(repo, "greenfield")
			},
			wantExit: exitOK,
			wantPlan: true,
		},
		{
			name: "first-adoption-preservation-with-complete-strict-decision-document-emits-plan",
			args: func(t *testing.T, repo string) []string {
				inspection, err := baseline.InspectRepository(context.Background(), repo, nil)
				if err != nil {
					t.Fatalf("inspect first-adoption repository: %v", err)
				}
				preservation, err := baseline.PlanRootPreservation(
					inspection,
					baseline.RootPreservationRequest{Mode: baseline.PreservationModePreservation},
				)
				if err != nil {
					t.Fatalf("build preservation Decision Document: %v", err)
				}
				if preservation.DecisionSkeleton == nil {
					t.Fatalf("preservation result has no strict Decision Document: %+v", preservation)
				}
				document := preservation.DecisionSkeleton.Document
				document.Decisions = append(
					[]baseline.DecisionValue{{ID: "preservation.mode", Value: "preservation"}},
					baselinePlanCharacterizationDecisions()...,
				)
				data, err := json.MarshalIndent(document, "", "  ")
				if err != nil {
					t.Fatalf("marshal strict Decision Document: %v", err)
				}
				decisionPath := filepath.Join(t.TempDir(), "decisions.json")
				if err := os.WriteFile(decisionPath, append(data, '\n'), 0o644); err != nil {
					t.Fatalf("write strict Decision Document: %v", err)
				}
				return []string{
					"baseline", "plan", "--repo", repo, "--profile", "go-cli-tui",
					"--decision-file", decisionPath, "--format=json",
				}
			},
			wantExit:                 exitOK,
			wantPlan:                 true,
			wantRetentionDisposition: "repository-rules",
		},
		{
			name: "decisions-absent-names-every-required-decision",
			args: func(_ *testing.T, repo string) []string {
				return []string{
					"baseline", "plan", "--repo", repo, "--profile", "go-cli-tui", "--format=json",
				}
			},
			wantExit:     exitUnverified,
			wantCategory: "decision",
			assertResult: func(t *testing.T, result baseline.Result) {
				catalog, err := baseline.LoadEmbeddedCatalog()
				if err != nil {
					t.Fatalf("load embedded catalog: %v", err)
				}
				profile, err := baseline.ResolveProfile("", "go-cli-tui", catalog)
				if err != nil {
					t.Fatalf("resolve characterization Profile: %v", err)
				}
				for _, decisionID := range profile.Decisions {
					if !strings.Contains(result.Message, decisionID) {
						t.Errorf("missing-decision result omitted %q", decisionID)
					}
				}
			},
		},
		{
			name: "decisions-supplied-without-preservation-mode-requires-action",
			args: func(_ *testing.T, repo string) []string {
				return baselinePlanCharacterizationArgs(repo, "")
			},
			wantExit:     exitUnverified,
			wantCategory: "decision",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repo := newHumanBaselineRepository(t)
			writeBaselinePlanTestFile(t, repo, "AGENTS.md", "retain this repository rule\n")
			commitBaselinePlanTestRepository(t, repo)
			if _, err := os.Stat(filepath.Join(repo, "docs", "agents", "setup-context.json")); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("first-adoption fixture unexpectedly has a Setup Manifest: %v", err)
			}
			before := baselinePlanTestTree(t, repo)
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			code := RunContext(context.Background(), test.args(t, repo), &stdout, &stderr)
			if code != test.wantExit || stderr.Len() != 0 {
				t.Fatalf(
					"exit identity = %d, want %d stdout=%s stderr=%s",
					code,
					test.wantExit,
					stdout.String(),
					stderr.String(),
				)
			}
			if after := baselinePlanTestTree(t, repo); after != before {
				t.Fatalf("non-interactive characterization mutated repository bytes")
			}

			if test.wantPlan {
				plan, err := baseline.ParsePlanDocument(stdout.Bytes())
				if err != nil {
					t.Fatalf("parse characterization Plan: %v", err)
				}
				if plan.SchemaVersion != baseline.PlanSchemaVersion || plan.Profile.ID != "go-cli-tui" {
					t.Fatalf(
						"plan identity = schema %q profile %q",
						plan.SchemaVersion,
						plan.Profile.ID,
					)
				}
				if test.wantRetentionDisposition == "" {
					if len(plan.Retention) != 0 {
						t.Fatalf("initial greenfield retention = %+v, want empty", plan.Retention)
					}
				} else if len(plan.Retention) != 1 ||
					plan.Retention[0].Disposition != test.wantRetentionDisposition {
					t.Fatalf(
						"initial preservation clause disposition = %+v, want %q",
						plan.Retention,
						test.wantRetentionDisposition,
					)
				}
				return
			}

			var result baseline.Result
			if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
				t.Fatalf("decode characterization Result: %v", err)
			}
			if result.SchemaVersion != baseline.ResultSchemaVersion ||
				result.Operation != "plan" ||
				result.State != "action_required" ||
				result.Category != test.wantCategory {
				t.Fatalf(
					"result identity = schema %q operation %q state %q category %q",
					result.SchemaVersion,
					result.Operation,
					result.State,
					result.Category,
				)
			}
			if test.assertResult != nil {
				test.assertResult(t, result)
			}
		})
	}
}

func baselinePlanCharacterizationArgs(repo, preservationMode string) []string {
	args := []string{"baseline", "plan", "--repo", repo, "--profile", "go-cli-tui"}
	for _, decision := range baselinePlanCharacterizationDecisions() {
		args = append(args, "--decision", fmt.Sprintf("%s=%v", decision.ID, decision.Value))
	}
	if preservationMode != "" {
		args = append(args, "--decision", "preservation.mode="+preservationMode)
	}
	return append(args, "--format=json")
}

func baselinePlanCharacterizationDecisions() []baseline.DecisionValue {
	return []baseline.DecisionValue{
		{ID: "language.generated", Value: "English"},
		{ID: "verification.gate", Value: "make verify"},
		{ID: "spec.scaffold", Value: true},
		{ID: "domain.layout", Value: "single-context"},
		{ID: "triage.external", Value: false},
		{ID: "autonomous.enabled", Value: true},
		{ID: "runtime.backend", Value: "codex gpt-5.5 xhigh"},
		{ID: "runtime.design", Value: "claude opus xhigh"},
		{ID: "secondbrain.enabled", Value: false},
		{ID: "repository.extension.enabled", Value: true},
	}
}

func newBaselinePlanTestRepository(t *testing.T) string {
	t.Helper()
	repo := t.TempDir()
	runBaselinePlanTestCommand(t, repo, "git", "init", "--quiet")
	gittest.Harden(t, repo)
	gittest.AppendConfig(t, repo, "[user]\n\tname = Roundfix Test\n\temail = roundfix@example.test\n[commit]\n\tgpgsign = false\n")
	return repo
}

func commitBaselinePlanTestRepository(t *testing.T, repo string) {
	t.Helper()
	runBaselinePlanTestCommand(t, repo, "git", "add", "--all")
	runBaselinePlanTestCommand(t, repo, "git", "commit", "--quiet", "-m", "seed")
}

func writeBaselinePlanTestFile(t *testing.T, repo, relative, content string) {
	t.Helper()
	target := filepath.Join(repo, filepath.FromSlash(relative))
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatalf("create %s parent: %v", relative, err)
	}
	if err := os.WriteFile(target, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", relative, err)
	}
}

func baselinePlanProfileFileFixture(
	t *testing.T,
) (string, baseline.ProfileDraftInput, []baseline.DecisionValue) {
	t.Helper()
	repository := newHumanBaselineRepository(t)
	writeBaselinePlanTestFile(
		t,
		repository,
		"package.json",
		`{"name":"root","packageManager":"bun@1.3.0","scripts":{"verify":"true"},"dependencies":{"hono":"latest","typescript":"latest","zod":"latest"}}`,
	)
	writeBaselinePlanTestFile(t, repository, "packages/backend/package.json", `{"name":"backend"}`)
	commitBaselinePlanTestRepository(t, repository)

	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	source, err := baseline.ResolveProfile("", "standard-typescript-monorepo", catalog)
	if err != nil {
		t.Fatal(err)
	}
	decisions := []baseline.DecisionValue{
		{ID: "language.generated", Value: "English"},
		{ID: "verification.gate", Value: "make verify"},
		{ID: "identifier.strategy", Value: map[string]any{"kind": "uuid-v7"}},
		{ID: "http.contract", Value: map[string]any{"mode": "Post-only"}},
		{
			ID: "auth.provider",
			Value: map[string]any{
				"kind": "better-auth",
				"routeException": map[string]any{
					"scope":   "/api/auth/*",
					"methods": []any{"GET", "POST"},
					"owner":   "Better Auth",
					"reason":  "Provider protocol routes require GET and POST semantics.",
				},
			},
		},
		{ID: "spec.scaffold", Value: true},
		{ID: "domain.layout", Value: "single-context"},
		{ID: "triage.external", Value: false},
		{ID: "autonomous.enabled", Value: false},
		{ID: "secondbrain.enabled", Value: false},
		{ID: "repository.extension.enabled", Value: false},
	}
	alignment, err := baseline.ResolveProfileAlignment(
		context.Background(),
		repository,
		baseline.ProfileAlignmentRequest{
			ProfileID:            source.ID,
			Decisions:            decisions,
			RemediationProfileID: source.ID,
		},
		catalog,
	)
	if err != nil {
		t.Fatal(err)
	}
	sourceCapabilities := make(map[string]struct{}, len(source.Capabilities))
	for _, capabilityID := range source.Capabilities {
		sourceCapabilities[capabilityID] = struct{}{}
	}
	removedCapabilities := make([]string, 0)
	for _, divergence := range alignment.Divergences {
		if !divergence.Blocking {
			continue
		}
		if _, profileSpecific := sourceCapabilities[divergence.ID]; profileSpecific {
			removedCapabilities = append(removedCapabilities, divergence.ID)
		}
	}
	input, err := baseline.NewProfileAdaptationDraft(
		source.ID,
		"guided-backend",
		[]string{"frontend", "autonomous-work"},
		removedCapabilities,
		catalog,
	)
	if err != nil {
		t.Fatal(err)
	}
	resolved, _, err := baseline.ResolveProfileDraft(repository, input, catalog)
	if err != nil {
		t.Fatal(err)
	}
	decisions = decisionsSelectedByProfile(decisions, resolved)
	return repository, input, decisions
}

func baselinePlanProfileFileArgs(
	repository string,
	draftPath string,
	decisions []baseline.DecisionValue,
) []string {
	args := []string{
		"baseline",
		"plan",
		"--repo",
		repository,
		"--profile-file",
		draftPath,
		"--decision",
		"preservation.mode=greenfield",
		"--format=json",
	}
	for _, decision := range decisions {
		value := ""
		switch typed := decision.Value.(type) {
		case string:
			value = typed
		default:
			data, _ := json.Marshal(typed)
			value = string(data)
		}
		args = append(args, "--decision", decision.ID+"="+value)
	}
	return args
}

func staleProfileDraftDocument(t *testing.T, document []byte) []byte {
	t.Helper()
	var raw map[string]any
	if err := json.Unmarshal(document, &raw); err != nil {
		t.Fatal(err)
	}
	raw["catalogSchema"] = "roundfix/baseline-catalog/stale"
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(data, '\n')
}

func runBaselinePlanTestCommand(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	command := exec.Command(name, args...)
	command.Dir = dir
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("run %s %s: %v\n%s", name, strings.Join(args, " "), err, output)
	}
}

func baselinePlanTestTree(t *testing.T, root string) string {
	t.Helper()
	var paths []string
	if err := filepath.WalkDir(root, func(filePath string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, filePath)
		if err != nil {
			return err
		}
		if relative == ".git" {
			return filepath.SkipDir
		}
		if relative != "." {
			paths = append(paths, filepath.ToSlash(relative))
		}
		return nil
	}); err != nil {
		t.Fatalf("list repository tree: %v", err)
	}
	sort.Strings(paths)
	digest := sha256.New()
	for _, relative := range paths {
		filePath := filepath.Join(root, filepath.FromSlash(relative))
		info, err := os.Lstat(filePath)
		if err != nil {
			t.Fatalf("inspect %s: %v", relative, err)
		}
		fmt.Fprintf(digest, "%d:%s:%d:", len(relative), relative, uint32(info.Mode()))
		switch {
		case info.Mode()&fs.ModeSymlink != 0:
			target, err := os.Readlink(filePath)
			if err != nil {
				t.Fatalf("read link %s: %v", relative, err)
			}
			fmt.Fprintf(digest, "%d:%s", len(target), target)
		case info.Mode().IsRegular():
			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("read %s: %v", relative, err)
			}
			fmt.Fprintf(digest, "%d:", len(data))
			if _, err := digest.Write(data); err != nil {
				t.Fatalf("hash %s: %v", relative, err)
			}
		}
	}
	return fmt.Sprintf("%x", digest.Sum(nil))
}
