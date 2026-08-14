package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"roundfix/internal/agent"
	"roundfix/internal/app"
	"roundfix/internal/baseline"
	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
	"roundfix/skills"
)

type doctorDependencies struct {
	loadConfig       func(roundconfig.LoadOptions) (roundconfig.Loaded, error)
	healthChecker    func(roundconfig.Loaded, string) HealthChecker
	profileReadiness func(context.Context, roundconfig.Config, []roundconfig.WorkCategory, string) profileProofResult
	resolveExternal  func(string) ([]string, bool, error)
	checkSkills      func(context.Context, string, []string) (skills.RepositoryReadiness, error)
	residue          func(context.Context, string) []CheckResult
}

func defaultDoctorDependencies() doctorDependencies {
	return doctorDependencies{
		loadConfig: roundconfig.Load,
		healthChecker: func(_ roundconfig.Loaded, codexPath string) HealthChecker {
			return defaultSetupDependencies().healthChecker(codexPath)
		},
		profileReadiness: func(ctx context.Context, config roundconfig.Config, categories []roundconfig.WorkCategory, workDir string) profileProofResult {
			return proveProfileSelections(ctx, config, categories, workDir, commandDependenciesForContext(ctx).newEngineCollaborators().runner)
		},
		resolveExternal: resolveExternalSkillRequirement,
		checkSkills:     skills.CheckRepositoryWithExternal,
		residue:         defaultDoctorResidueResults,
	}
}

