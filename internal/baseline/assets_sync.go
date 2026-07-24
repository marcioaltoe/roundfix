package baseline

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const AssetsSyncSchemaVersion = "setup-context-driven/audit-v1"

// AssetsSyncRequest selects a read-only drift check or canonical asset refresh.
type AssetsSyncRequest struct {
	SourceDir string
	Check     bool
}

// AssetsSyncSummary is the maintained finding count projection.
type AssetsSyncSummary struct {
	Errors    int `json:"errors"`
	Decisions int `json:"decisions"`
	Warnings  int `json:"warnings"`
	Info      int `json:"info"`
}

// AssetsSyncFinding is one stable source, drift, or refresh outcome.
type AssetsSyncFinding struct {
	Code      string `json:"code"`
	Severity  string `json:"severity"`
	Path      string `json:"path"`
	ManagedID string `json:"managedId"`
	Message   string `json:"message"`
	Action    string `json:"action"`
}

// AssetsSyncChange preserves the maintained audit result shape.
type AssetsSyncChange struct {
	Action    string `json:"action"`
	Path      string `json:"path"`
	ManagedID string `json:"managedId"`
}

// AssetsSyncPayload is the stable text and JSON command result.
type AssetsSyncPayload struct {
	SchemaVersion  string              `json:"schemaVersion"`
	OK             bool                `json:"ok"`
	Summary        AssetsSyncSummary   `json:"summary"`
	Findings       []AssetsSyncFinding `json:"findings"`
	PlannedChanges []AssetsSyncChange  `json:"plannedChanges"`
}

// AssetsSyncErrorCategory maps synchronization failures to Baseline CLI exits.
type AssetsSyncErrorCategory string

const (
	AssetsSyncInvalid   AssetsSyncErrorCategory = "invalid"
	AssetsSyncExecution AssetsSyncErrorCategory = "execution"
)

// AssetsSyncError is one typed synchronization refusal or failure.
type AssetsSyncError struct {
	Category AssetsSyncErrorCategory
	Finding  AssetsSyncFinding
	Err      error
}

func (err *AssetsSyncError) Error() string {
	return err.Finding.Message
}

func (err *AssetsSyncError) Unwrap() error {
	return err.Err
}

type assetsSyncDependencies struct {
	repository      string
	assetRoot       string
	transactionHook func(transactionFaultPoint) error
}

type assetsSyncSnapshot struct {
	SchemaVersion     string             `json:"schemaVersion"`
	ID                string             `json:"id"`
	Version           string             `json:"version"`
	Source            assetsSyncSource   `json:"source"`
	Digest            string             `json:"digest"`
	Skills            []assetsSyncSkill  `json:"skills"`
	ActivationBundles []assetsSyncBundle `json:"activationBundles,omitempty"`
}

type assetsSyncSource struct {
	Type       string `json:"type"`
	Repository string `json:"repository,omitempty"`
	Ref        string `json:"ref,omitempty"`
	Path       string `json:"path,omitempty"`
	Name       string `json:"name,omitempty"`
}

type assetsSyncSkill struct {
	Name          string           `json:"name"`
	Path          string           `json:"path"`
	Source        assetsSyncSource `json:"source"`
	TreeDigest    string           `json:"treeDigest,omitempty"`
	ContentDigest string           `json:"contentDigest,omitempty"`
}

type assetsSyncBundle struct {
	ID     string   `json:"id"`
	Skills []string `json:"skills"`
}

type assetsSyncCheckout struct {
	root       string
	repository string
	revision   string
}

type assetsSyncPlan struct {
	path    string
	setupID string
	content []byte
}

var assetsSyncRepoOwnedSkills = map[string]struct{}{
	"archive-spec":         {},
	"brainstorming":        {},
	"business-analyst":     {},
	"council":              {},
	"evidence-gate":        {},
	"implement-spec":       {},
	"implement-task":       {},
	"qa-gate":              {},
	"roundfix":             {},
	"setup-context-driven": {},
	"write-idea":           {},
	"write-prd":            {},
	"write-tasks":          {},
	"write-techspec":       {},
}

