package spec

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"roundfix/internal/authorization"
)

// These aliases preserve the internal/spec API while the reader lives in a
// lower-level package that spec and suiteguardcontract can both import.
type AuthorizationRole = authorization.AuthorizationRole

const (
	AuthorizationRoleSpec   = authorization.AuthorizationRoleSpec
	AuthorizationRoleLegacy = authorization.AuthorizationRoleLegacy
)

type AuthorizationOutcome = authorization.AuthorizationOutcome

const (
	AuthorizationGranted    = authorization.AuthorizationGranted
	AuthorizationRefused    = authorization.AuthorizationRefused
	AuthorizationUnresolved = authorization.AuthorizationUnresolved
)

type AuthorizationStatus = authorization.AuthorizationStatus

const (
	AuthorizationStatusApproved  = authorization.AuthorizationStatusApproved
	AuthorizationStatusProposed  = authorization.AuthorizationStatusProposed
	AuthorizationStatusWithdrawn = authorization.AuthorizationStatusWithdrawn
)

type AuthorizationOperation = authorization.AuthorizationOperation

const (
	AuthorizationOperationImplement   = authorization.AuthorizationOperationImplement
	AuthorizationOperationCommit      = authorization.AuthorizationOperationCommit
	AuthorizationOperationPush        = authorization.AuthorizationOperationPush
	AuthorizationOperationPullRequest = authorization.AuthorizationOperationPullRequest
	AuthorizationOperationMerge       = authorization.AuthorizationOperationMerge
	AuthorizationOperationRelease     = authorization.AuthorizationOperationRelease
)

type AuthorizationRegeneration = authorization.AuthorizationRegeneration
type AuthorizationSource = authorization.AuthorizationSource
type AuthorizationRecord = authorization.AuthorizationRecord
type AuthorizationReasonCode = authorization.AuthorizationReasonCode

const (
	AuthorizationReasonRole                = authorization.AuthorizationReasonRole
	AuthorizationReasonMalformedRecord     = authorization.AuthorizationReasonMalformedRecord
	AuthorizationReasonStatus              = authorization.AuthorizationReasonStatus
	AuthorizationReasonGranted             = authorization.AuthorizationReasonGranted
	AuthorizationReasonAction              = authorization.AuthorizationReasonAction
	AuthorizationReasonConsuming           = authorization.AuthorizationReasonConsuming
	AuthorizationReasonPaths               = authorization.AuthorizationReasonPaths
	AuthorizationReasonOperations          = authorization.AuthorizationReasonOperations
	AuthorizationReasonRegeneration        = authorization.AuthorizationReasonRegeneration
	AuthorizationReasonContradictory       = authorization.AuthorizationReasonContradictory
	AuthorizationReasonUnreadableRecord    = authorization.AuthorizationReasonUnreadableRecord
	AuthorizationReasonUnavailableRevision = authorization.AuthorizationReasonUnavailableRevision
)

type AuthorizationReason = authorization.AuthorizationReason
type AuthorizationResolution = authorization.AuthorizationResolution
type AuthorizationReadRequest = authorization.AuthorizationReadRequest

// AuthorizationLocation keeps the repository that supplies the authorization
// record separate from the project repository whose delivery it governs.
type AuthorizationLocation struct {
	SpecRepoRoot     string
	SpecRevision     string
	SpecRelativePath string
	ProjectRepoRoot  string
	DeliveryTarget   string
}

func AllAuthorizationOperations() []AuthorizationOperation {
	return authorization.AllAuthorizationOperations()
}

func ReadAuthorization(ctx context.Context, request AuthorizationReadRequest) AuthorizationResolution {
	return authorization.ReadAuthorization(ctx, request)
}

// AuthorizationRecordPath returns the canonical default-root record for a Spec.
func AuthorizationRecordPath(specSlug string) string {
	return path.Join("docs/specs", specSlug, "_authorization.md")
}

// ReadSpecAuthorization resolves the active record from the configured Spec
// Root. A Spec Root in another Git repository supplies its own revision.
func ReadSpecAuthorization(ctx context.Context, projectRepoRoot string, specsRoot string, specSlug string, deliveryTarget string) AuthorizationResolution {
	location, field, err := resolveAuthorizationLocation(ctx, projectRepoRoot, specsRoot, specSlug, deliveryTarget)
	if err != nil {
		return unresolvedSpecAuthorization(specsRoot, specSlug, field, err)
	}
	return ReadAuthorization(ctx, AuthorizationReadRequest{
		RepoRoot:   location.SpecRepoRoot,
		RecordPath: location.SpecRelativePath,
		Revision:   location.SpecRevision,
		Role:       AuthorizationRoleSpec,
		AskingSpec: specSlug,
	})
}

