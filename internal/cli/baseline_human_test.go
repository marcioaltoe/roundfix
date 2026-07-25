// Suite: Human Baseline workflow
// Invariant: one linear terminal workflow produces the automation-equivalent Plan and mutates only after exact final approval.
// Boundary IN: adoption/update state detection, numbered prompts, consolidated review, no-TTY refusal, and exact apply wiring.
// Boundary OUT: catalog semantics, repository inventory internals, and transaction phase failure injection.

package cli

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestHumanBaselineAdoption(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	before := baselinePlanTestTree(t, repo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader(humanBaselineAdoptionAnswers("1")),
			interactive: true,
		},
	)
	if code != exitOK {
		t.Fatalf("human adoption exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "Baseline apply: verified") ||
		!strings.Contains(stdout.String(), "Consolidated Change Plan review") {
		t.Fatalf("human adoption output does not show review and verified apply:\n%s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Final confirmation for Plan Digest sha256:") {
		t.Fatalf("human adoption did not bind confirmation to digest:\n%s", stderr.String())
	}
	manifest := filepath.Join(repo, filepath.FromSlash(baselineSetupManifestPath))
	if _, err := os.Stat(manifest); err != nil {
		t.Fatalf("human adoption did not write verified Setup Manifest: %v", err)
	}
	after := baselinePlanTestTree(t, repo)
	if before == after {
		t.Fatal("approved human adoption did not change repository bytes")
	}
}

func TestHumanBaselineUpdate(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	applyHumanBaselineFixturePlan(t, repo)
	before := baselinePlanTestTree(t, repo)

	answers := []string{"1", "1"}
	for range 10 {
		answers = append(answers, "1")
	}
	answers = append(answers, "2")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader(strings.Join(answers, "\n") + "\n"),
			interactive: true,
		},
	)
	if code != exitUnverified {
		t.Fatalf("human update decline exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"Baseline workflow: update",
		"Current Baseline Profile: go-cli-tui",
		"Baseline Plan was declined; no repository bytes were written",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("human update output missing %q:\n%s", want, stdout.String())
		}
	}
	if !strings.Contains(stderr.String(), "Change Baseline Profile") {
		t.Fatalf("human update did not offer explicit profile change:\n%s", stderr.String())
	}
	after := baselinePlanTestTree(t, repo)
	if before != after {
		t.Fatalf("declined human update changed repository bytes:\nbefore=%s\nafter=%s", before, after)
	}

	changeAnswers := []string{"1", "2", "2"}
	for range 10 {
		changeAnswers = append(changeAnswers, "1")
	}
	var changeReview bytes.Buffer
	var changePrompts bytes.Buffer
	humanPlan, err := driveHumanBaselinePlan(
		context.Background(),
		repo,
		&baselineHumanPrompt{
			reader: bufioReader(strings.Join(changeAnswers, "\n") + "\n"),
			writer: &changePrompts,
		},
		&changeReview,
	)
	if err != nil {
		t.Fatalf("build human profile-change Plan: %v\nreview=%s\nprompts=%s", err, changeReview.String(), changePrompts.String())
	}
	automation, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{
		Repository: repo,
		ProfileID:  "rust-cli",
		Decisions:  humanBaselineFixtureDecisions(),
		Preservation: baseline.RootPreservationRequest{
			Mode: baseline.PreservationModeGreenfield,
		},
	})
	if err != nil || automation.Plan == nil {
		t.Fatalf("build automation profile-change Plan: outcome=%+v error=%v", automation, err)
	}
	humanBytes, err := baseline.MarshalPlanDocument(humanPlan)
	if err != nil {
		t.Fatal(err)
	}
	automationBytes, err := baseline.MarshalPlanDocument(*automation.Plan)
	if err != nil {
		t.Fatal(err)
	}
	if humanPlan.Profile.ID != "rust-cli" || !bytes.Equal(humanBytes, automationBytes) {
		t.Fatalf(
			"profile-change parity failed: profile=%s human=%s automation=%s",
			humanPlan.Profile.ID,
			humanPlan.PlanDigest,
			automation.Plan.PlanDigest,
		)
	}
}

