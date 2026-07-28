package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/codex"
	roundconfig "roundfix/internal/config"
	"roundfix/skills"
)

func TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts(t *testing.T) {
	tests := []struct {
		name       string
		checker    *doctorFakeHealthChecker
		wantCode   int
		wantStdout string
	}{
		{
			name: "all checks pass",
			checker: newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK, Detail: "/home/roundfix/.local/bin/codex accepted"},
			),
			wantCode: exitOK,
			wantStdout: "node: ok (v25.6.1 >= " + setupNodeMinimumVersion + ")\n" +
				"acpx: ok (" + agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion + ")\n" +
				"adapter: ok (codex-acp)\n" +
				"profiles: ok (3 distinct tuples; 10 category references)\n" +
				"skills: ok (39 required: 14 Roundfix-owned, 25 external)\n" +
				"codex: ok (/home/roundfix/.local/bin/codex accepted)\n",
		},
		{
			name: "quarantined codex fails with reinstall action",
			checker: newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusFailed, Detail: "/tmp/codex is quarantined", NextAction: codex.ReinstallNextAction},
			),
			wantCode: exitRunFailed,
			wantStdout: "node: ok (v25.6.1 >= " + setupNodeMinimumVersion + ")\n" +
				"acpx: ok (" + agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion + ")\n" +
				"adapter: ok (codex-acp)\n" +
				"profiles: ok (3 distinct tuples; 10 category references)\n" +
				"skills: ok (39 required: 14 Roundfix-owned, 25 external)\n" +
				"codex: failed (/tmp/codex is quarantined; next: " + codex.ReinstallNextAction + ")\n",
		},
		{
			name: "codex not applicable does not fail",
			checker: newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusSkipped, Detail: "not-applicable on linux"},
			),
			wantCode: exitOK,
			wantStdout: "node: ok (v25.6.1 >= " + setupNodeMinimumVersion + ")\n" +
				"acpx: ok (" + agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion + ")\n" +
				"adapter: ok (codex-acp)\n" +
				"profiles: ok (3 distinct tuples; 10 category references)\n" +
				"skills: ok (39 required: 14 Roundfix-owned, 25 external)\n" +
				"codex: skipped (not-applicable on linux)\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir, repoDir := withCLIWorkspace(t)
			runner := &profileReadinessExactRunner{}
			withAgentRunner(t, runner)
			withDoctorFakeDeps(t, tt.checker)
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"doctor"}, &stdout, &stderr)

			if code != tt.wantCode {
				t.Fatalf("expected exit code %d, got %d", tt.wantCode, code)
			}
			if got := stdout.String(); got != tt.wantStdout {
				t.Fatalf("unexpected stdout:\n got: %q\nwant: %q", got, tt.wantStdout)
			}
			if stderr.Len() != 0 {
				t.Fatalf("expected no stderr, got %q", stderr.String())
			}
			if len(runner.exactRequests) != 3 {
				t.Fatalf("expected three distinct profile proofs, got %#v", runner.exactRequests)
			}
			wantModels := []string{"gpt-5.6-sol", "gpt-5.5", "opus"}
			for index, wantModel := range wantModels {
				if request := runner.exactRequests[index]; request.WorkDir != "/repo/project" || request.Runtime.Model != wantModel {
					t.Fatalf("profile proof %d = %#v, want model %q in repository", index, request, wantModel)
				}
			}
			if len(tt.checker.agentRequests) != 0 {
				t.Fatalf("Doctor must not run the legacy configured Agent probe, got %#v", tt.checker.agentRequests)
			}
			assertDoctorPathMissing(t, filepath.Join(homeDir, ".acpx"))
			assertDoctorPathMissing(t, filepath.Join(homeDir, ".roundfix"))
			assertDoctorPathMissing(t, filepath.Join(repoDir, ".roundfixrc.yml"))
		})
	}
}

