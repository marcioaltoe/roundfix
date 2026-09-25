package baseline

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

const SkillsReconcileSchemaVersion = "setup-context-driven/reconcile-v1"

const SkillsInstalledTreeRetained = "retained"

var hexadecimalGitCommit = regexp.MustCompile(`^[0-9a-fA-F]{40}$`)

// SkillsReconcileRequest selects one source repository at one immutable commit.
type SkillsReconcileRequest struct {
	Repository       string
	ProfileID        string
	SourceRepository string
	Commit           string
	SourceDir        string
	Confirmation     string
}

// SkillsLockDisposition records why one lock entry is retained or removable.
type SkillsLockDisposition string

const (
	SkillsLockPresent         SkillsLockDisposition = "present"
	SkillsLockMoved           SkillsLockDisposition = "moved"
	SkillsLockObsolete        SkillsLockDisposition = "obsolete"
	SkillsLockRequiredRemoved SkillsLockDisposition = "required-removed"
	SkillsLockUnrelated       SkillsLockDisposition = "unrelated"
)

// SkillsReconcileEntry is one ordered skills-lock.json classification.
type SkillsReconcileEntry struct {
	Skill         string                `json:"skill"`
	Disposition   SkillsLockDisposition `json:"disposition"`
	Source        string                `json:"source,omitempty"`
	SkillPath     string                `json:"skillPath,omitempty"`
	TargetPath    string                `json:"targetPath"`
	InstalledTree string                `json:"installedTree"`
	LockEdit      *RestoreLockEdit      `json:"lockEdit"`
}

// SkillsReconcilePayload is the stable preview and apply result.
type SkillsReconcilePayload struct {
	SchemaVersion  string                 `json:"schemaVersion"`
	OK             bool                   `json:"ok"`
	Applied        bool                   `json:"applied"`
	Profile        string                 `json:"profile"`
	Setup          *string                `json:"setup"`
	Acquisitions   []RestoreAcquisition   `json:"acquisitions"`
	Skills         []SkillsReconcileEntry `json:"skills"`
	PlannedChanges []RestorePlannedChange `json:"plannedChanges"`
	PlanDigest     *string                `json:"planDigest"`
	Finding        *RestoreFinding        `json:"finding,omitempty"`
}

type reconcilePlan struct {
	payload    SkillsReconcilePayload
	document   PlanDocument
	repository string
}

type reconcileDependencies struct {
	profile         restoreProfile
	transactionHook func(transactionFaultPoint) error
}

// ReconcileSkillsLock previews or applies removal of lock entries proven obsolete.
func ReconcileSkillsLock(
	ctx context.Context,
	request SkillsReconcileRequest,
) (SkillsReconcilePayload, error) {
	catalog, err := LoadEmbeddedCatalog()
	if err != nil {
		return failedReconcilePayload(request.ProfileID, "", restoreError(
			SkillsRestoreInvalid,
			"reconcile.assets-invalid",
			fmt.Sprintf("Embedded Baseline assets are invalid: %v.", err),
			"Fix the embedded Baseline catalog before reconciling the skills lock.",
			err,
		))
	}
	profile, err := loadRestoreProfile(catalog, request.ProfileID)
	if err != nil {
		return failedReconcilePayload(request.ProfileID, "", err)
	}
	return reconcileSkillsLock(ctx, request, reconcileDependencies{profile: profile})
}