func TestHumanBaselineIncompatibleManifestKeepsValidDefaults(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	applyHumanBaselineFixturePlan(t, repo)
	manifestPath := filepath.Join(repo, filepath.FromSlash(baselineSetupManifestPath))
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read Setup Manifest: %v", err)
	}
	var manifest baseline.SetupManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Fatalf("decode Setup Manifest: %v", err)
	}
	manifest.ProfileDigest = "sha256:" + strings.Repeat("0", 64)
	data, err = json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		t.Fatalf("encode incompatible Setup Manifest: %v", err)
	}
	if err := os.WriteFile(manifestPath, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write incompatible Setup Manifest: %v", err)
	}

	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load Baseline catalog: %v", err)
	}
	state, err := inspectBaselineHumanState(repo, catalog)
	if err != nil {
		t.Fatalf("inspect incompatible human state: %v", err)
	}
	if state.mode != "adoption" || state.incompatible == "" {
		t.Fatalf("incompatible state = %+v, want adoption with diagnostic", state)
	}
	if state.currentProfile == nil || state.currentProfile.ID != "go-cli-tui" {
		t.Fatalf("recovered profile = %+v, want go-cli-tui", state.currentProfile)
	}
	if got := state.currentDecisions["verification.gate"]; got != "make verify" {
		t.Fatalf("recovered verification.gate = %#v, want make verify", got)
	}

	var review bytes.Buffer
	var prompts bytes.Buffer
	profile, err := promptBaselineProfile(
		context.Background(),
		&baselineHumanPrompt{reader: bufioReader("\n"), writer: &prompts},
		&review,
		repo,
		catalog,
		state,
	)
	if err != nil {
		t.Fatalf("accept recovered profile default: %v", err)
	}
	if profile.ID != "go-cli-tui" {
		t.Fatalf("accepted profile = %q, want go-cli-tui", profile.ID)
	}
	if !strings.Contains(prompts.String(), "Reuse existing profile go-cli-tui (default)") {
		t.Fatalf("recovered profile default is not visible:\n%s", prompts.String())
	}
}

func TestHumanBaselineDecisionDefaults(t *testing.T) {
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load Baseline catalog: %v", err)
	}
	tests := []struct {
		name    string
		id      string
		current map[string]any
		want    any
	}{
		{
			name:    "existing manifest value wins",
			id:      "verification.gate",
			current: map[string]any{"verification.gate": "repository verify"},
			want:    "repository verify",
		},
		{
			name:    "invalid manifest value falls back to catalog",
			id:      "verification.gate",
			current: map[string]any{"verification.gate": "  "},
			want:    "rtk make verify",
		},
		{name: "language", id: "language.generated", want: "English"},
		{name: "verification", id: "verification.gate", want: "rtk make verify"},
		{name: "HTTP contract", id: "http.contract", want: map[string]any{"mode": "Post-only"}},
		{name: "spec scaffold", id: "spec.scaffold", want: true},
		{name: "domain layout", id: "domain.layout", want: "single-context"},
		{name: "external triage", id: "triage.external", want: false},
		{name: "autonomous work", id: "autonomous.enabled", want: true},
		{name: "Secondbrain", id: "secondbrain.enabled", want: true},
		{name: "repository extension", id: "repository.extension.enabled", want: true},
		{name: "backend runtime", id: "runtime.backend", want: "codex gpt-5.6-sol"},
		{name: "design runtime", id: "runtime.design", want: "claude fable high"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			got, err := promptBaselineDecision(
				context.Background(),
				&baselineHumanPrompt{reader: bufioReader("\n"), writer: &output},
				catalog,
				test.id,
				test.current,
			)
			if err != nil {
				t.Fatalf("accept %s default: %v", test.id, err)
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Fatalf("%s default = %#v, want %#v", test.id, got, test.want)
			}
			if !strings.Contains(output.String(), "default") {
				t.Fatalf("%s prompt does not expose its default:\n%s", test.id, output.String())
			}
		})
	}
}