func runDoctorCommand(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	ctx = contextWithCommandDependencies(ctx, environment.dependencies)
	dependencies := commandDependenciesForContext(ctx).doctor
	if commandWantsHelp(args) {
		fmt.Fprint(stdout, commandUsage("doctor"))
		return exitOK
	}
	if err := parseDoctorCommand(args); err != nil {
		printDoctorFailure(err, stderr)
		return exitPreflight
	}

	loadOptions, err := environment.loadOptions(stderr)
	if err != nil {
		printDoctorFailure(err, stderr)
		return exitRunFailed
	}
	loaded, err := dependencies.loadConfig(loadOptions)
	if err != nil {
		printDoctorFailure(err, stderr)
		return exitRunFailed
	}

	checker := dependencies.healthChecker(loaded, environment.codexPath)
	repositoryRoot := strings.TrimSpace(loaded.GitRoot)
	profileWorkDir := repositoryRoot
	if profileWorkDir == "" {
		profileWorkDir, err = environment.resolveWorkDir("resolve process working directory")
		if err != nil {
			printDoctorFailure(err, stderr)
			return exitRunFailed
		}
	}

	// Keep independent checks eager and ordered.
	results := make([]CheckResult, 0, 8)
	results = append(results, checker.Node(ctx))
	results = append(results, checker.ACPX(ctx))
	runtimes, runtimeErr := doctorAdapterRuntimes(loaded.Config)
	results = append(results, doctorAdapterCheck(ctx, checker, runtimes, runtimeErr))
	profileReadiness := dependencies.profileReadiness(ctx, loaded.Config, roundconfig.ConfiguredWorkCategories(loaded.Config), profileWorkDir)
	results = append(results, doctorProfileReadinessResult(profileReadiness))
	if repositoryRoot == "" {
		results = append(results, doctorMissingRepositoryRootResult())
	} else {
		external, manifestOK, requirementErr := dependencies.resolveExternal(repositoryRoot)
		if requirementErr != nil {
			results = append(results, doctorSkillRequirementResult(requirementErr))
		} else {
			skillReadiness, skillErr := dependencies.checkSkills(ctx, repositoryRoot, external)
			if manifestOK {
				results = append(results, doctorSkillReadinessResult(skillReadiness, skillErr))
			} else {
				results = append(results, doctorMissingSetupManifestResult(skillReadiness, skillErr))
			}
		}
	}
	results = append(results, dependencies.residue(ctx, loaded.HomeDir)...)
	results = append(results, checker.Codex(ctx))

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

// RunStore is the read-only Run Database surface needed by Residue.
type RunStore interface {
	ListRuns(context.Context, store.ListRunsQuery) ([]store.Run, error)
}

// ProcessLineage is the recorded ownership proof for one Roundfix spawn tree.
type ProcessLineage struct {
	OwnerPID      int
	OwnerIdentity string
}

// ProcessTable reads details only from a proven Roundfix spawn lineage.
type ProcessTable interface {
	ReadLineage(context.Context, ProcessLineage) ([]store.OwnedProcess, error)
}

// ResidualProcess is a live process whose originating Run is terminal.
type ResidualProcess struct {
	PID     int
	Started time.Time
	CPUTime time.Duration
	RunID   string
	Command string
}

type ownerProcessTable struct {
	controller *store.OwnerProcessControl
}

func (table ownerProcessTable) ReadLineage(ctx context.Context, lineage ProcessLineage) ([]store.OwnedProcess, error) {
	return table.controller.InspectTreeProcesses(ctx, lineage.OwnerPID, lineage.OwnerIdentity)
}

// Residue reports live processes from terminal Roundfix Run lineages. Active
// Runs are excluded before their process lineage is inspected, and the Run
// Database surface is read-only by construction.
func Residue(ctx context.Context, runStore RunStore, table ProcessTable) ([]ResidualProcess, error) {
	runs, err := runStore.ListRuns(ctx, store.ListRunsQuery{States: store.StatesAll})
	if err != nil {
		return nil, fmt.Errorf("read Run Database for process residue: %w", err)
	}

	residueByPID := make(map[int]ResidualProcess)
	var readErrors []error
	for _, run := range runs {
		if !store.IsTerminalState(run.State) || run.OwnerPID == nil || *run.OwnerPID <= 0 {
			continue
		}
		if run.OwnerIdentityUnproven || strings.TrimSpace(run.OwnerIdentity) == "" {
			readErrors = append(readErrors, fmt.Errorf(
				"could not read process table for Run %s: recorded owner identity is unproven",
				run.ID,
			))
			continue
		}
		processes, processErr := table.ReadLineage(ctx, ProcessLineage{
			OwnerPID:      *run.OwnerPID,
			OwnerIdentity: run.OwnerIdentity,
		})
		if processErr != nil {
			readErrors = append(readErrors, fmt.Errorf("could not read process table for Run %s: %w", run.ID, processErr))
		}
		for _, process := range processes {
			if process.PID <= 0 {
				continue
			}
			if _, alreadyReported := residueByPID[process.PID]; alreadyReported {
				continue
			}
			residueByPID[process.PID] = ResidualProcess{
				PID:     process.PID,
				Started: process.Started,
				CPUTime: process.CPUTime,
				RunID:   run.ID,
				Command: process.Command,
			}
		}
	}

	residue := make([]ResidualProcess, 0, len(residueByPID))
	for _, process := range residueByPID {
		residue = append(residue, process)
	}
	sort.Slice(residue, func(i, j int) bool {
		if residue[i].Started.Equal(residue[j].Started) {
			return residue[i].PID < residue[j].PID
		}
		return residue[i].Started.Before(residue[j].Started)
	})
	return residue, errors.Join(readErrors...)
}

func defaultDoctorResidueResults(ctx context.Context, homeDir string) []CheckResult {
	runStore, err := store.OpenReader(ctx, homeDir)
	if err != nil {
		return []CheckResult{{
			Name:   HealthCheckResidue,
			Status: CheckStatusPartial,
			Detail: fmt.Sprintf("could not read Run Database: %v", err),
		}}
	}
	results := doctorResidueResults(
		ctx,
		runStore,
		ownerProcessTable{controller: store.NewOwnerProcessController()},
		time.Now(),
	)
	if err := runStore.Close(); err != nil {
		results = append(results, CheckResult{
			Name:   HealthCheckResidue,
			Status: CheckStatusPartial,
			Detail: fmt.Sprintf("could not close Run Database reader: %v", err),
		})
	}
	return results
}

func doctorResidueResults(ctx context.Context, runStore RunStore, table ProcessTable, now time.Time) []CheckResult {
	processes, err := Residue(ctx, runStore, table)
	results := make([]CheckResult, 0, len(processes)+1)
	for _, process := range processes {
		age := now.Sub(process.Started)
		if age < 0 {
			age = 0
		}
		detail := fmt.Sprintf(
			"PID %d; age %s; CPU %s; originating Run %s; next: inspect the Run and terminate PID %d if it is no longer needed",
			process.PID,
			formatResidueDuration(age),
			formatResidueDuration(process.CPUTime),
			process.RunID,
			process.PID,
		)
		results = append(results, CheckResult{Name: HealthCheckResidue, Status: CheckStatusFound, Detail: detail})
	}
	if err != nil {
		results = append(results, CheckResult{Name: HealthCheckResidue, Status: CheckStatusPartial, Detail: err.Error()})
	}
	if len(results) == 0 {
		results = append(results, CheckResult{
			Name:   HealthCheckResidue,
			Status: CheckStatusOK,
			Detail: "no process residue found",
		})
	}
	return results
}

func formatResidueDuration(duration time.Duration) string {
	if duration <= 0 {
		return "0s"
	}
	return duration.Truncate(time.Second).String()
}

const doctorSetupManifestPath = "docs/agents/setup-context.json"

func resolveExternalSkillRequirement(root string) ([]string, bool, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return nil, false, nil
	}
	manifestPath := filepath.Join(root, filepath.FromSlash(doctorSetupManifestPath))
	info, err := os.Lstat(manifestPath)
	if err != nil || !info.Mode().IsRegular() {
		return nil, false, nil
	}
	manifestBytes, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, false, nil
	}
	var manifest struct {
		Modules []string `json:"modules"`
	}
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, false, nil
	}

	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		return nil, false, fmt.Errorf("resolve external skill requirement: %w", err)
	}
	ownedNames := skills.Names()
	owned := make(map[string]struct{}, len(ownedNames))
	for _, name := range ownedNames {
		owned[name] = struct{}{}
	}
	required := make(map[string]struct{})
	for _, moduleID := range manifest.Modules {
		module, ok := catalog.Module(moduleID)
		if !ok {
			return nil, false, fmt.Errorf(
				"resolve external skill requirement: Setup Manifest names unknown module %q",
				moduleID,
			)
		}
		var declaration struct {
			RequiredSkills []string `json:"requiredSkills"`
		}
		if err := json.Unmarshal(module.Data, &declaration); err != nil {
			return nil, false, fmt.Errorf(
				"resolve external skill requirement for module %q: %w",
				moduleID,
				err,
			)
		}
		for _, name := range declaration.RequiredSkills {
			if _, isOwned := owned[name]; !isOwned {
				required[name] = struct{}{}
			}
		}
	}

	external := make([]string, 0, len(required))
	for name := range required {
		external = append(external, name)
	}
	sort.Strings(external)
	return external, true, nil
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

