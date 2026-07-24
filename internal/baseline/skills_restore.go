package baseline

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

const (
	SkillsRestoreSchemaVersion = "setup-context-driven/restore-v1"

	lockHashCompatibilitySchema = "setup-context-driven/external-lock-hash-compatibility-v1"
	lockHashCompatibilityPath   = "lock-hash-compatibility-v1.json"
	skillsLockPath              = "skills-lock.json"
	restoreMaxFiles             = 2_000
	restoreMaxBytes             = 8 * 1024 * 1024
)

var lowercaseSHA256 = regexp.MustCompile(`^[0-9a-f]{64}$`)

// SkillsRestoreRequest selects one immutable Repository Skill Set restoration.
type SkillsRestoreRequest struct {
	Repository   string
	ProfileID    string
	Skills       []string
	SourceDir    string
	Confirmation string
}

// SkillsRestorePayload is the stable text and JSON command result.
type SkillsRestorePayload struct {
	SchemaVersion  string                 `json:"schemaVersion"`
	OK             bool                   `json:"ok"`
	Applied        bool                   `json:"applied"`
	Profile        string                 `json:"profile"`
	Setup          *string                `json:"setup"`
	Acquisitions   []RestoreAcquisition   `json:"acquisitions"`
	Skills         []RestoreSkill         `json:"skills"`
	PlannedChanges []RestorePlannedChange `json:"plannedChanges"`
	PlanDigest     *string                `json:"planDigest"`
	Finding        *RestoreFinding        `json:"finding,omitempty"`
}

// RestoreAcquisition identifies one exact immutable Git provenance fetch.
type RestoreAcquisition struct {
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
}