func TestRunDoctorProfileReadinessReportsLegacyAdapterThroughEffectiveProfile(t *testing.T) {
	config := roundconfig.Builtin()
	config.Defaults.Agent = "codex"
	config.Runtimes.Codex.Model = "legacy-model-default"
	proofs, err := buildProfileProofReports(config, roundconfig.RequiredWorkCategories())
	if err != nil {
		t.Fatalf("build profile proof reports: %v", err)
	}
	legacy := &agent.AdapterLineageError{
		Command: "codex-acp",
		Package: "@zed-industries/codex-acp",
		Version: "0.16.0",
	}
	applyProfileProofFailure(&proofs[0], legacy)
	readiness := profileProofResult{Proofs: proofs, Err: profileProofError{
		Selection:      proofs[0].Selection,
		References:     proofs[0].References,
		Classification: proofs[0].Classification,
		NextAction:     proofs[0].NextAction,
		Err:            legacy,
	}}
	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK, Detail: "/home/roundfix/.local/bin/codex accepted"},
	)
	checker.adapter = CheckResult{
		Name:   HealthCheckAdapter,
		Status: CheckStatusOK,
		Detail: "command=\"codex-acp\"; package=@zed-industries/codex-acp; version=0.16.0",
	}
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{
		Config:  config,
		GitRoot: "/repo/project",
		HomeDir: "/home/roundfix-test",
	}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return readiness
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"doctor"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected doctor adapter failure exit %d, got %d", exitRunFailed, code)
	}
	for _, want := range []string{
		"adapter: ok (command=\"codex-acp\"; package=@zed-industries/codex-acp; version=0.16.0)",
		"profiles: failed",
		`runtime="codex", model="gpt-5.6-sol", reasoning_effort="high"`,
		"affected categories: general preferred source=built-in, backend preferred source=built-in, frontend fallback[1] source=built-in, qa preferred source=built-in, review preferred source=built-in",
		"classification: adapter_lineage_unknown",
		"adapter evidence: command=\"codex-acp\", version=\"0.16.0\"",
		"next: run `" + agent.CodexAdapterInstallCommand() + "`",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("expected stdout to contain %q, got %q", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "legacy-model-default") || strings.Contains(stdout.String(), "model:") || strings.Contains(stdout.String(), "agent:") {
		t.Fatalf("Doctor reported legacy configured-runtime readiness: %q", stdout.String())
	}
	if len(checker.adapterRuntimes) != 1 || checker.adapterRuntimes[0].Model != "gpt-5.6-sol" || checker.adapterRuntimes[0].ReasoningEffort != "high" {
		t.Fatalf("adapter check did not use the effective general profile: %#v", checker.adapterRuntimes)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunDoctorProfileReadinessMatchesProfilesValidateFailureEvidence(t *testing.T) {
	_, repoDir := withCLIWorkspace(t)
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), `
profiles:
  backend:
    preferred:
      runtime: codex
      model: broken-backend
      reasoning_effort: high
    fallbacks:
      - runtime: claude
        model: backup-backend
        reasoning_effort: medium
`)
	failure := &agent.SelectionUnsupportedError{
		Kind:                agent.SelectionReasoningControlNotAdvertised,
		Runtime:             "codex",
		Model:               "broken-backend",
		ReasoningEffort:     "high",
		AdvertisedModels:    []string{"broken-backend", "gpt-5.5"},
		AdvertisedReasoning: []string{"low", "medium"},
	}
	runner := &profileReadinessExactRunner{prove: func(req agent.ProbeRequest) (agent.SelectionProof, error) {
		if req.Runtime.Model == "broken-backend" {
			return agent.SelectionProof{}, failure
		}
		return agent.SelectionProof{}, nil
	}}
	withAgentRunner(t, runner)
	var validateStdout bytes.Buffer
	var validateStderr bytes.Buffer
	validateCode := Run([]string{"profiles", "validate", "--category", "backend", "--json"}, &validateStdout, &validateStderr)
	if validateCode != exitPreflight {
		t.Fatalf("profiles validate exit = %d, want %d", validateCode, exitPreflight)
	}

	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK, Detail: "/home/roundfix/.local/bin/codex accepted"},
	)
	withDoctorLiveDeps(t, checker)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"doctor"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected doctor selection failure exit %d, got %d", exitRunFailed, code)
	}
	for _, evidence := range []struct {
		doctor   string
		validate string
	}{
		{`runtime="codex", model="broken-backend", reasoning_effort="high"`, `runtime "codex", model "broken-backend", reasoning_effort "high"`},
		{"affected categories: backend preferred source=project", "affected categories: backend preferred source=project"},
		{"classification: reasoning_control_not_advertised", "classification: reasoning_control_not_advertised"},
		{"advertised_models=broken-backend,gpt-5.5", "advertised Agent Models: broken-backend, gpt-5.5"},
		{"advertised_reasoning=low,medium", "advertised reasoning efforts: low, medium"},
		{"roundfix profiles configure --scope user|project", "roundfix profiles configure --scope user|project"},
	} {
		if !strings.Contains(stdout.String(), evidence.doctor) {
			t.Fatalf("Doctor output missing shared failure evidence %q: %q", evidence.doctor, stdout.String())
		}
		if !strings.Contains(validateStderr.String(), evidence.validate) {
			t.Fatalf("profiles validate output missing shared failure evidence %q: %q", evidence.validate, validateStderr.String())
		}
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunDoctorContinuesChecksAfterProfileReadinessFailure(t *testing.T) {
	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK, Detail: "v25.6.1 >= " + setupNodeMinimumVersion},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK, Detail: agent.MinimumACPXVersion + " >= " + agent.MinimumACPXVersion},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK, Detail: "/home/roundfix/.local/bin/codex accepted"},
	)
	proofCalls := 0
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{
		Config:  roundconfig.Builtin(),
		GitRoot: "/repo/project",
		HomeDir: "/home/roundfix-test",
	}, func(_ context.Context, _ roundconfig.Config, categories []roundconfig.WorkCategory, workDir string) profileProofResult {
		proofCalls++
		if got := formatWorkCategories(categories); got != "general, backend, frontend, qa, review" || workDir != "/repo/project" {
			t.Fatalf("profile readiness input categories=%q workDir=%q", got, workDir)
		}
		return profileProofResult{Err: errors.New("profile proof unavailable")}
	})
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"doctor"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("expected rejected model exit %d, got %d", exitRunFailed, code)
	}
	if proofCalls != 1 || checker.nodeCalls != 1 || checker.acpxCalls != 1 || checker.adapterCalls != 1 || checker.codexCalls != 1 {
		t.Fatalf("independent check calls profile=%d node=%d acpx=%d adapter=%d codex=%d", proofCalls, checker.nodeCalls, checker.acpxCalls, checker.adapterCalls, checker.codexCalls)
	}
	output := stdout.String()
	for _, want := range []string{"node: ok", "acpx: ok", "adapter: ok", "profiles: failed (profile proof unavailable)", "codex: ok"} {
		if !strings.Contains(output, want) {
			t.Fatalf("Doctor output missing %q after profile failure: %q", want, output)
		}
	}
	if strings.Index(output, "profiles:") > strings.Index(output, "codex:") {
		t.Fatalf("profile readiness must precede Codex so Repository Skill Set readiness can be inserted after it: %q", output)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunDoctorRepositorySkillReadiness(t *testing.T) {
	ownedCommand := "roundfix skills install --target project"
	externalCommand := "bunx skills experimental_install && bunx skills update -p -y"
	tests := []struct {
		name      string
		readiness skills.RepositoryReadiness
		err       error
		wantCode  int
		wantLine  string
	}{
		{
			name: "ready",
			readiness: skills.RepositoryReadiness{
				OwnedRequired:    14,
				ExternalRequired: 25,
			},
			wantCode: exitOK,
			wantLine: "skills: ok (39 required: 14 Roundfix-owned, 25 external)",
		},
		{
			name: "owned failure",
			readiness: skills.RepositoryReadiness{
				OwnedRequired: 14,
				MissingOwned:  []string{"write-prd"},
				OutdatedOwned: []string{"roundfix"},
			},
			wantCode: exitRunFailed,
			wantLine: "skills: failed (missing: write-prd; outdated: roundfix; next: " +
				ownedCommand + ")",
		},
		{
			name: "external failure",
			readiness: skills.RepositoryReadiness{
				ExternalRequired: 25,
				MissingExternal:  []string{"testing-boss"},
				OutdatedExternal: []string{"agentic-cli-design"},
			},
			wantCode: exitRunFailed,
			wantLine: "skills: failed (missing: testing-boss; outdated: agentic-cli-design; next: " +
				externalCommand + ")",
		},
		{
			name: "mixed failure is sorted with ordered remediation",
			readiness: skills.RepositoryReadiness{
				OwnedRequired:    14,
				ExternalRequired: 25,
				MissingOwned:     []string{"write-prd"},
				MissingExternal:  []string{"agentic-cli-design"},
				OutdatedOwned:    []string{"roundfix"},
				OutdatedExternal: []string{"testing-boss"},
			},
			wantCode: exitRunFailed,
			wantLine: "skills: failed (missing: agentic-cli-design, write-prd; outdated: roundfix, testing-boss; next: " +
				ownedCommand + " && " + externalCommand + ")",
		},
		{
			name: "symlinked lock keeps external remediation",
			err: &skills.RepositoryReadinessError{
				Ownership: skills.RepositoryOwnershipExternal,
				Operation: "inspect skills lock",
				Path:      "/repo/project/skills-lock.json",
				Err:       errors.New("path must be a regular file"),
			},
			wantCode: exitRunFailed,
			wantLine: "skills: failed (inspect skills lock \"/repo/project/skills-lock.json\": path must be a regular file; next: " +
				externalCommand + ")",
		},
		{
			name:     "unclassified checker error",
			err:      errors.New("repository skill check failed"),
			wantCode: exitRunFailed,
			wantLine: "skills: failed (repository skill check failed; next: " +
				ownedCommand + " && " + externalCommand + ")",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var calls []string
			recordCall := func(name string) {
				calls = append(calls, name)
			}
			checker := newDoctorFakeHealthChecker(
				CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
				CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
			)
			checker.recordCall = recordCall
			withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{
				Config:  roundconfig.Builtin(),
				GitRoot: "  /repo/project\n",
			}, func(_ context.Context, _ roundconfig.Config, _ []roundconfig.WorkCategory, workDir string) profileProofResult {
				recordCall(HealthCheckProfiles)
				if workDir != "/repo/project" {
					t.Fatalf("profile readiness work directory = %q, want trimmed Git root", workDir)
				}
				return profileProofResult{}
			})
			skillCalls := 0
			doctorDeps.checkSkills = func(_ context.Context, root string) (skills.RepositoryReadiness, error) {
				skillCalls++
				recordCall(HealthCheckSkills)
				if root != "/repo/project" {
					t.Fatalf("skill readiness root = %q, want /repo/project", root)
				}
				return test.readiness, test.err
			}
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"doctor"}, &stdout, &stderr)

			if code != test.wantCode {
				t.Fatalf("exit code = %d, want %d; stderr=%q", code, test.wantCode, stderr.String())
			}
			lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
			if len(lines) != 6 || lines[4] != test.wantLine {
				t.Fatalf("unexpected Doctor output lines:\n%q\nwant skills line %q at index 4", lines, test.wantLine)
			}
			if skillCalls != 1 || checker.nodeCalls != 1 || checker.acpxCalls != 1 || checker.adapterCalls != 1 || checker.codexCalls != 1 {
				t.Fatalf("independent check calls skills=%d node=%d acpx=%d adapter=%d codex=%d",
					skillCalls, checker.nodeCalls, checker.acpxCalls, checker.adapterCalls, checker.codexCalls)
			}
			wantCalls := []string{
				HealthCheckNode,
				HealthCheckACPX,
				HealthCheckAdapter,
				HealthCheckProfiles,
				HealthCheckSkills,
				HealthCheckCodex,
			}
			if !reflect.DeepEqual(calls, wantCalls) {
				t.Fatalf("check order = %v, want %v", calls, wantCalls)
			}
			if stderr.Len() != 0 {
				t.Fatalf("expected no stderr, got %q", stderr.String())
			}
		})
	}
}