func doctorAdapterRuntimes(config roundconfig.Config) ([]agent.RuntimeSpec, error) {
	proofs, err := buildProfileProofReports(config, roundconfig.ConfiguredWorkCategories(config))
	if err != nil {
		return nil, fmt.Errorf("resolve effective configured Agent Selection Profiles for adapter readiness: %w", err)
	}

	selectionsByRuntime := make(map[string]roundconfig.AgentSelection)
	for _, proof := range proofs {
		runtimeID := strings.TrimSpace(proof.Selection.Runtime)
		if _, exists := selectionsByRuntime[runtimeID]; exists {
			continue
		}
		selection := proof.Selection
		selection.Runtime = runtimeID
		selectionsByRuntime[runtimeID] = selection
	}

	runtimeIDs := make([]string, 0, len(selectionsByRuntime))
	for runtimeID := range selectionsByRuntime {
		runtimeIDs = append(runtimeIDs, runtimeID)
	}
	sort.Strings(runtimeIDs)

	runtimes := make([]agent.RuntimeSpec, 0, len(runtimeIDs))
	for _, runtimeID := range runtimeIDs {
		runtime, err := runtimeForProfileSelection(selectionsByRuntime[runtimeID])
		if err != nil {
			return nil, fmt.Errorf("resolve ACP Runtime %q for adapter readiness: %w", runtimeID, err)
		}
		runtimes = append(runtimes, runtime)
	}
	return runtimes, nil
}

