package spec

import (
	"context"
	"path"

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

func AllAuthorizationOperations() []AuthorizationOperation {
	return authorization.AllAuthorizationOperations()
}

func ReadAuthorization(ctx context.Context, request AuthorizationReadRequest) AuthorizationResolution {
	return authorization.ReadAuthorization(ctx, request)
}

// AuthorizationRecordPath returns the canonical active record for a Spec.
func AuthorizationRecordPath(specSlug string) string {
	return path.Join("docs/specs", specSlug, "_authorization.md")
}

// ReadSpecAuthorization resolves the canonical active record for specSlug.
func ReadSpecAuthorization(ctx context.Context, repoRoot string, specSlug string, revision string) AuthorizationResolution {
	return ReadAuthorization(ctx, AuthorizationReadRequest{
		RepoRoot:   repoRoot,
		RecordPath: AuthorizationRecordPath(specSlug),
		Revision:   revision,
		Role:       AuthorizationRoleSpec,
		AskingSpec: specSlug,
	})
}

// RequireOperation asks an already-resolved grant for one operation.
func RequireOperation(resolution AuthorizationResolution, operation AuthorizationOperation) error {
	return authorization.RequireOperation(resolution, operation)
}
