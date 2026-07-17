package releaseplan

import (
	"errors"
	"fmt"
)

var (
	ErrMalformedStableVersion = errors.New("malformed stable version")
	ErrPrereleaseVersion      = errors.New("pre-release version")
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