// SyncAssets checks or refreshes the Go-owned canonical Baseline setup
// snapshots from an explicit clean Git checkout.
func SyncAssets(ctx context.Context, request AssetsSyncRequest) (AssetsSyncPayload, error) {
	if ctx == nil {
		return failedAssetsSyncPayload(assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			"assets/setups",
			"setups",
			"Asset synchronization requires a live context.",
			"Rerun roundfix baseline assets sync.",
			errors.New("context is required"),
		))
	}
	root, _, err := inspectRepositoryIdentity(ctx, ".", nil)
	if err != nil {
		return failedAssetsSyncPayload(assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			"internal/baseline/assets",
			"setups",
			fmt.Sprintf("Canonical Baseline assets require the Roundfix Git worktree: %v.", err),
			"Run the command from the Roundfix source checkout.",
			err,
		))
	}
	return syncAssets(ctx, request, assetsSyncDependencies{
		repository: root,
		assetRoot:  filepath.Join(root, "internal", "baseline", "assets"),
	})
}

func syncAssets(
	ctx context.Context,
	request AssetsSyncRequest,
	dependencies assetsSyncDependencies,
) (AssetsSyncPayload, error) {
	if ctx == nil {
		return failedAssetsSyncPayload(assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			"assets/setups",
			"setups",
			"Asset synchronization requires a live context.",
			"Rerun roundfix baseline assets sync.",
			errors.New("context is required"),
		))
	}
	if err := ctx.Err(); err != nil {
		return emptyAssetsSyncPayload(), err
	}
	request.SourceDir = strings.TrimSpace(request.SourceDir)
	if request.SourceDir == "" {
		return failedAssetsSyncPayload(assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			"assets/setups",
			"setups",
			"Canonical setups directory is required.",
			"Pass --source-dir with an explicit canonical setups directory.",
			errors.New("source directory is required"),
		))
	}
	sourceDir, err := filepath.Abs(request.SourceDir)
	if err != nil {
		return failedAssetsSyncPayload(assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			request.SourceDir,
			"setups",
			fmt.Sprintf("Canonical setups directory cannot be resolved: %v.", err),
			"Pass a readable canonical setups directory.",
			err,
		))
	}
	sourceDir = filepath.Clean(sourceDir)
	if resolved, resolveErr := filepath.EvalSymlinks(sourceDir); resolveErr == nil {
		sourceDir = resolved
	}
	if info, statErr := os.Stat(sourceDir); statErr != nil || !info.IsDir() {
		if statErr == nil {
			statErr = errors.New("path is not a directory")
		}
		return failedAssetsSyncPayload(assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			filepath.ToSlash(sourceDir),
			"setups",
			fmt.Sprintf("Canonical setups directory is not a directory: %s.", sourceDir),
			"Pass a readable canonical setups directory.",
			statErr,
		))
	}
	if strings.TrimSpace(dependencies.repository) == "" ||
		strings.TrimSpace(dependencies.assetRoot) == "" {
		return failedAssetsSyncPayload(assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			"internal/baseline/assets",
			"setups",
			"Canonical Baseline asset target is unavailable.",
			"Run the command from the Roundfix source checkout.",
			errors.New("asset synchronization target is required"),
		))
	}
	if resolved, resolveErr := filepath.EvalSymlinks(dependencies.repository); resolveErr == nil {
		dependencies.repository = resolved
	}
	if resolved, resolveErr := filepath.EvalSymlinks(dependencies.assetRoot); resolveErr == nil {
		dependencies.assetRoot = resolved
	}
	if _, err := LoadCatalog(os.DirFS(dependencies.assetRoot)); err != nil {
		return failedAssetsSyncPayload(assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			"internal/baseline/assets",
			"setups",
			fmt.Sprintf("Go-owned canonical Baseline assets are invalid: %v.", err),
			"Fix the canonical Baseline catalog before synchronizing.",
			err,
		))
	}

	current, err := loadAssetsSyncSnapshots(dependencies.assetRoot)
	if err != nil {
		return failedAssetsSyncPayload(err)
	}
	checkout, err := inspectAssetsSyncCheckout(ctx, sourceDir)
	if err != nil {
		return failedAssetsSyncPayload(err)
	}

	findings := []AssetsSyncFinding{}
	plans := []assetsSyncPlan{}
	overrides := map[string][]byte{}
	for _, setupID := range sortedMapKeys(current) {
		if err := ctx.Err(); err != nil {
			return assetsSyncPayload(findings), err
		}
		snapshot, snapshotFindings := buildAssetsSyncSnapshot(
			ctx,
			setupID,
			sourceDir,
			current[setupID],
			checkout,
		)
		findings = append(findings, snapshotFindings...)
		if snapshot == nil {
			continue
		}
		data, marshalErr := json.MarshalIndent(snapshot, "", "  ")
		if marshalErr != nil {
			findings = append(findings, assetsSyncInvalidFinding(
				filepath.Join(sourceDir, setupID+".json"),
				setupID,
				fmt.Sprintf("Canonical setup snapshot cannot be serialized: %v.", marshalErr),
			))
			continue
		}
		data = append(data, '\n')
		currentPath := filepath.Join(dependencies.assetRoot, "setups", setupID+".json")
		currentBytes, readErr := os.ReadFile(currentPath)
		if readErr != nil {
			findings = append(findings, assetsSyncInvalidFinding(
				currentPath,
				setupID,
				fmt.Sprintf("Bundled setup snapshot cannot be read: %v.", readErr),
			))
			continue
		}
		if bytes.Equal(data, currentBytes) {
			continue
		}
		assetPath := path.Join("setups", setupID+".json")
		overrides[assetPath] = append([]byte(nil), data...)
		plans = append(plans, assetsSyncPlan{
			path:    assetPath,
			setupID: setupID,
			content: data,
		})
	}
	if hasAssetsSyncErrors(findings) {
		return failedAssetsSyncFindings(findings, AssetsSyncInvalid)
	}
	if len(overrides) != 0 {
		if _, err := LoadCatalog(assetsOverlayFS{
			base:      os.DirFS(dependencies.assetRoot),
			overrides: overrides,
		}); err != nil {
			finding := assetsSyncInvalidFinding(
				dependencies.assetRoot,
				"setups",
				fmt.Sprintf("Generated setup snapshots are incompatible with the Baseline catalog: %v.", err),
			)
			return failedAssetsSyncFindings([]AssetsSyncFinding{finding}, AssetsSyncInvalid)
		}
	}
	if request.Check {
		for _, plan := range plans {
			findings = append(findings, assetsSyncDriftFinding(plan.setupID))
		}
		if len(findings) != 0 {
			return failedAssetsSyncFindings(findings, AssetsSyncExecution)
		}
		return emptyAssetsSyncPayload(), nil
	}
	if len(plans) == 0 {
		return emptyAssetsSyncPayload(), nil
	}
	if err := applyAssetsSyncPlans(ctx, dependencies, plans); err != nil {
		finding := AssetsSyncFinding{
			Code:      "skills.setup-snapshot.drift",
			Severity:  "error",
			Path:      "assets/setups",
			ManagedID: "setups",
			Message:   fmt.Sprintf("Could not update setup snapshots atomically: %v.", err),
			Action:    "Fix filesystem permissions and rerun roundfix baseline assets sync.",
		}
		return failedAssetsSyncFindings([]AssetsSyncFinding{finding}, AssetsSyncExecution)
	}
	for _, plan := range plans {
		findings = append(findings, AssetsSyncFinding{
			Code:      "skills.setup-snapshot.updated",
			Severity:  "info",
			Path:      "assets/setups/" + plan.setupID + ".json",
			ManagedID: "setup." + plan.setupID,
			Message:   "Setup snapshot was synchronized from the canonical source.",
			Action:    "Review the snapshot diff and run asset validation.",
		})
	}
	return assetsSyncPayload(findings), nil
}