// RestoreSource identifies one external skill's immutable source subtree.
type RestoreSource struct {
	Provider   string `json:"provider"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	Path       string `json:"path"`
}

// RestoreFileChange is one exact skill-tree file change.
type RestoreFileChange struct {
	Action       string  `json:"action"`
	Path         string  `json:"path"`
	Skill        string  `json:"skill"`
	BeforeDigest *string `json:"beforeDigest"`
	AfterDigest  *string `json:"afterDigest"`
}

// RestoreLockEdit records one portable skills-lock.json entry update.
type RestoreLockEdit struct {
	Action string         `json:"action"`
	Path   string         `json:"path"`
	Skill  string         `json:"skill"`
	Before map[string]any `json:"before"`
	After  map[string]any `json:"after"`
}

// RestorePlannedChange is one flattened file or lock operation.
type RestorePlannedChange struct {
	Action       string         `json:"action"`
	Path         string         `json:"path"`
	Skill        string         `json:"skill"`
	BeforeDigest *string        `json:"-"`
	AfterDigest  *string        `json:"-"`
	Before       map[string]any `json:"-"`
	After        map[string]any `json:"-"`
}

// RestoreSkill is one selected external Repository Skill Set member.
type RestoreSkill struct {
	Skill          string              `json:"skill"`
	TargetPath     string              `json:"targetPath"`
	Source         RestoreSource       `json:"source"`
	ExpectedDigest string              `json:"expectedDigest"`
	ObservedDigest *string             `json:"observedDigest"`
	Changes        []RestoreFileChange `json:"changes"`
	LockEdit       *RestoreLockEdit    `json:"lockEdit"`
}

// RestoreFinding is one stable actionable command outcome.
type RestoreFinding struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Action  string `json:"action"`
}

// MarshalJSON preserves the maintained operation-specific restoration shape:
// file changes carry byte digests, while lock edits carry entry documents.
func (change RestorePlannedChange) MarshalJSON() ([]byte, error) {
	if change.Action == "update-lock-entry" {
		return json.Marshal(struct {
			Action string         `json:"action"`
			Path   string         `json:"path"`
			Skill  string         `json:"skill"`
			Before map[string]any `json:"before"`
			After  map[string]any `json:"after"`
		}{
			Action: change.Action,
			Path:   change.Path,
			Skill:  change.Skill,
			Before: change.Before,
			After:  change.After,
		})
	}
	return json.Marshal(struct {
		Action       string  `json:"action"`
		Path         string  `json:"path"`
		Skill        string  `json:"skill"`
		BeforeDigest *string `json:"beforeDigest"`
		AfterDigest  *string `json:"afterDigest"`
	}{
		Action:       change.Action,
		Path:         change.Path,
		Skill:        change.Skill,
		BeforeDigest: change.BeforeDigest,
		AfterDigest:  change.AfterDigest,
	})
}

// SkillsRestoreErrorCategory maps engine failures to the Baseline CLI exits.
type SkillsRestoreErrorCategory string

const (
	SkillsRestoreInvalid   SkillsRestoreErrorCategory = "invalid"
	SkillsRestoreAction    SkillsRestoreErrorCategory = "action_required"
	SkillsRestoreExecution SkillsRestoreErrorCategory = "execution"
)

// SkillsRestoreError is one typed restoration refusal or failure.
type SkillsRestoreError struct {
	Category SkillsRestoreErrorCategory
	Finding  RestoreFinding
	Err      error
}

func (err *SkillsRestoreError) Error() string {
	return err.Finding.Message
}

func (err *SkillsRestoreError) Unwrap() error {
	return err.Err
}

type restoreProfile struct {
	ID     string
	Setup  string
	Skills map[string]restoreSkillContract
}

type restoreSkillContract struct {
	Name       string
	Source     RestoreSource
	TreeDigest string
}

type restoreFile struct {
	Path    string
	Content []byte
}

type restorePlan struct {
	payload     SkillsRestorePayload
	document    PlanDocument
	repository  string
	sourceFiles map[string][]restoreFile
}

type restoreProvenance struct {
	Provider   string
	Repository string
	Ref        string
}

type restoreDependencies struct {
	profile         restoreProfile
	adapterFixture  []byte
	transactionHook func(transactionFaultPoint) error
}

// RestoreSkills previews or applies one exact immutable skill restoration.
func RestoreSkills(ctx context.Context, request SkillsRestoreRequest) (SkillsRestorePayload, error) {
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		return failedRestorePayload(request.ProfileID, "", restoreError(
			SkillsRestoreInvalid,
			"restore.assets-invalid",
			fmt.Sprintf("Embedded Baseline assets are invalid: %v.", err),
			"Fix the embedded Baseline catalog before restoring skills.",
			err,
		))
	}
	profile, err := loadRestoreProfile(catalog, request.ProfileID)
	if err != nil {
		return failedRestorePayload(request.ProfileID, "", err)
	}
	fixture, ok := catalog.Asset(lockHashCompatibilityPath)
	if !ok {
		err := restoreError(
			SkillsRestoreExecution,
			"lock.adapter-incompatible",
			"The lock compatibility fixture is missing from the embedded Baseline catalog.",
			"Update the isolated lock adapter before restoring skills.",
			errors.New("lock compatibility fixture is missing"),
		)
		return failedRestorePayload(request.ProfileID, profile.Setup, err)
	}
	return restoreSkills(ctx, request, restoreDependencies{
		profile:        profile,
		adapterFixture: fixture.Data,
	})
}

func restoreSkills(
	ctx context.Context,
	request SkillsRestoreRequest,
	dependencies restoreDependencies,
) (SkillsRestorePayload, error) {
	if ctx == nil {
		err := restoreError(
			SkillsRestoreInvalid,
			"restore.context-invalid",
			"Restoration requires a live context.",
			"Rerun roundfix baseline skills restore.",
			errors.New("context is required"),
		)
		return failedRestorePayload(request.ProfileID, dependencies.profile.Setup, err)
	}
	if err := ctx.Err(); err != nil {
		return failedRestorePayload(request.ProfileID, dependencies.profile.Setup, err)
	}
	request.Repository = strings.TrimSpace(request.Repository)
	request.ProfileID = strings.TrimSpace(request.ProfileID)
	request.SourceDir = strings.TrimSpace(request.SourceDir)
	request.Confirmation = strings.TrimSpace(request.Confirmation)
	if request.Repository == "" {
		err := restoreError(
			SkillsRestoreInvalid,
			"restore.repo-invalid",
			"Repository path cannot be empty.",
			"Pass an existing Git worktree with --repo.",
			errors.New("repository path is empty"),
		)
		return failedRestorePayload(request.ProfileID, dependencies.profile.Setup, err)
	}
	if request.Confirmation != "" && !lowercaseSHA256.MatchString(request.Confirmation) {
		err := restoreError(
			SkillsRestoreInvalid,
			"plan.confirmation.invalid",
			"Plan confirmation must be a lowercase SHA-256 digest.",
			"Pass the exact planDigest returned by the restoration preview.",
			errors.New("confirmation digest is malformed"),
		)
		return failedRestorePayload(request.ProfileID, dependencies.profile.Setup, err)
	}

	plan, err := buildSkillsRestorePlan(ctx, request, dependencies)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			err = restoreError(
				SkillsRestoreExecution,
				"restore.canceled",
				"Restoration was canceled before mutation completed.",
				"Rerun roundfix baseline skills restore when the operation can complete.",
				err,
			)
		}
		return failedRestorePayload(request.ProfileID, dependencies.profile.Setup, err)
	}
	if len(plan.payload.Skills) == 0 {
		plan.payload.OK = true
		return plan.payload, nil
	}
	if plan.payload.PlanDigest == nil || request.Confirmation != *plan.payload.PlanDigest {
		code := "plan.confirmation.stale"
		message := "The supplied confirmation does not match the current restoration Change Plan."
		if request.Confirmation == "" {
			code = "plan.confirmation.required"
			message = "Restoration requires confirmation of this exact Change Plan."
		}
		err := restoreError(
			SkillsRestoreAction,
			code,
			message,
			"Review plannedChanges and rerun with --confirm-plan planDigest.",
			errors.New("restoration plan is not confirmed"),
		)
		plan.payload.Finding = &err.Finding
		return plan.payload, err
	}
	if err := applySkillsRestorePlan(ctx, plan, dependencies.transactionHook); err != nil {
		var restoreErr *SkillsRestoreError
		if errors.As(err, &restoreErr) {
			plan.payload.Finding = &restoreErr.Finding
		}
		return plan.payload, err
	}
	plan.payload.OK = true
	plan.payload.Applied = true
	plan.payload.Finding = &RestoreFinding{
		Code:    "restore.completed",
		Message: "Selected external Repository Skill Set members match their immutable snapshots.",
		Action:  "Run roundfix baseline plan to verify the selected Baseline Profile.",
	}
	return plan.payload, nil
}

func buildSkillsRestorePlan(
	ctx context.Context,
	request SkillsRestoreRequest,
	dependencies restoreDependencies,
) (restorePlan, error) {
	if err := validateLockAdapter(dependencies.adapterFixture); err != nil {
		return restorePlan{}, err
	}
	if request.ProfileID != dependencies.profile.ID {
		return restorePlan{}, restoreError(
			SkillsRestoreInvalid,
			"restore.profile-unknown",
			fmt.Sprintf("Unknown built-in Baseline Profile %q.", request.ProfileID),
			"Choose a profile id from the embedded Baseline catalog.",
			errors.New("profile does not match the selected catalog"),
		)
	}
	selected, err := selectRestoreSkills(dependencies.profile, request.Skills)
	if err != nil {
		return restorePlan{}, err
	}
	if request.SourceDir != "" {
		info, statErr := os.Stat(request.SourceDir)
		if statErr != nil || !info.IsDir() {
			return restorePlan{}, restoreError(
				SkillsRestoreInvalid,
				"restore.source-dir-invalid",
				fmt.Sprintf("Offline Git object store is not a directory: %s.", request.SourceDir),
				"Pass an existing Git checkout or bare object store to --source-dir.",
				statErr,
			)
		}
	}

	root, identity, err := inspectRepositoryIdentity(ctx, request.Repository, nil)
	if err != nil {
		return restorePlan{}, restoreError(
			SkillsRestoreInvalid,
			"restore.repo-invalid",
			fmt.Sprintf("Repository root is not a usable Git worktree: %v.", err),
			"Pass an existing Git worktree with --repo.",
			err,
		)
	}
	lock, err := loadSkillsLock(filepath.Join(root, skillsLockPath))
	if err != nil {
		return restorePlan{}, err
	}

	groups := make(map[restoreProvenance][]restoreSkillContract)
	for _, contract := range selected {
		key := restoreProvenance{
			Provider:   contract.Source.Provider,
			Repository: contract.Source.Repository,
			Ref:        contract.Source.Ref,
		}
		groups[key] = append(groups[key], contract)
	}
	provenances := make([]restoreProvenance, 0, len(groups))
	for provenance := range groups {
		provenances = append(provenances, provenance)
	}
	sort.Slice(provenances, func(i, j int) bool {
		left, right := provenances[i], provenances[j]
		if left.Provider != right.Provider {
			return left.Provider < right.Provider
		}
		if left.Repository != right.Repository {
			return left.Repository < right.Repository
		}
		return left.Ref < right.Ref
	})

	acquired := make(map[string][]restoreFile, len(selected))
	acquisitions := make([]RestoreAcquisition, 0, len(provenances))
	totalFiles := 0
	totalBytes := 0
	for _, provenance := range provenances {
		contracts := groups[provenance]
		sort.Slice(contracts, func(i, j int) bool {
			return contracts[i].Name < contracts[j].Name
		})
		filesBySkill, acquireErr := acquireRestoreGroup(ctx, provenance, contracts, request.SourceDir)
		if acquireErr != nil {
			return restorePlan{}, acquireErr
		}
		acquisitions = append(acquisitions, RestoreAcquisition{
			Provider: provenance.Provider, Repository: provenance.Repository, Ref: provenance.Ref,
		})
		for name, files := range filesBySkill {
			totalFiles += len(files)
			for _, file := range files {
				totalBytes += len(file.Content)
			}
			if totalFiles > restoreMaxFiles || totalBytes > restoreMaxBytes {
				return restorePlan{}, restoreError(
					SkillsRestoreExecution,
					"source.limit-exceeded",
					"Acquired skill content exceeds the restoration file-count or byte limit.",
					"Reduce the declared immutable source subtree before retrying.",
					errors.New("aggregate source limit exceeded"),
				)
			}
			acquired[name] = files
		}
	}

	lockAfter := lock.clone()
	var skillPlans []RestoreSkill
	var plannedChanges []RestorePlannedChange
	sourceFiles := make(map[string][]restoreFile)
	for _, contract := range selected {
		target := path.Join(".agents/skills", contract.Name)
		current, exists, inspectErr := inspectRestoreTarget(root, target)
		if inspectErr != nil {
			return restorePlan{}, inspectErr
		}
		files := acquired[contract.Name]
		observed := nullableDigest(exists, portableRestoreDigest(current))
		changes := restoreChanges(contract.Name, target, current, files)
		expectedEntry, entryErr := expectedSkillsLockEntry(contract, files)
		if entryErr != nil {
			return restorePlan{}, entryErr
		}
		beforeEntry, beforeObject := lockAfter.objectEntry("skills", contract.Name)
		var lockEdit *RestoreLockEdit
		if !beforeObject || !reflect.DeepEqual(beforeEntry, expectedEntry) {
			lockEdit = &RestoreLockEdit{
				Action: "update-lock-entry",
				Path:   skillsLockPath,
				Skill:  contract.Name,
				Before: beforeEntry,
				After:  expectedEntry,
			}
			lockAfter.setSkillsEntry(contract.Name, expectedEntry)
		}
		if len(changes) == 0 && lockEdit == nil {
			continue
		}
		skill := RestoreSkill{
			Skill:          contract.Name,
			TargetPath:     target,
			Source:         contract.Source,
			ExpectedDigest: contract.TreeDigest,
			ObservedDigest: observed,
			Changes:        changes,
			LockEdit:       lockEdit,
		}
		skillPlans = append(skillPlans, skill)
		sourceFiles[contract.Name] = files
		for _, change := range changes {
			plannedChanges = append(plannedChanges, RestorePlannedChange{
				Action:       change.Action,
				Path:         change.Path,
				Skill:        change.Skill,
				BeforeDigest: change.BeforeDigest,
				AfterDigest:  change.AfterDigest,
			})
		}
		if lockEdit != nil {
			plannedChanges = append(plannedChanges, RestorePlannedChange{
				Action: lockEdit.Action,
				Path:   lockEdit.Path,
				Skill:  lockEdit.Skill,
				Before: lockEdit.Before,
				After:  lockEdit.After,
			})
		}
	}

	lockAfterBytes := lock.before
	if hasLockEdits(skillPlans) {
		lockAfterBytes, err = lockAfter.marshalIndent()
		if err != nil {
			return restorePlan{}, restoreError(
				SkillsRestoreExecution,
				"lock.write-failed",
				fmt.Sprintf("Could not serialize skills-lock.json: %v.", err),
				"Fix the lock document and preview restoration again.",
				err,
			)
		}
	}
	digest, err := restorePlanDigest(
		dependencies.profile,
		acquisitions,
		skillPlans,
		plannedChanges,
		lock.before,
		lockAfterBytes,
	)
	if err != nil {
		return restorePlan{}, err
	}
	payload := SkillsRestorePayload{
		SchemaVersion:  SkillsRestoreSchemaVersion,
		Profile:        dependencies.profile.ID,
		Setup:          stringPointer(dependencies.profile.Setup),
		Acquisitions:   acquisitions,
		Skills:         skillPlans,
		PlannedChanges: plannedChanges,
		PlanDigest:     stringPointer(digest),
	}
	document, err := buildRestoreTransactionDocument(
		root,
		identity,
		digest,
		skillPlans,
		sourceFiles,
		lockAfterBytes,
	)
	if err != nil {
		return restorePlan{}, err
	}
	return restorePlan{
		payload: payload, document: document, repository: root, sourceFiles: sourceFiles,
	}, nil
}

func loadRestoreProfile(catalog *Catalog, profileID string) (restoreProfile, error) {
	profileID = strings.TrimSpace(profileID)
	entry, ok := catalog.Profile(profileID)
	if !ok {
		return restoreProfile{}, restoreError(
			SkillsRestoreInvalid,
			"restore.profile-unknown",
			fmt.Sprintf("Unknown built-in Baseline Profile %q.", profileID),
			"Choose a profile id from the embedded Baseline catalog.",
			errors.New("profile is not embedded"),
		)
	}
	var profileDocument struct {
		ID    string `json:"id"`
		Setup string `json:"setup"`
	}
	if err := json.Unmarshal(entry.Data, &profileDocument); err != nil {
		return restoreProfile{}, restoreError(
			SkillsRestoreInvalid,
			"restore.assets-invalid",
			fmt.Sprintf("Embedded Baseline Profile %q is invalid.", profileID),
			"Fix the embedded Baseline catalog before restoring skills.",
			err,
		)
	}
	setupEntry, ok := catalog.Setup(profileDocument.Setup)
	if !ok {
		return restoreProfile{}, restoreError(
			SkillsRestoreInvalid,
			"restore.assets-invalid",
			fmt.Sprintf("Embedded Baseline Profile %q has no Setup Snapshot.", profileID),
			"Fix the embedded Baseline catalog before restoring skills.",
			errors.New("profile setup is missing"),
		)
	}
	var setupDocument struct {
		ID     string `json:"id"`
		Skills []struct {
			Name       string `json:"name"`
			TreeDigest string `json:"treeDigest"`
			Source     struct {
				Type       string `json:"type"`
				Repository string `json:"repository"`
				Ref        string `json:"ref"`
				Path       string `json:"path"`
			} `json:"source"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(setupEntry.Data, &setupDocument); err != nil {
		return restoreProfile{}, restoreError(
			SkillsRestoreInvalid,
			"restore.assets-invalid",
			fmt.Sprintf("Embedded Setup Snapshot %q is invalid.", profileDocument.Setup),
			"Fix the embedded Baseline catalog before restoring skills.",
			err,
		)
	}
	contracts := make(map[string]restoreSkillContract)
	for _, skill := range setupDocument.Skills {
		if skill.Source.Type != "github" {
			continue
		}
		contracts[skill.Name] = restoreSkillContract{
			Name: skill.Name,
			Source: RestoreSource{
				Provider:   skill.Source.Type,
				Repository: skill.Source.Repository,
				Ref:        skill.Source.Ref,
				Path:       skill.Source.Path,
			},
			TreeDigest: skill.TreeDigest,
		}
	}
	return restoreProfile{ID: profileDocument.ID, Setup: setupDocument.ID, Skills: contracts}, nil
}