func reconcileSkillsLock(
	ctx context.Context,
	request SkillsReconcileRequest,
	dependencies reconcileDependencies,
) (SkillsReconcilePayload, error) {
	request.Repository = strings.TrimSpace(request.Repository)
	request.ProfileID = strings.TrimSpace(request.ProfileID)
	request.SourceRepository = strings.TrimSpace(request.SourceRepository)
	request.Commit = strings.TrimSpace(request.Commit)
	if hexadecimalGitCommit.MatchString(request.Commit) {
		request.Commit = strings.ToLower(request.Commit)
	}
	request.SourceDir = strings.TrimSpace(request.SourceDir)
	request.Confirmation = strings.TrimSpace(request.Confirmation)
	if err := validateSkillsReconcileRequest(ctx, request, dependencies.profile); err != nil {
		return failedReconcilePayload(request.ProfileID, dependencies.profile.Setup, err)
	}

	plan, err := buildSkillsReconcilePlan(ctx, request, dependencies.profile)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			err = restoreError(
				SkillsRestoreExecution,
				"reconcile.canceled",
				"Lock reconciliation was canceled before mutation completed.",
				"Rerun roundfix baseline skills reconcile when the operation can complete.",
				err,
			)
		}
		return failedReconcilePayload(request.ProfileID, dependencies.profile.Setup, err)
	}
	if plan.payload.Finding != nil && plan.payload.Finding.Code == "reconcile.required-removed" {
		err := restoreError(
			SkillsRestoreAction,
			plan.payload.Finding.Code,
			plan.payload.Finding.Message,
			plan.payload.Finding.Action,
			errors.New("required profile skill was removed upstream"),
		)
		return plan.payload, err
	}
	if len(plan.payload.PlannedChanges) == 0 {
		plan.payload.OK = true
		return plan.payload, nil
	}
	if plan.payload.PlanDigest == nil || request.Confirmation != *plan.payload.PlanDigest {
		code := "plan.confirmation.stale"
		message := "The supplied confirmation does not match the current lock reconciliation Change Plan."
		if request.Confirmation == "" {
			code = "plan.confirmation.required"
			message = "Lock reconciliation requires confirmation of this exact Change Plan."
		}
		err := restoreError(
			SkillsRestoreAction,
			code,
			message,
			"Review plannedChanges and rerun with --confirm-plan planDigest.",
			errors.New("lock reconciliation plan is not confirmed"),
		)
		plan.payload.Finding = &err.Finding
		return plan.payload, err
	}
	if err := applySkillsRestorePlan(ctx, restorePlan{
		payload:     SkillsRestorePayload{},
		document:    plan.document,
		repository:  plan.repository,
		sourceFiles: map[string][]restoreFile{},
	}, dependencies.transactionHook); err != nil {
		var restoreErr *SkillsRestoreError
		if errors.As(err, &restoreErr) {
			plan.payload.Finding = &restoreErr.Finding
		}
		return plan.payload, err
	}
	plan.payload.OK = true
	plan.payload.Applied = true
	plan.payload.Finding = &RestoreFinding{
		Code:    "reconcile.completed",
		Message: "Obsolete skills-lock.json entries were removed; installed skill trees were retained.",
		Action:  "Run roundfix doctor to inspect the repository's current skill readiness.",
	}
	return plan.payload, nil
}