func TestProjectDecisionPrompts(t *testing.T) {
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatalf("load Baseline catalog: %v", err)
	}

	t.Run("UUID version 7 suggestion is visible", func(t *testing.T) {
		var output bytes.Buffer
		got, err := promptBaselineDecision(
			context.Background(),
			&baselineHumanPrompt{reader: bufioReader("\n"), writer: &output},
			catalog,
			"identifier.strategy",
			nil,
		)
		if err != nil {
			t.Fatalf("accept identifier suggestion: %v", err)
		}
		want := map[string]any{"kind": "uuid-v7"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("identifier suggestion = %#v, want %#v", got, want)
		}
		if !strings.Contains(output.String(), "UUID version 7") ||
			!strings.Contains(output.String(), `"kind":"uuid-v7"`) {
			t.Fatalf("identifier suggestion is not visible:\n%s", output.String())
		}
	})

	t.Run("complete Better Auth suggestion is visible", func(t *testing.T) {
		var output bytes.Buffer
		got, err := promptBaselineDecision(
			context.Background(),
			&baselineHumanPrompt{reader: bufioReader("\n"), writer: &output},
			catalog,
			"auth.provider",
			nil,
		)
		if err != nil {
			t.Fatalf("accept Better Auth suggestion: %v", err)
		}
		encoded, err := json.Marshal(got)
		if err != nil {
			t.Fatalf("encode Better Auth suggestion: %v", err)
		}
		for _, fragment := range []string{
			`"kind":"better-auth"`,
			`"scope":"/api/auth/*"`,
			`"methods":["GET","POST"]`,
			`"owner":"Better Auth"`,
			`"reason":`,
		} {
			if !strings.Contains(output.String(), fragment) ||
				!strings.Contains(string(encoded), fragment) {
				t.Fatalf("Better Auth suggestion is incomplete for %q:\n%s", fragment, output.String())
			}
		}
	})

	t.Run("compatible stored object is reused", func(t *testing.T) {
		current := map[string]any{
			"identifier.strategy": map[string]any{
				"kind":     "repository-defined",
				"guidance": "Use immutable aggregate sequence identifiers.",
			},
		}
		var output bytes.Buffer
		got, err := promptBaselineDecision(
			context.Background(),
			&baselineHumanPrompt{reader: bufioReader("\n"), writer: &output},
			catalog,
			"identifier.strategy",
			current,
		)
		if err != nil {
			t.Fatalf("reuse stored identifier strategy: %v", err)
		}
		if !reflect.DeepEqual(got, current["identifier.strategy"]) {
			t.Fatalf("reused identifier = %#v, want %#v", got, current["identifier.strategy"])
		}
		if !strings.Contains(output.String(), "Keep identifier.strategy=") {
			t.Fatalf("stored identifier keep-or-change prompt is missing:\n%s", output.String())
		}
	})

	t.Run("invalid stored object is not reused", func(t *testing.T) {
		current := map[string]any{
			"identifier.strategy": map[string]any{
				"kind":    "uuid-v7",
				"unknown": true,
			},
		}
		var output bytes.Buffer
		got, err := promptBaselineDecision(
			context.Background(),
			&baselineHumanPrompt{reader: bufioReader("\n"), writer: &output},
			catalog,
			"identifier.strategy",
			current,
		)
		if err != nil {
			t.Fatalf("replace invalid stored identifier strategy: %v", err)
		}
		if !reflect.DeepEqual(got, map[string]any{"kind": "uuid-v7"}) {
			t.Fatalf("replacement identifier = %#v", got)
		}
		if strings.Contains(output.String(), "Keep identifier.strategy=") {
			t.Fatalf("invalid stored identifier was offered for reuse:\n%s", output.String())
		}
	})
}