func TestRunDoctorPassesCommandContextToRepositorySkillReadiness(t *testing.T) {
	type contextKey struct{}
	const marker = "doctor-command"

	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
	)
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{
		Config:  roundconfig.Builtin(),
		GitRoot: "/repo/project",
	}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return profileProofResult{}
	})
	doctorDeps.checkSkills = func(ctx context.Context, root string) (skills.RepositoryReadiness, error) {
		if got := ctx.Value(contextKey{}); got != marker {
			t.Fatalf("repository checker context marker = %v, want %q", got, marker)
		}
		if root != "/repo/project" {
			t.Fatalf("repository checker root = %q, want /repo/project", root)
		}
		return skills.RepositoryReadiness{
			OwnedRequired:    14,
			ExternalRequired: 25,
		}, nil
	}
	commandContext := context.WithValue(t.Context(), contextKey{}, marker)
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := runDoctorCommand(commandContext, nil, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("exit code = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
}

func TestRunDoctorMissingRepositoryRoot(t *testing.T) {
	processDir := t.TempDir()
	t.Chdir(processDir)
	var calls []string
	recordCall := func(name string) {
		calls = append(calls, name)
	}
	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
	)
	checker.recordCall = recordCall
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{
		Config: roundconfig.Builtin(),
	}, func(_ context.Context, _ roundconfig.Config, _ []roundconfig.WorkCategory, workDir string) profileProofResult {
		recordCall(HealthCheckProfiles)
		if workDir != processDir {
			t.Fatalf("profile readiness work directory = %q, want process working directory %q", workDir, processDir)
		}
		return profileProofResult{}
	})
	skillCalls := 0
	doctorDeps.checkSkills = func(_ context.Context, root string) (skills.RepositoryReadiness, error) {
		skillCalls++
		return skills.RepositoryReadiness{}, fmt.Errorf("unexpected repository check for %q", root)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"doctor"}, &stdout, &stderr)

	if code != exitRunFailed {
		t.Fatalf("exit code = %d, want %d", code, exitRunFailed)
	}
	wantStdout := "node: ok\n" +
		"acpx: ok\n" +
		"adapter: ok (codex-acp)\n" +
		"profiles: ok (0 distinct tuples; 0 category references)\n" +
		"skills: failed (Repository Skill Set readiness requires a Git repository; next: run roundfix doctor from a Git repository)\n" +
		"codex: ok\n"
	if got := stdout.String(); got != wantStdout {
		t.Fatalf("unexpected stdout:\n got: %q\nwant: %q", got, wantStdout)
	}
	if skillCalls != 0 {
		t.Fatalf("repository checker calls = %d, want 0", skillCalls)
	}
	wantCalls := []string{
		HealthCheckNode,
		HealthCheckACPX,
		HealthCheckAdapter,
		HealthCheckProfiles,
		HealthCheckCodex,
	}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Fatalf("independent check order = %v, want %v", calls, wantCalls)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
}