func selectRestoreSkills(profile restoreProfile, requested []string) ([]restoreSkillContract, error) {
	names := make([]string, 0, len(profile.Skills))
	if len(requested) == 0 {
		for name := range profile.Skills {
			names = append(names, name)
		}
	} else {
		seen := make(map[string]struct{})
		var unknown []string
		for _, raw := range requested {
			name := strings.TrimSpace(raw)
			if _, ok := profile.Skills[name]; !ok {
				unknown = append(unknown, name)
				continue
			}
			if _, duplicate := seen[name]; duplicate {
				continue
			}
			seen[name] = struct{}{}
			names = append(names, name)
		}
		if len(unknown) != 0 {
			sort.Strings(unknown)
			return nil, restoreError(
				SkillsRestoreInvalid,
				"restore.skill-invalid",
				"Selected skills are not external members of this profile: "+strings.Join(unknown, ", ")+".",
				"Choose repeatable --skill values from this profile's external Repository Skill Set.",
				errors.New("selected skill is not in the profile"),
			)
		}
	}
	sort.Strings(names)
	contracts := make([]restoreSkillContract, len(names))
	for index, name := range names {
		contracts[index] = profile.Skills[name]
	}
	return contracts, nil
}

func validateLockAdapter(data []byte) error {
	var fixture struct {
		SchemaVersion string `json:"schemaVersion"`
		Version       int    `json:"version"`
		Files         []struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		} `json:"files"`
		ExpectedSHA256 string `json:"expectedSha256"`
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&fixture); err != nil {
		return lockAdapterError(fmt.Errorf("decode fixture: %w", err))
	}
	if err := requireJSONEOF(decoder); err != nil {
		return lockAdapterError(err)
	}
	if fixture.SchemaVersion != lockHashCompatibilitySchema ||
		fixture.Version != 1 ||
		len(fixture.Files) == 0 ||
		!lowercaseSHA256.MatchString(fixture.ExpectedSHA256) {
		return lockAdapterError(errors.New("fixture identity is invalid"))
	}
	files := make([]restoreFile, len(fixture.Files))
	seen := make(map[string]struct{}, len(fixture.Files))
	for index, item := range fixture.Files {
		if !safeRestoreRelative(item.Path) {
			return lockAdapterError(fmt.Errorf("fixture path %q is unsafe", item.Path))
		}
		if _, duplicate := seen[item.Path]; duplicate {
			return lockAdapterError(fmt.Errorf("fixture path %q is duplicated", item.Path))
		}
		seen[item.Path] = struct{}{}
		files[index] = restoreFile{Path: item.Path, Content: []byte(item.Content)}
	}
	if observed := externalSkillsLockDigest(files); observed != fixture.ExpectedSHA256 {
		return lockAdapterError(fmt.Errorf(
			"computed digest %s does not match pinned digest %s",
			observed,
			fixture.ExpectedSHA256,
		))
	}
	return nil
}