func doctorAdapterCheck(ctx context.Context, checker HealthChecker, runtimes []agent.RuntimeSpec, runtimeErr error) CheckResult {
	if runtimeErr != nil {
		return CheckResult{
			Name:   HealthCheckAdapter,
			Status: CheckStatusFailed,
			Detail: runtimeErr.Error(),
			Err:    runtimeErr,
		}
	}
	if len(runtimes) == 0 {
		err := errors.New("effective required Agent Selection Profiles reference no ACP Runtime")
		return CheckResult{
			Name:   HealthCheckAdapter,
			Status: CheckStatusFailed,
			Detail: err.Error(),
			Err:    err,
		}
	}

	result := CheckResult{
		Name:   HealthCheckAdapter,
		Status: CheckStatusOK,
	}
	details := make([]string, 0, len(runtimes))
	nextActions := make([]string, 0, len(runtimes))
	for _, runtime := range runtimes {
		runtimeResult := checker.Adapter(ctx, runtime)
		detail := strings.TrimSpace(runtimeResult.Detail)
		if detail == "" {
			detail = string(runtimeResult.Status)
		}
		if runtimeResult.Status == CheckStatusFailed {
			result.Status = CheckStatusFailed
			if classification := doctorAdapterFailureClassification(runtimeResult.Err); classification != "" {
				detail += "; classification: " + classification
			}
			if nextAction := strings.TrimSpace(runtimeResult.NextAction); nextAction != "" {
				nextActions = append(nextActions, nextAction)
			}
			if runtimeResult.Err != nil {
				if result.Err == nil {
					result.Err = runtimeResult.Err
				} else {
					result.Err = errors.Join(result.Err, runtimeResult.Err)
				}
			}
		}
		details = append(details, runtime.ID+": "+detail)
	}
	result.Detail = strings.Join(details, " | ")
	result.NextAction = strings.Join(nextActions, " && ")
	return result
}

func doctorAdapterFailureClassification(err error) string {
	var classified interface {
		error
		Classification() string
	}
	if !errors.As(err, &classified) {
		return ""
	}
	return strings.TrimSpace(classified.Classification())
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
	ownedSkillsNextAction      = skills.OwnedSkillsUpgradeAction
	externalSkillsNextAction   = "bunx skills experimental_install && bunx skills update -p -y"
	baselineAdoptionNextAction = "roundfix baseline"
	missingGitRootNextAction   = "run roundfix doctor from a Git repository"
	missingSetupManifestDetail = "Setup Manifest is absent or unreadable"
)

func doctorMissingRepositoryRootResult() CheckResult {
	return CheckResult{
		Name:       HealthCheckSkills,
		Status:     CheckStatusFailed,
		Detail:     "Repository Skill Set readiness requires a Git repository",
		NextAction: missingGitRootNextAction,
	}
}

func doctorSkillRequirementResult(err error) CheckResult {
	return CheckResult{
		Name:   HealthCheckSkills,
		Status: CheckStatusFailed,
		Detail: err.Error(),
		Err:    err,
	}
}

func doctorMissingSetupManifestResult(readiness skills.RepositoryReadiness, checkErr error) CheckResult {
	result := doctorSkillReadinessResult(readiness, checkErr)
	result.Status = CheckStatusFailed
	if result.Detail == "" {
		result.Detail = missingSetupManifestDetail
	} else {
		result.Detail += "; " + missingSetupManifestDetail
	}
	if !(checkErr == nil && readiness.Ready()) {
		result.Detail += "; 0 external required"
	}

	next := make([]string, 0, 2)
	if doctorOwnedSkillRequirementFailed(readiness, checkErr) {
		next = append(next, ownedSkillsNextAction)
	}
	next = append(next, baselineAdoptionNextAction)
	result.NextAction = strings.Join(next, " && ")
	return result
}

func doctorOwnedSkillRequirementFailed(readiness skills.RepositoryReadiness, checkErr error) bool {
	if len(readiness.MissingOwned) > 0 || len(ownedSkillReadinessWithState(readiness, skills.ReadinessBelow)) > 0 {
		return true
	}
	if checkErr == nil {
		return false
	}
	var readinessErr *skills.RepositoryReadinessError
	if !errors.As(checkErr, &readinessErr) {
		return true
	}
	return readinessErr.Ownership == "" || readinessErr.Ownership == skills.RepositoryOwnershipOwned
}

