// Suite: Public Baseline release journeys
// Invariant: a built Roundfix binary applies only one exact reviewed Baseline Plan and composes with repository-owned tools without hidden execution.
// Boundary IN: real process plan/apply handoff, greenfield, preservation, update, profile change, stale and cross-clone plans, unsafe carriers, external composition, and empty reapply.
// Boundary OUT: transaction fault injection (internal/baseline/release_gate_test.go), sealed ACP supervision (internal/baselineacp/release_gate_test.go), and the separately authorized live Fluxus repository.

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
	"sort"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

func TestBaselineMacroJourneysPublicCLI(t *testing.T) {
	binary := buildBaselineReleaseBinary(t)

	t.Run("greenfield automation apply and empty reapply", func(t *testing.T) {
		repo := newBaselineReleaseRepository(t, "go-cli-tui")
		plan, planPath := baselineReleasePlan(t, binary, repo, "go-cli-tui")
		first := baselineReleaseApply(t, binary, repo, plan, planPath)
		if first.State != "verified" {
			t.Fatalf("first apply state = %q, want verified", first.State)
		}
		before := baselineReleaseManagedState(t, repo, plan)
		second := baselineReleaseApply(t, binary, repo, plan, planPath)
		if second.State != "verified" || !strings.Contains(second.Message, "already applied") {
			t.Fatalf("empty reapply result = %+v", second)
		}
		if after := baselineReleaseManagedState(t, repo, plan); after != before {
			t.Fatalf("empty reapply changed managed state:\nbefore=%s\nafter=%s", before, after)
		}
	})

	t.Run("preservation", func(t *testing.T) {
		const repositoryRule = "Keep this repository-specific release rule.\n"
		repo := newBaselineReleaseRepository(t, "go-cli-tui")
		writeBaselinePlanTestFile(t, repo, "AGENTS.md", repositoryRule)
		commitBaselinePlanTestRepository(t, repo)
		decisionPath := baselineReleasePreservationDecisionFile(t, repo)
		plan, planPath := baselineReleasePlanWithDecisionFile(
			t, binary, repo, "go-cli-tui", decisionPath,
		)
		baselineReleaseApply(t, binary, repo, plan, planPath)

		rules, err := os.ReadFile(filepath.Join(repo, "docs", "agents", "repository-rules.md"))
		if err != nil || !bytes.Contains(rules, []byte(strings.TrimSpace(repositoryRule))) {
			t.Fatalf("preserved repository rule = %q error=%v", rules, err)
		}
		var backup baseline.ManagedEntry
		for _, entry := range plan.ManagedEntries {
			if entry.Kind == "backup" {
				backup = entry
				break
			}
		}
		content, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(backup.Path)))
		if err != nil || string(content) != repositoryRule {
			t.Fatalf("immutable preservation backup %q = %q error=%v", backup.Path, content, err)
		}
	})

	t.Run("update and profile change", func(t *testing.T) {
		repo := newBaselineReleaseRepository(t, "go-cli-tui")
		goPlan, goPlanPath := baselineReleasePlan(t, binary, repo, "go-cli-tui")
		baselineReleaseApply(t, binary, repo, goPlan, goPlanPath)

		updatePlan, updatePath := baselineReleasePlan(t, binary, repo, "go-cli-tui")
		baselineReleaseApply(t, binary, repo, updatePlan, updatePath)

		rustPlan, rustPlanPath := baselineReleasePlan(t, binary, repo, "rust-cli")
		if rustPlan.Profile.ID != "rust-cli" || rustPlan.PlanDigest == goPlan.PlanDigest {
			t.Fatalf("profile-change plan = profile %q digest %q, original %q",
				rustPlan.Profile.ID, rustPlan.PlanDigest, goPlan.PlanDigest)
		}
		baselineReleaseApply(t, binary, repo, rustPlan, rustPlanPath)
		manifest, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(baselineSetupManifestPath)))
		if err != nil || !bytes.Contains(manifest, []byte(`"profile": "rust-cli"`)) {
			t.Fatalf("profile-change manifest = %q error=%v", manifest, err)
		}
	})

	t.Run("stale plan", func(t *testing.T) {
		repo := newBaselineReleaseRepository(t, "go-cli-tui")
		plan, planPath := baselineReleasePlan(t, binary, repo, "go-cli-tui")
		writeBaselinePlanTestFile(t, repo, "Makefile", "verify:\n\t@echo changed\n")
		before := baselineReleaseVisibleState(t, repo)
		code, stdout, _ := runBaselineReleaseCLI(
			t, binary,
			"baseline", "apply",
			"--repo", repo,
			"--plan", planPath,
			"--confirm-plan", plan.PlanDigest,
			"--format=json",
		)
		if code != exitUnverified {
			t.Fatalf("stale apply exit = %d stdout=%s, want %d", code, stdout, exitUnverified)
		}
		result, err := baseline.ParseResult(stdout)
		if err != nil || result.Category != "stale" {
			t.Fatalf("stale apply result = %+v error=%v", result, err)
		}
		if after := baselineReleaseVisibleState(t, repo); after != before {
			t.Fatalf("stale apply changed repository:\nbefore=%s\nafter=%s", before, after)
		}
	})

	t.Run("cross-clone apply", func(t *testing.T) {
		source := newBaselineReleaseRepository(t, "go-cli-tui")
		plan, planPath := baselineReleasePlan(t, binary, source, "go-cli-tui")
		parent := t.TempDir()
		clone := filepath.Join(parent, "matching")
		runBaselinePlanTestCommand(t, parent, "git", "clone", "--quiet", source, clone)
		result := baselineReleaseApply(t, binary, clone, plan, planPath)
		if result.State != "verified" {
			t.Fatalf("cross-clone apply result = %+v", result)
		}
	})

	t.Run("unsafe carrier", func(t *testing.T) {
		repo := newBaselinePlanTestRepository(t)
		if err := os.Symlink("../outside.md", filepath.Join(repo, "AGENTS.md")); err != nil {
			t.Fatal(err)
		}
		commitBaselinePlanTestRepository(t, repo)
		before := baselineReleaseVisibleState(t, repo)
		code, stdout, _ := runBaselineReleaseCLI(
			t, binary, "baseline", "plan", "--repo", repo, "--format=json",
		)
		if code != exitPreflight {
			t.Fatalf("unsafe plan exit = %d stdout=%s, want %d", code, stdout, exitPreflight)
		}
		result, err := baseline.ParseResult(stdout)
		if err != nil || result.Category != "preflight" {
			t.Fatalf("unsafe carrier result = %+v error=%v", result, err)
		}
		if after := baselineReleaseVisibleState(t, repo); after != before {
			t.Fatalf("unsafe planning changed repository:\nbefore=%s\nafter=%s", before, after)
		}
	})

	t.Run("interactive and automation parity", TestHumanAutomationPlanParity)
	t.Run("consolidated preservation review", TestConsolidatedReview)
	t.Run("rejected-plan revision", TestRejectedPlanRevision)
}