func TestHumanBaselinePreservationDefaultFollowsInstructionInventory(t *testing.T) {
	tests := []struct {
		name            string
		hasInstructions bool
		want            baseline.PreservationMode
	}{
		{name: "instructions use preservation", hasInstructions: true, want: baseline.PreservationModePreservation},
		{name: "no instructions use greenfield", hasInstructions: false, want: baseline.PreservationModeGreenfield},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var output bytes.Buffer
			mode, err := promptPreservationMode(
				context.Background(),
				&baselineHumanPrompt{reader: bufioReader("\n"), writer: &output},
				test.hasInstructions,
			)
			if err != nil {
				t.Fatalf("accept preservation default: %v", err)
			}
			if mode != test.want {
				t.Fatalf("preservation default = %q, want %q", mode, test.want)
			}
			if !strings.Contains(output.String(), "(default)") {
				t.Fatalf("preservation prompt does not expose its default:\n%s", output.String())
			}
		})
	}
}

func TestHumanBaselineConfirmationRequiresExplicitChoice(t *testing.T) {
	var output bytes.Buffer
	selected, err := (&baselineHumanPrompt{
		reader: bufioReader("\n2\n"),
		writer: &output,
	}).selectOne(context.Background(), "Apply reviewed Plan", []string{"Apply", "Decline"})
	if err != nil {
		t.Fatalf("read explicit confirmation: %v", err)
	}
	if selected != 1 {
		t.Fatalf("confirmation selection = %d, want explicit decline", selected)
	}
	if !strings.Contains(output.String(), "Enter a number from 1 to 2") {
		t.Fatalf("blank confirmation was accepted without an explicit choice:\n%s", output.String())
	}
}

func TestConsolidatedReview(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	writeBaselinePlanTestFile(t, repo, "AGENTS.md", "Preserve this repository-specific rule.\n")
	commitBaselinePlanTestRepository(t, repo)
	before := baselinePlanTestTree(t, repo)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader(humanBaselinePreservationAnswers()),
			interactive: true,
		},
	)
	if code != exitUnverified {
		t.Fatalf("preservation review exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	for _, want := range []string{
		"Consolidated editable classification review:",
		"normative-clause -> repository-rules",
		"Consolidated Change Plan review",
		"Complete managed-entry ledger:",
		"Complete Upgrade Retention Contract ledger:",
		"Plan Digest: sha256:",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("consolidated review missing %q:\n%s", want, stdout.String())
		}
	}
	fileChanges := strings.Index(stdout.String(), "File changes:")
	managedLedger := strings.Index(stdout.String(), "Complete managed-entry ledger:")
	retentionLedger := strings.Index(stdout.String(), "Complete Upgrade Retention Contract ledger:")
	if fileChanges < 0 || managedLedger < fileChanges || retentionLedger < managedLedger {
		t.Fatalf("consolidated review order is not fileChanges then ledgers:\n%s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "Edit classifications and dispositions") {
		t.Fatalf("classification proposal was not editable:\n%s", stderr.String())
	}
	after := baselinePlanTestTree(t, repo)
	if before != after {
		t.Fatalf("declined preservation review changed repository bytes:\nbefore=%s\nafter=%s", before, after)
	}
}

