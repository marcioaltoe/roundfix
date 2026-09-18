package releaseplan

import (
	"sort"
	"strconv"
	"strings"
)

// Version is a stable semantic version without pre-release or build metadata.
type Version struct {
	major    int
	minor    int
	patch    int
	Prefixed bool // the tag carried a leading "v"
}

// ParseStableVersion parses the supported MAJOR.MINOR.PATCH and
// vMAJOR.MINOR.PATCH release tag forms. It rejects pre-release and malformed
// values without normalization.
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
		return Version{
			major:    major,
			minor:    minor,
			patch:    patch,
			Prefixed: strings.HasPrefix(tag, "v"),
		}, nil
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

// SelectHighestVersion returns the ref with the highest semantic version.
// It refuses when more than one ref reaches that version.
func SelectHighestVersion(refs []VersionRef) (VersionRef, error) {
	if len(refs) == 0 {
		return VersionRef{}, ErrNoStableReleaseTag
	}

	highest := []VersionRef{refs[0]}
	for _, ref := range refs[1:] {
		switch compareStableVersion(ref.Version, highest[0].Version) {
		case 1:
			highest = []VersionRef{ref}
		case 0:
			highest = append(highest, ref)
		}
	}
	if len(highest) > 1 {
		sort.Slice(highest, func(left, right int) bool {
			if highest[left].Tag != highest[right].Tag {
				return highest[left].Tag < highest[right].Tag
			}
			return highest[left].CommitSHA < highest[right].CommitSHA
		})
		return VersionRef{}, AmbiguousHighestVersionError{Refs: highest}
	}
	return highest[0], nil
}

func compareStableVersion(left Version, right Version) int {
	for _, pair := range [][2]int{
		{left.major, right.major},
		{left.minor, right.minor},
		{left.patch, right.patch},
	} {
		if pair[0] > pair[1] {
			return 1
		}
		if pair[0] < pair[1] {
			return -1
		}
	}
	return 0
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
	core := strings.TrimPrefix(tag, "v")
	parts := strings.Split(core, ".")
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