func lockAdapterError(err error) error {
	return restoreError(
		SkillsRestoreExecution,
		"lock.adapter-incompatible",
		"The lock compatibility fixture cannot prove agreement with the maintained lock contract.",
		"Update the isolated lock adapter before restoring skills.",
		err,
	)
}

func requireJSONEOF(decoder *json.Decoder) error {
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("JSON document contains trailing values")
		}
		return err
	}
	return nil
}

func expectedSkillsLockEntry(
	contract restoreSkillContract,
	files []restoreFile,
) (map[string]any, error) {
	hasSkill := false
	for _, file := range files {
		if file.Path == "SKILL.md" {
			hasSkill = true
			break
		}
	}
	if !hasSkill {
		return nil, restoreError(
			SkillsRestoreExecution,
			"source.skill-file-missing",
			"The verified source subtree does not contain SKILL.md.",
			"Fix the immutable source snapshot before restoring this skill.",
			errors.New("SKILL.md is missing"),
		)
	}
	return map[string]any{
		"source":       contract.Source.Repository,
		"ref":          contract.Source.Ref,
		"sourceType":   contract.Source.Provider,
		"skillPath":    path.Join(contract.Source.Path, "SKILL.md"),
		"computedHash": externalSkillsLockDigest(files),
	}, nil
}