func TestConsolidatedReviewEditsManagedClassification(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	writeBaselinePlanTestFile(
		t,
		repo,
		"AGENTS.md",
		"<!-- setup-context-driven:begin id=legacy.rule version=1 -->\nPreserve me.\n<!-- setup-context-driven:end id=legacy.rule -->\n",
	)
	commitBaselinePlanTestRepository(t, repo)

	answers := strings.TrimSuffix(humanBaselineAdoptionAnswers(""), "\n")
	answers = "2" + strings.TrimPrefix(answers, "1")
	answers += "\n2\n1\nmanaged setup evidence is not repository policy\n2\n"
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader(answers),
			interactive: true,
		},
	)
	if code != exitUnverified {
		t.Fatalf("edited classification exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "Consolidated Change Plan review") ||
		strings.Contains(stdout.String(), "Category: decision") {
		t.Fatalf("edited managed classification did not produce a valid review:\n%s", stdout.String())
	}
}

func TestHumanAutomationPlanParity(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	var review bytes.Buffer
	var prompts bytes.Buffer
	humanPlan, err := driveHumanBaselinePlan(
		context.Background(),
		repo,
		&baselineHumanPrompt{
			reader: bufioReader(humanBaselineAdoptionAnswers("")),
			writer: &prompts,
		},
		&review,
	)
	if err != nil {
		t.Fatalf("build human Baseline Plan: %v\nreview=%s\nprompts=%s", err, review.String(), prompts.String())
	}
	automation, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  humanBaselineFixtureDecisions(),
		Preservation: baseline.RootPreservationRequest{
			Mode: baseline.PreservationModeGreenfield,
		},
	})
	if err != nil || automation.Plan == nil {
		t.Fatalf("build automation Baseline Plan: outcome=%+v error=%v", automation, err)
	}
	humanBytes, err := baseline.MarshalPlanDocument(humanPlan)
	if err != nil {
		t.Fatal(err)
	}
	automationBytes, err := baseline.MarshalPlanDocument(*automation.Plan)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(humanBytes, automationBytes) {
		t.Fatalf(
			"human and automation plans differ\ndigests: human=%s automation=%s",
			humanPlan.PlanDigest,
			automation.Plan.PlanDigest,
		)
	}
}