func loadAssetsSyncSnapshots(assetRoot string) (map[string]assetsSyncSnapshot, error) {
	paths, err := filepath.Glob(filepath.Join(assetRoot, "setups", "*.json"))
	if err != nil {
		return nil, assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			"assets/setups",
			"setups",
			fmt.Sprintf("Bundled setup snapshots cannot be listed: %v.", err),
			"Fix the bundled setup snapshots before synchronizing.",
			err,
		)
	}
	if len(paths) == 0 {
		return nil, assetsSyncError(
			AssetsSyncInvalid,
			"skills.setup-snapshot.drift",
			"error",
			"assets/setups",
			"setups",
			"Bundled setup snapshots are missing.",
			"Restore the Go-owned canonical setup snapshots before synchronizing.",
			errors.New("setup snapshot collection is empty"),
		)
	}
	result := make(map[string]assetsSyncSnapshot, len(paths))
	for _, snapshotPath := range paths {
		data, readErr := os.ReadFile(snapshotPath)
		if readErr != nil {
			return nil, assetsSyncError(
				AssetsSyncInvalid,
				"skills.setup-snapshot.drift",
				"error",
				filepath.ToSlash(snapshotPath),
				strings.TrimSuffix(filepath.Base(snapshotPath), ".json"),
				fmt.Sprintf("Bundled setup snapshot cannot be read: %v.", readErr),
				"Fix the bundled setup snapshot before synchronizing.",
				readErr,
			)
		}
		var snapshot assetsSyncSnapshot
		if err := json.Unmarshal(data, &snapshot); err != nil {
			return nil, assetsSyncError(
				AssetsSyncInvalid,
				"skills.setup-snapshot.drift",
				"error",
				filepath.ToSlash(snapshotPath),
				strings.TrimSuffix(filepath.Base(snapshotPath), ".json"),
				fmt.Sprintf("Bundled setup snapshot is not valid JSON: %v.", err),
				"Fix the bundled setup snapshot before synchronizing.",
				err,
			)
		}
		setupID := strings.TrimSuffix(filepath.Base(snapshotPath), ".json")
		if snapshot.ID != setupID {
			return nil, assetsSyncError(
				AssetsSyncInvalid,
				"skills.setup-snapshot.drift",
				"error",
				filepath.ToSlash(snapshotPath),
				setupID,
				"Bundled setup snapshot id does not match its filename.",
				"Fix the bundled setup snapshot before synchronizing.",
				errors.New("setup snapshot id mismatch"),
			)
		}
		result[setupID] = snapshot
	}
	return result, nil
}