func inspectRestoreTarget(root, target string) ([]restoreFile, bool, error) {
	absolute := filepath.Join(root, filepath.FromSlash(target))
	info, err := os.Lstat(absolute)
	if errors.Is(err, fs.ErrNotExist) {
		if err := validateRestoreParents(root, target); err != nil {
			return nil, false, err
		}
		return []restoreFile{}, false, nil
	}
	if err != nil {
		return nil, false, restoreError(
			SkillsRestoreExecution,
			"target.read-failed",
			fmt.Sprintf("Could not inspect restoration target %s: %v.", target, err),
			"Fix repository permissions and rerun the restoration preview.",
			err,
		)
	}
	if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		return nil, false, restoreError(
			SkillsRestoreInvalid,
			"target.unsafe-tree",
			fmt.Sprintf("Restoration target is not a regular directory: %s.", target),
			"Replace the unsafe target manually, then rerun the restoration preview.",
			errors.New("target root is unsafe"),
		)
	}
	if err := validateRestoreParents(root, target); err != nil {
		return nil, false, err
	}
	files := []restoreFile{}
	totalBytes := 0
	err = filepath.WalkDir(absolute, func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if name == absolute {
			return nil
		}
		relative, err := filepath.Rel(absolute, name)
		if err != nil {
			return err
		}
		portable := filepath.ToSlash(relative)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return fmt.Errorf("entry %q is a symbolic link", portable)
		}
		if entry.IsDir() {
			return nil
		}
		if !info.Mode().IsRegular() || fileLinkCount(info) != 1 {
			return fmt.Errorf("entry %q is not an independent regular file", portable)
		}
		content, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		files = append(files, restoreFile{Path: portable, Content: content})
		totalBytes += len(content)
		if len(files) > restoreMaxFiles || totalBytes > restoreMaxBytes {
			return errors.New("target limit exceeded")
		}
		return nil
	})
	if err != nil {
		code := "target.read-failed"
		action := "Fix repository permissions and rerun the restoration preview."
		if strings.Contains(err.Error(), "symbolic link") ||
			strings.Contains(err.Error(), "regular file") {
			code = "target.unsafe-tree"
			action = "Remove the unsafe target entry manually, then rerun the restoration preview."
		}
		if strings.Contains(err.Error(), "limit exceeded") {
			code = "target.limit-exceeded"
			action = "Reduce the target directory before retrying."
		}
		return nil, false, restoreError(
			SkillsRestoreInvalid,
			code,
			fmt.Sprintf("Could not safely inspect restoration target %s: %v.", target, err),
			action,
			err,
		)
	}
	sortRestoreFiles(files)
	return files, true, nil
}

