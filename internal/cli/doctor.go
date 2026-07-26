package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	roundconfig "roundfix/internal/config"
	"roundfix/skills"
)

var doctorDeps = defaultDoctorDependencies()

type doctorDependencies struct {
	loadConfig       func(roundconfig.LoadOptions) (roundconfig.Loaded, error)
	healthChecker    func(roundconfig.Loaded) HealthChecker
	profileReadiness func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult
	checkSkills      func(string) (skills.RepositoryReadiness, error)
}

func defaultDoctorDependencies() doctorDependencies {
	return doctorDependencies{
		loadConfig: roundconfig.Load,
		healthChecker: func(roundconfig.Loaded) HealthChecker {
			return setupDeps.healthChecker()
		},
		profileReadiness: func(ctx context.Context, config roundconfig.Config, categories []roundconfig.WorkCategory, workDir string) profileProofResult {
			return proveProfileSelections(ctx, config, categories, workDir, newEngineCollaborators().runner)
		},
		checkSkills: skills.CheckRepository,
	}
}

func runDoctorCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("doctor"))
		return exitOK
	}
	if err := parseDoctorCommand(args); err != nil {
		printDoctorFailure(err, stderr)
		return exitPreflight
	}

	loaded, err := doctorDeps.loadConfig(roundconfig.LoadOptions{Stderr: stderr})
	if err != nil {
		printDoctorFailure(err, stderr)
		return exitRunFailed
	}

	checker := doctorDeps.healthChecker(loaded)
	workDir := strings.TrimSpace(loaded.GitRoot)
	if workDir == "" {
		workDir, err = os.Getwd()
		if err != nil {
			printDoctorFailure(fmt.Errorf("resolve Doctor working directory: %w", err), stderr)
			return exitRunFailed
		}
	}

	profileReadiness := doctorDeps.profileReadiness(ctx, loaded.Config, roundconfig.RequiredWorkCategories(), workDir)
	profileResult := doctorProfileReadinessResult(profileReadiness)
	skillReadiness, skillErr := doctorDeps.checkSkills(workDir)
	skillResult := doctorSkillReadinessResult(skillReadiness, skillErr)
	runtime, runtimeErr := doctorAdapterRuntime(loaded.Config)

	// Keep independent checks eager and ordered.
	results := []CheckResult{
		checker.Node(ctx),
		checker.ACPX(ctx),
		doctorAdapterCheck(ctx, checker, runtime, runtimeErr),
		profileResult,
		skillResult,
		checker.Codex(ctx),
	}

	failed := false
	for _, result := range results {
		printDoctorResult(stdout, result)
		if result.Status == CheckStatusFailed {
			failed = true
		}
	}
	if failed {
		return exitRunFailed
	}
	return exitOK
}

