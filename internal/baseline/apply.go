package baseline

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"sort"
	"strings"
)

// ApplyErrorKind classifies an apply refusal without coupling the Baseline
// engine to CLI exit codes.
type ApplyErrorKind string

const (
	ApplyErrorInvalid      ApplyErrorKind = "invalid"
	ApplyErrorApproval     ApplyErrorKind = "approval"
	ApplyErrorStale        ApplyErrorKind = "stale"
	ApplyErrorExecution    ApplyErrorKind = "execution"
	ApplyErrorVerification ApplyErrorKind = "verification"
)

// ApplyError is one strict apply refusal or failure.
type ApplyError struct {
	Kind       ApplyErrorKind
	NextAction string
	Err        error
}

func (err *ApplyError) Error() string {
	return err.Err.Error()
}

func (err *ApplyError) Unwrap() error {
	return err.Err
}

// ApplyPlan validates and applies exactly one approved portable Plan
// Document. It never recalculates or substitutes another plan.
func ApplyPlan(
	ctx context.Context,
	repository string,
	document PlanDocument,
	confirmation string,
) (Result, error) {
	if ctx == nil {
		return Result{}, applyError(
			ApplyErrorInvalid,
			"provide a live context and rerun roundfix baseline apply",
			errors.New("apply Baseline Plan: context is required"),
		)
	}
	if err := ctx.Err(); err != nil {
		return Result{}, err
	}
	if err := ValidatePlanDocument(document); err != nil {
		return Result{}, applyError(
			ApplyErrorInvalid,
			"generate a new plan with roundfix baseline plan and approve that exact document",
			fmt.Errorf("validate supplied Baseline Plan: %w", err),
		)
	}
	if confirmation != document.PlanDigest {
		return Result{}, applyError(
			ApplyErrorApproval,
			"review the supplied plan and rerun with --confirm-plan "+document.PlanDigest,
			fmt.Errorf(
				"confirmed Plan Digest %q does not match supplied plan %q",
				confirmation,
				document.PlanDigest,
			),
		)
	}

	state, err := inspectApplyPreimage(ctx, repository, document)
	if err != nil {
		return Result{}, err
	}
	if state == applyPreimagePostimage {
		return verifiedApplyResult(document, document.Postimages, true), nil
	}

	transaction, err := BeginTransaction(ctx, repository, document)
	if err != nil {
		// A mixed pre/postimage can be an interrupted transaction. BeginTransaction
		// recovers it before returning; classify any remaining mismatch as stale.
		if strings.Contains(strings.ToLower(err.Error()), "stale") ||
			strings.Contains(err.Error(), "repository identity does not match") {
			return Result{}, applyError(
				ApplyErrorStale,
				"run roundfix baseline plan again and approve the new Plan Digest",
				fmt.Errorf("apply supplied Baseline Plan: %w", err),
			)
		}
		return Result{}, applyError(
			ApplyErrorExecution,
			"resolve the reported transaction failure and rerun roundfix baseline apply",
			fmt.Errorf("begin Baseline apply transaction: %w", err),
		)
	}
	evidence, applyErr := transaction.Apply(ctx)
	closeErr := transaction.Close()
	if applyErr != nil {
		kind := ApplyErrorExecution
		next := "repair the reported failure and rerun roundfix baseline apply"
		if strings.Contains(strings.ToLower(applyErr.Error()), "verify") {
			kind = ApplyErrorVerification
			next = "inspect the recovery result, repair the reported Baseline state, and rerun with a new plan when required"
		}
		if closeErr != nil {
			applyErr = errors.Join(applyErr, closeErr)
		}
		return Result{}, applyError(kind, next, fmt.Errorf("apply supplied Baseline Plan: %w", applyErr))
	}
	if closeErr != nil {
		return Result{}, applyError(
			ApplyErrorExecution,
			"inspect the transaction lock and recovery state before retrying",
			fmt.Errorf("close Baseline apply transaction: %w", closeErr),
		)
	}
	return verifiedApplyResult(document, evidence.VerifiedPostimages, false), nil
}