func validateRestoreParents(root, target string) error {
	current := root
	parts := strings.Split(path.Dir(target), "/")
	for _, part := range parts {
		if part == "." || part == "" {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		if err != nil {
			return restoreError(
				SkillsRestoreExecution,
				"target.read-failed",
				fmt.Sprintf("Could not inspect restoration target parent %s: %v.", current, err),
				"Fix repository permissions and rerun the restoration preview.",
				err,
			)
		}
		if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
			return restoreError(
				SkillsRestoreInvalid,
				"target.unsafe-path",
				fmt.Sprintf("Restoration target parent is not a repository directory: %s.", current),
				"Replace the unsafe parent before retrying.",
				errors.New("target parent is unsafe"),
			)
		}
	}
	return nil
}

func restoreChanges(
	skillName, target string,
	before, after []restoreFile,
) []RestoreFileChange {
	beforeByPath := restoreFilesByPath(before)
	afterByPath := restoreFilesByPath(after)
	paths := make([]string, 0, len(beforeByPath)+len(afterByPath))
	seen := make(map[string]struct{}, cap(paths))
	for relative := range beforeByPath {
		seen[relative] = struct{}{}
		paths = append(paths, relative)
	}
	for relative := range afterByPath {
		if _, ok := seen[relative]; !ok {
			paths = append(paths, relative)
		}
	}
	sort.Strings(paths)
	changes := make([]RestoreFileChange, 0, len(paths))
	for _, relative := range paths {
		beforeContent, beforeExists := beforeByPath[relative]
		afterContent, afterExists := afterByPath[relative]
		if beforeExists && afterExists && bytes.Equal(beforeContent, afterContent) {
			continue
		}
		action := "refresh"
		if !beforeExists {
			action = "create"
		}
		if !afterExists {
			action = "remove"
		}
		changes = append(changes, RestoreFileChange{
			Action:       action,
			Path:         path.Join(target, relative),
			Skill:        skillName,
			BeforeDigest: optionalBytesDigest(beforeContent, beforeExists),
			AfterDigest:  optionalBytesDigest(afterContent, afterExists),
		})
	}
	return changes
}