func buildAssetsSyncSnapshot(
	ctx context.Context,
	setupID, sourceDir string,
	current assetsSyncSnapshot,
	checkout assetsSyncCheckout,
) (*assetsSyncSnapshot, []AssetsSyncFinding) {
	sourceDocument, sourcePath, err := loadAssetsSyncSourceDocument(setupID, sourceDir)
	if err != nil {
		return nil, []AssetsSyncFinding{assetsSyncInvalidFinding(sourcePath, setupID, err.Error())}
	}
	if err := verifyAssetsSyncCommittedFile(ctx, checkout, sourcePath); err != nil {
		return nil, []AssetsSyncFinding{assetsSyncInvalidFinding(sourcePath, setupID, err.Error())}
	}
	if declared, ok := sourceDocument["source"]; ok {
		relative, relativeErr := filepath.Rel(checkout.root, sourcePath)
		if relativeErr != nil {
			return nil, []AssetsSyncFinding{assetsSyncInvalidFinding(sourcePath, setupID, "Canonical setup source file is outside the Git checkout.")}
		}
		if message := validateAssetsSyncDeclaredSource(
			ctx,
			declared,
			checkout,
			filepath.ToSlash(relative),
		); message != "" {
			return nil, []AssetsSyncFinding{assetsSyncInvalidFinding(sourcePath, setupID, message)}
		}
	}

	currentByPath := make(map[string]assetsSyncSkill, len(current.Skills))
	currentByName := make(map[string]assetsSyncSkill, len(current.Skills))
	for _, skill := range current.Skills {
		currentByPath[skill.Path] = skill
		currentByName[skill.Name] = skill
	}
	rawSkills, ok := sourceDocument["skills"].([]any)
	if !ok || len(rawSkills) == 0 {
		return nil, []AssetsSyncFinding{assetsSyncInvalidFinding(
			sourcePath,
			setupID,
			"Canonical setup must contain a non-empty skills list.",
		)}
	}
	skills := make([]assetsSyncSkill, 0, len(rawSkills))
	findings := []AssetsSyncFinding{}
	seenNames := map[string]struct{}{}
	seenPaths := map[string]struct{}{}
	for _, rawSkill := range rawSkills {
		skill, finding := normalizeAssetsSyncSkill(
			ctx,
			rawSkill,
			sourceDir,
			sourcePath,
			checkout,
			currentByPath,
			currentByName,
		)
		if finding != nil {
			findings = append(findings, *finding)
			continue
		}
		if _, duplicate := seenNames[skill.Name]; duplicate {
			findings = append(findings, assetsSyncInvalidFinding(
				sourcePath,
				setupID,
				fmt.Sprintf("Canonical setup contains duplicate skill name %s.", skill.Name),
			))
		}
		if _, duplicate := seenPaths[skill.Path]; duplicate {
			findings = append(findings, assetsSyncInvalidFinding(
				sourcePath,
				setupID,
				fmt.Sprintf("Canonical setup contains duplicate skill path %s.", skill.Path),
			))
		}
		seenNames[skill.Name] = struct{}{}
		seenPaths[skill.Path] = struct{}{}
		skills = append(skills, *skill)
	}
	if len(findings) != 0 {
		return nil, findings
	}
	relative, err := filepath.Rel(checkout.root, sourcePath)
	if err != nil || !safeRelative(filepath.ToSlash(relative)) {
		return nil, []AssetsSyncFinding{assetsSyncInvalidFinding(
			sourcePath,
			setupID,
			"Canonical setup source file is outside the Git checkout.",
		)}
	}
	normalizedSkills := make([]any, len(skills))
	for index, skill := range skills {
		normalizedSkills[index] = assetsSyncSkillDocument(skill)
	}
	digestPayload := any(normalizedSkills)
	if current.ActivationBundles != nil {
		normalizedBundles := make([]any, len(current.ActivationBundles))
		for index, bundle := range current.ActivationBundles {
			normalizedBundles[index] = map[string]any{
				"id":     bundle.ID,
				"skills": bundle.Skills,
			}
		}
		digestPayload = map[string]any{
			"activationBundles": normalizedBundles,
			"skills":            normalizedSkills,
		}
	}
	digest, err := canonicalSHA256(digestPayload)
	if err != nil {
		return nil, []AssetsSyncFinding{assetsSyncInvalidFinding(
			sourcePath,
			setupID,
			fmt.Sprintf("Canonical setup digest cannot be calculated: %v.", err),
		)}
	}
	return &assetsSyncSnapshot{
		SchemaVersion: "setup-context-driven/setup-snapshot/0.0.1",
		ID:            setupID,
		Version:       "0.0.1",
		Source: assetsSyncSource{
			Type:       "github",
			Repository: checkout.repository,
			Ref:        checkout.revision,
			Path:       filepath.ToSlash(relative),
		},
		Digest:            digest,
		Skills:            skills,
		ActivationBundles: append([]assetsSyncBundle(nil), current.ActivationBundles...),
	}, nil
}

