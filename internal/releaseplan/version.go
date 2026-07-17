package releaseplan

import (
	"strconv"
	"strings"
)

// Version is a stable semantic version without pre-release or build metadata.
type Version struct {
	major int
	minor int
	patch int
}

// ParseStableVersion parses only the supported vMAJOR.MINOR.PATCH release tag
// form. It rejects pre-release and malformed values without normalization.
func ParseStableVersion(tag string) (Version, error) {
	if parts, ok := stableVersionParts(tag); ok {
		major, err := strconv.Atoi(parts[0])
		if err != nil {
			return Version{}, malformedStableVersion(tag)
		}
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			return Version{}, malformedStableVersion(tag)
		}
		patch, err := strconv.Atoi(parts[2])
		if err != nil {
			return Version{}, malformedStableVersion(tag)
		}
		return Version{major: major, minor: minor, patch: patch}, nil
	}
	if isPrereleaseVersion(tag) {
		return Version{}, StableVersionError{
			Input:      tag,
			Reason:     "pre-release tags are not supported",
			NextAction: "use a stable vMAJOR.MINOR.PATCH tag",
			Err:        ErrPrereleaseVersion,
		}
	}
	return Version{}, malformedStableVersion(tag)
}

func malformedStableVersion(tag string) StableVersionError {
	return StableVersionError{
		Input:      tag,
		Reason:     "expected a stable vMAJOR.MINOR.PATCH tag",
		NextAction: "use a tag like v1.2.3",
		Err:        ErrMalformedStableVersion,
	}
}

func stableVersionParts(tag string) ([]string, bool) {
	if !strings.HasPrefix(tag, "v") {
		return nil, false
	}
	parts := strings.Split(strings.TrimPrefix(tag, "v"), ".")
	if len(parts) != 3 {
		return nil, false
	}
	for _, part := range parts {
		if !isCanonicalNumericIdentifier(part) {
			return nil, false
		}
	}
	return parts, true
}

func isPrereleaseVersion(tag string) bool {
	core, _, hasPrerelease := strings.Cut(tag, "-")
	if !hasPrerelease {
		return false
	}
	_, ok := stableVersionParts(core)
	return ok
}

func isCanonicalNumericIdentifier(value string) bool {
	if value == "" {
		return false
	}
	if len(value) > 1 && value[0] == '0' {
		return false
	}
	for _, char := range value {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

func (version Version) String() string {
	return "v" + strconv.Itoa(version.major) + "." + strconv.Itoa(version.minor) + "." + strconv.Itoa(version.patch)
}

func (version Version) Major() int {
	return version.major
}

func (version Version) Minor() int {
	return version.minor
}

func (version Version) Patch() int {
	return version.patch
}

func (version Version) IncrementPatch() Version {
	return Version{major: version.major, minor: version.minor, patch: version.patch + 1}
}

func (version Version) IncrementMinor() Version {
	return Version{major: version.major, minor: version.minor + 1}
}

func (version Version) IncrementMajor() Version {
	return Version{major: version.major + 1}
}