func doctorSkillReadinessResult(readiness skills.RepositoryReadiness, checkErr error) CheckResult {
	result := CheckResult{Name: HealthCheckSkills}
	if checkErr == nil && readiness.Ready() {
		unversioned := ownedSkillReadinessWithState(readiness, skills.ReadinessUnversioned)
		result.Status = CheckStatusOK
		if len(unversioned) > 0 {
			result.Status = CheckStatusUnversioned
		}
		result.Detail = fmt.Sprintf(
			"%d required: %d Roundfix-owned, %d external",
			readiness.OwnedRequired+readiness.ExternalRequired,
			readiness.OwnedRequired,
			readiness.ExternalRequired,
		)
		if len(unversioned) > 0 {
			result.Detail += "; unversioned: " + strings.Join(ownedSkillNames(unversioned), ", ")
		}
		return result
	}

	result.Status = CheckStatusFailed
	missing := append([]string{}, readiness.MissingOwned...)
	missing = append(missing, readiness.MissingExternal...)
	outdated := append([]string{}, readiness.OutdatedExternal...)
	sort.Strings(missing)
	sort.Strings(outdated)

	var details []string
	if len(missing) > 0 {
		details = append(details, "missing: "+strings.Join(missing, ", "))
	}
	if len(outdated) > 0 {
		details = append(details, "outdated: "+strings.Join(outdated, ", "))
	}
	below := ownedSkillReadinessWithState(readiness, skills.ReadinessBelow)
	for _, owned := range below {
		details = append(details, owned.ComparisonDetail())
	}
	unversioned := ownedSkillReadinessWithState(readiness, skills.ReadinessUnversioned)
	if len(unversioned) > 0 {
		details = append(details, "unversioned: "+strings.Join(ownedSkillNames(unversioned), ", "))
	}
	if checkErr != nil {
		details = append(details, checkErr.Error())
		result.Err = checkErr
	}
	result.Detail = strings.Join(details, "; ")

	ownedFailure := len(readiness.MissingOwned) > 0 || len(below) > 0
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
	if checkErr != nil && !ownedFailure && !externalFailure {
		ownedFailure = true
		externalFailure = true
	}
	var next []string
	if ownedFailure {
		next = append(next, ownedSkillsNextAction)
	}
	if externalFailure {
		fromError := []string(nil)
		if readinessErr != nil {
			fromError = readinessErr.MissingExternal
		}
		externalActions := doctorExternalSkillNextActions(readiness, fromError)
		if len(externalActions) == 0 {
			externalActions = append(externalActions, externalSkillsNextAction)
		}
		next = append(next, externalActions...)
	}
	result.NextAction = strings.Join(next, " && ")
	return result
}

func ownedSkillReadinessWithState(readiness skills.RepositoryReadiness, state skills.ReadinessState) []skills.OwnedSkillReadiness {
	matched := make([]skills.OwnedSkillReadiness, 0, len(readiness.Owned))
	for _, owned := range readiness.Owned {
		if owned.State == state {
			matched = append(matched, owned)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Skill < matched[j].Skill
	})
	return matched
}

func ownedSkillNames(readiness []skills.OwnedSkillReadiness) []string {
	names := make([]string, 0, len(readiness))
	for _, owned := range readiness {
		names = append(names, owned.Skill)
	}
	return names
}

// doctorExternalSkillNextActions renders one install command per external
// skill. A lock-entry failure reports its gap through the readiness error
// rather than the readiness struct, so both sources feed the same remediation.
func doctorExternalSkillNextActions(readiness skills.RepositoryReadiness, fromError []string) []string {
	names := append([]string{}, readiness.MissingExternal...)
	names = append(names, readiness.OutdatedExternal...)
	names = append(names, fromError...)
	sort.Strings(names)

	actions := make([]string, 0, len(names))
	previous := ""
	for _, name := range names {
		if name == previous {
			continue
		}
		actions = append(actions, "bunx skills add marcioaltoe/skills@"+name)
		previous = name
	}
	return actions
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
