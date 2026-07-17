package releaseplan

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMalformedStableVersion = errors.New("malformed stable version")
	ErrPrereleaseVersion      = errors.New("pre-release version")
	ErrDirtyWorktree          = errors.New("dirty worktree")
	ErrNoStableReleaseTag     = errors.New("no stable release tag")
	ErrUnresolvedRevision     = errors.New("unresolved revision")
	ErrNonCommitRevision      = errors.New("non-commit revision")
	ErrInvalidReleaseRange    = errors.New("invalid release range")
	ErrManualImpactRequired   = errors.New("manual impact required")
	ErrManualReasonRequired   = errors.New("manual impact reason required")
	ErrManualImpactTooLow     = errors.New("manual impact below automatic minimum")
)

// StableVersionError reports why a release base could not be parsed as the
// supported vMAJOR.MINOR.PATCH stable tag form.
type StableVersionError struct {
	Input      string
	Reason     string
	NextAction string
	Err        error
}

func (err StableVersionError) Error() string {
	return fmt.Sprintf("release base %q: %s; %s", err.Input, err.Reason, err.NextAction)
}

func (err StableVersionError) Unwrap() error {
	return err.Err
}

// UnknownImpactError reports an impact value outside the Release Plan schema.
type UnknownImpactError struct {
	Impact Impact
}

func (err UnknownImpactError) Error() string {
	return fmt.Sprintf("release impact %q is not supported; use one of none, patch, minor, major", err.Impact)
}

// GitSourceError reports a failed local Git release-range operation with the
// next corrective action a CLI can show to the user.
type GitSourceError struct {
	Operation  string
	Ref        string
	Paths      []string
	NextAction string
	Err        error
}

func (err GitSourceError) Error() string {
	message := err.Operation
	if err.Ref != "" {
		message += fmt.Sprintf(" %q", err.Ref)
	}
	if err.Err != nil {
		message += ": " + err.Err.Error()
	}
	if len(err.Paths) > 0 {
		message += fmt.Sprintf(": paths: %s", strings.Join(err.Paths, ", "))
	}
	if err.NextAction != "" {
		message += "; " + err.NextAction
	}
	return message
}

func (err GitSourceError) Unwrap() error {
	return err.Err
}

// ManualImpactError reports a rejected manual classification input.
type ManualImpactError struct {
	Impact     Impact
	Minimum    Impact
	Reason     string
	NextAction string
	Err        error
}

func (err ManualImpactError) Error() string {
	if err.Minimum != "" {
		return fmt.Sprintf("manual impact %q is invalid for automatic minimum %q: %s", err.Impact, err.Minimum, err.NextAction)
	}
	return fmt.Sprintf("manual impact %q is invalid: %s", err.Impact, err.NextAction)
}

func (err ManualImpactError) Unwrap() error {
	return err.Err
}