func assetsSyncSkillDocument(skill assetsSyncSkill) map[string]any {
	source := map[string]any{"type": skill.Source.Type}
	if skill.Source.Repository != "" {
		source["repository"] = skill.Source.Repository
	}
	if skill.Source.Ref != "" {
		source["ref"] = skill.Source.Ref
	}
	if skill.Source.Path != "" {
		source["path"] = skill.Source.Path
	}
	if skill.Source.Name != "" {
		source["name"] = skill.Source.Name
	}
	document := map[string]any{
		"name":   skill.Name,
		"path":   skill.Path,
		"source": source,
	}
	if skill.TreeDigest != "" {
		document["treeDigest"] = skill.TreeDigest
	}
	if skill.ContentDigest != "" {
		document["contentDigest"] = skill.ContentDigest
	}
	return document
}

func loadAssetsSyncSourceDocument(setupID, sourceDir string) (map[string]any, string, error) {
	jsonPath := filepath.Join(sourceDir, setupID+".json")
	if info, err := os.Stat(jsonPath); err == nil && info.Mode().IsRegular() {
		data, readErr := os.ReadFile(jsonPath)
		if readErr != nil {
			return nil, jsonPath, fmt.Errorf("Canonical setup JSON cannot be read: %w", readErr)
		}
		var document map[string]any
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.UseNumber()
		if err := decoder.Decode(&document); err != nil {
			return nil, jsonPath, fmt.Errorf("Canonical setup JSON is invalid: %w", err)
		}
		if err := requireJSONEOF(decoder); err != nil {
			return nil, jsonPath, fmt.Errorf("Canonical setup JSON is invalid: %w", err)
		}
		if document == nil {
			return nil, jsonPath, errors.New("Canonical setup JSON must be an object.")
		}
		return document, jsonPath, nil
	}
	for _, extension := range []string{".txt", ".md"} {
		sourcePath := filepath.Join(sourceDir, setupID+extension)
		data, err := os.ReadFile(sourcePath)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, sourcePath, fmt.Errorf("Canonical setup source cannot be read: %w", err)
		}
		skills := []any{}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			skills = append(skills, map[string]any{"path": line})
		}
		return map[string]any{"skills": skills}, sourcePath, nil
	}
	return nil, jsonPath, errors.New("Canonical setup source file is missing.")
}