type applyPreimageState int

const (
	applyPreimageApproved applyPreimageState = iota
	applyPreimageMixed
	applyPreimagePostimage
)

func inspectApplyPreimage(
	ctx context.Context,
	repository string,
	document PlanDocument,
) (applyPreimageState, error) {
	root, identity, err := inspectRepositoryIdentity(ctx, repository, nil)
	if err != nil {
		return applyPreimageApproved, applyError(
			ApplyErrorInvalid,
			"repair the repository and rerun roundfix baseline apply",
			fmt.Errorf("inspect Baseline apply repository: %w", err),
		)
	}
	if identity.Digest != document.Repository.Digest {
		return applyPreimageApproved, applyError(
			ApplyErrorStale,
			"apply the plan only in a matching clone or generate a new plan for this repository",
			errors.New("Baseline Plan repository identity does not match"),
		)
	}
	paths := make([]string, len(document.Preimages))
	for index, preimage := range document.Preimages {
		paths[index] = preimage.Path
	}
	snapshot, err := inspectRepositorySnapshot(root, InventoryRequest{MutablePaths: paths})
	if err != nil {
		return applyPreimageApproved, applyError(
			ApplyErrorInvalid,
			"repair the unsafe repository state and rerun roundfix baseline apply",
			fmt.Errorf("inspect Baseline apply preimage: %w", err),
		)
	}
	if len(snapshot.Blocking) != 0 {
		return applyPreimageApproved, applyError(
			ApplyErrorInvalid,
			"repair every blocking carrier and generate a new Baseline Plan",
			fmt.Errorf("Baseline apply repository has blocking carrier %q: %s",
				snapshot.Blocking[0].Path, snapshot.Blocking[0].Message),
		)
	}

	current := preimagesByPath(snapshot.Preimages)
	postimages := make(map[string]Postimage, len(document.Postimages))
	for _, postimage := range document.Postimages {
		postimages[postimage.Path] = postimage
	}
	allApproved := true
	allPostimages := true
	for _, approved := range document.Preimages {
		observed, ok := current[approved.Path]
		matchesApproved := ok && samePreimage(approved, observed)
		postimage, mutable := postimages[approved.Path]
		if !mutable {
			if !matchesApproved {
				return applyPreimageApproved, applyError(
					ApplyErrorStale,
					"run roundfix baseline plan again and approve the new Plan Digest",
					fmt.Errorf("Baseline Plan preimage is stale at %q", approved.Path),
				)
			}
			continue
		}
		matchesPostimage := preimageMatchesPostimage(observed, postimage)
		allApproved = allApproved && matchesApproved
		allPostimages = allPostimages && matchesPostimage
		if !matchesApproved && !matchesPostimage {
			return applyPreimageApproved, applyError(
				ApplyErrorStale,
				"run roundfix baseline plan again and approve the new Plan Digest",
				fmt.Errorf("Baseline Plan preimage is stale at %q", approved.Path),
			)
		}
	}
	if allPostimages {
		return applyPreimagePostimage, nil
	}
	if allApproved {
		return applyPreimageApproved, nil
	}
	return applyPreimageMixed, nil
}

func preimageMatchesPostimage(preimage Preimage, postimage Postimage) bool {
	if preimage.Path != postimage.Path ||
		preimage.Kind != postimage.Kind ||
		preimage.Mode != postimage.Mode {
		return false
	}
	switch postimage.Kind {
	case PreimageMissing:
		return !preimage.Exists
	case PreimageRegular:
		return preimage.Exists &&
			preimage.Bytes == int64(len(postimage.Content)) &&
			preimage.ContentIdentity == postimage.ContentIdentity
	default:
		return false
	}
}