func resolveAuthorizationLocation(ctx context.Context, projectRepoRoot string, specsRoot string, specSlug string, deliveryTarget string) (AuthorizationLocation, string, error) {
	specRepoRoot, specPrefix, specCommonDir, err := authorizationRepositoryLocation(ctx, specsRoot)
	if err != nil {
		return AuthorizationLocation{}, "spec_root", fmt.Errorf("resolve Spec Root %q: %w", specsRoot, err)
	}
	_, _, projectCommonDir, err := authorizationRepositoryLocation(ctx, projectRepoRoot)
	if err != nil {
		return AuthorizationLocation{}, "project_repo_root", fmt.Errorf("resolve project repository root %q: %w", projectRepoRoot, err)
	}

	revision := "HEAD"
	if sameAuthorizationRepository(specCommonDir, projectCommonDir) && strings.TrimSpace(deliveryTarget) != "" {
		revision = strings.TrimSpace(deliveryTarget)
	}
	resolvedRevision, err := authorizationGitOutput(
		ctx,
		specRepoRoot,
		"rev-parse",
		"--verify",
		"--end-of-options",
		revision+"^{commit}",
	)
	if err != nil {
		return AuthorizationLocation{}, "spec_root", fmt.Errorf("resolve Spec Root revision %q: %w", revision, err)
	}

	return AuthorizationLocation{
		SpecRepoRoot:     specRepoRoot,
		SpecRevision:     resolvedRevision,
		SpecRelativePath: path.Join(specPrefix, specSlug, "_authorization.md"),
		ProjectRepoRoot:  projectRepoRoot,
		DeliveryTarget:   deliveryTarget,
	}, "", nil
}

func authorizationRepositoryLocation(ctx context.Context, workDir string) (string, string, string, error) {
	repoRoot, err := authorizationGitOutput(ctx, workDir, "rev-parse", "--show-toplevel")
	if err != nil {
		return "", "", "", err
	}
	prefix, err := authorizationGitOutputAllowEmpty(ctx, workDir, "rev-parse", "--show-prefix")
	if err != nil {
		return "", "", "", err
	}
	commonDir, err := authorizationGitOutput(ctx, workDir, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", "", "", err
	}
	return filepath.Clean(repoRoot), strings.TrimSuffix(filepath.ToSlash(prefix), "/"), canonicalAuthorizationPath(commonDir), nil
}

func authorizationGitOutput(ctx context.Context, workDir string, args ...string) (string, error) {
	return authorizationGitOutputValue(ctx, workDir, false, args...)
}

func authorizationGitOutputAllowEmpty(ctx context.Context, workDir string, args ...string) (string, error) {
	return authorizationGitOutputValue(ctx, workDir, true, args...)
}

func authorizationGitOutputValue(ctx context.Context, workDir string, allowEmpty bool, args ...string) (string, error) {
	gitArgs := append([]string{"-C", workDir, "-c", "core.fsmonitor=false"}, args...)
	command := exec.CommandContext(ctx, "git", gitArgs...)
	command.Env = authorizationGitEnvironment()
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return "", fmt.Errorf("git %s: %s: %w", strings.Join(args, " "), detail, err)
	}
	value := strings.TrimSpace(stdout.String())
	if (!allowEmpty && value == "") || strings.Contains(value, "\n") {
		return "", fmt.Errorf("git %s returned invalid value %q", strings.Join(args, " "), value)
	}
	return value, nil
}

func authorizationGitEnvironment() []string {
	environment := make([]string, 0, len(os.Environ())+1)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "GIT_OPTIONAL_LOCKS=") {
			continue
		}
		environment = append(environment, entry)
	}
	return append(environment, "GIT_OPTIONAL_LOCKS=0")
}

func canonicalAuthorizationPath(value string) string {
	resolved, err := filepath.EvalSymlinks(value)
	if err == nil {
		return filepath.Clean(resolved)
	}
	absolute, err := filepath.Abs(value)
	if err == nil {
		return filepath.Clean(absolute)
	}
	return filepath.Clean(value)
}

func sameAuthorizationRepository(first string, second string) bool {
	return canonicalAuthorizationPath(first) == canonicalAuthorizationPath(second)
}

func unresolvedSpecAuthorization(specsRoot string, specSlug string, field string, err error) AuthorizationResolution {
	return AuthorizationResolution{
		Outcome: AuthorizationUnresolved,
		Record: AuthorizationRecord{
			Role:   AuthorizationRoleSpec,
			Source: AuthorizationSource{Path: path.Join(specSlug, "_authorization.md")},
		},
		Reason: AuthorizationReason{
			Code:   AuthorizationReasonUnreadableRecord,
			Field:  field,
			Value:  specsRoot,
			Detail: err.Error(),
		},
	}
}

// RequireOperation asks an already-resolved grant for one operation.
func RequireOperation(resolution AuthorizationResolution, operation AuthorizationOperation) error {
	return authorization.RequireOperation(resolution, operation)
}