func TestRunDoctorRealRepositoryCheckDoesNotMutateState(t *testing.T) {
	homeDir := t.TempDir()
	repoDir := t.TempDir()
	t.Setenv("HOME", homeDir)
	t.Chdir(repoDir)
	mustMkdir(t, filepath.Join(repoDir, ".git"))
	writeDoctorReadyRepositoryFixture(t, repoDir)

	userConfigPath := filepath.Join(homeDir, ".roundfix", "config.yml")
	mustMkdir(t, filepath.Dir(userConfigPath))
	mustWrite(t, userConfigPath, roundconfig.DefaultConfigYAML())
	mustWrite(t, filepath.Join(homeDir, ".roundfix", "state-marker"), "user state\n")
	mustWrite(t, filepath.Join(repoDir, ".roundfixrc.yml"), roundconfig.DefaultConfigYAML())
	mustMkdir(t, filepath.Join(repoDir, ".roundfix"))
	mustWrite(t, filepath.Join(repoDir, ".roundfix", "run-marker"), "repository state\n")

	snapshotRoots := map[string]string{
		"repository":      repoDir,
		"roundfix home":   filepath.Join(homeDir, ".roundfix"),
		"repository data": filepath.Join(repoDir, ".roundfix"),
		"skill state":     filepath.Join(repoDir, ".agents", "skills"),
	}
	beforeRoots := make(map[string]map[string]doctorPathSnapshot, len(snapshotRoots))
	for name, root := range snapshotRoots {
		beforeRoots[name] = snapshotDoctorPath(t, root)
	}
	beforeUserConfig := mustReadBytes(t, userConfigPath)
	lockPath := filepath.Join(repoDir, "skills-lock.json")
	beforeLock := mustReadBytes(t, lockPath)

	checker := newDoctorFakeHealthChecker(
		CheckResult{Name: HealthCheckNode, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckACPX, Status: CheckStatusOK},
		CheckResult{Name: HealthCheckCodex, Status: CheckStatusOK},
	)
	withDoctorFakeLoadedAndReadiness(t, checker, roundconfig.Loaded{
		Config:            roundconfig.Builtin(),
		GitRoot:           repoDir,
		HomeDir:           homeDir,
		UserConfigPath:    userConfigPath,
		ProjectConfigPath: filepath.Join(repoDir, ".roundfixrc.yml"),
	}, func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult {
		return profileProofResult{}
	})
	doctorDeps.checkSkills = skills.CheckRepository
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"doctor"}, &stdout, &stderr)

	if code != exitOK {
		t.Fatalf("exit code = %d, want %d; stdout=%q stderr=%q", code, exitOK, stdout.String(), stderr.String())
	}
	wantStdout := "node: ok\n" +
		"acpx: ok\n" +
		"adapter: ok (codex-acp)\n" +
		"profiles: ok (0 distinct tuples; 0 category references)\n" +
		"skills: ok (39 required: 14 Roundfix-owned, 25 external)\n" +
		"codex: ok\n"
	if got := stdout.String(); got != wantStdout {
		t.Fatalf("unexpected stdout:\n got: %q\nwant: %q", got, wantStdout)
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no stderr, got %q", stderr.String())
	}
	for name, root := range snapshotRoots {
		if after := snapshotDoctorPath(t, root); !reflect.DeepEqual(after, beforeRoots[name]) {
			t.Fatalf("Doctor mutated %s:\nbefore=%#v\nafter=%#v", name, beforeRoots[name], after)
		}
	}
	if after := mustReadBytes(t, userConfigPath); !bytes.Equal(after, beforeUserConfig) {
		t.Fatalf("Doctor mutated User Config:\nbefore=%q\nafter=%q", beforeUserConfig, after)
	}
	if after := mustReadBytes(t, lockPath); !bytes.Equal(after, beforeLock) {
		t.Fatalf("Doctor mutated skills lock:\nbefore=%q\nafter=%q", beforeLock, after)
	}
}

func TestRunDoctorRejectsArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"doctor", "extra"}, &stdout, &stderr)

	if code != exitPreflight {
		t.Fatalf("expected exit code %d, got %d", exitPreflight, code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("expected no stdout, got %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unexpected argument \"extra\"") {
		t.Fatalf("expected argument diagnostic, got %q", stderr.String())
	}
}

func newDoctorFakeHealthChecker(node, acpx, codex CheckResult) *doctorFakeHealthChecker {
	return &doctorFakeHealthChecker{
		node:    node,
		acpx:    acpx,
		adapter: CheckResult{Name: HealthCheckAdapter, Status: CheckStatusOK, Detail: "codex-acp"},
		codex:   codex,
	}
}

type doctorFakeHealthChecker struct {
	node            CheckResult
	acpx            CheckResult
	adapter         CheckResult
	agentResult     CheckResult
	codex           CheckResult
	agentRequests   []agent.ProbeRequest
	adapterRuntimes []agent.RuntimeSpec
	nodeCalls       int
	acpxCalls       int
	adapterCalls    int
	codexCalls      int
	recordCall      func(string)
}

func (checker *doctorFakeHealthChecker) Node(context.Context) CheckResult {
	checker.nodeCalls++
	checker.record(HealthCheckNode)
	return checker.node
}

func (checker *doctorFakeHealthChecker) ACPX(context.Context) CheckResult {
	checker.acpxCalls++
	checker.record(HealthCheckACPX)
	return checker.acpx
}

func (checker *doctorFakeHealthChecker) Adapter(_ context.Context, runtime agent.RuntimeSpec) CheckResult {
	checker.adapterCalls++
	checker.adapterRuntimes = append(checker.adapterRuntimes, runtime)
	checker.record(HealthCheckAdapter)
	return checker.adapter
}

func (checker *doctorFakeHealthChecker) Agent(_ context.Context, req agent.ProbeRequest) CheckResult {
	checker.agentRequests = append(checker.agentRequests, req)
	return checker.agentResult
}

func (checker *doctorFakeHealthChecker) Codex(context.Context) CheckResult {
	checker.codexCalls++
	checker.record(HealthCheckCodex)
	return checker.codex
}

func (checker *doctorFakeHealthChecker) record(name string) {
	if checker.recordCall != nil {
		checker.recordCall(name)
	}
}

type doctorRepositoryLockFixture struct {
	Version int                                         `json:"version"`
	Skills  map[string]doctorRepositoryLockSkillFixture `json:"skills"`
}

type doctorRepositoryLockSkillFixture struct {
	ComputedHash string `json:"computedHash"`
}

func writeDoctorReadyRepositoryFixture(t *testing.T, root string) {
	t.Helper()
	skillsRoot := filepath.Join(root, ".agents", "skills")
	files, err := skills.Files()
	if err != nil {
		t.Fatalf("read embedded Repository Skill Set: %v", err)
	}
	for _, file := range files {
		path := filepath.Join(skillsRoot, filepath.FromSlash(file.Path))
		mustMkdir(t, filepath.Dir(path))
		if err := os.WriteFile(path, file.Data, 0o644); err != nil {
			t.Fatalf("write embedded skill artifact %q: %v", path, err)
		}
	}

	lock := doctorRepositoryLockFixture{
		Version: 1,
		Skills:  make(map[string]doctorRepositoryLockSkillFixture, len(skills.Recommended())),
	}
	for _, name := range skills.Recommended() {
		skillRoot := filepath.Join(skillsRoot, name)
		mustMkdir(t, skillRoot)
		mustWrite(t, filepath.Join(skillRoot, "SKILL.md"), name+"\n")
		hash, err := skills.SkillFolderHash(t.Context(), skillRoot)
		if err != nil {
			t.Fatalf("hash external skill fixture %q: %v", name, err)
		}
		lock.Skills[name] = doctorRepositoryLockSkillFixture{ComputedHash: hash}
	}
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		t.Fatalf("encode Doctor skills lock fixture: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "skills-lock.json"), append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write Doctor skills lock fixture: %v", err)
	}
}