func verifiedApplyResult(document PlanDocument, verified []Postimage, alreadyApplied bool) Result {
	recommendations := make([]string, 0, len(document.SetupManifest.Verification))
	for _, verification := range document.SetupManifest.Verification {
		if verification.Command != "" {
			recommendations = append(recommendations, verification.Command)
		}
	}
	sort.Strings(recommendations)
	message := "approved Baseline Plan applied and verified"
	if alreadyApplied {
		message = "approved Baseline Plan is already applied and verified"
	}
	return Result{
		SchemaVersion:      ResultSchemaVersion,
		Operation:          "apply",
		State:              "verified",
		Message:            message,
		PlanDigest:         document.PlanDigest,
		VerifiedPostimages: clonePostimages(verified),
		Warnings:           cloneFindings(document.Warnings),
		Recommendations:    recommendations,
	}
}

func clonePostimages(postimages []Postimage) []Postimage {
	cloned := make([]Postimage, len(postimages))
	for index, postimage := range postimages {
		cloned[index] = clonePostimage(postimage)
	}
	return cloned
}

func applyError(kind ApplyErrorKind, next string, err error) error {
	return &ApplyError{Kind: kind, NextAction: next, Err: err}
}

func validatePlanApplyContract(document PlanDocument) error {
	postimages := make(map[string]Postimage, len(document.Postimages))
	preimages := preimagesByPath(document.Preimages)
	entriesByKind := make(map[string][]ManagedEntry)
	for _, postimage := range document.Postimages {
		postimages[postimage.Path] = postimage
	}
	for _, entry := range document.ManagedEntries {
		entriesByKind[entry.Kind] = append(entriesByKind[entry.Kind], entry)
	}

	manifestPostimage, ok := postimages[manifestPath]
	if !ok || manifestPostimage.Kind != PreimageRegular {
		return errors.New("Baseline Plan has no regular Setup Manifest postimage")
	}
	manifestBytes, err := json.MarshalIndent(document.SetupManifest, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize Setup Manifest validation bytes: %w", err)
	}
	manifestBytes = append(manifestBytes, '\n')
	if !bytes.Equal(manifestPostimage.Content, manifestBytes) {
		return errors.New("Setup Manifest postimage does not match the plan identity")
	}

	expectedBackups, err := expectedRootBackups(document.Preimages)
	if err != nil {
		return err
	}
	actualBackups := make(map[string]ManagedEntry, len(entriesByKind["backup"]))
	for _, entry := range entriesByKind["backup"] {
		source := strings.TrimPrefix(entry.ID, "backup:")
		if source == entry.ID || !safeRelative(source) {
			return fmt.Errorf("invalid immutable backup entry %q", entry.ID)
		}
		if _, duplicate := actualBackups[source]; duplicate {
			return fmt.Errorf("duplicate immutable backup for %q", source)
		}
		actualBackups[source] = entry
		expectedPath, expected := expectedBackups[source]
		if !expected || entry.Path != expectedPath {
			return fmt.Errorf("immutable backup %q has no matching root carrier relationship", entry.Path)
		}
		sourcePreimage := preimages[source]
		postimage, exists := postimages[entry.Path]
		if !exists || postimage.Kind != PreimageRegular ||
			postimage.ContentIdentity != sourcePreimage.ContentIdentity ||
			entry.ContentIdentity != sourcePreimage.ContentIdentity {
			return fmt.Errorf("immutable backup %q does not preserve exact source bytes", entry.Path)
		}
		digest := strings.TrimPrefix(postimage.ContentIdentity, "sha256:")
		if len(digest) != sha256HexLength ||
			path.Base(entry.Path) != strings.TrimSuffix(path.Base(expectedPath), ".md")+".md" ||
			!strings.Contains(path.Base(entry.Path), "."+digest+".") {
			return fmt.Errorf("immutable backup %q does not match its full digest name", entry.Path)
		}
		backupPreimage := preimages[entry.Path]
		if backupPreimage.Exists &&
			(backupPreimage.Kind != PreimageRegular ||
				backupPreimage.ContentIdentity != postimage.ContentIdentity) {
			return fmt.Errorf("immutable backup %q collides with different bytes", entry.Path)
		}
	}
	for source := range expectedBackups {
		if _, exists := actualBackups[source]; !exists {
			return fmt.Errorf("root carrier source %q has no immutable backup", source)
		}
	}

	seenRetention := make(map[string]struct{}, len(document.Retention))
	for _, retention := range document.Retention {
		if _, duplicate := seenRetention[retention.FromClause]; duplicate {
			return fmt.Errorf("duplicate retention record for %q", retention.FromClause)
		}
		seenRetention[retention.FromClause] = struct{}{}
		if retention.Disposition == "rejected" && len(retention.Targets) != 0 {
			return fmt.Errorf("rejected retention record %q has targets", retention.FromClause)
		}
		if retention.Disposition != "rejected" && len(retention.Targets) == 0 {
			return fmt.Errorf("retention record %q has no target", retention.FromClause)
		}
	}
	return nil
}