func validateSkillsReconcileRequest(
	ctx context.Context,
	request SkillsReconcileRequest,
	profile restoreProfile,
) error {
	if ctx == nil {
		return restoreError(
			SkillsRestoreInvalid,
			"reconcile.context-invalid",
			"Lock reconciliation requires a live context.",
			"Rerun roundfix baseline skills reconcile.",
			errors.New("context is required"),
		)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if request.Repository == "" {
		return restoreError(
			SkillsRestoreInvalid,
			"reconcile.repo-invalid",
			"Repository path cannot be empty.",
			"Pass an existing Git worktree with --repo.",
			errors.New("repository path is empty"),
		)
	}
	if request.ProfileID != profile.ID {
		return restoreError(
			SkillsRestoreInvalid,
			"restore.profile-unknown",
			fmt.Sprintf("Unknown built-in Baseline Profile %q.", request.ProfileID),
			"Choose a profile id from the embedded Baseline catalog.",
			errors.New("profile does not match the selected catalog"),
		)
	}
	if request.SourceRepository == "" {
		return restoreError(
			SkillsRestoreInvalid,
			"reconcile.source-invalid",
			"Source repository cannot be empty.",
			"Pass the skills-lock.json source repository with --source.",
			errors.New("source repository is empty"),
		)
	}
	if !hexadecimalGitCommit.MatchString(request.Commit) {
		return restoreError(
			SkillsRestoreInvalid,
			"reconcile.commit-invalid",
			"Reconciliation revision must be an exact 40-hex commit.",
			"Resolve the source revision to an immutable commit and pass it with --revision.",
			errors.New("reconciliation revision is mutable or malformed"),
		)
	}
	if request.Confirmation != "" && !lowercaseSHA256.MatchString(request.Confirmation) {
		return restoreError(
			SkillsRestoreInvalid,
			"plan.confirmation.invalid",
			"Plan confirmation must be a lowercase SHA-256 digest.",
			"Pass the exact planDigest returned by the reconciliation preview.",
			errors.New("confirmation digest is malformed"),
		)
	}
	if request.SourceDir != "" {
		info, err := os.Stat(request.SourceDir)
		if err != nil || !info.IsDir() {
			return restoreError(
				SkillsRestoreInvalid,
				"restore.source-dir-invalid",
				fmt.Sprintf("Offline Git object store is not a directory: %s.", request.SourceDir),
				"Pass an existing Git checkout or bare object store to --source-dir.",
				err,
			)
		}
	}
	return nil
}

func buildSkillsReconcilePlan(
	ctx context.Context,
	request SkillsReconcileRequest,
	profile restoreProfile,
) (reconcilePlan, error) {
	root, identity, err := inspectRepositoryIdentity(ctx, request.Repository, nil)
	if err != nil {
		return reconcilePlan{}, restoreError(
			SkillsRestoreInvalid,
			"reconcile.repo-invalid",
			fmt.Sprintf("Repository root is not a usable Git worktree: %v.", err),
			"Pass an existing Git worktree with --repo.",
			err,
		)
	}
	lock, err := loadSkillsLock(filepath.Join(root, skillsLockPath))
	if err != nil {
		return reconcilePlan{}, err
	}
	provenance := restoreProvenance{
		Provider: "github", Repository: request.SourceRepository, Ref: request.Commit,
	}
	var sourcePaths []string
	if err := withAcquiredRestoreCommit(
		ctx,
		provenance,
		request.SourceDir,
		func(objectStore string) error {
			var inspectErr error
			sourcePaths, inspectErr = readReconcileGitPaths(ctx, objectStore, provenance)
			return inspectErr
		},
	); err != nil {
		return reconcilePlan{}, err
	}

	pathSet := make(map[string]struct{}, len(sourcePaths))
	movedSkills := make(map[string]struct{})
	for _, sourcePath := range sourcePaths {
		pathSet[sourcePath] = struct{}{}
		if path.Base(sourcePath) == "SKILL.md" {
			movedSkills[path.Base(path.Dir(sourcePath))] = struct{}{}
		}
	}
	lockAfter := lock.clone()
	entries := make([]SkillsReconcileEntry, 0, len(lock.skillEntries()))
	plannedChanges := []RestorePlannedChange{}
	requiredRemoved := []string{}
	for _, locked := range lock.skillEntries() {
		source, _ := locked.value.stringField("source")
		skillPath, _ := locked.value.stringField("skillPath")
		disposition := SkillsLockUnrelated
		if source == request.SourceRepository {
			switch {
			case hasStringKey(pathSet, skillPath):
				disposition = SkillsLockPresent
			case hasStringKey(movedSkills, locked.name):
				disposition = SkillsLockMoved
			case hasStringKey(profile.RequiredSkills, locked.name):
				disposition = SkillsLockRequiredRemoved
				requiredRemoved = append(requiredRemoved, locked.name)
			default:
				disposition = SkillsLockObsolete
			}
		}
		entry := SkillsReconcileEntry{
			Skill:         locked.name,
			Disposition:   disposition,
			Source:        source,
			SkillPath:     skillPath,
			TargetPath:    path.Join(".agents/skills", locked.name),
			InstalledTree: SkillsInstalledTreeRetained,
		}
		if disposition == SkillsLockObsolete {
			before, _ := locked.value.interfaceValue().(map[string]any)
			entry.LockEdit = &RestoreLockEdit{
				Action: "remove-lock-entry", Path: skillsLockPath, Skill: locked.name, Before: before,
			}
			lockAfter.removeSkillsEntry(locked.name)
			plannedChanges = append(plannedChanges, RestorePlannedChange{
				Action: "remove-lock-entry", Path: skillsLockPath, Skill: locked.name, Before: before,
			})
		}
		entries = append(entries, entry)
	}
	payload := SkillsReconcilePayload{
		SchemaVersion: SkillsReconcileSchemaVersion,
		Profile:       profile.ID,
		Setup:         stringPointer(profile.Setup),
		Acquisitions: []RestoreAcquisition{{
			Provider: provenance.Provider, Repository: provenance.Repository, Ref: provenance.Ref,
		}},
		Skills: entries,
	}
	if len(requiredRemoved) != 0 {
		for index := range payload.Skills {
			payload.Skills[index].LockEdit = nil
		}
		payload.Finding = &RestoreFinding{
			Code: "reconcile.required-removed",
			Message: fmt.Sprintf(
				"The selected commit removes Profile-required skills: %s.",
				strings.Join(requiredRemoved, ", "),
			),
			Action: "Choose a Profile that does not require them or restore them upstream before reconciling.",
		}
		return reconcilePlan{payload: payload, repository: root}, nil
	}

	lockAfterBytes := lock.before
	if len(plannedChanges) != 0 {
		lockAfterBytes, err = lockAfter.marshalIndent()
		if err != nil {
			return reconcilePlan{}, restoreError(
				SkillsRestoreExecution,
				"lock.write-failed",
				fmt.Sprintf("Could not serialize skills-lock.json: %v.", err),
				"Fix the lock document and preview reconciliation again.",
				err,
			)
		}
	}
	digest, err := skillsReconcilePlanDigest(
		profile,
		payload.Acquisitions,
		entries,
		plannedChanges,
		lock.before,
		lockAfterBytes,
	)
	if err != nil {
		return reconcilePlan{}, err
	}
	payload.PlannedChanges = plannedChanges
	payload.PlanDigest = stringPointer(digest)
	plan := reconcilePlan{payload: payload, repository: root}
	if len(plannedChanges) != 0 {
		plan.document, err = buildSkillsReconcileTransactionDocument(
			root, identity, digest, lockAfterBytes,
		)
		if err != nil {
			return reconcilePlan{}, err
		}
	}
	return plan, nil
}

func readReconcileGitPaths(
	ctx context.Context,
	objectStore string,
	provenance restoreProvenance,
) ([]string, error) {
	output, err := runRestoreGit(
		ctx,
		"--git-dir", objectStore,
		"ls-tree", "-r", "-z", "--full-tree", "--name-only",
		provenance.Ref,
	)
	if err != nil {
		return nil, restoreError(
			SkillsRestoreExecution,
			"source.commit-unavailable",
			fmt.Sprintf("Could not inspect %s@%s: %v.", provenance.Repository, provenance.Ref, err),
			"Verify the exact commit and offline source before retrying.",
			err,
		)
	}
	parts := strings.Split(string(output), "\x00")
	paths := make([]string, 0, len(parts))
	for _, item := range parts {
		if item != "" {
			paths = append(paths, item)
		}
	}
	return paths, nil
}

func skillsReconcilePlanDigest(
	profile restoreProfile,
	acquisitions []RestoreAcquisition,
	entries []SkillsReconcileEntry,
	changes []RestorePlannedChange,
	lockBefore, lockAfter []byte,
) (string, error) {
	payload := map[string]any{
		"kind":             "reconcile-skills-lock",
		"profile":          profile.ID,
		"setup":            profile.Setup,
		"acquisitions":     acquisitions,
		"skills":           entries,
		"plannedChanges":   changes,
		"lockBeforeDigest": optionalBytesDigest(lockBefore, lockBefore != nil),
		"lockAfterDigest":  optionalBytesDigest(lockAfter, lockAfter != nil),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", restoreError(
			SkillsRestoreExecution,
			"reconcile.plan-failed",
			fmt.Sprintf("Could not serialize the lock reconciliation Change Plan: %v.", err),
			"Fix the reconciliation contract before retrying.",
			err,
		)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func buildSkillsReconcileTransactionDocument(
	root string,
	identity RepositoryIdentity,
	digest string,
	lockAfter []byte,
) (PlanDocument, error) {
	anchored, err := os.OpenRoot(root)
	if err != nil {
		return PlanDocument{}, restoreError(
			SkillsRestoreInvalid,
			"reconcile.repo-invalid",
			fmt.Sprintf("Could not anchor the reconciliation repository: %v.", err),
			"Repair the repository and rerun the preview.",
			err,
		)
	}
	defer anchored.Close()
	if err := validatePathParents(anchored, skillsLockPath); err != nil {
		return PlanDocument{}, restoreError(
			SkillsRestoreInvalid,
			"lock.unsafe-path",
			fmt.Sprintf("Reconciliation target %s is unsafe: %v.", skillsLockPath, err),
			"Replace the unsafe target parent before retrying.",
			err,
		)
	}
	state, err := captureTransactionState(anchored, skillsLockPath)
	if err != nil {
		return PlanDocument{}, restoreError(
			SkillsRestoreExecution,
			"lock.read-failed",
			fmt.Sprintf("Could not capture the skills lock preimage: %v.", err),
			"Fix repository permissions and rerun the preview.",
			err,
		)
	}
	return PlanDocument{
		Repository: identity,
		Preimages: []Preimage{
			preimageFromTransactionState(skillsLockPath, state),
		},
		Postimages: []Postimage{{
			Path:            skillsLockPath,
			Kind:            PreimageRegular,
			Mode:            0o600,
			Content:         append([]byte(nil), lockAfter...),
			ContentIdentity: transactionContentIdentity(lockAfter),
		}},
		PlanDigest: "sha256:" + digest,
	}, nil
}

func failedReconcilePayload(
	profile string,
	setup string,
	err error,
) (SkillsReconcilePayload, error) {
	payload := SkillsReconcilePayload{
		SchemaVersion:  SkillsReconcileSchemaVersion,
		Profile:        profile,
		Acquisitions:   []RestoreAcquisition{},
		Skills:         []SkillsReconcileEntry{},
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

func hasStringKey(values map[string]struct{}, key string) bool {
	_, ok := values[key]
	return ok
}