type doctorPathSnapshot struct {
	Mode os.FileMode
	Data string
}

func snapshotDoctorPath(t *testing.T, root string) map[string]doctorPathSnapshot {
	t.Helper()
	snapshot := make(map[string]doctorPathSnapshot)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		item := doctorPathSnapshot{Mode: info.Mode()}
		if info.Mode().IsRegular() {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			item.Data = string(data)
		}
		snapshot[filepath.ToSlash(relative)] = item
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot Doctor path %q: %v", root, err)
	}
	return snapshot
}

func mustReadBytes(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}

func withDoctorFakeDeps(t *testing.T, checker HealthChecker) {
	t.Helper()
	withDoctorFakeLoaded(t, checker, roundconfig.Loaded{
		Config:  roundconfig.Builtin(),
		GitRoot: "/repo/project",
		HomeDir: "/home/roundfix-test",
	})
}

func withDoctorFakeLoaded(t *testing.T, checker HealthChecker, loaded roundconfig.Loaded) {
	withDoctorFakeLoadedAndReadiness(t, checker, loaded, defaultDoctorDependencies().profileReadiness)
}

func withDoctorFakeLoadedAndReadiness(t *testing.T, checker HealthChecker, loaded roundconfig.Loaded, readiness func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult) {
	t.Helper()
	old := doctorDeps
	doctorDeps = doctorDependencies{
		loadConfig: func(roundconfig.LoadOptions) (roundconfig.Loaded, error) {
			return loaded, nil
		},
		healthChecker: func(roundconfig.Loaded) HealthChecker {
			return checker
		},
		profileReadiness: readiness,
		checkSkills: func(context.Context, string) (skills.RepositoryReadiness, error) {
			return skills.RepositoryReadiness{
				OwnedRequired:    14,
				ExternalRequired: 25,
			}, nil
		},
	}
	t.Cleanup(func() {
		doctorDeps = old
	})
}

func withDoctorLiveDeps(t *testing.T, checker HealthChecker) {
	t.Helper()
	old := doctorDeps
	doctorDeps = defaultDoctorDependencies()
	doctorDeps.healthChecker = func(roundconfig.Loaded) HealthChecker { return checker }
	doctorDeps.checkSkills = func(context.Context, string) (skills.RepositoryReadiness, error) {
		return skills.RepositoryReadiness{
			OwnedRequired:    14,
			ExternalRequired: 25,
		}, nil
	}
	t.Cleanup(func() {
		doctorDeps = old
	})
}

func assertDoctorPathMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("expected %s to be absent, got error %v", path, err)
	}
}
