package releaseplan

import (
	"errors"
	"fmt"
)

var (
	ErrMalformedStableVersion = errors.New("malformed stable version")
	ErrPrereleaseVersion      = errors.New("pre-release version")
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