const sha256HexLength = 64

func expectedRootBackups(preimages []Preimage) (map[string]string, error) {
	byPath := preimagesByPath(preimages)
	type carrier struct {
		path   string
		source string
		digest string
		direct bool
	}
	var carriers []carrier
	for _, carrierPath := range []string{"AGENTS.md", "CLAUDE.md"} {
		preimage, ok := byPath[carrierPath]
		if !ok || !preimage.Exists {
			continue
		}
		switch preimage.Kind {
		case PreimageRegular:
			if preimage.ContentIdentity == "" {
				return nil, fmt.Errorf("root carrier %q has no content identity", carrierPath)
			}
			carriers = append(carriers, carrier{
				path: carrierPath, source: carrierPath,
				digest: preimage.ContentIdentity, direct: true,
			})
		case PreimageSymlink:
			source := path.Clean(path.Join(path.Dir(carrierPath), preimage.LinkTarget))
			sourcePreimage, exists := byPath[source]
			if !exists || sourcePreimage.Kind != PreimageRegular ||
				sourcePreimage.ContentIdentity == "" {
				return nil, fmt.Errorf("root carrier %q has no safe bounded source", carrierPath)
			}
			carriers = append(carriers, carrier{
				path: carrierPath, source: source,
				digest: sourcePreimage.ContentIdentity,
			})
		default:
			return nil, fmt.Errorf("root carrier %q has unsupported kind %q", carrierPath, preimage.Kind)
		}
	}
	sort.Slice(carriers, func(i, j int) bool {
		if carriers[i].source != carriers[j].source {
			return carriers[i].source < carriers[j].source
		}
		if carriers[i].direct != carriers[j].direct {
			return carriers[i].direct
		}
		return carriers[i].path < carriers[j].path
	})
	expected := make(map[string]string)
	for _, carrier := range carriers {
		if _, exists := expected[carrier.source]; exists {
			continue
		}
		digest := strings.TrimPrefix(carrier.digest, "sha256:")
		base := strings.TrimSuffix(carrier.path, path.Ext(carrier.path))
		expected[carrier.source] = base + "." + digest + ".md"
	}
	return expected, nil
}

func verifyAppliedPlanState(root string, document PlanDocument) error {
	if err := validatePlanApplyContract(document); err != nil {
		return err
	}
	paths := make([]string, len(document.Preimages))
	for index, preimage := range document.Preimages {
		paths[index] = preimage.Path
	}
	snapshot, err := inspectRepositorySnapshot(root, InventoryRequest{MutablePaths: paths})
	if err != nil {
		return fmt.Errorf("verify applied Baseline repository: %w", err)
	}
	if len(snapshot.Blocking) != 0 {
		return fmt.Errorf("verify applied Baseline repository: blocking carrier %q: %s",
			snapshot.Blocking[0].Path, snapshot.Blocking[0].Message)
	}
	current := preimagesByPath(snapshot.Preimages)
	for _, postimage := range document.Postimages {
		if !preimageMatchesPostimage(current[postimage.Path], postimage) {
			observed := current[postimage.Path]
			return fmt.Errorf(
				"verify applied Baseline postimage %q: got exists=%t kind=%q mode=%#o bytes=%d identity=%q",
				postimage.Path,
				observed.Exists,
				observed.Kind,
				observed.Mode,
				observed.Bytes,
				observed.ContentIdentity,
			)
		}
	}
	return nil
}