func normalizeAssetsSyncSkill(
	ctx context.Context,
	raw any,
	sourceDir, sourcePath string,
	checkout assetsSyncCheckout,
	currentByPath, currentByName map[string]assetsSyncSkill,
) (*assetsSyncSkill, *AssetsSyncFinding) {
	var data map[string]any
	switch value := raw.(type) {
	case string:
		data = map[string]any{"path": value}
	case map[string]any:
		data = value
	default:
		finding := assetsSyncInvalidFinding(
			sourcePath,
			"setup",
			"Each canonical setup skill must be a string or object.",
		)
		return nil, &finding
	}
	rawPath := firstAssetsSyncString(data, "path", "skillPath", "name")
	if strings.TrimSpace(rawPath) == "" {
		finding := assetsSyncInvalidFinding(
			sourcePath,
			"setup",
			"Canonical setup skill is missing a path.",
		)
		return nil, &finding
	}
	normalizedPath := normalizeAssetsSyncSkillPath(rawPath)
	if normalizedPath == "" {
		finding := assetsSyncInvalidFinding(
			sourcePath,
			"setup",
			fmt.Sprintf("Canonical setup skill path %q is not portable.", rawPath),
		)
		return nil, &finding
	}
	name, _ := data["name"].(string)
	name = strings.TrimSpace(name)
	if name == "" {
		name = inferAssetsSyncSkillName(normalizedPath)
	}
	current := currentByPath[normalizedPath]
	if current.Name == "" {
		current = currentByName[name]
	}
	_, repoOwnedByName := assetsSyncRepoOwnedSkills[name]
	if repoOwnedByName || current.Source.Type == "repo" {
		digest := current.ContentDigest
		if !lowercaseSHA256.MatchString(digest) {
			digest, _ = data["contentDigest"].(string)
		}
		if !lowercaseSHA256.MatchString(digest) {
			finding := assetsSyncInvalidFinding(
				sourcePath,
				name,
				fmt.Sprintf("Repository-owned skill %s is missing a valid content digest.", name),
			)
			return nil, &finding
		}
		return &assetsSyncSkill{
			Name: name,
			Path: normalizedPath,
			Source: assetsSyncSource{
				Type: "repo",
				Name: "roundfix",
			},
			ContentDigest: digest,
		}, nil
	}

	effectiveRevision := checkout.revision
	if declared, ok := data["source"]; ok {
		if message := validateAssetsSyncDeclaredSource(
			ctx,
			declared,
			checkout,
			normalizedPath,
		); message != "" {
			finding := assetsSyncInvalidFinding(sourcePath, name, message)
			return nil, &finding
		}
		effectiveRevision = declared.(map[string]any)["ref"].(string)
	}
	treeRoot := resolveAssetsSyncSkillTree(sourceDir, rawPath)
	if treeRoot == "" {
		finding := assetsSyncInvalidFinding(
			sourcePath,
			name,
			fmt.Sprintf("Canonical setup skill %s source directory is missing.", name),
		)
		return nil, &finding
	}
	relative, err := filepath.Rel(checkout.root, treeRoot)
	if err != nil || filepath.ToSlash(relative) != normalizedPath {
		finding := assetsSyncInvalidFinding(
			sourcePath,
			name,
			fmt.Sprintf("Canonical setup skill %s source path does not match %s.", name, normalizedPath),
		)
		return nil, &finding
	}
	workingDigest, err := assetsSyncFilesystemTreeDigest(treeRoot)
	if err != nil {
		finding := assetsSyncInvalidFinding(
			sourcePath,
			name,
			fmt.Sprintf("Canonical setup skill %s cannot be proven: %v.", name, err),
		)
		return nil, &finding
	}
	committedDigest, err := assetsSyncCommittedTreeDigest(
		ctx,
		checkout.root,
		effectiveRevision,
		normalizedPath,
	)
	if err != nil {
		finding := assetsSyncInvalidFinding(
			sourcePath,
			name,
			fmt.Sprintf("Canonical setup skill %s cannot be proven: %v.", name, err),
		)
		return nil, &finding
	}
	if workingDigest != committedDigest {
		finding := assetsSyncInvalidFinding(
			sourcePath,
			name,
			fmt.Sprintf("Canonical setup skill %s bytes do not match commit %s.", name, effectiveRevision),
		)
		return nil, &finding
	}
	if declaredDigest, ok := data["treeDigest"].(string); ok && declaredDigest != workingDigest {
		finding := assetsSyncInvalidFinding(
			sourcePath,
			name,
			fmt.Sprintf("Canonical setup skill %s tree digest does not match its source bytes.", name),
		)
		return nil, &finding
	}
	return &assetsSyncSkill{
		Name: name,
		Path: normalizedPath,
		Source: assetsSyncSource{
			Type:       "github",
			Repository: checkout.repository,
			Ref:        effectiveRevision,
			Path:       normalizedPath,
		},
		TreeDigest: workingDigest,
	}, nil
}