func TestBaselineHumanProfileAdaptation(t *testing.T) {
	repository, _, _ := baselinePlanProfileFileFixture(t)
	before := baselinePlanTestTree(t, repository)
	answers := "\n3\n" +
		strings.Repeat("\n", 32) +
		"2\n" +
		"1\n" +
		"guided-human-backend\n"
	var review bytes.Buffer
	var prompts bytes.Buffer
	var request baseline.PlanRequest
	humanPlan, err := driveHumanBaselinePlanWithRequest(
		context.Background(),
		repository,
		&baselineHumanPrompt{
			reader: bufioReader(answers),
			writer: &prompts,
		},
		&review,
		&request,
	)
	if err != nil {
		t.Fatalf(
			"guide human Profile adaptation: %v\nreview=%s\nprompts=%s",
			err,
			review.String(),
			prompts.String(),
		)
	}
	for _, want := range []string{
		"Baseline Profile alignment: action_required",
		"blocking capability.workspace.frontend",
		"advisory capability.firecrawl",
		"Repository-owned Profile adaptation proposal",
		"Modules removed",
		"Capabilities removed",
		"Baseline Profile alignment: ready",
		"Consolidated Change Plan review",
	} {
		if !strings.Contains(review.String(), want) {
			t.Fatalf("guided adaptation review missing %q:\n%s", want, review.String())
		}
	}
	if !strings.Contains(prompts.String(), "Change Baseline Profile") ||
		!strings.Contains(prompts.String(), "repository-owned Profile adaptation") ||
		!strings.Contains(prompts.String(), "Decline without writing") {
		t.Fatalf("Profile divergence choices are incomplete:\n%s", prompts.String())
	}
	if request.ProfileDraft == nil || request.ProfileID != "" {
		t.Fatalf("human Profile adaptation request = %+v", request)
	}

	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	automationInput, err := baseline.ProfileDraftInputFromDocument(
		request.ProfileDraft.Document,
		catalog,
	)
	if err != nil {
		t.Fatalf("normalize automation Profile draft: %v", err)
	}
	automation, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{
		Repository:   repository,
		ProfileDraft: &automationInput,
		Decisions:    request.Decisions,
		Preservation: request.Preservation,
	})
	if err != nil || automation.Plan == nil {
		t.Fatalf("build automation Profile draft Plan: outcome=%+v error=%v", automation, err)
	}
	if !reflect.DeepEqual(humanPlan.Profile, automation.Plan.Profile) ||
		!reflect.DeepEqual(humanPlan.Postimages, automation.Plan.Postimages) ||
		humanPlan.PlanDigest != automation.Plan.PlanDigest {
		t.Fatalf(
			"human and automation Profile drafts differ: human=%s automation=%s",
			humanPlan.PlanDigest,
			automation.Plan.PlanDigest,
		)
	}
	if after := baselinePlanTestTree(t, repository); before != after {
		t.Fatal("human Profile adaptation planning wrote repository bytes")
	}

	t.Run("decline writes nothing", func(t *testing.T) {
		source, resolveErr := baseline.ResolveProfile(
			repository,
			"standard-typescript-monorepo",
			catalog,
		)
		if resolveErr != nil {
			t.Fatal(resolveErr)
		}
		treeBefore := baselinePlanTestTree(t, repository)
		var declineReview bytes.Buffer
		var declinePrompts bytes.Buffer
		_, _, _, declineErr := promptBaselineProfileAlignment(
			context.Background(),
			&baselineHumanPrompt{
				reader: bufioReader("3\n"),
				writer: &declinePrompts,
			},
			&declineReview,
			repository,
			catalog,
			baselineHumanState{},
			source,
			[]baseline.DecisionValue{
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
			},
		)
		var actionErr *baselineHumanActionError
		if !errors.As(declineErr, &actionErr) ||
			!strings.Contains(actionErr.result.Message, "declined") {
			t.Fatalf("decline error = %v", declineErr)
		}
		if treeAfter := baselinePlanTestTree(t, repository); treeBefore != treeAfter {
			t.Fatal("declined Profile adaptation changed repository bytes")
		}
	})

	t.Run("review output failure writes nothing", func(t *testing.T) {
		treeBefore := baselinePlanTestTree(t, repository)
		var outputErr bytes.Buffer
		exit := runBaselineHumanCommandWithIO(
			context.Background(),
			[]string{"--repo", repository},
			failingWriter{err: errors.New("injected human review output failure")},
			&outputErr,
			baselineHumanCommandIO{
				input:       strings.NewReader(answers + "1\n"),
				interactive: true,
			},
		)
		if exit != exitRunFailed ||
			!strings.Contains(outputErr.String(), "injected human review output failure") {
			t.Fatalf("output failure exit=%d stderr=%s", exit, outputErr.String())
		}
		if treeAfter := baselinePlanTestTree(t, repository); treeBefore != treeAfter {
			t.Fatal("human review output failure changed repository bytes")
		}
	})
}