func TestBaselineFindingRegressionsHumanReview(t *testing.T) {
	t.Run("decisions use one consolidated review", TestConsolidatedReview)
	t.Run("file projection precedes canonical ledgers", TestRejectedPlanRevision)
}

func TestBaselineDocumentationContractExamples(t *testing.T) {
	t.Run("public examples parse", TestBaselineExamplesParse)
	t.Run("Decision Documents parse", TestBaselineDecisionExamples)
}

func TestBaselineFormatterComposition(t *testing.T) {
	binary := buildBaselineReleaseBinary(t)
	tests := []struct {
		profile   string
		formatter func(*testing.T, string)
	}{
		{
			profile: "go-cli-tui",
			formatter: func(t *testing.T, repo string) {
				runBaselineReleaseExternal(t, repo, "gofmt", "-w", "main.go")
			},
		},
		{
			profile: "rust-cli",
			formatter: func(t *testing.T, repo string) {
				runBaselineReleaseExternal(t, repo, "rustfmt", "src/main.rs")
			},
		},
		{
			profile: "standard-typescript-monorepo",
			formatter: func(t *testing.T, repo string) {
				runBaselineReleaseExternal(
					t,
					repo,
					"bunx", "--no-install", "oxfmt@0.59.0", "--check", "AGENTS.md", "docs/agents",
				)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.profile, func(t *testing.T) {
			repo := newBaselineReleaseRepository(t, test.profile)
			plan, planPath := baselineReleasePlan(t, binary, repo, test.profile)
			baselineReleaseApply(t, binary, repo, plan, planPath)

			for _, marker := range []string{".qa-verification-ran", ".qa-profile-command-ran"} {
				if _, err := os.Stat(filepath.Join(repo, marker)); !errors.Is(err, os.ErrNotExist) {
					t.Fatalf("Baseline executed repository command marker %s: %v", marker, err)
				}
			}

			managedBefore := baselineReleaseManagedState(t, repo, plan)
			test.formatter(t, repo)
			runBaselineReleaseExternal(t, repo, "make", "verify")
			if _, err := os.Stat(filepath.Join(repo, ".qa-verification-ran")); err != nil {
				t.Fatalf("external repository Verification did not run: %v", err)
			}
			if managedAfter := baselineReleaseManagedState(t, repo, plan); managedAfter != managedBefore {
				t.Fatalf("external composition changed managed state:\nbefore=%s\nafter=%s",
					managedBefore, managedAfter)
			}

			baselineReleaseApply(t, binary, repo, plan, planPath)
			if managedAfter := baselineReleaseManagedState(t, repo, plan); managedAfter != managedBefore {
				t.Fatalf("empty reapply changed managed state:\nbefore=%s\nafter=%s",
					managedBefore, managedAfter)
			}
		})
	}
}

func buildBaselineReleaseBinary(t *testing.T) string {
	t.Helper()
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	binary := filepath.Join(t.TempDir(), "roundfix")
	command := exec.Command("go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
	command.Dir = projectRoot
	command.Env = append(os.Environ(), "GOCACHE="+filepath.Join(t.TempDir(), "go-cache"))
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("build release-gate Roundfix binary: %v\n%s", err, output)
	}
	return binary
}

func baselineReleasePlan(
	t *testing.T,
	binary string,
	repo string,
	profile string,
) (baseline.PlanDocument, string) {
	t.Helper()
	args := []string{"baseline", "plan", "--repo", repo, "--profile", profile}
	for _, decision := range baselineReleaseDecisionArgs(profile, "greenfield") {
		args = append(args, "--decision", decision)
	}
	args = append(args, "--format=json")
	code, stdout, stderr := runBaselineReleaseCLI(t, binary, args...)
	if code != exitOK {
		t.Fatalf("release plan exit = %d stdout=%s stderr=%s", code, stdout, stderr)
	}
	return baselineReleasePlanFile(t, stdout)
}

func baselineReleasePlanWithDecisionFile(
	t *testing.T,
	binary string,
	repo string,
	profile string,
	decisionPath string,
) (baseline.PlanDocument, string) {
	t.Helper()
	code, stdout, stderr := runBaselineReleaseCLI(
		t,
		binary,
		"baseline", "plan",
		"--repo", repo,
		"--profile", profile,
		"--decision-file", decisionPath,
		"--format=json",
	)
	if code != exitOK {
		t.Fatalf("preservation plan exit = %d stdout=%s stderr=%s", code, stdout, stderr)
	}
	return baselineReleasePlanFile(t, stdout)
}

func baselineReleasePlanFile(t *testing.T, data []byte) (baseline.PlanDocument, string) {
	t.Helper()
	plan, err := baseline.ParsePlanDocument(data)
	if err != nil {
		t.Fatalf("parse release-gate Plan: %v\n%s", err, data)
	}
	planPath := filepath.Join(t.TempDir(), "baseline-plan.json")
	if err := os.WriteFile(planPath, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return plan, planPath
}

func baselineReleaseApply(
	t *testing.T,
	binary string,
	repo string,
	plan baseline.PlanDocument,
	planPath string,
) baseline.Result {
	t.Helper()
	code, stdout, stderr := runBaselineReleaseCLI(
		t,
		binary,
		"baseline", "apply",
		"--repo", repo,
		"--plan", planPath,
		"--confirm-plan", plan.PlanDigest,
		"--format=json",
	)
	if code != exitOK {
		t.Fatalf("release apply exit = %d stdout=%s stderr=%s", code, stdout, stderr)
	}
	result, err := baseline.ParseResult(stdout)
	if err != nil {
		t.Fatalf("parse release apply result: %v\n%s", err, stdout)
	}
	return result
}

func baselineReleasePreservationDecisionFile(t *testing.T, repo string) string {
	t.Helper()
	inspection, err := baseline.InspectRepository(context.Background(), repo, nil)
	if err != nil {
		t.Fatal(err)
	}
	preservation, err := baseline.PlanRootPreservation(
		inspection,
		baseline.RootPreservationRequest{Mode: baseline.PreservationModePreservation},
	)
	if err != nil || preservation.DecisionSkeleton == nil {
		t.Fatalf("build preservation Decision Document: plan=%+v error=%v", preservation, err)
	}
	document := preservation.DecisionSkeleton.Document
	for _, raw := range baselineReleaseDecisionArgs("go-cli-tui", "preservation") {
		id, value, ok := strings.Cut(raw, "=")
		if !ok {
			t.Fatalf("invalid release decision %q", raw)
		}
		document.Decisions = append(document.Decisions, baseline.DecisionValue{
			ID: id, Value: baselineReleaseDecisionValue(value),
		})
	}
	data, err := json.MarshalIndent(document, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "preservation-decisions.json")
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func baselineReleaseDecisionValue(raw string) any {
	switch raw {
	case "true":
		return true
	case "false":
		return false
	}
	if strings.HasPrefix(raw, "{") {
		var value any
		if json.Unmarshal([]byte(raw), &value) == nil {
			return value
		}
	}
	return raw
}

func baselineReleaseDecisionArgs(profile, preservation string) []string {
	decisions := []string{
		"preservation.mode=" + preservation,
		"language.generated=English",
		"verification.gate=make verify",
		"spec.scaffold=true",
		"domain.layout=single-context",
		"triage.external=false",
		"autonomous.enabled=false",
		"secondbrain.enabled=false",
		"repository.extension.enabled=false",
	}
	if profile == "standard-typescript-monorepo" {
		decisions = append(decisions, `http.contract={"mode":"REST"}`)
	}
	return decisions
}

func newBaselineReleaseRepository(t *testing.T, profile string) string {
	t.Helper()
	repo := newBaselinePlanTestRepository(t)
	writeBaselinePlanTestFile(t, repo, ".agents/skills/context7/SKILL.md", "# context7\n")
	writeBaselinePlanTestFile(t, repo, ".agents/skills/exa-web-search/SKILL.md", "# exa\n")
	writeBaselinePlanTestFile(t, repo, "Makefile", strings.Join([]string{
		"verify:",
		"\t@touch .qa-verification-ran",
		"",
	}, "\n"))

	switch profile {
	case "go-cli-tui":
		writeBaselinePlanTestFile(t, repo, "main.go", "package main\nfunc main( ){ }\n")
	case "rust-cli":
		writeBaselinePlanTestFile(t, repo, "src/main.rs", "fn main( ) {println!(\"ok\");}\n")
	case "standard-typescript-monorepo":
		writeBaselinePlanTestFile(t, repo, "package.json", baselineReleaseTypeScriptPackageJSON())
		writeBaselinePlanTestFile(t, repo, "packages/frontend/package.json", `{"name":"frontend"}`)
		writeBaselinePlanTestFile(
			t,
			repo,
			"packages/backend/package.json",
			`{"name":"backend","dependencies":{"postgres":"latest","drizzle-orm":"latest"}}`,
		)
		writeBaselinePlanTestFile(
			t,
			repo,
			"packages/backend/src/infra/controllers/http/app.ts",
			"app.get('/health', handler)\napp.post('/orders', handler)\n",
		)
		writeBaselinePlanTestFile(
			t,
			repo,
			"DATABASE.md",
			"# Database\n\nPostgreSQL is the repository database contract.\n",
		)
	default:
		t.Fatalf("unknown maintained profile %q", profile)
	}
	commitBaselinePlanTestRepository(t, repo)
	return repo
}

func baselineReleaseTypeScriptPackageJSON() string {
	dependencies := []string{
		`"@logtape/logtape":"latest"`,
		`"@tanstack/react-query":"latest"`,
		`"@tanstack/react-router":"latest"`,
		`"better-auth":"latest"`,
		`"drizzle-orm":"latest"`,
		`"hono":"latest"`,
		`"oxfmt":"0.59.0"`,
		`"oxlint":"latest"`,
		`"postgres":"latest"`,
		`"react":"latest"`,
		`"shadcn":"latest"`,
		`"tailwindcss":"latest"`,
		`"turbo":"latest"`,
		`"typescript":"latest"`,
		`"vite":"latest"`,
		`"vitest":"latest"`,
		`"zod":"latest"`,
	}
	return `{"name":"root","packageManager":"bun@1.3.0","scripts":{` +
		`"format":"touch .qa-profile-command-ran","lint":"touch .qa-profile-command-ran",` +
		`"test":"touch .qa-profile-command-ran","build":"touch .qa-profile-command-ran",` +
		`"verify":"touch .qa-profile-command-ran"},` +
		`"dependencies":{` + strings.Join(dependencies, ",") + `}}`
}

func runBaselineReleaseCLI(
	t *testing.T,
	binary string,
	args ...string,
) (int, []byte, []byte) {
	t.Helper()
	command := exec.Command(binary, args...)
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	err := command.Run()
	if err == nil {
		return 0, stdout.Bytes(), stderr.Bytes()
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatalf("run Roundfix release binary: %v", err)
	}
	return exitErr.ExitCode(), stdout.Bytes(), stderr.Bytes()
}

func runBaselineReleaseExternal(t *testing.T, repo, name string, args ...string) {
	t.Helper()
	command := exec.Command(name, args...)
	command.Dir = repo
	command.Env = append(os.Environ(), "GIT_CONFIG_NOSYSTEM=1", "GIT_OPTIONAL_LOCKS=0")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("run external %s %s: %v\n%s", name, strings.Join(args, " "), err, output)
	}
}

func baselineReleaseManagedState(
	t *testing.T,
	repo string,
	plan baseline.PlanDocument,
) string {
	t.Helper()
	digest := sha256.New()
	for _, postimage := range plan.Postimages {
		content, err := os.ReadFile(filepath.Join(repo, filepath.FromSlash(postimage.Path)))
		if err != nil {
			t.Fatalf("read managed path %s: %v", postimage.Path, err)
		}
		fmt.Fprintf(digest, "%s\x00%x\x00", postimage.Path, sha256.Sum256(content))
	}
	return fmt.Sprintf("%x", digest.Sum(nil))
}

func baselineReleaseVisibleState(t *testing.T, root string) string {
	t.Helper()
	var rows []string
	err := filepath.WalkDir(root, func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		relative = filepath.ToSlash(relative)
		if relative == ".git" {
			return filepath.SkipDir
		}
		if relative == "." {
			return nil
		}
		info, err := os.Lstat(name)
		if err != nil {
			return err
		}
		row := fmt.Sprintf("%s:%s:%#o", relative, info.Mode().Type(), info.Mode().Perm())
		switch {
		case info.Mode()&os.ModeSymlink != 0:
			target, err := os.Readlink(name)
			if err != nil {
				return err
			}
			row += ":" + target
		case info.Mode().IsRegular():
			content, err := os.ReadFile(name)
			if err != nil {
				return err
			}
			row += fmt.Sprintf(":%x", sha256.Sum256(content))
		}
		rows = append(rows, row)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(rows)
	return strings.Join(rows, "\n")
}