func parseDoctorCommand(args []string) error {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	if err := fs.Parse(args); err != nil {
		return validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	return nil
}

func doctorAdapterRuntime(config roundconfig.Config) (agent.RuntimeSpec, error) {
	resolved, err := roundconfig.ResolveProfile(config, roundconfig.CategoryGeneral, nil)
	if err != nil {
		return agent.RuntimeSpec{}, fmt.Errorf("resolve effective general Agent Selection Profile for adapter readiness: %w", err)
	}
	return runtimeForProfileSelection(resolved.Profile.Preferred)
}

func doctorAdapterCheck(ctx context.Context, checker HealthChecker, runtime agent.RuntimeSpec, runtimeErr error) CheckResult {
	if runtimeErr != nil {
		return CheckResult{
			Name:   HealthCheckAdapter,
			Status: CheckStatusFailed,
			Detail: runtimeErr.Error(),
		}
	}
	return checker.Adapter(ctx, runtime)
}

func doctorProfileReadinessResult(readiness profileProofResult) CheckResult {
	result := CheckResult{Name: HealthCheckProfiles}
	if readiness.Err == nil {
		result.Status = CheckStatusOK
		result.Detail = fmt.Sprintf("%d distinct tuples; %d category references", len(readiness.Proofs), profileProofReferenceCount(readiness.Proofs))
		return result
	}

	result.Status = CheckStatusFailed
	failed, ok := firstFailedProfileProof(readiness.Proofs)
	if !ok {
		result.Detail = readiness.Err.Error()
		return result
	}
	parts := []string{
		fmt.Sprintf("runtime=%q, model=%q, reasoning_effort=%q", failed.Selection.Runtime, failed.Selection.Model, failed.Selection.ReasoningEffort),
		"affected categories: " + formatProfileProofReferences(failed.References),
	}
	if classification := strings.TrimSpace(failed.Classification); classification != "" {
		parts = append(parts, "classification: "+classification)
	}
	parts = append(parts, "adapter evidence: "+doctorProfileAdapterEvidence(failed))
	result.Detail = strings.Join(parts, "; ")
	result.NextAction = strings.TrimSpace(failed.NextAction)
	return result
}

const (
	ownedSkillsNextAction    = "roundfix skills install --target project"
	externalSkillsNextAction = "bunx skills experimental_install && bunx skills update -p -y"
)

func doctorSkillReadinessResult(readiness skills.RepositoryReadiness, checkErr error) CheckResult {
	result := CheckResult{Name: HealthCheckSkills}
	if checkErr == nil && readiness.Ready() {
		result.Status = CheckStatusOK
		result.Detail = fmt.Sprintf(
			"%d required: %d Roundfix-owned, %d external",
			readiness.OwnedRequired+readiness.ExternalRequired,
			readiness.OwnedRequired,
			readiness.ExternalRequired,
		)
		return result
	}

	result.Status = CheckStatusFailed
	missing := append([]string{}, readiness.MissingOwned...)
	missing = append(missing, readiness.MissingExternal...)
	outdated := append([]string{}, readiness.OutdatedOwned...)
	outdated = append(outdated, readiness.OutdatedExternal...)
	sort.Strings(missing)
	sort.Strings(outdated)

	var details []string
	if len(missing) > 0 {
		details = append(details, "missing: "+strings.Join(missing, ", "))
	}
	if len(outdated) > 0 {
		details = append(details, "outdated: "+strings.Join(outdated, ", "))
	}
	if checkErr != nil {
		details = append(details, checkErr.Error())
		result.Err = checkErr
	}
	result.Detail = strings.Join(details, "; ")

	ownedFailure := len(readiness.MissingOwned) > 0 || len(readiness.OutdatedOwned) > 0
	externalFailure := len(readiness.MissingExternal) > 0 || len(readiness.OutdatedExternal) > 0
	var readinessErr *skills.RepositoryReadinessError
	if errors.As(checkErr, &readinessErr) {
		switch readinessErr.Ownership {
		case skills.RepositoryOwnershipOwned:
			ownedFailure = true
		case skills.RepositoryOwnershipExternal:
			externalFailure = true
		}
	}
	var next []string
	if ownedFailure {
		next = append(next, ownedSkillsNextAction)
	}
	if externalFailure {
		next = append(next, externalSkillsNextAction)
	}
	result.NextAction = strings.Join(next, "; ")
	return result
}

func profileProofReferenceCount(proofs []profileProofReport) int {
	count := 0
	for _, proof := range proofs {
		count += len(proof.References)
	}
	return count
}

func firstFailedProfileProof(proofs []profileProofReport) (profileProofReport, bool) {
	for _, proof := range proofs {
		if proof.Status == "failed" {
			return proof, true
		}
	}
	return profileProofReport{}, false
}

func doctorProfileAdapterEvidence(proof profileProofReport) string {
	evidence := make([]string, 0, 5)
	if command := strings.TrimSpace(proof.AdapterCommand); command != "" {
		evidence = append(evidence, fmt.Sprintf("command=%q", command))
	}
	if version := strings.TrimSpace(proof.AdapterVersion); version != "" {
		evidence = append(evidence, fmt.Sprintf("version=%q", version))
	}
	if len(proof.AdvertisedModels) > 0 {
		evidence = append(evidence, "advertised_models="+strings.Join(proof.AdvertisedModels, ","))
	}
	if len(proof.AdvertisedReasoning) > 0 {
		evidence = append(evidence, "advertised_reasoning="+strings.Join(proof.AdvertisedReasoning, ","))
	}
	if detail := strings.TrimSpace(proof.Error); detail != "" {
		evidence = append(evidence, fmt.Sprintf("error=%q", detail))
	}
	if len(evidence) == 0 {
		return "unavailable"
	}
	return strings.Join(evidence, ", ")
}

func printDoctorResult(stdout io.Writer, result CheckResult) {
	detail := strings.TrimSpace(result.Detail)
	if result.Status == CheckStatusFailed {
		nextAction := strings.TrimSpace(result.NextAction)
		if nextAction != "" {
			if detail == "" {
				detail = "next: " + nextAction
			} else {
				detail += "; next: " + nextAction
			}
		}
	}
	if detail == "" {
		fmt.Fprintf(stdout, "%s: %s\n", result.Name, result.Status)
		return
	}
	fmt.Fprintf(stdout, "%s: %s (%s)\n", result.Name, result.Status, detail)
}

func printDoctorFailure(err error, stderr io.Writer) {
	fmt.Fprintf(stderr, "%s: doctor failed: %v\n", app.Name, err)
	fmt.Fprintf(stderr, "Run '%s doctor --help' for usage.\n", app.Name)
}