func normalizeAssetsSyncSkillPath(raw string) string {
	value := strings.ReplaceAll(strings.TrimSpace(raw), `\`, "/")
	if value == "" {
		return ""
	}
	if !strings.Contains(value, "/") {
		value = ".agents/skills/" + value + "/SKILL.md"
	}
	if !safeRelative(value) {
		return ""
	}
	return value
}

func resolveAssetsSyncSkillTree(sourceDir, raw string) string {
	value := strings.ReplaceAll(strings.TrimSpace(raw), `\`, "/")
	if value == "" || path.IsAbs(value) || containsParent(value) || isDrivePath(value) {
		return ""
	}
	var candidate string
	if strings.HasPrefix(value, "skills/") {
		candidate = filepath.Join(filepath.Dir(sourceDir), filepath.FromSlash(value))
	} else {
		candidate = filepath.Join(sourceDir, filepath.FromSlash(value))
	}
	info, err := os.Stat(candidate)
	if err != nil {
		return ""
	}
	if info.IsDir() {
		return filepath.Clean(candidate)
	}
	if filepath.Base(candidate) == "SKILL.md" && info.Mode().IsRegular() {
		return filepath.Dir(candidate)
	}
	return filepath.Clean(candidate)
}

func inferAssetsSyncSkillName(skillPath string) string {
	if path.Base(skillPath) == "SKILL.md" {
		return path.Base(path.Dir(skillPath))
	}
	return strings.TrimSuffix(path.Base(skillPath), path.Ext(skillPath))
}

func firstAssetsSyncString(values map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := values[key].(string); ok && strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func applyAssetsSyncPlans(
	ctx context.Context,
	dependencies assetsSyncDependencies,
	plans []assetsSyncPlan,
) error {
	_, identity, err := inspectRepositoryIdentity(ctx, dependencies.repository, nil)
	if err != nil {
		return fmt.Errorf("inspect asset target repository: %w", err)
	}
	anchored, err := os.OpenRoot(dependencies.repository)
	if err != nil {
		return fmt.Errorf("anchor asset target repository: %w", err)
	}
	defer anchored.Close()
	postimages := make([]Postimage, len(plans))
	preimages := make([]Preimage, len(plans))
	digestInput := make(map[string]string, len(plans))
	for index, plan := range plans {
		target := filepath.Join(dependencies.assetRoot, filepath.FromSlash(plan.path))
		relative, err := filepath.Rel(dependencies.repository, target)
		if err != nil {
			return fmt.Errorf("resolve canonical asset path: %w", err)
		}
		relative = filepath.ToSlash(relative)
		if !safeRestoreRelative(relative) {
			return fmt.Errorf("canonical asset path %q is unsafe", relative)
		}
		if err := validatePathParents(anchored, relative); err != nil {
			return fmt.Errorf("validate canonical asset path %q: %w", relative, err)
		}
		state, err := captureTransactionState(anchored, relative)
		if err != nil {
			return fmt.Errorf("capture canonical asset preimage %q: %w", relative, err)
		}
		preimages[index] = preimageFromTransactionState(relative, state)
		postimages[index] = Postimage{
			Path:            relative,
			Kind:            PreimageRegular,
			Mode:            0o644,
			Content:         append([]byte(nil), plan.content...),
			ContentIdentity: transactionContentIdentity(plan.content),
		}
		digestInput[relative] = transactionContentIdentity(plan.content)
	}
	digest, err := canonicalSHA256(digestInput)
	if err != nil {
		return fmt.Errorf("digest canonical asset refresh: %w", err)
	}
	document := PlanDocument{
		Repository: identity,
		Preimages:  preimages,
		Postimages: postimages,
		PlanDigest: "sha256:" + digest,
	}
	transaction, err := beginFileTransaction(
		ctx,
		dependencies.repository,
		document,
		dependencies.transactionHook,
		validateRestoreTransactionPlan,
		verifyRestoreTransactionState,
	)
	if err != nil {
		return fmt.Errorf("begin canonical asset transaction: %w", err)
	}
	_, applyErr := transaction.Apply(ctx)
	closeErr := transaction.Close()
	if applyErr != nil || closeErr != nil {
		return fmt.Errorf(
			"apply canonical asset transaction: %w",
			errors.Join(applyErr, closeErr),
		)
	}
	return nil
}

func assetsSyncInvalidFinding(sourcePath, setupID, message string) AssetsSyncFinding {
	return AssetsSyncFinding{
		Code:      "skills.setup-snapshot.drift",
		Severity:  "error",
		Path:      filepath.ToSlash(sourcePath),
		ManagedID: "setup." + setupID,
		Message:   message,
		Action:    "Fix the canonical setup source before synchronizing snapshots.",
	}
}

func assetsSyncDriftFinding(setupID string) AssetsSyncFinding {
	return AssetsSyncFinding{
		Code:      "skills.setup-snapshot.drift",
		Severity:  "error",
		Path:      "assets/setups/" + setupID + ".json",
		ManagedID: "setup." + setupID,
		Message:   "Bundled setup snapshot differs from the canonical source.",
		Action:    "Run roundfix baseline assets sync without --check to refresh the bundled snapshot.",
	}
}

func assetsSyncError(
	category AssetsSyncErrorCategory,
	code, severity, assetPath, managedID, message, action string,
	err error,
) *AssetsSyncError {
	return &AssetsSyncError{
		Category: category,
		Finding: AssetsSyncFinding{
			Code:      code,
			Severity:  severity,
			Path:      assetPath,
			ManagedID: managedID,
			Message:   message,
			Action:    action,
		},
		Err: err,
	}
}

func failedAssetsSyncPayload(err error) (AssetsSyncPayload, error) {
	var syncErr *AssetsSyncError
	if errors.As(err, &syncErr) {
		return assetsSyncPayload([]AssetsSyncFinding{syncErr.Finding}), err
	}
	return emptyAssetsSyncPayload(), err
}

func failedAssetsSyncFindings(
	findings []AssetsSyncFinding,
	category AssetsSyncErrorCategory,
) (AssetsSyncPayload, error) {
	sortAssetsSyncFindings(findings)
	first := findings[0]
	return assetsSyncPayload(findings), &AssetsSyncError{
		Category: category,
		Finding:  first,
		Err:      errors.New(strings.TrimSuffix(first.Message, ".")),
	}
}

func assetsSyncPayload(findings []AssetsSyncFinding) AssetsSyncPayload {
	sortAssetsSyncFindings(findings)
	summary := AssetsSyncSummary{}
	for _, finding := range findings {
		switch finding.Severity {
		case "error":
			summary.Errors++
		case "decision":
			summary.Decisions++
		case "warning":
			summary.Warnings++
		case "info":
			summary.Info++
		}
	}
	return AssetsSyncPayload{
		SchemaVersion:  AssetsSyncSchemaVersion,
		OK:             summary.Errors == 0 && summary.Decisions == 0,
		Summary:        summary,
		Findings:       append([]AssetsSyncFinding(nil), findings...),
		PlannedChanges: []AssetsSyncChange{},
	}
}

func emptyAssetsSyncPayload() AssetsSyncPayload {
	return assetsSyncPayload([]AssetsSyncFinding{})
}

func hasAssetsSyncErrors(findings []AssetsSyncFinding) bool {
	for _, finding := range findings {
		if finding.Severity == "error" || finding.Severity == "decision" {
			return true
		}
	}
	return false
}

func sortAssetsSyncFindings(findings []AssetsSyncFinding) {
	order := map[string]int{"error": 0, "decision": 1, "warning": 2, "info": 3}
	sort.Slice(findings, func(i, j int) bool {
		left, right := findings[i], findings[j]
		if order[left.Severity] != order[right.Severity] {
			return order[left.Severity] < order[right.Severity]
		}
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		if left.Path != right.Path {
			return left.Path < right.Path
		}
		if left.ManagedID != right.ManagedID {
			return left.ManagedID < right.ManagedID
		}
		return left.Message < right.Message
	})
}

type assetsOverlayFS struct {
	base      fs.FS
	overrides map[string][]byte
}

func (overlay assetsOverlayFS) Open(name string) (fs.File, error) {
	if data, ok := overlay.overrides[name]; ok {
		return &assetsMemoryFile{
			Reader: bytes.NewReader(data),
			name:   path.Base(name),
			size:   int64(len(data)),
		}, nil
	}
	return overlay.base.Open(name)
}

type assetsMemoryFile struct {
	*bytes.Reader
	name string
	size int64
}

func (file *assetsMemoryFile) Close() error {
	return nil
}

func (file *assetsMemoryFile) Stat() (fs.FileInfo, error) {
	return assetsMemoryInfo{name: file.name, size: file.size}, nil
}

type assetsMemoryInfo struct {
	name string
	size int64
}

func (info assetsMemoryInfo) Name() string       { return info.name }
func (info assetsMemoryInfo) Size() int64        { return info.size }
func (info assetsMemoryInfo) Mode() fs.FileMode  { return 0o444 }
func (info assetsMemoryInfo) ModTime() time.Time { return time.Time{} }
func (info assetsMemoryInfo) IsDir() bool        { return false }
func (info assetsMemoryInfo) Sys() any           { return nil }