func TestRejectedPlanRevision(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	initial := humanBaselineFixturePlan(t, repo)
	specIndex := humanDecisionIndex(t, initial.Decisions, "spec.scaffold")
	answers := humanBaselineAdoptionAnswers("") +
		strings.Join([]string{
			"3", // reject and revise
			"3", // divergences
			"1", // structured manual revision
			strconv.Itoa(specIndex + 1),
			"2", // change current value
			"2", // no
			"2", // decline the revised digest without writing
		}, "\n") + "\n"
	analyzer := &countingBaselineRevisionAnalyzer{}
	before := baselinePlanTestTree(t, repo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:            strings.NewReader(answers),
			interactive:      true,
			revisionAnalyzer: analyzer,
		},
	)
	if code != exitUnverified {
		t.Fatalf("rejected-plan revision exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	if analyzer.calls != 0 {
		t.Fatalf("direct structured revision started %d ACP sessions", analyzer.calls)
	}
	digests := humanReviewDigests(stdout.String())
	if len(digests) != 2 || digests[0] == digests[1] {
		t.Fatalf("review digests = %v, want distinct original and recomputed approvals\n%s", digests, stdout.String())
	}
	for _, want := range []string{
		"Rejected Plan revision accepted",
		"File changes:",
		"Complete managed-entry ledger:",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("revised review missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Count(stderr.String(), "Final confirmation for Plan Digest") != 2 {
		t.Fatalf("revised plan did not require a new exact approval:\n%s", stderr.String())
	}
	if after := baselinePlanTestTree(t, repo); after != before {
		t.Fatalf("rejected and declined revision changed repository bytes")
	}
}

func TestRepeatedPlanRevisionDeterminism(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	run := func(t *testing.T, repo string) []string {
		t.Helper()
		initial := humanBaselineFixturePlan(t, repo)
		specIndex := humanDecisionIndex(t, initial.Decisions, "spec.scaffold")
		secondbrainIndex := humanDecisionIndex(t, initial.Decisions, "secondbrain.enabled")
		answers := humanBaselineAdoptionAnswers("") +
			strings.Join([]string{
				"3", "3", "1", strconv.Itoa(specIndex + 1), "2", "2",
				"3", "4", "1", strconv.Itoa(secondbrainIndex + 1), "2", "1",
				"2",
			}, "\n") + "\n"
		var stdout bytes.Buffer
		var stderr bytes.Buffer
		code := runBaselineHumanCommandWithIO(
			context.Background(),
			[]string{"--repo", repo},
			&stdout,
			&stderr,
			baselineHumanCommandIO{input: strings.NewReader(answers), interactive: true},
		)
		if code != exitUnverified {
			t.Fatalf("repeated revision exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
		}
		if strings.Count(stdout.String(), "Baseline workflow: adoption") != 1 {
			t.Fatalf("repeated revision restarted repository adoption:\n%s", stdout.String())
		}
		digests := humanReviewDigests(stdout.String())
		if len(digests) != 3 || digests[0] == digests[1] || digests[1] == digests[2] {
			t.Fatalf("repeated revision digests = %v", digests)
		}
		return digests
	}
	first := run(t, repo)
	second := run(t, repo)
	if strings.Join(first, "\n") != strings.Join(second, "\n") {
		t.Fatalf("equivalent repeated revisions are not deterministic:\nfirst=%v\nsecond=%v", first, second)
	}
}

func TestBaselineNoTTY(t *testing.T) {
	repo := newHumanBaselineRepository(t)
	before := baselinePlanTestTree(t, repo)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	code := runBaselineHumanCommandWithIO(
		context.Background(),
		[]string{"--repo", repo, "--format=json"},
		&stdout,
		&stderr,
		baselineHumanCommandIO{
			input:       strings.NewReader("1\n"),
			interactive: false,
		},
	)
	if code != exitUnverified {
		t.Fatalf("no-TTY baseline exit = %d stdout=%s stderr=%s", code, stdout.String(), stderr.String())
	}
	result, err := baseline.ParseResult(stdout.Bytes())
	if err != nil {
		t.Fatalf("parse no-TTY result: %v\n%s", err, stdout.String())
	}
	if result.Category != "interactive_input" ||
		!strings.Contains(result.NextAction, "roundfix baseline plan") ||
		!strings.Contains(result.NextAction, "roundfix baseline apply") {
		t.Fatalf("no-TTY result is not actionable: %+v", result)
	}
	if strings.Contains(stderr.String(), "Prompt ") {
		t.Fatalf("no-TTY baseline emitted a hidden prompt:\n%s", stderr.String())
	}
	if after := baselinePlanTestTree(t, repo); before != after {
		t.Fatalf("no-TTY baseline changed repository bytes:\nbefore=%s\nafter=%s", before, after)
	}

	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "roundfix")
	build := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
	build.Dir = projectRoot
	build.Env = append(os.Environ(), "GOCACHE="+filepath.Join(t.TempDir(), "go-cache"))
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build real roundfix CLI: %v\n%s", err, output)
	}
	command := exec.Command(binary, "baseline", "--repo", repo, "--format=json")
	command.Stdin = strings.NewReader("1\n1\n1\n")
	var realOut bytes.Buffer
	var realErr bytes.Buffer
	command.Stdout = &realOut
	command.Stderr = &realErr
	err = command.Run()
	exitErr, ok := err.(*exec.ExitError)
	if !ok || exitErr.ExitCode() != exitUnverified {
		t.Fatalf("redirected real baseline error=%v stdout=%s stderr=%s", err, realOut.String(), realErr.String())
	}
	if strings.Contains(realErr.String(), "Prompt ") {
		t.Fatalf("redirected real baseline emitted a prompt:\n%s", realErr.String())
	}
}