func restorePlanDigest(
	profile restoreProfile,
	acquisitions []RestoreAcquisition,
	skills []RestoreSkill,
	changes []RestorePlannedChange,
	lockBefore, lockAfter []byte,
) (string, error) {
	payload := map[string]any{
		"kind":             "restore-skills",
		"profile":          profile.ID,
		"setup":            profile.Setup,
		"acquisitions":     acquisitions,
		"skills":           skills,
		"plannedChanges":   changes,
		"lockBeforeDigest": optionalBytesDigest(lockBefore, lockBefore != nil),
		"lockAfterDigest":  optionalBytesDigest(lockAfter, lockAfter != nil),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", restoreError(
			SkillsRestoreExecution,
			"restore.plan-failed",
			fmt.Sprintf("Could not serialize the restoration Change Plan: %v.", err),
			"Fix the embedded restoration contract before retrying.",
			err,
		)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func buildRestoreTransactionDocument(
	root string,
	identity RepositoryIdentity,
	digest string,
	skills []RestoreSkill,
	sourceFiles map[string][]restoreFile,
	lockAfter []byte,
) (PlanDocument, error) {
	postimages := make(map[string]Postimage)
	for _, skill := range skills {
		changedFiles := make(map[string]struct{}, len(skill.Changes))
		for _, change := range skill.Changes {
			if change.Action == "remove" {
				postimages[change.Path] = Postimage{Path: change.Path, Kind: PreimageMissing}
				continue
			}
			changedFiles[change.Path] = struct{}{}
		}
		for _, file := range sourceFiles[skill.Skill] {
			target := path.Join(skill.TargetPath, file.Path)
			if _, changed := changedFiles[target]; !changed {
				continue
			}
			postimages[target] = Postimage{
				Path:            target,
				Kind:            PreimageRegular,
				Mode:            0o644,
				Content:         append([]byte(nil), file.Content...),
				ContentIdentity: transactionContentIdentity(file.Content),
			}
		}
		if skill.LockEdit != nil {
			postimages[skillsLockPath] = Postimage{
				Path:            skillsLockPath,
				Kind:            PreimageRegular,
				Mode:            0o600,
				Content:         append([]byte(nil), lockAfter...),
				ContentIdentity: transactionContentIdentity(lockAfter),
			}
		}
	}
	paths := sortedKeys(postimages)
	preimages := make([]Preimage, len(paths))
	anchored, err := os.OpenRoot(root)
	if err != nil {
		return PlanDocument{}, restoreError(
			SkillsRestoreInvalid,
			"restore.repo-invalid",
			fmt.Sprintf("Could not anchor the restoration repository: %v.", err),
			"Repair the repository and rerun the preview.",
			err,
		)
	}
	defer anchored.Close()
	for index, relative := range paths {
		if err := validatePathParents(anchored, relative); err != nil {
			return PlanDocument{}, restoreError(
				SkillsRestoreInvalid,
				"target.unsafe-path",
				fmt.Sprintf("Restoration target %s is unsafe: %v.", relative, err),
				"Replace the unsafe target parent before retrying.",
				err,
			)
		}
		state, err := captureTransactionState(anchored, relative)
		if err != nil {
			return PlanDocument{}, restoreError(
				SkillsRestoreExecution,
				"target.read-failed",
				fmt.Sprintf("Could not capture restoration preimage %s: %v.", relative, err),
				"Fix repository permissions and rerun the preview.",
				err,
			)
		}
		preimages[index] = preimageFromTransactionState(relative, state)
	}
	orderedPostimages := make([]Postimage, len(paths))
	for index, relative := range paths {
		orderedPostimages[index] = postimages[relative]
	}
	return PlanDocument{
		Repository: identity,
		Preimages:  preimages,
		Postimages: orderedPostimages,
		PlanDigest: "sha256:" + digest,
	}, nil
}

func preimageFromTransactionState(relative string, state transactionFileState) Preimage {
	return Preimage{
		Path:            relative,
		Exists:          state.Exists,
		Kind:            state.Kind,
		Mode:            state.Mode,
		LinkTarget:      state.LinkTarget,
		Bytes:           state.Bytes,
		ContentIdentity: state.ContentIdentity,
	}
}

func applySkillsRestorePlan(
	ctx context.Context,
	plan restorePlan,
	hook func(transactionFaultPoint) error,
) error {
	transaction, err := beginFileTransaction(
		ctx,
		plan.repository,
		plan.document,
		hook,
		validateRestoreTransactionPlan,
		verifyRestoreTransactionState,
	)
	if err != nil {
		return classifyRestoreApplyError(fmt.Errorf("begin restoration transaction: %w", err))
	}
	_, applyErr := transaction.Apply(ctx)
	closeErr := transaction.Close()
	if applyErr != nil {
		return classifyRestoreApplyError(fmt.Errorf(
			"apply restoration transaction: %w",
			errors.Join(applyErr, closeErr),
		))
	}
	if closeErr != nil {
		return restoreError(
			SkillsRestoreExecution,
			"restore.close-failed",
			fmt.Sprintf("Restoration postimages were applied, but the transaction did not close cleanly: %v.", closeErr),
			"Inspect the Git-private Baseline transaction lock and recovery state before retrying.",
			closeErr,
		)
	}
	return nil
}

func validateRestoreTransactionPlan(
	ctx context.Context,
	repository string,
	document PlanDocument,
) error {
	if !strings.HasPrefix(document.PlanDigest, "sha256:") ||
		!lowercaseSHA256.MatchString(strings.TrimPrefix(document.PlanDigest, "sha256:")) {
		return errors.New("restoration Plan Digest is invalid")
	}
	root, identity, err := inspectRepositoryIdentity(ctx, repository, nil)
	if err != nil {
		return err
	}
	if identity.Digest != document.Repository.Digest {
		return errors.New("restoration repository identity does not match")
	}
	anchored, err := os.OpenRoot(root)
	if err != nil {
		return err
	}
	defer anchored.Close()
	preimages := make(map[string]Preimage, len(document.Preimages))
	for _, approved := range document.Preimages {
		if _, duplicate := preimages[approved.Path]; duplicate {
			return fmt.Errorf("duplicate restoration preimage %q", approved.Path)
		}
		preimages[approved.Path] = approved
		current, err := captureTransactionState(anchored, approved.Path)
		if err != nil {
			return err
		}
		if !transactionStateMatchesPreimage(current, approved) {
			return fmt.Errorf("restoration Plan preimage is stale at %q", approved.Path)
		}
	}
	seenPostimages := make(map[string]struct{}, len(document.Postimages))
	for _, postimage := range document.Postimages {
		if _, ok := preimages[postimage.Path]; !ok {
			return fmt.Errorf("restoration postimage %q has no exact preimage", postimage.Path)
		}
		if _, duplicate := seenPostimages[postimage.Path]; duplicate {
			return fmt.Errorf("duplicate restoration postimage %q", postimage.Path)
		}
		seenPostimages[postimage.Path] = struct{}{}
	}
	return nil
}

func verifyRestoreTransactionState(root string, document PlanDocument) error {
	anchored, err := os.OpenRoot(root)
	if err != nil {
		return err
	}
	defer anchored.Close()
	for _, postimage := range document.Postimages {
		if err := verifyTransactionPostimage(anchored, postimage); err != nil {
			return err
		}
	}
	return nil
}

func classifyRestoreApplyError(err error) error {
	lower := strings.ToLower(err.Error())
	if strings.Contains(lower, "stale") ||
		strings.Contains(lower, "identity does not match") {
		return restoreError(
			SkillsRestoreAction,
			"plan.confirmation.stale",
			"The repository changed after restoration planning.",
			"Preview the current Change Plan and confirm its new planDigest.",
			err,
		)
	}
	code := "restore.apply-failed"
	message := fmt.Sprintf("Restoration apply failed and the transaction restored the exact preimage: %v.", err)
	action := "Fix repository filesystem permissions and rerun the current preview."
	var rollbackErr *IncompleteRollbackError
	if errors.As(err, &rollbackErr) {
		code = "restore.rollback-incomplete"
		message = fmt.Sprintf("Restoration apply failed and rollback is incomplete: %v.", err)
		action = "Recover the Git-private Baseline transaction before retrying."
	}
	return restoreError(SkillsRestoreExecution, code, message, action, err)
}

func failedRestorePayload(profile, setup string, err error) (SkillsRestorePayload, error) {
	payload := SkillsRestorePayload{
		SchemaVersion:  SkillsRestoreSchemaVersion,
		Profile:        profile,
		Acquisitions:   []RestoreAcquisition{},
		Skills:         []RestoreSkill{},
		PlannedChanges: []RestorePlannedChange{},
	}
	if setup != "" {
		payload.Setup = stringPointer(setup)
	}
	var restoreErr *SkillsRestoreError
	if errors.As(err, &restoreErr) {
		payload.Finding = &restoreErr.Finding
	}
	return payload, err
}

func restoreError(
	category SkillsRestoreErrorCategory,
	code, message, action string,
	err error,
) *SkillsRestoreError {
	return &SkillsRestoreError{
		Category: category,
		Finding: RestoreFinding{
			Code: code, Message: message, Action: action,
		},
		Err: err,
	}
}

func portableRestoreDigest(files []restoreFile) string {
	sortRestoreFiles(files)
	digest := sha256.New()
	var length [8]byte
	for _, file := range files {
		binary.BigEndian.PutUint64(length[:], uint64(len(file.Path)))
		_, _ = digest.Write(length[:])
		_, _ = digest.Write([]byte(file.Path))
		binary.BigEndian.PutUint64(length[:], uint64(len(file.Content)))
		_, _ = digest.Write(length[:])
		_, _ = digest.Write(file.Content)
	}
	return hex.EncodeToString(digest.Sum(nil))
}

func externalSkillsLockDigest(files []restoreFile) string {
	sortRestoreFiles(files)
	digest := sha256.New()
	for _, file := range files {
		_, _ = digest.Write([]byte(file.Path))
		_, _ = digest.Write(file.Content)
	}
	return hex.EncodeToString(digest.Sum(nil))
}

func sortRestoreFiles(files []restoreFile) {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})
}

func restoreFilesByPath(files []restoreFile) map[string][]byte {
	result := make(map[string][]byte, len(files))
	for _, file := range files {
		result[file.Path] = file.Content
	}
	return result
}

func optionalBytesDigest(content []byte, exists bool) *string {
	if !exists {
		return nil
	}
	sum := sha256.Sum256(content)
	value := hex.EncodeToString(sum[:])
	return &value
}

func nullableDigest(exists bool, digest string) *string {
	if !exists {
		return nil
	}
	return &digest
}

func stringPointer(value string) *string {
	return &value
}

func hasLockEdits(skills []RestoreSkill) bool {
	for _, skill := range skills {
		if skill.LockEdit != nil {
			return true
		}
	}
	return false
}

func safeRestoreRelative(relative string) bool {
	return utf8.ValidString(relative) &&
		safeRelative(relative) &&
		!strings.HasPrefix(relative, ".git/") &&
		!strings.Contains(relative, "/.git/") &&
		!strings.HasPrefix(relative, "node_modules/") &&
		!strings.Contains(relative, "/node_modules/")
}