func newHumanBaselineRepository(t *testing.T) string {
	t.Helper()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeBaselinePlanTestFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeBaselinePlanTestFile(t, repo, "Makefile", "verify:\n\t@true\n")
	commitBaselinePlanTestRepository(t, repo)
	return repo
}

func applyHumanBaselineFixturePlan(t *testing.T, repo string) {
	t.Helper()
	outcome, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  humanBaselineFixtureDecisions(),
		Preservation: baseline.RootPreservationRequest{
			Mode: baseline.PreservationModeGreenfield,
		},
	})
	if err != nil || outcome.Plan == nil {
		t.Fatalf("build update fixture plan: outcome=%+v error=%v", outcome, err)
	}
	if _, err := baseline.ApplyPlan(context.Background(), repo, *outcome.Plan, outcome.Plan.PlanDigest); err != nil {
		t.Fatalf("apply update fixture plan: %v", err)
	}
}

func humanBaselineFixtureDecisions() []baseline.DecisionValue {
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
		{ID: "repository.extension.enabled", Value: false},
	}
}

func humanBaselineFixturePlan(t *testing.T, repo string) baseline.PlanDocument {
	t.Helper()
	outcome, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{
		Repository: repo,
		ProfileID:  "go-cli-tui",
		Decisions:  humanBaselineFixtureDecisions(),
		Preservation: baseline.RootPreservationRequest{
			Mode: baseline.PreservationModeGreenfield,
		},
	})
	if err != nil || outcome.Plan == nil {
		t.Fatalf("build human revision fixture: outcome=%+v error=%v", outcome, err)
	}
	return *outcome.Plan
}

func humanDecisionIndex(t *testing.T, decisions []baseline.DecisionValue, id string) int {
	t.Helper()
	for index, decision := range decisions {
		if decision.ID == id {
			return index
		}
	}
	t.Fatalf("missing decision %q", id)
	return -1
}

func humanReviewDigests(output string) []string {
	var digests []string
	for _, line := range strings.Split(output, "\n") {
		if digest, found := strings.CutPrefix(line, "Plan Digest: "); found {
			digests = append(digests, digest)
		}
	}
	return digests
}

type countingBaselineRevisionAnalyzer struct {
	calls int
}

func (analyzer *countingBaselineRevisionAnalyzer) Revise(
	context.Context,
	baseline.RevisionSnapshot,
) (baseline.RevisionProposal, error) {
	analyzer.calls++
	return baseline.RevisionProposal{}, errors.New("unexpected semantic revision")
}

func humanBaselineAdoptionAnswers(final string) string {
	answers := []string{
		"1",
		"1",
		"1",
		"make verify",
		"1",
		"1",
		"2",
		"1",
		"2",
		"2",
		"codex gpt-5.5 xhigh",
		"claude opus xhigh",
	}
	if final != "" {
		answers = append(answers, final)
	}
	return strings.Join(answers, "\n") + "\n"
}

func humanBaselinePreservationAnswers() string {
	answers := strings.Split(strings.TrimSuffix(humanBaselineAdoptionAnswers(""), "\n"), "\n")
	answers[0] = "2"
	answers[9] = "1"
	return strings.Join(answers, "\n") + "\n1\n2\n"
}

func bufioReader(input string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(input))
}

func decodeHumanBaselineResult(t *testing.T, data []byte) baseline.Result {
	t.Helper()
	var result baseline.Result
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	return result
}
